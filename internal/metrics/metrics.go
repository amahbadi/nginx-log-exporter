package metrics

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Metrics struct {
	UpstreamDuration *prometheus.HistogramVec
}

func NewMetrics(buckets []float64) *Metrics {
	upstreamDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "upstream_response_duration_seconds",
			Help:    "Duration of upstream responses by URI",
			Buckets: buckets,
		},
		[]string{"uri"},
	)

	prometheus.MustRegister(upstreamDuration)

	return &Metrics{
		UpstreamDuration: upstreamDuration,
	}
}

func ServeMetrics(port string) {
	http.Handle("/metrics", promhttp.Handler())
	log.Printf("Starting metrics server on :%s/metrics", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func (m *Metrics) ObserveUpstreamDuration(uri string, duration float64) {
	m.UpstreamDuration.WithLabelValues(uri).Observe(duration)
}
