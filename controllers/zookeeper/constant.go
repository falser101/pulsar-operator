package zookeeper

const (
	// PulsarMemData Mem
	PulsarMemData = "\" -Xms100m -Xmx256m \""

	// PulsarGCData GC
	PulsarGCData = "\" -XX:+UseG1GC -XX:MaxGCPauseMillis=10\""

	// ContainerZookeeperServerList Container zookeeper server list env key
	ContainerZookeeperServerList = "ZOOKEEPER_SERVERS"

	// ContainerDataVolumeName Container data volume name
	ContainerDataVolumeName = "datadir"

	// ContainerDataPath Container zookeeper data path
	ContainerDataPath = "/pulsar/data"

	// ReadinessProbeScript ReadinessProbe script
	ReadinessProbeScript = "bin/pulsar-zookeeper-ruok.sh"

	// LivenessProbeScript LivenessProbe script
	LivenessProbeScript = "bin/pulsar-zookeeper-ruok.sh"
)

// ServiceAnnotations Annotations
var ServiceAnnotations map[string]string
var StatefulSetAnnotations map[string]string

func init() {
	// Init Service Annotations
	ServiceAnnotations = make(map[string]string)
	ServiceAnnotations["service.alpha.kubernetes.io/tolerate-unready-endpoints"] = "true"

	// Init StatefulSet Annotations
	StatefulSetAnnotations = make(map[string]string)
	StatefulSetAnnotations["pod.alpha.kubernetes.io/initialized"] = "true"
	StatefulSetAnnotations["prometheus.io/scrape"] = "true"
	StatefulSetAnnotations["prometheus.io/port"] = "8000"
}
