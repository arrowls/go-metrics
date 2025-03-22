package middleware

import "net/http"

type Middleware func(next http.Handler) http.Handler

func Wrap(handler http.Handler, middlewareList []Middleware) http.Handler {
	for i := len(middlewareList) - 1; i >= 0; i-- {
		handler = middlewareList[i](handler)
	}
	return handler
}
