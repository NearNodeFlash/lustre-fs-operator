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
		key                types.NamespacedName
		created, retrieved *LustreFileSystem
	)

	BeforeEach(func() {
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
	})

	Context("Create", func() {
		It("should create an object successfully, with IP", func() {
			By("creating an object")
			Expect(k8sClient.Create(context.TODO(), created)).To(Succeed())

			retrieved = &LustreFileSystem{}
			Expect(k8sClient.Get(context.TODO(), key, retrieved)).To(Succeed())
			Expect(retrieved).To(Equal(created))

			By("deleting the object")
			Expect(k8sClient.Delete(context.TODO(), created)).To(Succeed())
			Expect(k8sClient.Get(context.TODO(), key, created)).ToNot(Succeed())
		})

		It("should create an object successfully, with hostname", func() {
			By("creating an object")
			created.Spec.MgsNid = "localhost@tcp"
			Expect(k8sClient.Create(context.TODO(), created)).To(Succeed())

			retrieved = &LustreFileSystem{}
			Expect(k8sClient.Get(context.TODO(), key, retrieved)).To(Succeed())
			Expect(retrieved).To(Equal(created))

			By("deleting the object")
			Expect(k8sClient.Delete(context.TODO(), created)).To(Succeed())
			Expect(k8sClient.Get(context.TODO(), key, created)).ToNot(Succeed())
		})
	})

	Context("Negatives", func() {

		It("should fail with empty name attribute", func() {
			created.Spec.Name = ""
			Expect(k8sClient.Create(context.TODO(), created)).NotTo(Succeed())
		})

		It("should fail with an overflowing name attribute", func() {
			created.Spec.Name = "some_really_long_and_invalid_name"
			Expect(k8sClient.Create(context.TODO(), created)).NotTo(Succeed())
		})

		It("should fail with an invalid 'mgsNid' attribute", func() {
			By("invalid format")
			created.Spec.MgsNid = "this_format_is_missing_an_ampersand"
			Expect(k8sClient.Create(context.TODO(), created)).NotTo(Succeed())
		})

		It("should fail with an invalid 'mountRoot' attribute", func() {
			created.Spec.MountRoot = "mangled\npath\r"
			Expect(k8sClient.Create(context.TODO(), created)).NotTo(Succeed())
		})
	})

	Context("Update", func() {
		It("should create and update an object successfully", func() {

		})
	})
})
