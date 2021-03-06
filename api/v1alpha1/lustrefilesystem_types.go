/*
 * Copyright 2021, 2022 Hewlett Packard Enterprise Development LP
 * Other additional copyright holders may be indicated within.
 *
 * The entirety of this work is licensed under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 *
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// LustreFileSystemSpec defines the desired state of LustreFileSystem
type LustreFileSystemSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Name is the name of the Lustre file system.
	// +kubebuilder:validation:MaxLength:=8
	// +kubebuilder:validation:MinLength:=1
	Name string `json:"name"`

	// MgsNids lists the NID of the MGS to use for accessing the Lustre file system. The MGS NIDs are combined with the
	// Name to establish a connection to a global Lustre file system. Connections occur in the order of the listed NIDs.
	// +kubebuilder:validation:MinItems:=1
	MgsNids []string `json:"mgsNids"`

	// MountRoot is the mount path used to access the Lustre file system from a host. Data Movement directives can
	// reference this field when performing data movement from or to the Lustre file system.
	MountRoot string `json:"mountRoot"`

	// StorageClassName refers to the StorageClass to use for this
	// file system.
	// +kubebuilder:default="nnf-lustre-fs"
	StorageClassName string `json:"storageClassName,omitempty"`
}

// LustreFileSystemStatus defines the observed state of LustreFileSystem
type LustreFileSystemStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="FSNAME",type="string",JSONPath=".spec.name",description="Lustre file system name"
//+kubebuilder:printcolumn:name="MgsNID",type="string",JSONPath=".spec.mgsNids[0]",description="MGS NID"
//+kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
//+kubebuilder:printcolumn:name="MountRoot",type="string",JSONPath=".spec.mountRoot",priority=1,description="Mount path used to mount filesystem"
//+kubebuilder:printcolumn:name="StorageClass",type="string",JSONPath=".spec.storageClassName",priority=1,description="StorageClass to use"

// LustreFileSystem is the Schema for the lustrefilesystems API
type LustreFileSystem struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   LustreFileSystemSpec   `json:"spec,omitempty"`
	Status LustreFileSystemStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// LustreFileSystemList contains a list of LustreFileSystem
type LustreFileSystemList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []LustreFileSystem `json:"items"`
}

func init() {
	SchemeBuilder.Register(&LustreFileSystem{}, &LustreFileSystemList{})
}
