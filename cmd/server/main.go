package main

import (
	"net/http"

	"github.com/devlucas-java/klyp-shop/configs"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {

	cfg := configs.InitConfigDev()
	_ = configs.InitDBDev()

	r := chi.NewRouter()

	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)

	err := http.ListenAndServe(":"+cfg.GetWebServerPort(), r)
	if err != nil {
		panic(err)
	}

}
