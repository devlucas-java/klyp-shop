package http

import (
	"fmt"
	"net/http"

	"github.com/devlucas-java/klyp-shop/pkg/logger"
)

// Server encapsula o servidor HTTP e suas dependências de ciclo de vida.
type Server struct {
	port   string
	router http.Handler
	log    *logger.Logger
}

// NewServer cria um novo Server com o router e porta fornecidos.
func NewServer(port string, router http.Handler, log *logger.Logger) *Server {
	return &Server{
		port:   port,
		router: router,
		log:    log,
	}
}

// Run inicia o servidor HTTP e bloqueia até receber um erro fatal.
func (s *Server) Run() error {
	addr := fmt.Sprintf(":%s", s.port)
	s.log.Infof("Server is running on port %s", s.port)
	return http.ListenAndServe(addr, s.router)
}
