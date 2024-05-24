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
var StatefulSetAnnotations map[string]string

func init() {
	StatefulSetAnnotations = make(map[string]string)
	StatefulSetAnnotations["prometheus.io/scrape"] = "true"
	StatefulSetAnnotations["prometheus.io/port"] = "8080"
}
