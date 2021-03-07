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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// KubeKillerServerSpec defines the desired state of KubeKillerServer
type KubeKillerServerSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of KubeKillerServer. Edit KubeKillerServer_types.go to remove/update
	Name    string `json:"name"`
	Image   string `json:"image"`
	Replica int32  `json:"replica"`
}

// KubeKillerServerStatus defines the observed state of KubeKillerServer
type KubeKillerServerStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	DeploymentStatus string `json:"deploymentStatus"`
	ServiceStatus    string `json:"serviceStatus"`
	IngressStatus    string `json:"ingressStatus"`
}

// +kubebuilder:object:root=true

// KubeKillerServer is the Schema for the kubekillerservers API
type KubeKillerServer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KubeKillerServerSpec   `json:"spec,omitempty"`
	Status KubeKillerServerStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// KubeKillerServerList contains a list of KubeKillerServer
type KubeKillerServerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KubeKillerServer `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KubeKillerServer{}, &KubeKillerServerList{})
}
