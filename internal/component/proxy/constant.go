package proxy

const (
	MemData = "-Xms64m -Xmx256m -XX:MaxDirectMemorySize=256m"

	PulsarGC = `-XX:+UseG1GC -XX:MaxGCPauseMillis=10 -XX:+ParallelRefProcEnabled -XX:+UnlockExperimentalVMOptions -XX:+DoEscapeAnalysis -XX:ParallelGCThreads=4 -XX:ConcGCThreads=4 -XX:G1NewSizePercent=50 -XX:+DisableExplicitGC -XX:-ResizePLAB -XX:+ExitOnOutOfMemoryError -XX:+PerfDisableSharedMem -verbosegc`

	PulsarMem = "-Xms128m -Xmx256m -XX:MaxDirectMemorySize=256m"

	StatsProviderClass = "org.apache.bookkeeper.stats.prometheus.PrometheusMetricsProvider"

	JournalDataMountPath = "/pulsar/data/bookkeeper/journal"

	LedgersDataMountPath = "/pulsar/data/bookkeeper/ledgers"
)

// DeploymentAnnotations Annotations
var DeploymentAnnotations map[string]string

func init() {
	DeploymentAnnotations = make(map[string]string)
	DeploymentAnnotations["prometheus.io/scrape"] = "true"
	DeploymentAnnotations["prometheus.io/port"] = "8000"
}
