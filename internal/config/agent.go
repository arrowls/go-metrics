// Package config содержит конфигурацию для сервера и агента метрик.
package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/caarlos0/env/v7"
)

// AgentConfig содержит конфигурацию для агента сбора метрик.
// Настройки определяют интервалы сбора и отправки метрик, а также параметры подключения к серверу.
type AgentConfig struct {
	// ReportInterval - интервал в секундах для отправки метрик на сервер
	ReportInterval int `env:"REPORT_INTERVAL"`
	// PollInterval - интервал в секундах для сбора метрик с системы
	PollInterval int `env:"POLL_INTERVAL"`
	// ServerEndpoint - адрес сервера для отправки метрик
	ServerEndpoint string `env:"ADDRESS"`
	// Key - ключ для подписи метрик при отправке
	Key string `env:"KEY"`
	// RateLimit - максимальное количество одновременных запросов к серверу
	RateLimit int `env:"RATE_LIMIT"`
}

// NewAgentConfig создает и инициализирует конфигурацию агента.
// Читает настройки из флагов командной строки и переменных окружения.
// Переменные окружения имеют приоритет над флагами командной строки.
func NewAgentConfig() AgentConfig {
	flag.IntVar(&agentConfig.ReportInterval, "r", reportIntervalDefault, "report interval in seconds")
	flag.IntVar(&agentConfig.PollInterval, "p", pollIntervalDefault, "collection interval in seconds")
	flag.StringVar(&agentConfig.ServerEndpoint, "a", serverEndpointDefault, "server endpoint url")
	flag.StringVar(&serverConfig.Key, "k", "", "encoding key")
	flag.IntVar(&agentConfig.RateLimit, "l", 0, "max requests to server")

	flag.Parse()

	if err := env.Parse(&agentConfig); err != nil {
		fmt.Printf("Failed to parse env: %v\n", err)
		os.Exit(1)
	}

	return agentConfig
}
