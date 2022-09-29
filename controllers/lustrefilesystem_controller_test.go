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

package controllers

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/NearNodeFlash/lustre-fs-operator/api/v1alpha1"
)

var _ = Describe("LustreFileSystem Controller", func() {
	var (
		key                types.NamespacedName
		created, retrieved *v1alpha1.LustreFileSystem
	)

	BeforeEach(func() {
		key = types.NamespacedName{
			Name:      "lustre-fs-example",
			Namespace: corev1.NamespaceDefault,
		}

		created = &v1alpha1.LustreFileSystem{
			ObjectMeta: metav1.ObjectMeta{
				Name:      key.Name,
				Namespace: key.Namespace,
			},
			Spec: v1alpha1.LustreFileSystemSpec{
				Name:             "test",
				MgsNids:          "172.0.0.1@tcp",
				MountRoot:        "/lus/test",
				StorageClassName: "nnf-lustre-fs",
			},
		}
	})

	Context("Create", func() {
		It("Should create successfully", func() {
			By("creating the lustre file system")
			Expect(k8sClient.Create(context.TODO(), created)).To(Succeed())

			By("retrieving created lustre file system")
			retrieved = &v1alpha1.LustreFileSystem{}
			Eventually(func() error {
				return k8sClient.Get(context.TODO(), key, retrieved)
			}, "1s").Should(Succeed())

			//Expect(retrieved).To(Equal(created)) // retrieved will have TypeMeta which is not provided on create call, but is filled in by kubernetes
		})

		It("Should have created a persistent volume", func() {
			pvkey := types.NamespacedName{
				Name: key.Name + PersistentVolumeSuffix,
				//Namespace: key.Namespace, // Cluster-scoped resource cannot have a namespace (even if "default")
			}

			pv := &corev1.PersistentVolume{}

			By("get persistent volume")
			Eventually(func() error {
				return k8sClient.Get(context.TODO(), pvkey, pv)
			}, "3s").Should(Succeed())

			By("and it must have the storage class name set")
			Expect(pv.Spec.StorageClassName).To(Equal(created.Spec.StorageClassName))
			By("and it must be reserved for the matching PVC")
			Expect(pv.Spec.ClaimRef.Name).To(Equal(key.Name + PersistentVolumeClaimSuffix))
			Expect(pv.Spec.ClaimRef.Namespace).To(Equal(key.Namespace))
		})

		It("Should have created a persistent volume claim", func() {
			pvckey := types.NamespacedName{
				Name:      key.Name + PersistentVolumeClaimSuffix,
				Namespace: key.Namespace,
			}

			pvc := &corev1.PersistentVolumeClaim{}

			By("get persistent volume claim")
			Eventually(func() error {
				return k8sClient.Get(context.TODO(), pvckey, pvc)
			}, "3s").Should(Succeed())

			By("and it must have the volume name set")
			Expect(pvc.Spec.VolumeName).To(Equal(key.Name + PersistentVolumeSuffix))
			By("and it must have the storage class name set")
			Expect(*pvc.Spec.StorageClassName).To(Equal(created.Spec.StorageClassName))
		})
	})

	Context("Delete", func() {
		It("Deletion of LustreFileSystem should delete the PV", func() {
			fsDeleter := &v1alpha1.LustreFileSystem{
				ObjectMeta: metav1.ObjectMeta{
					Name:      key.Name,
					Namespace: key.Namespace,
				},
			}
			// In finalizer cleanup, the LustreFileSystem resource
			// will explicitly delete the PV.  The PV does not
			// have an owner ref, so it isn't cleaned up by
			// garbage collection.
			Expect(k8sClient.Delete(context.TODO(), fsDeleter)).To(Succeed())

			fsKey := types.NamespacedName{
				Name:      key.Name + PersistentVolumeClaimSuffix,
				Namespace: key.Namespace,
			}
			fs := &v1alpha1.LustreFileSystem{}
			By("Wait for the LustreFileSystem to be deleted")
			Eventually(func() error {
				return k8sClient.Get(context.TODO(), fsKey, fs)
			}, "3s").ShouldNot(Succeed())

			By("Wait for reconciler to run")
			time.Sleep(10 * time.Second)
		})

		PIt("Lifecycle of the PVC", func() {

			// The PVC is owned by the LustreFileSystem resource,
			// so when running live it will be garbage-collected.
			// In the test env that garbage collection doesn't
			// happen.  That means the PV is still hanging around
			// because it is still bound to the PVC.

			// Unfortunately, in the test environment the PVC
			// is still hanging around even after that explicit
			// deletion.

			pvcGarbage := &corev1.PersistentVolumeClaim{
				ObjectMeta: metav1.ObjectMeta{
					Name:      key.Name + PersistentVolumeClaimSuffix,
					Namespace: key.Namespace,
				},
			}
			By("Fake garbage collection of the PVC")
			Expect(k8sClient.Delete(context.TODO(), pvcGarbage)).To(Succeed())
			time.Sleep(10 * time.Second)
			By("Confirm the deletion of the PVC")
			pvckey := types.NamespacedName{
				Name:      key.Name + PersistentVolumeClaimSuffix,
				Namespace: key.Namespace,
			}

			pvc := &corev1.PersistentVolumeClaim{}
			Eventually(func() error {
				return k8sClient.Get(context.TODO(), pvckey, pvc)
			}, "10s", "3s").ShouldNot(Succeed())
		})

		PIt("Lifecycle of the PV", func() {
			// Now that we've done the fake garbage-collection
			// of the PVC, the PV should disappear.

			// Unfortunately, in the test environment the PV is
			// hanging on.

			pvkey := types.NamespacedName{
				Name: key.Name + PersistentVolumeSuffix,
			}
			pv := &corev1.PersistentVolume{}
			By("Verify that the PV has been deleted")
			Eventually(func() error {
				return k8sClient.Get(context.TODO(), pvkey, pv)
			}, "10s", "3s").ShouldNot(Succeed())
		})
	})
})
