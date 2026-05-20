package mocks

import (
	"github.com/devlucas-java/klyp-shop/internal/infrastructure/observability/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

// NewTestMetric retorna um *metrics.Metric com registry isolado, sem registrar
// no registry global do Prometheus. Seguro para uso em testes paralelos.
func NewTestMetric() *metrics.Metric {
	reg := prometheus.NewRegistry()
	return metrics.NewMetric(reg)
}
