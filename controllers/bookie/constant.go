package bookie

const (
	MemData = "-Xms64m -Xmx256m -XX:MaxDirectMemorySize=256m"

	StatsProviderClass = "org.apache.bookkeeper.stats.prometheus.PrometheusMetricsProvider"

	JournalDataMountPath = "/pulsar/data/bookkeeper/journal"

	LedgersDataMountPath = "/pulsar/data/bookkeeper/ledgers"
)

// StatefulSetAnnotations Annotations
var StatefulSetAnnotations map[string]string

func init() {
	StatefulSetAnnotations = make(map[string]string)
	StatefulSetAnnotations["prometheus.io/scrape"] = "true"
	StatefulSetAnnotations["prometheus.io/port"] = "8000"
}
