package v1alpha1

// Broker defines the desired state of Broker
type Broker struct {
	// Image is the  container image. default is apachepulsar/pulsar-all:latest
	Image ContainerImage `json:"image,omitempty"`

	// Labels specifies the labels to attach to pods the operator creates for
	// the broker cluster.
	Labels map[string]string `json:"labels,omitempty"`

	// Replicas is the expected size of the broker cluster.
	Replicas int32 `json:"replicas,omitempty"`

	// Pod defines the policy to create pod for the broker cluster.
	//
	// Updating the pod does not take effect on any existing pods.
	Pod PodPolicy `json:"pod,omitempty"`

	ConfigData map[string]string `json:"configData,omitempty"`

	Ports BrokerPorts `json:"ports,omitempty"`
}

type BrokerPorts struct {
	Http      string `json:"http,omitempty"`
	Https     string `json:"https,omitempty"`
	Pulsar    string `json:"pulsar,omitempty"`
	Pulsarssl string `json:"pulsarssl,omitempty"`
}

func (b *Broker) SetDefault(cluster *PulsarCluster) bool {
	changed := false

	if b.Image.SetDefault(cluster, BrokerComponent) {
		changed = true
	}

	if b.Pod.SetDefault(cluster, BrokerComponent) {
		changed = true
	}

	if b.setConfigDataDefault(cluster) {
		changed = true
	}

	if b.setPortsDefault(cluster) {
		changed = true
	}
	return changed
}

func (b *Broker) setPortsDefault(cluster *PulsarCluster) (changed bool) {
	if b.Ports.Http == "" {
		b.Ports.Http = "8080"
		changed = true
	}

	if b.Ports.Https == "" {
		b.Ports.Https = "8443"
		changed = true
	}

	if b.Ports.Pulsar == "" {
		b.Ports.Pulsar = "6650"
		changed = true
	}

	if b.Ports.Pulsarssl == "" {
		b.Ports.Pulsarssl = "6651"
		changed = true
	}
	return
}

func (b *Broker) setConfigDataDefault(c *PulsarCluster) (changed bool) {
	if b.ConfigData == nil {
		b.ConfigData = make(map[string]string)
		b.ConfigData["managedLedgerDefaultEnsembleSize"] = "1"
		b.ConfigData["managedLedgerDefaultWriteQuorum"] = "1"
		b.ConfigData["managedLedgerDefaultAckQuorum"] = "1"
		b.ConfigData["PULSAR_MEM"] = `>
			-Xms128m -Xmx256m -XX:MaxDirectMemorySize=256m`
		b.ConfigData["PULSAR_GC"] = `>
			-XX:+UseG1GC
			-XX:MaxGCPauseMillis=10
			-Dio.netty.leakDetectionLevel=disabled
			-Dio.netty.recycler.linkCapacity=1024
			-XX:+ParallelRefProcEnabled
			-XX:+UnlockExperimentalVMOptions
			-XX:+DoEscapeAnalysis
			-XX:ParallelGCThreads=4
			-XX:ConcGCThreads=4
			-XX:G1NewSizePercent=50
			-XX:+DisableExplicitGC
			-XX:-ResizePLAB
			-XX:+ExitOnOutOfMemoryError
			-XX:+PerfDisableSharedMem`
		changed = true
	}
	return
}
