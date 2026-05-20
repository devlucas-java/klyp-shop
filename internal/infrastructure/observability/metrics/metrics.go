package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

// Metric agrupa todos os collectors do Prometheus registrados no registry da aplicação.
type Metric struct {
	// HTTP
	ApiRequests *prometheus.CounterVec
	ApiLatency  *prometheus.HistogramVec

	// Negócio — pedidos
	OrdersCreated   prometheus.Counter
	OrdersCancelled prometheus.Counter

	// Negócio — pagamentos
	PaymentsCreated prometheus.Counter
	PaymentsSettled prometheus.Counter
	PaymentsFailed  prometheus.Counter

	// WebSocket
	WebSocketConnections prometheus.Gauge
}

func NewMetric(reg prometheus.Registerer) *Metric {
	m := &Metric{
		// ── HTTP ──────────────────────────────────────────────────────────────
		ApiRequests: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "klyp_shop",
				Name:      "api_requests_total",
				Help:      "Total de requisições HTTP recebidas.",
			},
			[]string{"route", "method", "status"},
		),
		ApiLatency: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: "klyp_shop",
				Name:      "api_request_latency_seconds",
				Help:      "Latência das requisições HTTP em segundos.",
				Buckets:   prometheus.DefBuckets,
			},
			[]string{"route", "method"},
		),

		// ── Pedidos ───────────────────────────────────────────────────────────
		OrdersCreated: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "klyp_shop",
			Name:      "orders_created_total",
			Help:      "Total de pedidos criados.",
		}),
		OrdersCancelled: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "klyp_shop",
			Name:      "orders_cancelled_total",
			Help:      "Total de pedidos cancelados.",
		}),

		// ── Pagamentos ────────────────────────────────────────────────────────
		PaymentsCreated: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "klyp_shop",
			Name:      "payments_created_total",
			Help:      "Total de invoices Bitcoin criadas.",
		}),
		PaymentsSettled: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "klyp_shop",
			Name:      "payments_settled_total",
			Help:      "Total de pagamentos confirmados (InvoiceSettled).",
		}),
		PaymentsFailed: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "klyp_shop",
			Name:      "payments_failed_total",
			Help:      "Total de pagamentos expirados ou inválidos.",
		}),

		// ── WebSocket ─────────────────────────────────────────────────────────
		WebSocketConnections: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "klyp_shop",
			Name:      "websocket_connections_active",
			Help:      "Número de conexões WebSocket ativas no hub de chat.",
		}),
	}

	// Collectors padrão do Go runtime e processo
	reg.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)

	// Collectors da aplicação
	reg.MustRegister(
		m.ApiRequests,
		m.ApiLatency,
		m.OrdersCreated,
		m.OrdersCancelled,
		m.PaymentsCreated,
		m.PaymentsSettled,
		m.PaymentsFailed,
		m.WebSocketConnections,
	)

	return m
}
