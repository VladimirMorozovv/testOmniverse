package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"service_api/internal/config"
	"time"
)

type Status string

const (
	OkStatus Status = "ok"
)

var (
	// Количество обработанных/не обработанных входящих запросов
	requestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name:      "total_requests",
		Namespace: config.Namespace,
		Help:      "Number of processed incoming requests",
	}, []string{"status", "path"})

	requestDurationSeconds = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:      "duration_request",
		Namespace: config.Namespace,
		Help:      "Total request processing time",
	}, []string{"path"})

	postgresQueryDuration = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:      "duration_postgres_query",
		Namespace: config.Namespace,
		Help:      "Postgres query execution time",
	})
)

func IncRequestsTotal(status Status, path string) {
	requestsTotal.With(prometheus.Labels{"status": string(status), "path": path}).Inc()
}
func ObserveRequestDurationSeconds(path string) func() {
	start := time.Now()
	return func() {
		duration := time.Now().Sub(start).Seconds()
		requestDurationSeconds.With(prometheus.Labels{"path": path}).Observe(duration)
	}
}
func ObservePostgresQueryDuration() func() {
	start := time.Now()
	return func() {
		duration := time.Now().Sub(start).Seconds()
		postgresQueryDuration.Observe(duration)
	}
}
