package bookie

const (
	ComponentName = "init-bookie-job"
)
const (
	MemData = "-Xms64m -Xmx256m -XX:MaxDirectMemorySize=256m"

	PulsarGC = `-XX:+UseG1GC -XX:MaxGCPauseMillis=10 -XX:+ParallelRefProcEnabled -XX:+UnlockExperimentalVMOptions -XX:+DoEscapeAnalysis -XX:ParallelGCThreads=4 -XX:ConcGCThreads=4 -XX:G1NewSizePercent=50 -XX:+DisableExplicitGC -XX:-ResizePLAB -XX:+ExitOnOutOfMemoryError -XX:+PerfDisableSharedMem -verbosegc`

	PulsarMem = "-Xms128m -Xmx256m -XX:MaxDirectMemorySize=256m"

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
