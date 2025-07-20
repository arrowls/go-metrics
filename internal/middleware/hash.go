// Package middleware содержит HTTP middleware компоненты.
package middleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"

	"github.com/arrowls/go-metrics/internal/apperrors"
	"github.com/arrowls/go-metrics/internal/config"
	"github.com/sirupsen/logrus"
)

// responseWriter оборачивает http.ResponseWriter для перехвата тела ответа.
// Позволяет получить тело ответа для вычисления хеша подписи.
type responseWriter struct {
	http.ResponseWriter
	// body - накопленное тело ответа для вычисления хеша
	body []byte
}

// Write перехватывает запись в ResponseWriter и сохраняет данные в body.
// Позволяет получить полное тело ответа для подписи.
func (rw *responseWriter) Write(b []byte) (int, error) {
	rw.body = append(rw.body, b...)
	return rw.ResponseWriter.Write(b)
}

// NewHashingMiddleware создает middleware для проверки и добавления хеш-подписи.
// Проверяет HMAC-SHA256 подпись входящих запросов и добавляет подпись к ответам.
// Использует ключ из конфигурации сервера для вычисления подписи.
// Если ключ не задан, middleware пропускает запросы без изменений.
func NewHashingMiddleware(serverConfig config.ServerConfig, logger *logrus.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if serverConfig.Key == "" {
				next.ServeHTTP(w, r)
				return
			}

			bodyBytes, err := io.ReadAll(r.Body)
			defer func() {
				errClose := r.Body.Close()
				if errClose != nil {
					logger.Errorf("error closing request body: %v", err)
				}
			}()

			r.Body = io.NopCloser(bytes.NewReader(bodyBytes))

			if err != nil {
				logger.Errorf("error reading request body: %v", err)
				return
			}
			// если есть тело запроса
			if len(bodyBytes) != 0 {
				logger.Info("begin checking hash from request body")

				hash := r.Header.Get(config.HashHeaderName)
				if hash == "" {
					next.ServeHTTP(w, r)
					return
				}

				hasher := hmac.New(sha256.New, []byte(serverConfig.Key))
				hasher.Write(bodyBytes)
				sum := hex.EncodeToString(hasher.Sum(nil))

				if sum != hash {
					w.Header().Set("Content-Type", "application/json")
					_, err = w.Write(apperrors.ErrorResponse("invalid hash"))
					if err != nil {
						logger.Errorf("error writing error response: %v", err)
					}
					w.WriteHeader(http.StatusBadRequest)
					return
				}
			}

			writer := &responseWriter{
				w,
				[]byte{},
			}
			next.ServeHTTP(writer, r)

			if len(writer.body) > 0 {
				hasher := hmac.New(sha256.New, []byte(serverConfig.Key))
				hasher.Write(writer.body)
				sum := hex.EncodeToString(hasher.Sum(nil))
				w.Header().Set(config.HashHeaderName, sum)
			}
		})
	}
}
