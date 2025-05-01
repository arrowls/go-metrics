package controller

import (
	"net/http"
)

type Metric interface {
	HandleNew(rw http.ResponseWriter, r *http.Request)
	HandleItem(rw http.ResponseWriter, r *http.Request)
	HandleNewFromBody(rw http.ResponseWriter, r *http.Request)
	HandleGetItemFromBody(rw http.ResponseWriter, r *http.Request)
	HandleCreateBatch(rw http.ResponseWriter, r *http.Request)
}

type Public interface {
	HandlePublic(rw http.ResponseWriter, r *http.Request)
	HandleIndex(rw http.ResponseWriter, r *http.Request)
	Ping(rw http.ResponseWriter, r *http.Request)
}

type ErrorHandler interface {
	Handle(w http.ResponseWriter, err error)
}

type Controller struct {
	Metric Metric
	Public Public
}
