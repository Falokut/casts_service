package config

import (
	"sync"
	"time"

	"github.com/Falokut/casts_service/internal/repository"
	"github.com/Falokut/casts_service/pkg/jaeger"
	"github.com/Falokut/casts_service/pkg/logging"
	"github.com/Falokut/casts_service/pkg/metrics"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	LogLevel        string `yaml:"log_level" env:"LOG_LEVEL"`
	HealthcheckPort string `yaml:"healthcheck_port" env:"HEALTHCHECK_PORT"`
	Listen          struct {
		Host string `yaml:"host" env:"HOST"`
		Port string `yaml:"port" env:"PORT"`
		Mode string `yaml:"server_mode" env:"SERVER_MODE"` // support GRPC, REST, BOTH
	} `yaml:"listen"`

	PrometheusConfig struct {
		Name         string                      `yaml:"service_name" ENV:"PROMETHEUS_SERVICE_NAME"`
		ServerConfig metrics.MetricsServerConfig `yaml:"server_config"`
	} `yaml:"prometheus"`

	CastsCache struct {
		Network  string        `yaml:"network" env:"CASTS_CACHE_NETWORK"`
		Addr     string        `yaml:"addr" env:"CASTS_CACHE_ADDR"`
		Password string        `yaml:"password" env:"CASTS_CACHE_PASSWORD"`
		DB       int           `yaml:"db" env:"CASTS_CACHE_DB"`
		CastTTL  time.Duration `yaml:"cast_ttl"`
	} `yaml:"casts_cache"`

	ProfessionsCache struct {
		Network        string        `yaml:"network" env:"PROFESSIONS_CACHE_NETWORK"`
		Addr           string        `yaml:"addr" env:"PROFESSIONS_CACHE_ADDR"`
		Password       string        `yaml:"password" env:"PROFESSIONS_CACHE_PASSWORD"`
		DB             int           `yaml:"db" env:"PROFESSIONS_CACHE_DB"`
		ProfessionsTTL time.Duration `yaml:"professions_ttl"`
	} `yaml:"professions_cache"`

	DBConfig     repository.DBConfig `yaml:"db_config"`
	JaegerConfig jaeger.Config       `yaml:"jaeger"`
}

var instance *Config
var once sync.Once

const configsPath = "configs/"

func GetConfig() *Config {
	once.Do(func() {
		logger := logging.GetLogger()
		instance = &Config{}

		if err := cleanenv.ReadConfig(configsPath+"config.yml", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			logger.Fatal(help, " ", err)
		}
	})

	return instance
}
