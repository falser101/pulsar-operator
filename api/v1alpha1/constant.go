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
)

const (
	// ZookeeperComponent Zookeeper component
	ZookeeperComponent = "zookeeper"

	// BrokerComponent broker
	BrokerComponent = "broker"

	// BookieComponent bookie
	BookieComponent = "bookie"

	// ProxyComponent bookie
	ProxyComponent = "proxy"

	// AutoRecoveryComponent autoRecovery
	AutoRecoveryComponent = "autorecovery"

	ToolsetComponent = "toolset"
)

const (
	// DefaultAllPulsarContainerRepository is the default docker repo for components
	DefaultAllPulsarContainerRepository = "apachepulsar/pulsar-all"

	// DefaultAllPulsarContainerVersion is the default tag used for components
	DefaultAllPulsarContainerVersion = "latest"

	// DefaultContainerPolicy is the default container pull policy used
	DefaultContainerPolicy = "IfNotPresent"
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
	labels[LabelCluster] = c.Name
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

	// PulsarBookieClientPort Bookie client port
	PulsarBookieClientPort = 8000
)

// Storage default capacity
const (
	// JournalStorageDefaultCapacity journal storage default capacity
	JournalStorageDefaultCapacity = "1Gi"

	// LedgersStorageDefaultCapacity ledgers storage default capacity
	LedgersStorageDefaultCapacity = "10Gi"

	ZookeeperClusterDefaultStorageCapacity = "10Gi"
)

const (
	GrafanaDefaultAdminUser     = "test"
	GrafanaDefaultAdminPassword = "test"
)

const (
	// ServiceDomain service domain
	ServiceDomain = "svc.cluster.local"
)
