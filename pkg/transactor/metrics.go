package transactor

import "github.com/prometheus/client_golang/prometheus"

const (
	defaultNamespace = "mev_commit_oracle"
	subsystem        = "settler"
)

type metrics struct {
	LastConfirmedNonce prometheus.Gauge
	LastUsedNonce      prometheus.Gauge
	LastSentNonce      prometheus.Gauge
	LastConfirmedBlock prometheus.Gauge
}

func newMetrics() *metrics {
	m := &metrics{}
	m.LastConfirmedNonce = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: defaultNamespace,
			Subsystem: subsystem,
			Name:      "last_confirmed_nonce",
			Help:      "Last confirmed nonce",
		},
	)
	m.LastUsedNonce = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: defaultNamespace,
			Subsystem: subsystem,
			Name:      "last_used_nonce",
			Help:      "Last used nonce",
		},
	)
	m.LastConfirmedBlock = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: defaultNamespace,
			Subsystem: subsystem,
			Name:      "last_confirmed_block",
			Help:      "Last confirmed block",
		},
	)
	m.LastSentNonce = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: defaultNamespace,
			Subsystem: subsystem,
			Name:      "last_sent_nonce",
			Help:      "Last sent nonce",
		},
	)
	return m
}

func (m *metrics) Collectors() []prometheus.Collector {
	return []prometheus.Collector{
		m.LastConfirmedNonce,
		m.LastUsedNonce,
		m.LastConfirmedBlock,
		m.LastSentNonce,
	}
}
