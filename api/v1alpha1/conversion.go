/*
 * Copyright 2022 Hewlett Packard Enterprise Development LP
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
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	lusv1beta1 "github.com/NearNodeFlash/lustre-fs-operator/api/v1beta1"
	utilconversion "github.com/NearNodeFlash/lustre-fs-operator/github/cluster-api/util/conversion"
)

var convertlog = logf.Log.V(2).WithName("convert-v1alpha1")

func (src *LustreFileSystem) ConvertTo(dstRaw conversion.Hub) error {
	convertlog.Info("Convert LustreFileSystem To Hub", "name", src.GetName(), "namespace", src.GetNamespace())
	dst := dstRaw.(*lusv1beta1.LustreFileSystem)

	if err := Convert_v1alpha1_LustreFileSystem_To_v1beta1_LustreFileSystem(src, dst, nil); err != nil {
		return err
	}

	// Manually restore data.
	restored := &lusv1beta1.LustreFileSystem{}
	if ok, err := utilconversion.UnmarshalData(src, restored); err != nil || !ok {
		return err
	}
	// EDIT THIS FUNCTION! If the annotation is holding anything that is
	// hub-specific then copy it into 'dst' from 'restored'.
	// Otherwise, you may comment out UnmarshalData() until it's needed.

	return nil
}

func (dst *LustreFileSystem) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*lusv1beta1.LustreFileSystem)
	convertlog.Info("Convert LustreFileSystem From Hub", "name", src.GetName(), "namespace", src.GetNamespace())

	if err := Convert_v1beta1_LustreFileSystem_To_v1alpha1_LustreFileSystem(src, dst, nil); err != nil {
		return err
	}

	// Preserve Hub data on down-conversion except for metadata.
	return utilconversion.MarshalData(src, dst)
}

// The List-based ConvertTo/ConvertFrom routines are never used by the
// conversion webhook, but the conversion-verifier tool wants to see them.
// The conversion-gen tool generated the Convert_X_to_Y routines, should they
// ever be needed.

func resource(resource string) schema.GroupResource {
	return schema.GroupResource{Group: "lus", Resource: resource}
}

func (src *LustreFileSystemList) ConvertTo(dstRaw conversion.Hub) error {
	return apierrors.NewMethodNotSupported(resource("LustreFileSystemList"), "ConvertTo")
}

func (dst *LustreFileSystemList) ConvertFrom(srcRaw conversion.Hub) error {
	return apierrors.NewMethodNotSupported(resource("LustreFileSystemList"), "ConvertFrom")
}
