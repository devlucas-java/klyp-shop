package adapter

import (
	"net/http"

	"github.com/devlucas-java/klyp-shop/internal/delivery/http/response"
)

type AdapterHandler func(w http.ResponseWriter, r *http.Request) error

func Adapt(h AdapterHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		err := h(w, r)

		if err != nil {
			response.ResponseError(w, err)
			return
		}
	}
}
