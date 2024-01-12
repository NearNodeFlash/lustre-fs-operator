/*
 * Copyright 2021-2023 Hewlett Packard Enterprise Development LP
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

package controller

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/retry"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	lusv1beta1 "github.com/NearNodeFlash/lustre-fs-operator/api/v1beta1"
)

var _ = Describe("LustreFileSystem Controller", func() {

	var fs *lusv1beta1.LustreFileSystem

	BeforeEach(func() {
		fs = &lusv1beta1.LustreFileSystem{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "controller",
				Namespace: corev1.NamespaceDefault,
			},
			Spec: lusv1beta1.LustreFileSystemSpec{
				Name:             "test",
				MgsNids:          "172.0.0.1@tcp",
				MountRoot:        "/lus/test",
				StorageClassName: "nnf-lustre-fs",
			},
		}
	})

	JustBeforeEach(func() {
		Expect(k8sClient.Create(ctx, fs)).Should(Succeed())
	})

	JustAfterEach(func() {
		if fs != nil {
			Expect(k8sClient.Delete(ctx, fs)).Should(Succeed())
			Eventually(func() error {
				return k8sClient.Get(ctx, client.ObjectKeyFromObject(fs), fs)
			}).ShouldNot(Succeed())
		}
	})

	Context("creates successfully with no namespaces", func() {
		It("has no accesses", func() {
			Eventually(func() error {
				return k8sClient.Get(ctx, client.ObjectKeyFromObject(fs), fs)
			}).Should(Succeed())

			Expect(fs.Spec.Namespaces).To(HaveLen(0))
			Expect(fs.Status.Namespaces).To(HaveLen(0))
		})
	})

	Context("creates successfully with namespace but no mode", func() {
		const namespace = "dummy-namespace"

		BeforeEach(func() {
			fs.Spec.Namespaces = map[string]lusv1beta1.LustreFileSystemNamespaceSpec{
				namespace: {},
			}
		})

		It("has mode but no namespace accesses", func() {
			Eventually(func() error {
				return k8sClient.Get(ctx, client.ObjectKeyFromObject(fs), fs)
			}).Should(Succeed())

			Expect(fs.Spec.Namespaces).To(HaveKey(namespace))
			Expect(fs.Spec.Namespaces[namespace].Modes).To(HaveLen(0))
		})
	})

	Context("with a dummy namespace", Ordered, func() {
		const namespace = "dummy-namespace"
		const mode = corev1.ReadWriteMany

		/*
			For some reason envtest never actually deletes a namespace, so instead of managing the
			namespace for each test case, it is done once in the BeforeAll() below.

			BeforeEach(func() {
				Expect(k8sClient.Create(ctx, ns)).Should(Succeed())
			})

			AfterEach(func() {
				Expect(k8sClient.Delete(ctx, ns)).Should(Succeed())
				Eventually(k8sClient.Get(ctx, client.ObjectKeyFromObject(ns), ns)).ShouldNot(Succeed())
			})
		*/

		BeforeAll(func() {
			ns := &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{
				Name: namespace,
			}}

			Expect(k8sClient.Create(ctx, ns)).Should(Succeed())
		})

		validateCreateOccurredFn := func() {
			By("verifying namespaces are ready")
			Eventually(func(g Gomega) lusv1beta1.LustreFileSystemNamespaceAccessStatus {
				g.Expect(k8sClient.Get(ctx, client.ObjectKeyFromObject(fs), fs)).Should(Succeed())
				g.Expect(fs.Status.Namespaces).To(HaveKey(namespace))
				g.Expect(fs.Status.Namespaces[namespace].Modes).To(HaveKey(mode))

				return fs.Status.Namespaces[namespace].Modes[mode]
			}).Should(MatchAllFields(Fields{
				"State":                    Equal(lusv1beta1.NamespaceAccessReady),
				"PersistentVolumeRef":      Not(BeNil()),
				"PersistentVolumeClaimRef": Not(BeNil()),
			}))

			By("verifying PV exists")
			pv := &corev1.PersistentVolume{
				ObjectMeta: metav1.ObjectMeta{
					Name: fs.PersistentVolumeName(namespace, mode),
				},
			}

			Expect(k8sClient.Get(ctx, client.ObjectKeyFromObject(pv), pv)).Should(Succeed())

			By("verifying PVC exists")
			pvc := &corev1.PersistentVolumeClaim{
				ObjectMeta: metav1.ObjectMeta{
					Name:      fs.PersistentVolumeClaimName(namespace, mode),
					Namespace: namespace,
				},
			}

			Expect(k8sClient.Get(ctx, client.ObjectKeyFromObject(pvc), pvc)).Should(Succeed())
		}

		Context("with namespace and mode on create", func() {

			BeforeEach(func() {
				fs.Spec.Namespaces = map[string]lusv1beta1.LustreFileSystemNamespaceSpec{
					namespace: {
						Modes: []corev1.PersistentVolumeAccessMode{
							mode,
						},
					},
				}
			})

			It("creates pv/pvc and goes ready", func() {
				validateCreateOccurredFn()
			})

			It("does not delete until only our finalizer is left", func() {
				const finalizer = "test-finalizer"

				By("adding outside finalizer")
				Expect(retry.RetryOnConflict(retry.DefaultBackoff, func() error {
					if err := k8sClient.Get(ctx, client.ObjectKeyFromObject(fs), fs); err != nil {
						return err
					}
					controllerutil.AddFinalizer(fs, finalizer)
					return k8sClient.Update(ctx, fs)
				})).To(Succeed())

				validateCreateOccurredFn()

				By("deleting fs and expecting it to do nothing")
				Expect(k8sClient.Delete(ctx, fs)).To(Succeed())
				Consistently(func(g Gomega) error {
					g.Expect(k8sClient.Get(ctx, client.ObjectKeyFromObject(fs), fs)).Should(Succeed())
					validateCreateOccurredFn()
					return nil
				}, "3s").Should(Succeed(), "fs was deleted prior to removing outside finalizer")

				By("removing the finalizer and verify its removal")
				Expect(retry.RetryOnConflict(retry.DefaultBackoff, func() error {
					if err := k8sClient.Get(ctx, client.ObjectKeyFromObject(fs), fs); err != nil {
						return err
					}
					controllerutil.RemoveFinalizer(fs, finalizer)
					return k8sClient.Update(ctx, fs)
				})).To(Succeed())

				Eventually(func(g Gomega) error {
					return k8sClient.Get(ctx, client.ObjectKeyFromObject(fs), fs)
				}).ShouldNot(Succeed())

				fs = nil // we already cleaned up
			})
		})

		Context("adding a namespace post create", func() {
			const mode = corev1.ReadWriteMany

			JustBeforeEach(func() {
				Eventually(func(g Gomega) error {
					g.Expect(k8sClient.Get(ctx, client.ObjectKeyFromObject(fs), fs)).Should(Succeed())
					Expect(fs.Spec.Namespaces).To(BeEmpty())

					fs.Spec.Namespaces = map[string]lusv1beta1.LustreFileSystemNamespaceSpec{
						namespace: {
							Modes: []corev1.PersistentVolumeAccessMode{
								mode,
							},
						},
					}

					return k8sClient.Update(ctx, fs)
				}).Should(Succeed())
			})

			It("creates pv/pvc and goes ready", func() {
				validateCreateOccurredFn()
			})

			It("removes pv/pvc", func() {
				validateCreateOccurredFn()

				Eventually(func(g Gomega) error {
					g.Expect(k8sClient.Get(ctx, client.ObjectKeyFromObject(fs), fs)).Should(Succeed())
					g.Expect(fs.Spec.Namespaces).To(HaveKey(namespace))

					delete(fs.Spec.Namespaces, namespace)

					return k8sClient.Update(ctx, fs)
				}).Should(Succeed())

				Eventually(func(g Gomega) map[string]lusv1beta1.LustreFileSystemNamespaceStatus {
					g.Expect(k8sClient.Get(ctx, client.ObjectKeyFromObject(fs), fs)).Should(Succeed())
					return fs.Status.Namespaces
				}).ShouldNot(HaveKey(namespace))

				// envtest doesn't support deletion of PV/PVC resources
				/*
					Eventually(func() error {
						pvc := &corev1.PersistentVolumeClaim{
							ObjectMeta: metav1.ObjectMeta{
								Name:      fs.PersistentVolumeClaimName(namespace, mode),
								Namespace: namespace,
							},
						}

						return k8sClient.Get(ctx, client.ObjectKeyFromObject(pvc), pvc)
					}).ShouldNot(Succeed())

					Eventually(func() error {
						pv := &corev1.PersistentVolume{
							ObjectMeta: metav1.ObjectMeta{
								Name: fs.PersistentVolumeName(namespace, mode),
							},
						}

						return k8sClient.Get(ctx, client.ObjectKeyFromObject(pv), pv)
					}).ShouldNot(Succeed())
				*/
			})
		})
	})
})
