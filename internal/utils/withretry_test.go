package utils

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWithRetry_SuccessOnFirstTry(t *testing.T) {
	attempts := 0
	retryFunc := func() (bool, error) {
		attempts++
		return false, nil // Успех, не нужно повторять
	}

	err := WithRetry(retryFunc)

	assert.NoError(t, err)
	assert.Equal(t, 1, attempts)
}

func TestWithRetry_SuccessAfterRetries(t *testing.T) {
	attempts := 0
	retryFunc := func() (bool, error) {
		attempts++
		if attempts < 3 {
			return true, errors.New("temporary error")
		}
		return false, nil // Успех на третьей попытке
	}

	err := WithRetry(retryFunc)

	assert.NoError(t, err)
	assert.Equal(t, 3, attempts)
}

func TestWithRetry_NoRetryOnPermanentError(t *testing.T) {
	attempts := 0
	testError := errors.New("permanent error")
	retryFunc := func() (bool, error) {
		attempts++
		return false, testError // Постоянная ошибка, не повторяем
	}

	err := WithRetry(retryFunc)

	assert.Equal(t, testError, err)
	assert.Equal(t, 1, attempts)
}

func TestWithRetry_AllRetriesExhausted(t *testing.T) {
	attempts := 0
	testError := errors.New("temporary error")
	retryFunc := func() (bool, error) {
		attempts++
		return true, testError // Всегда временная ошибка
	}

	err := WithRetry(retryFunc)

	assert.Equal(t, testError, err)
	// Ожидаем 4 попытки: 1 начальная + 3 повторные
	assert.Equal(t, 4, attempts)
}

func TestWithRetry_RetryDelays(t *testing.T) {
	startTime := time.Now()
	attempts := 0
	retryFunc := func() (bool, error) {
		attempts++
		if attempts < 3 {
			return true, errors.New("temporary error")
		}
		return false, nil
	}

	err := WithRetry(retryFunc)

	assert.NoError(t, err)

	elapsed := time.Since(startTime)

	// Проверяем, что прошло достаточно времени для задержек
	// 1 сек + 3 сек = 4 секунды минимум
	assert.GreaterOrEqual(t, elapsed, 4*time.Second)
}

func TestWithRetry_ErrorChanges(t *testing.T) {
	attempts := 0
	testErrors := []error{
		errors.New("error 1"),
		errors.New("error 2"),
		errors.New("error 3"),
		nil,
	}

	retryFunc := func() (bool, error) {
		if attempts < len(testErrors) {
			err := testErrors[attempts]
			attempts++
			if err == nil {
				return false, nil
			}
			return true, err
		}
		return false, errors.New("unexpected")
	}

	err := WithRetry(retryFunc)

	assert.NoError(t, err)
	assert.Equal(t, 4, attempts)
}

func TestWithRetry_RetryFlagChanges(t *testing.T) {
	attempts := 0
	retryFlags := []bool{true, true, false}

	retryFunc := func() (bool, error) {
		if attempts < len(retryFlags) {
			shouldRetry := retryFlags[attempts]
			attempts++
			return shouldRetry, errors.New("test error")
		}
		return false, nil
	}

	err := WithRetry(retryFunc)

	assert.Error(t, err)
	assert.Equal(t, 3, attempts)
}
