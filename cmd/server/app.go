package main

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"

	"github.com/Falokut/casts_service/internal/config"
	"github.com/Falokut/casts_service/internal/repository"
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
		if err := metrics.RunMetricServer(cfg.PrometheusConfig.ServerConfig); err != nil {
			logger.Errorf("Shutting down, error while running metrics server %v", err)
			shutdown <- err
			return
		}
	}()

	logger.Info("Database initializing")
	database, err := repository.NewPostgreDB(cfg.DBConfig)
	if err != nil {
		logger.Errorf("Shutting down, connection to the database is not established: %s", err.Error())
		return
	}

	logger.Info("Repository initializing")
	repo := repository.NewCastsRepository(database)
	defer repo.Shutdown()

	castsCache, err := repository.NewCastsCache(logger.Logger, getCastsCacheOptions(cfg))
	if err != nil {
		logger.Errorf("Shutting down, connection to the cache is not established: %s", err.Error())
		return
	}
	defer castsCache.Shutdown()

	professionsCache, err := repository.NewProfessionsCache(logger.Logger,
		getProfessionsCacheOptions(cfg))
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
			return
		}
	}()

	repoManager := repository.NewCastsRepositoryManager(logger.Logger, repo,
		castsCache, cfg.CastsCache.CastTTL, professionsCache, cfg.ProfessionsCache.ProfessionsTTL, metric)
	logger.Info("Service initializing")
	service := service.NewCastsService(logger.Logger, repoManager)
	logger.Info("Server initializing")
	s := server.NewServer(logger.Logger, service)
	go func() {
		if err := s.Run(getListenServerConfig(cfg), metric, nil, nil); err != nil {
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

	s.Shutdown()
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

			return casts_service.RegisterCastsServiceV1HandlerServer(context.Background(),
				mux, serv)
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
