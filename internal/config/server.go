// Package config содержит конфигурацию для сервера и агента метрик.
package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/caarlos0/env/v7"
)

const (
	reportIntervalDefault  = 10
	pollIntervalDefault    = 2
	serverEndpointDefault  = "localhost:8080"
	storeIntervalDefault   = 300
	storageFilePathDefault = "metrics.json"
	restoreDefault         = false
)

// ServerConfig содержит конфигурацию для HTTP сервера метрик.
// Настройки могут быть заданы через переменные окружения или флаги командной строки.
type ServerConfig struct {
	// ServerEndpoint - адрес и порт для прослушивания HTTP запросов
	ServerEndpoint string `env:"ADDRESS"`
	// StoreInterval - интервал в секундах для сохранения метрик в файл
	StoreInterval int `env:"STORE_INTERVAL"`
	// StorageFilePath - путь к файлу для сохранения метрик
	StorageFilePath string `env:"FILE_STORAGE_PATH"`
	// Restore - флаг для восстановления метрик из файла при запуске
	Restore bool `env:"RESTORE"`
	// DatabaseDSN - строка подключения к базе данных PostgreSQL
	DatabaseDSN string `env:"DATABASE_DSN"`
	// Key - ключ для подписи и проверки целостности метрик
	Key string `env:"KEY"`
}

var serverConfig ServerConfig
var agentConfig AgentConfig

// newServerConfig создает и инициализирует конфигурацию сервера.
// Читает настройки из флагов командной строки и переменных окружения.
// Переменные окружения имеют приоритет над флагами командной строки.
func newServerConfig() ServerConfig {
	flag.StringVar(&serverConfig.ServerEndpoint, "a", serverEndpointDefault, "server endpoint url")
	flag.IntVar(&serverConfig.StoreInterval, "i", storeIntervalDefault, "interval to write metrics to file")
	flag.StringVar(&serverConfig.StorageFilePath, "f", storageFilePathDefault, "file to write metrics backup")
	flag.BoolVar(&serverConfig.Restore, "r", restoreDefault, "restore on startup")
	flag.StringVar(&serverConfig.DatabaseDSN, "d", "", "database connection string")
	flag.StringVar(&serverConfig.Key, "k", "", "encoding key")

	flag.Parse()

	if err := env.Parse(&serverConfig); err != nil {
		fmt.Printf("Failed to parse env: %v\n", err)
		os.Exit(1)
	}

	return serverConfig
}
