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
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// These tests are written in BDD-style using Ginkgo framework. Refer to
// http://onsi.github.io/ginkgo to learn more.

var _ = Describe("LustreFileSystemWebhook", func() {
	var (
		key                    types.NamespacedName
		createdFS, retrievedFS *LustreFileSystem
	)

	BeforeEach(func() {
		key = types.NamespacedName{
			Name:      "lustre-fs-example",
			Namespace: "default",
		}

		createdFS = &LustreFileSystem{
			ObjectMeta: metav1.ObjectMeta{
				Name:      key.Name,
				Namespace: key.Namespace,
			},
			Spec: LustreFileSystemSpec{
				Name:      "foo",
				MgsNids:   "127.0.0.1@tcp",
				MountRoot: "/lus/foo",
			},
		}

		retrievedFS = &LustreFileSystem{}
	})

	AfterEach(func() {
		if createdFS != nil {
			Expect(k8sClient.Delete(context.TODO(), createdFS)).To(Succeed())
			expectedFS := &LustreFileSystem{}
			Eventually(func() error { // Delete can still return the cached object. Wait until the object is no longer present
				return k8sClient.Get(context.TODO(), client.ObjectKeyFromObject(createdFS), expectedFS)
			}).ShouldNot(Succeed())
		}
	})

	Context("Create", func() {
		It("should create an object successfully, with IP", func() {
			By("creating an object")
			Expect(k8sClient.Create(context.TODO(), createdFS)).To(Succeed())
		})

		It("should create an object successfully, with hostname", func() {
			By("creating an object")
			createdFS.Spec.MgsNids = "localhost@tcp"
			Expect(k8sClient.Create(context.TODO(), createdFS)).To(Succeed())
		})

		It("should create an object successfully, with complex nid list", func() {
			By("creating an object")
			createdFS.Spec.MgsNids = "localhost@tcp,kingkong@tcp:red@tcp,blue@tcp:cat@tcp"
			Expect(k8sClient.Create(context.TODO(), createdFS)).To(Succeed())
		})

		It("should allow an update to the metadata", func() {
			By("creating an object")
			Expect(k8sClient.Create(context.TODO(), createdFS)).To(Succeed())

			Expect(k8sClient.Get(context.TODO(), key, retrievedFS)).To(Succeed())

			By("updating the object")
			// A finalizer or ownerRef will interfere with
			// deletion, so set a label, instead.
			labels := retrievedFS.GetLabels()
			if labels == nil {
				labels = make(map[string]string)
			}
			labels["fs-label"] = "fs-label"
			retrievedFS.SetLabels(labels)
			Expect(k8sClient.Update(context.TODO(), retrievedFS)).To(Succeed())
		})
	})

	Context("Negatives", func() {

		It("should fail with empty name attribute", func() {
			createdFS.Spec.Name = ""
			Expect(k8sClient.Create(context.TODO(), createdFS)).NotTo(Succeed())
			createdFS = nil
		})

		It("should fail with an overflowing name attribute", func() {
			createdFS.Spec.Name = "some_really_long_and_invalid_name"
			Expect(k8sClient.Create(context.TODO(), createdFS)).NotTo(Succeed())
			createdFS = nil
		})

		It("should fail with an invalid 'mgsNid' attribute", func() {
			By("invalid format")
			createdFS.Spec.MgsNids = "this_format_is_missing_an_ampersand"
			Expect(k8sClient.Create(context.TODO(), createdFS)).NotTo(Succeed())
			createdFS = nil
		})

		It("should fail with an invalid nid in a complex nid list", func() {
			createdFS.Spec.MgsNids = "localhost@tcp,kingkong@tcp:red@tcp,this_format_is_missing_an_ampersand"
			Expect(k8sClient.Create(context.TODO(), createdFS)).NotTo(Succeed())
			createdFS = nil
		})

		It("should fail with an invalid 'mountRoot' attribute", func() {
			createdFS.Spec.MountRoot = "mangled\npath\r"
			Expect(k8sClient.Create(context.TODO(), createdFS)).NotTo(Succeed())
			createdFS = nil
		})

		It("should fail to update the spec", func() {
			By("creating an object")
			Expect(k8sClient.Create(context.TODO(), createdFS)).To(Succeed())

			retrievedFS = &LustreFileSystem{}
			Expect(k8sClient.Get(context.TODO(), key, retrievedFS)).To(Succeed())

			By("updating the object")
			retrievedFS.Spec.MountRoot = "/lus/other_mount"
			Expect(k8sClient.Update(context.TODO(), retrievedFS)).ToNot(Succeed())
		})
	})
})
