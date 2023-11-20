/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are requ ired.  Any new fields you add must have json tags for the fields to be serialized.

// PulsarSpec defines the desired state of Pulsar
type PulsarSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// The following markers will use OpenAPI v3 schema to validate the value
	// More info: https://book.kubebuilder.io/reference/markers/crd-validation.html
	// Broker defines the desired state of Broker
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Broker Broker `json:"broker,omitempty"`

	// Zookeeper defines the desired state of Zookeeper
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Zookeeper Zookeeper `json:"zookeeper,omitempty"`

	// Bookie defines the desired state of Bookie
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Bookie Bookie `json:"bookie,omitempty"`

	// Proxy defines the desired state of Proxy
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Proxy Proxy `json:"proxy,omitempty"`

	// AutoRecovery defines the desired state of AutoRecovery
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	AutoRecovery AutoRecovery `json:"autoRecovery,omitempty"`

	// Manager defines the desired state of Manager
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Manager Manager `json:"manager,omitempty"`

	// Monitor defines the desired state of Monitor
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Monitor Monitor `json:"monitor,omitempty"`
}

func (s *PulsarSpec) SetDefault(c *Pulsar) bool {
	changed := false
	if s.Broker.SetDefault(c) {
		changed = true
	}
	if s.Zookeeper.SetDefault(c) {
		changed = true
	}
	if s.Bookie.SetDefault(c) {
		changed = true
	}
	if s.AutoRecovery.SetDefault(c) {
		changed = true
	}
	if s.Manager.SetDefault(c) {
		changed = true
	}
	if s.Monitor.SetDefault(c) {
		changed = true
	}
	if s.Proxy.SetDefault(c) {
		changed = true
	}
	return changed
}

// PulsarStatus defines the observed state of PulsarCluster
type PulsarStatus struct {
	// Represents the observations of a Memcached's current state.
	// Memcached.status.conditions.type are: "Available", "Progressing", and "Degraded"
	// Memcached.status.conditions.status are one of True, False, Unknown.
	// Memcached.status.conditions.reason the value should be a CamelCase string and producers of specific
	// condition types may define expected values and meanings for this field, and whether the values
	// are considered a guaranteed API.
	// Memcached.status.conditions.Message is a human readable message indicating details about the transition.
	// For further information see: https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md#typical-status-properties

	// Conditions store the status conditions of the Memcached instances
	// +operator-sdk:csv:customresourcedefinitions:type=status
	Phase string `json:"phase,omitempty"`
}

func (s *PulsarStatus) SetDefault(c *Pulsar) bool {
	changed := false
	if s.Phase == "" {
		s.Phase = PulsarClusterInitializingPhase
		changed = true
	}
	return changed
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Pulsar is the Schema for the pulsar API
type Pulsar struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PulsarSpec   `json:"spec,omitempty"`
	Status PulsarStatus `json:"status,omitempty"`
}

func (c *Pulsar) SpecSetDefault() bool {
	return c.Spec.SetDefault(c)
}

func (c *Pulsar) StatusSetDefault() bool {
	return c.Status.SetDefault(c)
}

//+kubebuilder:object:root=true

// PulsarList contains a list of Pulsar
type PulsarList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Pulsar `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Pulsar{}, &PulsarList{})
}
