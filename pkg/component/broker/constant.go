package broker

const (
	PulsarMemData = "-Xms1024m -Xmx1024m -XX:MaxDirectMemorySize=1024m"

	ManagedLedgerDefaultEnsembleSize = "1"

	ManagedLedgerDefaultWriteQuorum = "1"

	ManagedLedgerDefaultAckQuorum = "1"

	FunctionsWorkerEnabled = "false"

	AdvertisedAddress = "advertisedAddress"
)

// Annotations
var DeploymentAnnotations map[string]string

func init() {
	DeploymentAnnotations = make(map[string]string)
	DeploymentAnnotations["prometheus.io/scrape"] = "true"
	DeploymentAnnotations["prometheus.io/port"] = "8080"
}
