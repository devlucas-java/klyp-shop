package adapter

import (
	"net/http"

	"github.com/devlucas-java/klyp-shop/internal/delivery/http/response"
	"github.com/devlucas-java/klyp-shop/pkg/logger"
)

type AdapterHandler func(w http.ResponseWriter, r *http.Request) error

type Adapter struct {
	log *logger.Logger
}

func NewAdapter(log *logger.Logger) *Adapter {
	return &Adapter{log: log}
}

func (a *Adapter) Adapt(h AdapterHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			response.ResponseError(w, r, err, a.log)
		}
	}
}
