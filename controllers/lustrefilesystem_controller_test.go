/*
Copyright 2021 Hewlett Packard Enterprise Development LP
*/
package controllers

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.hpe.com/hpe/hpc-rabsw-lustre-fs-operator/api/v1alpha1"
)

var _ = Describe("LustreFileSystem Controller", func() {
	var (
		key                types.NamespacedName
		created, retrieved *v1alpha1.LustreFileSystem
	)

	Context("Create", func() {
		It("Should create successfully", func() {
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
					Name:      "test",
					MgsNid:    "172.0.0.1@tcp",
					MountRoot: "/lus/test",
				},
			}

			By("creating the lustre file system")
			Expect(k8sClient.Create(context.TODO(), created)).To(Succeed())

			By("retrieving created lustre file system")
			retrieved = &v1alpha1.LustreFileSystem{}
			Eventually(func() error {
				return k8sClient.Get(context.TODO(), key, retrieved)
			}, "1s").Should(Succeed())

			//Expect(retrieved).To(Equal(created)) // retrieved will have TypeMeta which is not provided on create call, but is filled in by kubernetes

			By("expect persistent volume", func() {
				pvkey := types.NamespacedName{
					Name: key.Name + PersistentVolumeSuffix,
					//Namespace: key.Namespace, // Cluster-scoped resource cannot have a namespace (even if "default")
				}

				pv := &corev1.PersistentVolume{}

				By("get persistent volume")
				Eventually(func() error {
					return k8sClient.Get(context.TODO(), pvkey, pv)
				}, "3s").Should(Succeed())
			})

			By("expect persistent volume claim", func() {
				pvckey := types.NamespacedName{
					Name:      key.Name + PersistentVolumeClaimSuffix,
					Namespace: key.Namespace,
				}

				pvc := &corev1.PersistentVolumeClaim{}

				By("get persistent volume claim")
				Eventually(func() error {
					return k8sClient.Get(context.TODO(), pvckey, pvc)
				}, "3s").Should(Succeed())
			})
		})
	})
})
