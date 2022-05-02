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
	"context"

	. "github.com/onsi/ginkgo/v2"
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
