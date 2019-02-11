/*
Copyright 2019 B23 LLC.

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

// NiFiFnSpec defines the desired state of NiFiFn
type NiFiFnSpec struct {

	// RegistryVariables map[string,string]
	// TTLSecondsAfterFinished int32

	RegistryURL string `json:"registryUrl"`

	// +kubebuilder:validation:MaxLength=36
	// +kubebuilder:validation:MinLength=36
	Bucket string `json:"bucket"`

	// +kubebuilder:validation:MaxLength=36
	// +kubebuilder:validation:MinLength=36
	Flow string `json:"flow"`

	// +kubebuilder:validation:Minimum=-1
	FlowVersion int32 `json:"flowVersion"`

	// +kubebuilder:validation:MinItems=1
	FlowFiles []string `json:"flowFiles"`

	// +kubebuilder:validation:Pattern=.+:.+
	Image string `json:"image"`
}

// NiFiFnStatus defines the observed state of NiFiFn
type NiFiFnStatus struct {
	// TODO
	// CurrentVersion int32
	// Flow string
	// Bucket string
	// Queued Files ?

	Flow string `json:"flow"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NiFiFn is the Schema for the nififns API
// +k8s:openapi-gen=true
// +kubebuilder:printcolumn:name="Flow",type="string",JSONPath=".spec.flow",description="The UUID of the Flow in NiFi-Registry"
// +kubebuilder:printcolumn:name="Version",type="integer",JSONPath=".spec.flowVersion",description="The version of the NiFiFlow"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:resource:shortName=nifn
type NiFiFn struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NiFiFnSpec   `json:"spec,omitempty"`
	Status NiFiFnStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NiFiFnList contains a list of NiFiFn
type NiFiFnList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NiFiFn `json:"items"`
}

func init() {
	SchemeBuilder.Register(&NiFiFn{}, &NiFiFnList{})
}
