//go:build !ignore_autogenerated

/*
 * Copyright 2024 Hewlett Packard Enterprise Development LP
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

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LustreFileSystem) DeepCopyInto(out *LustreFileSystem) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LustreFileSystem.
func (in *LustreFileSystem) DeepCopy() *LustreFileSystem {
	if in == nil {
		return nil
	}
	out := new(LustreFileSystem)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *LustreFileSystem) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LustreFileSystemList) DeepCopyInto(out *LustreFileSystemList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]LustreFileSystem, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LustreFileSystemList.
func (in *LustreFileSystemList) DeepCopy() *LustreFileSystemList {
	if in == nil {
		return nil
	}
	out := new(LustreFileSystemList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *LustreFileSystemList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LustreFileSystemNamespaceAccessStatus) DeepCopyInto(out *LustreFileSystemNamespaceAccessStatus) {
	*out = *in
	if in.PersistentVolumeRef != nil {
		in, out := &in.PersistentVolumeRef, &out.PersistentVolumeRef
		*out = new(v1.LocalObjectReference)
		**out = **in
	}
	if in.PersistentVolumeClaimRef != nil {
		in, out := &in.PersistentVolumeClaimRef, &out.PersistentVolumeClaimRef
		*out = new(v1.LocalObjectReference)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LustreFileSystemNamespaceAccessStatus.
func (in *LustreFileSystemNamespaceAccessStatus) DeepCopy() *LustreFileSystemNamespaceAccessStatus {
	if in == nil {
		return nil
	}
	out := new(LustreFileSystemNamespaceAccessStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LustreFileSystemNamespaceSpec) DeepCopyInto(out *LustreFileSystemNamespaceSpec) {
	*out = *in
	if in.Modes != nil {
		in, out := &in.Modes, &out.Modes
		*out = make([]v1.PersistentVolumeAccessMode, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LustreFileSystemNamespaceSpec.
func (in *LustreFileSystemNamespaceSpec) DeepCopy() *LustreFileSystemNamespaceSpec {
	if in == nil {
		return nil
	}
	out := new(LustreFileSystemNamespaceSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LustreFileSystemNamespaceStatus) DeepCopyInto(out *LustreFileSystemNamespaceStatus) {
	*out = *in
	if in.Modes != nil {
		in, out := &in.Modes, &out.Modes
		*out = make(map[v1.PersistentVolumeAccessMode]LustreFileSystemNamespaceAccessStatus, len(*in))
		for key, val := range *in {
			(*out)[key] = *val.DeepCopy()
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LustreFileSystemNamespaceStatus.
func (in *LustreFileSystemNamespaceStatus) DeepCopy() *LustreFileSystemNamespaceStatus {
	if in == nil {
		return nil
	}
	out := new(LustreFileSystemNamespaceStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LustreFileSystemSpec) DeepCopyInto(out *LustreFileSystemSpec) {
	*out = *in
	if in.Namespaces != nil {
		in, out := &in.Namespaces, &out.Namespaces
		*out = make(map[string]LustreFileSystemNamespaceSpec, len(*in))
		for key, val := range *in {
			(*out)[key] = *val.DeepCopy()
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LustreFileSystemSpec.
func (in *LustreFileSystemSpec) DeepCopy() *LustreFileSystemSpec {
	if in == nil {
		return nil
	}
	out := new(LustreFileSystemSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LustreFileSystemStatus) DeepCopyInto(out *LustreFileSystemStatus) {
	*out = *in
	if in.Namespaces != nil {
		in, out := &in.Namespaces, &out.Namespaces
		*out = make(map[string]LustreFileSystemNamespaceStatus, len(*in))
		for key, val := range *in {
			(*out)[key] = *val.DeepCopy()
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LustreFileSystemStatus.
func (in *LustreFileSystemStatus) DeepCopy() *LustreFileSystemStatus {
	if in == nil {
		return nil
	}
	out := new(LustreFileSystemStatus)
	in.DeepCopyInto(out)
	return out
}
