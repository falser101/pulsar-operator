package v1alpha1

type Provider string

const (
	JWT Provider = "jwt"
)

type Authentication struct {
	// Authentication is the authentication policy for the broker cluster.
	Enabled  bool     `json:"enabled,omitempty"`
	Provider Provider `json:"provider,omitempty"`
}
