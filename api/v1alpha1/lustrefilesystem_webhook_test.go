/*
Copyright 2021 Hewlett Packard Enterprise Development LP
*/

package v1alpha1

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// These tests are written in BDD-style using Ginkgo framework. Refer to
// http://onsi.github.io/ginkgo to learn more.

var _ = Describe("LustreFileSystemWebhook", func() {
	var (
		key                                types.NamespacedName
		created, retrieved, valid, invalid *LustreFileSystem
	)

	Context("Create", func() {

		It("should create an object successfully", func() {
			key = types.NamespacedName{
				Name:      "lustre-fs-example",
				Namespace: "default",
			}

			created = &LustreFileSystem{
				ObjectMeta: metav1.ObjectMeta{
					Name:      key.Name,
					Namespace: key.Namespace,
				},
				Spec: LustreFileSystemSpec{
					Name:      "foo",
					MgsNid:    "127.0.0.1@tcp",
					MountRoot: "/lus/foo",
				},
			}

			By("creating an object")
			Expect(k8sClient.Create(context.TODO(), created)).To(Succeed())

			retrieved = &LustreFileSystem{}
			Expect(k8sClient.Get(context.TODO(), key, retrieved)).To(Succeed())
			Expect(retrieved).To(Equal(created))

			By("deleting the object")
			Expect(k8sClient.Delete(context.TODO(), created)).To(Succeed())
			Expect(k8sClient.Get(context.TODO(), key, created)).ToNot(Succeed())
		})

		It("should fail validating admission webhooks with invalid values", func() {
			key = types.NamespacedName{
				Name:      "lustre-fs-invalid-example",
				Namespace: "default",
			}

			valid = &LustreFileSystem{
				ObjectMeta: metav1.ObjectMeta{
					Name:      key.Name,
					Namespace: key.Namespace,
				},
				Spec: LustreFileSystemSpec{
					Name:      "foo",
					MgsNid:    "172.0.0.1@tcp",
					MountRoot: "lus/foo",
				},
			}

			makeInvalid := func(mutate func(fs *LustreFileSystem)) *LustreFileSystem {
				invalid = &LustreFileSystem{}
				valid.DeepCopyInto(invalid)

				mutate(invalid)

				return invalid
			}

			By("validate 'name' attribute", func() {
				By("empty attribute")
				Expect(k8sClient.Create(context.TODO(), makeInvalid(func(fs *LustreFileSystem) {
					fs.Spec.Name = ""
				}))).NotTo(Succeed())

				By("overflowing attribute")
				Expect(k8sClient.Create(context.TODO(), makeInvalid(func(fs *LustreFileSystem) {
					fs.Spec.Name = "some_really_long_and_invalid_name"
				}))).NotTo(Succeed())
			})

			By("validate 'mgsNid' attribute", func() {
				By("invalid format")
				Expect(k8sClient.Create(context.TODO(), makeInvalid(func(fs *LustreFileSystem) {
					fs.Spec.MgsNid = "this_format_is_missing_an_ampersand"
				}))).NotTo(Succeed())

				By("invalid IP address format")
				Expect(k8sClient.Create(context.TODO(), makeInvalid(func(fs *LustreFileSystem) {
					fs.Spec.MgsNid = "172.0@tcp"
				}))).NotTo(Succeed())

				By("invalid hostname format")
				Expect(k8sClient.Create(context.TODO(), makeInvalid(func(fs *LustreFileSystem) {
					fs.Spec.MgsNid = "this_is_an_invalid_$hostname@tcp"
				}))).NotTo(Succeed())
			})

			By("validate 'mountRoot' attribute", func() {
				By("invalid path")
				Expect(k8sClient.Create(context.TODO(), makeInvalid(func(fs *LustreFileSystem) {
					fs.Spec.MountRoot = "mangled\npath\r"
				}))).NotTo(Succeed())
			})
		})
	})

	Context("Update", func() {
		It("should create and update an object successfully", func() {

		})
	})
})
