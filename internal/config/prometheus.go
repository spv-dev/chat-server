package config

import (
	"errors"
	"net"
	"os"
	"strconv"

	"github.com/spv-dev/chat-server/internal/serviceerror"
)

const (
	prometheusHostEnvName          = "PROMETHEUS_HOST"
	prometheusPortEnvName          = "PROMETHEUS_PORT"
	prometheusHeaderTimeoutEnvName = "PROMETHEUS_TIMEOUT"
)

// PrometheusConfig интерфейс для конфигурации gRPC
type PrometheusConfig interface {
	Address() string
	HeaderTimeout() int64
}

type prometheusConfig struct {
	host          string
	port          string
	headerTimeout int64
}

// NewPrometheusConfig получает конфигурацию gRPC
func NewPrometheusConfig() (*prometheusConfig, error) {
	host := os.Getenv(prometheusHostEnvName)
	if len(host) == 0 {
		return nil, errors.New(serviceerror.PrometheusHostNotFound)
	}

	port := os.Getenv(prometheusPortEnvName)
	if len(port) == 0 {
		return nil, errors.New(serviceerror.PrometheusPortNotFound)
	}

	sTimeout := os.Getenv(prometheusHeaderTimeoutEnvName)
	if len(sTimeout) == 0 {
		return nil, errors.New(serviceerror.PrometheusPortNotFound)
	}

	timeoutStr := os.Getenv(prometheusHeaderTimeoutEnvName)
	if len(timeoutStr) == 0 {
		return nil, errors.New(serviceerror.PrometheusTimeoutNotFound)
	}

	timeout, err := strconv.ParseInt(timeoutStr, 10, 64)
	if err != nil {
		return nil, errors.New(serviceerror.FailedToConvertPrometheusTimeout)
	}

	return &prometheusConfig{
		host:          host,
		port:          port,
		headerTimeout: timeout,
	}, nil
}

// Address возвращает адрес сервиса
func (cfg *prometheusConfig) Address() string {
	return net.JoinHostPort(cfg.host, cfg.port)
}

// Address возвращает адрес сервиса
func (cfg *prometheusConfig) HeaderTimeout() int64 {
	return cfg.headerTimeout
}
