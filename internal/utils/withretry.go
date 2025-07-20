// Package utils содержит вспомогательные функции для приложения.
// Включает утилиты для повторных попыток, логирования и других общих операций.
package utils

import (
	"time"
)

// WithRetry выполняет функцию с повторными попытками при ошибках.
// Принимает функцию, которая возвращает флаг необходимости повтора и ошибку.
// Использует фиксированные задержки: 1, 3, 5 секунд между попытками.
// Возвращает последнюю ошибку если все попытки неудачны.
func WithRetry(retryFunc func() (bool, error)) error {
	var retryDelays = []time.Duration{
		1 * time.Second,
		3 * time.Second,
		5 * time.Second,
	}

	shouldRetry, err := retryFunc()
	if err == nil {
		return nil
	}

	if !shouldRetry {
		return err
	}

	for _, delay := range retryDelays {
		time.Sleep(delay)
		shouldRetry, err = retryFunc()
		if err == nil {
			return nil
		}

		if !shouldRetry {
			return err
		}

	}

	return err
}
