package updater

import "github.com/prometheus/client_golang/prometheus"

const (
	defaultNamespace = "mev_commit_oracle"
	subsystem        = "updater"
)

type metrics struct {
	CommitmentsCount          prometheus.Counter
	RewardsCount              prometheus.Counter
	SlashesCount              prometheus.Counter
	EncryptedCommitmentsCount prometheus.Counter
}

func newMetrics() *metrics {
	m := &metrics{}
	m.CommitmentsCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: defaultNamespace,
			Subsystem: subsystem,
			Name:      "commitments_count",
			Help:      "Number of commitments processed",
		},
	)
	m.RewardsCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: defaultNamespace,
			Subsystem: subsystem,
			Name:      "rewards_count",
			Help:      "Number of rewards",
		},
	)
	m.SlashesCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: defaultNamespace,
			Subsystem: subsystem,
			Name:      "slashes_count",
			Help:      "Number of slashes",
		},
	)
	m.EncryptedCommitmentsCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: defaultNamespace,
			Subsystem: subsystem,
			Name:      "encrypted_commitments_count",
			Help:      "Number of encrypted commitments",
		},
	)
	return m
}

func (m *metrics) Collectors() []prometheus.Collector {
	return []prometheus.Collector{
		m.CommitmentsCount,
		m.RewardsCount,
		m.SlashesCount,
		m.EncryptedCommitmentsCount,
	}
}
