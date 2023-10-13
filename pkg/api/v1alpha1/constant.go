package v1alpha1

// Service pulsar
const (
	Service = "pulsar"
)

const (
	PulsarClusterInitializingPhase = "Initializing"

	// PulsarClusterLaunchingPhase Launching phase
	PulsarClusterLaunchingPhase = "Launching"

	// PulsarClusterRunningPhase Running phase
	PulsarClusterRunningPhase = "Running"
)

// default number
const (
	// ZookeeperClusterDefaultNodeNum number default num is 3
	ZookeeperClusterDefaultNodeNum = 3

	// BrokerClusterDefaultNodeNum number default num is 3
	BrokerClusterDefaultNodeNum = 3

	// BookieClusterDefaultNodeNum number default num is 3
	BookieClusterDefaultNodeNum = 3

	// ProxyClusterDefaultNodeNum number default num is 3
	ProxyClusterDefaultNodeNum = 3

	ManagerDefaultNodeNum = 1
)

const (
	// ZookeeperComponent Zookeeper component
	ZookeeperComponent = "zookeeper"

	// BrokerComponent broker
	BrokerComponent = "broker"

	// BookieComponent bookie
	BookieComponent = "bookie"

	// AutoRecoveryComponent autoRecovery
	AutoRecoveryComponent = "autoRecovery"

	// ManagerComponent component
	ManagerComponent = "manager"

	// MonitorComponent component
	MonitorComponent = "monitor"
)

// monitor child component
const (
	// MonitorPrometheusComponent prometheus component
	MonitorPrometheusComponent = "monitor-prometheus"

	// MonitorGrafanaComponent grafana component
	MonitorGrafanaComponent = "monitor-grafana"
)

const (
	// DefaultAllPulsarContainerRepository is the default docker repo for components
	DefaultAllPulsarContainerRepository = "apachepulsar/pulsar-all"

	// DefaultAllPulsarContainerVersion is the default tag used for components
	DefaultAllPulsarContainerVersion = "latest"

	// DefaultPulsarManagerContainerRepository is default docker image name of pulsar manager
	DefaultPulsarManagerContainerRepository = "tlq-cn-console"

	// DefaultPulsarManagerContainerVersion is
	DefaultPulsarManagerContainerVersion = "v10.0.0.1"

	// DefaultPrometheusContainerRepository prometheus
	DefaultPrometheusContainerRepository = "prom/prometheus"

	// DefaultPrometheusContainerVersion version
	DefaultPrometheusContainerVersion = "v2.17.2"

	// DefaultContainerPolicy is the default container pull policy used
	DefaultContainerPolicy = "Always"

	DefaultMonitorGrafanaContainerRepository = "tlq-cn-grafana"

	DefaultMonitorGrafanaContainerTag = "v10.0.0.1"
)

// Labels
const (
	// LabelService App
	LabelService = "app"

	// LabelCluster LabelCluster
	LabelCluster = "cluster"

	// LabelComponent LabelComponent
	LabelComponent = "component"

	// LabelChildComponent child-component
	LabelChildComponent = "child-component"
)

func MakeComponentLabels(c *PulsarCluster, component string) map[string]string {
	return MakeAllLabels(c, component, "")
}

func MakeAllLabels(c *PulsarCluster, component string, childComponent string) map[string]string {
	labels := make(map[string]string)
	labels[LabelService] = Service
	labels[LabelCluster] = c.GetName()
	labels[LabelComponent] = component
	if childComponent != "" {
		labels[LabelChildComponent] = childComponent
	}
	return labels
}

// All component ports
const (
	// ZookeeperContainerClientDefaultPort Container client default port
	ZookeeperContainerClientDefaultPort = 2181

	// ZookeeperContainerServerDefaultPort Container server default port
	ZookeeperContainerServerDefaultPort = 2888

	// ZookeeperContainerLeaderElectionPort Container leader election port
	ZookeeperContainerLeaderElectionPort = 3888

	// PulsarBrokerPulsarServerPort Broker server port
	PulsarBrokerPulsarServerPort = 6650

	// PulsarBrokerHttpServerPort Broker http server port
	PulsarBrokerHttpServerPort = 8080

	// PulsarBookieServerPort Bookie server port
	PulsarBookieServerPort = 3181

	// PulsarGrafanaServerPort Grafana server port
	PulsarGrafanaServerPort = 3000

	// PulsarPrometheusServerPort Prometheus server port
	PulsarPrometheusServerPort = 9090

	// PulsarManagerBackendPort server port
	PulsarManagerBackendPort = 7750

	// PulsarManagerBackNodePort nodePort
	PulsarManagerBackNodePort = 30750

	// PulsarManagerFrontendNodePort nodePort
	PulsarManagerFrontendNodePort = 30527

	// PulsarManagerFrontendPort frontend port
	PulsarManagerFrontendPort = 9527
)

// Storage default capacity
const (
	// JournalStorageDefaultCapacity journal storage default capacity
	JournalStorageDefaultCapacity = 1

	// LedgersStorageDefaultCapacity ledgers storage default capacity
	LedgersStorageDefaultCapacity = 10

	PrometheusStorageDefaultCapacity = 10
	GrafanaStorageDefaultCapacity    = 10
	ZookeeperStorageDefaultCapacity  = 10
)

const (
	GrafanaDefaultAdminUser     = "test"
	GrafanaDefaultAdminPassword = "test"
)

const (
	// ServiceDomain service domain
	ServiceDomain = "svc.cluster.local"
)
