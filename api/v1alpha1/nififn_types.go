/*

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

// Package v1alpha1 defines NiFiFn types for the controller
// +kubebuilder:validation:Optional
package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SSLConfig defines an SSL context for securing NiFi communication
type SSLConfig struct {
	// +kubebuilder:validation:Required
	KeystoreFile string `json:"keystoreFile"`
	// +kubebuilder:validation:Required
	KeystorePass string `json:"keystorePass"`
	// +kubebuilder:validation:Required
	KeyPass string `json:"keyPass"`
	// +kubebuilder:validation:Required
	KeystoreType string `json:"keystoreType"`
	// +kubebuilder:validation:Required
	TruststoreFile string `json:"truststoreFile"`
	// +kubebuilder:validation:Required
	TruststorePass string `json:"truststorePass"`
	// +kubebuilder:validation:Required
	TruststoreType string `json:"truststoreType"`
}

// NiFiFnSpec defines the desired state of NiFiFn
type NiFiFnSpec struct {

	// FailureRetry
	// FailureBackoff
	// TTLSecondsAfterFinished int32

	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Enum=registry;xml
	RunFrom string `json:"runFrom"`

	RegistryURL string `json:"registryUrl,omitempty"`

	// +kubebuilder:validation:MaxLength=36
	// +kubebuilder:validation:MinLength=36
	BucketID string `json:"bucketId,omitempty"`

	// +kubebuilder:validation:MaxLength=36
	// +kubebuilder:validation:MinLength=36
	FlowID string `json:"flowId,omitempty"`

	FlowVersion int32 `json:"flowVersion,omitempty"`

	FlowXMLPath string `json:"flowXmlPath,omitempty"`

	// +kubebuilder:validation:Pattern=.+:.+
	Image string `json:"-"`

	MaterializeContent bool `json:"materializeContent,omitempty"`

	FailurePortIDs []string `json:"failurePortIds,omitempty"`

	SSLConfig `json:"ssl,omitempty"`

	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinItems=1
	FlowFiles []map[string]string `json:"flowFiles"`

	Parameters map[string]string `json:"parameters,omitempty"`
}

// NiFiFnStatus defines the observed state of NiFiFn
type NiFiFnStatus struct {
	// TODO
	// CurrentVersion int32
	// Flow string
	// Bucket string
	// Queued Files ?
	// Processed Files ?

	Flow string `json:"flow"`
}

// +kubebuilder:object:root=true

// NiFiFn is the Schema for the nififns API
// +kubebuilder:printcolumn:name="Flow",type="string",JSONPath=".spec.flow",description="The UUID of the Flow in NiFi-Registry"
// +kubebuilder:printcolumn:name="Version",type="integer",JSONPath=".spec.flowVersion",description="The version of the NiFiFlow"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// TODO: +kubebuilder:subresource:status
// TODO: +kubebuilder:subresource:scale
type NiFiFn struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NiFiFnSpec   `json:"spec,omitempty"`
	Status NiFiFnStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// NiFiFnList contains a list of NiFiFn
type NiFiFnList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NiFiFn `json:"items"`
}

func init() {
	SchemeBuilder.Register(&NiFiFn{}, &NiFiFnList{})
}
