package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/devlucas-java/klyp-shop/internal/infrastructure/observability/metrics"
	"github.com/go-chi/chi"
)

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

// RecordMetricsMiddleware registra contagem e latência de cada request usando
// o padrão de rota do chi (ex: /api/v1/product/{id}) para evitar cardinalidade
// alta com IDs dinâmicos nos labels.
func RecordMetricsMiddleware(m *metrics.Metric) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Não instrumenta o próprio endpoint de scraping
			if r.URL.Path == "/metrics" {
				next.ServeHTTP(w, r)
				return
			}

			start := time.Now()
			rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}
			next.ServeHTTP(rw, r)
			duration := time.Since(start).Seconds()

			// Usa o padrão da rota registrada no chi para evitar cardinalidade alta.
			// Ex: "/api/v1/product/abc123" vira "/api/v1/product/{productID}"
			route := chi.RouteContext(r.Context()).RoutePattern()
			if route == "" {
				route = r.URL.Path
			}

			statusStr := strconv.Itoa(rw.status)

			m.ApiRequests.WithLabelValues(route, r.Method, statusStr).Inc()
			m.ApiLatency.WithLabelValues(route, r.Method).Observe(duration)
		})
	}
}
