package main

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"

	"github.com/Falokut/casts_service/internal/config"
	"github.com/Falokut/casts_service/internal/repository"
	"github.com/Falokut/casts_service/internal/repository/postgresrepository"
	"github.com/Falokut/casts_service/internal/repository/rediscache"

	"github.com/Falokut/casts_service/internal/handler"
	"github.com/Falokut/casts_service/internal/service"

	casts_service "github.com/Falokut/casts_service/pkg/casts_service/v1/protos"
	jaegerTracer "github.com/Falokut/casts_service/pkg/jaeger"
	"github.com/Falokut/casts_service/pkg/metrics"
	server "github.com/Falokut/grpc_rest_server"
	"github.com/Falokut/healthcheck"
	logging "github.com/Falokut/online_cinema_ticket_office.loggerwrapper"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/opentracing/opentracing-go"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

func main() {
	logging.NewEntry(logging.ConsoleOutput)
	logger := logging.GetLogger()

	cfg := config.GetConfig()
	logLevel, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Logger.SetLevel(logLevel)

	tracer, closer, err := jaegerTracer.InitJaeger(cfg.JaegerConfig)
	if err != nil {
		logger.Errorf("Shutting down, error while creating tracer %v", err)
		return
	}
	logger.Info("Jaeger connected")
	defer closer.Close()

	opentracing.SetGlobalTracer(tracer)

	logger.Info("Metrics initializing")
	metric, err := metrics.CreateMetrics(cfg.PrometheusConfig.Name)
	if err != nil {
		logger.Errorf("Shutting down, error while creating metrics %v", err)
		return
	}

	shutdown := make(chan error, 1)
	go func() {
		logger.Info("Metrics server running")
		if serr := metrics.RunMetricServer(cfg.PrometheusConfig.ServerConfig); serr != nil {
			logger.Errorf("Shutting down, error while running metrics server %v", serr)
			shutdown <- serr
			return
		}
	}()

	logger.Info("Database initializing")
	database, err := postgresrepository.NewPostgreDB(&cfg.DBConfig)
	if err != nil {
		logger.Errorf("Shutting down, connection to the database is not established: %s", err.Error())
		return
	}

	logger.Info("Repository initializing")
	repo := postgresrepository.NewCastsRepository(database, logger.Logger)
	defer repo.Shutdown()

	castsCache, err := rediscache.NewCastsCache(logger.Logger, getCastsCacheOptions(cfg), metric)
	if err != nil {
		logger.Errorf("Shutting down, connection to the cache is not established: %s", err.Error())
		return
	}
	defer castsCache.Shutdown()

	professionsCache, err := rediscache.NewProfessionsCache(logger.Logger,
		getProfessionsCacheOptions(cfg), metric)
	if err != nil {
		logger.Errorf("Shutting down, connection to the professions cache is not established: %s", err.Error())
		return
	}
	defer professionsCache.Shutdown()

	logger.Info("Healthcheck initializing")
	healthcheckManager := healthcheck.NewHealthManager(logger.Logger,
		[]healthcheck.HealthcheckResource{database, castsCache, professionsCache}, cfg.HealthcheckPort, nil)
	go func() {
		logger.Info("Healthcheck server running")
		if err := healthcheckManager.RunHealthcheckEndpoint(); err != nil {
			logger.Errorf("Shutting down, error while running healthcheck endpoint %s", err.Error())
			shutdown <- err
		}
	}()

	repoManager := repository.NewCastsRepository(
		logger.Logger,
		repo,
		castsCache,
		cfg.CastsCache.CastTTL,
		professionsCache,
		cfg.ProfessionsCache.ProfessionsTTL)
	logger.Info("Service initializing")
	s := service.NewCastsService(logger.Logger, repoManager)
	logger.Info("Handler initializing")
	h := handler.NewCastsServiceHandler(s)

	logger.Info("Server initializing")
	serv := server.NewServer(logger.Logger, h)
	go func() {
		if err := serv.Run(getListenServerConfig(cfg), metric, nil, nil); err != nil {
			logger.Errorf("Shutting down, error while running server %s", err.Error())
			shutdown <- err
			return
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGTERM)

	select {
	case <-quit:
		break
	case <-shutdown:
		break
	}

	serv.Shutdown()
}

func getListenServerConfig(cfg *config.Config) server.Config {
	return server.Config{
		Mode:        cfg.Listen.Mode,
		Host:        cfg.Listen.Host,
		Port:        cfg.Listen.Port,
		ServiceDesc: &casts_service.CastsServiceV1_ServiceDesc,
		RegisterRestHandlerServer: func(ctx context.Context, mux *runtime.ServeMux, service any) error {
			serv, ok := service.(casts_service.CastsServiceV1Server)
			if !ok {
				return errors.New("can't convert")
			}

			return casts_service.RegisterCastsServiceV1HandlerServer(ctx, mux, serv)
		},
	}
}

func getCastsCacheOptions(cfg *config.Config) *redis.Options {
	return &redis.Options{
		Network:  cfg.CastsCache.Network,
		Addr:     cfg.CastsCache.Addr,
		Password: cfg.CastsCache.Password,
		DB:       cfg.CastsCache.DB,
	}
}

func getProfessionsCacheOptions(cfg *config.Config) *redis.Options {
	return &redis.Options{
		Network:  cfg.ProfessionsCache.Network,
		Addr:     cfg.ProfessionsCache.Addr,
		Password: cfg.ProfessionsCache.Password,
		DB:       cfg.ProfessionsCache.DB,
	}
}
