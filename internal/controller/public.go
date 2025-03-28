package controller

import (
	"html/template"
	"net/http"

	"github.com/arrowls/go-metrics/internal/service"
)

type PublicController struct {
	service *service.Service
}

func NewPublicController(service *service.Service) *PublicController {
	return &PublicController{service}
}

func (c *PublicController) HandlePublic(w http.ResponseWriter, r *http.Request) {
	http.FileServer(http.Dir("./frontend/dist")).ServeHTTP(w, r)
}

func (c *PublicController) HandleIndex(w http.ResponseWriter, r *http.Request) {
	// SSR?
	tmpl := template.Must(template.ParseFiles("./frontend/dist/index.html"))
	data := c.service.Metric.GetList()

	err := tmpl.Execute(w, *data)

	if err != nil {
		// page 500
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
