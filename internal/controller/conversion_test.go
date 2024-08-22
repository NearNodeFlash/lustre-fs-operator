/*
 * Copyright 2023-2024 Hewlett Packard Enterprise Development LP
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
	"context"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	lusv1alpha1 "github.com/NearNodeFlash/lustre-fs-operator/api/v1alpha1"
	lusv1beta1 "github.com/NearNodeFlash/lustre-fs-operator/api/v1beta1"
	utilconversion "github.com/NearNodeFlash/lustre-fs-operator/github/cluster-api/util/conversion"
)

var _ = Describe("Conversion Webhook Test", func() {

	// Don't get deep into verifying the conversion.
	// We have api/<spoke_ver>/conversion_test.go that is digging deep.
	// We're just verifying that the conversion webhook is hooked up.

	// Note: if a resource is accessed by its spoke API, then it should
	// have the utilconversion.DataAnnotation annotation.  It will not
	// have that annotation when it is accessed by its hub API.

	Context("LustreFileSystem", func() {
		var resHub *lusv1beta1.LustreFileSystem

		BeforeEach(func() {
			id := uuid.NewString()[0:8]
			resHub = &lusv1beta1.LustreFileSystem{
				ObjectMeta: metav1.ObjectMeta{
					Name:      id,
					Namespace: corev1.NamespaceDefault,
				},
				Spec: lusv1beta1.LustreFileSystemSpec{
					Name:      "w0",
					MgsNids:   "172.0.0.0@tcp",
					MountRoot: "/lus/w0",
				},
			}

			Expect(k8sClient.Create(context.TODO(), resHub)).To(Succeed())
		})

		AfterEach(func() {
			if resHub != nil {
				Expect(k8sClient.Delete(context.TODO(), resHub)).To(Succeed())
				expected := &lusv1beta1.LustreFileSystem{}
				Eventually(func() error { // Delete can still return the cached object. Wait until the object is no longer present.
					return k8sClient.Get(context.TODO(), client.ObjectKeyFromObject(resHub), expected)
				}).ShouldNot(Succeed())
			}
		})

		It("reads LustreFileSystem resource via hub and via spoke v1alpha1", func() {
			// Spoke should have annotation.
			resSpoke := &lusv1alpha1.LustreFileSystem{}
			Eventually(func(g Gomega) {
				g.Expect(k8sClient.Get(context.TODO(), client.ObjectKeyFromObject(resHub), resSpoke)).To(Succeed())
				anno := resSpoke.GetAnnotations()
				g.Expect(anno).To(HaveLen(1))
				g.Expect(anno).Should(HaveKey(utilconversion.DataAnnotation))
			}).Should(Succeed())

			// Hub should not have annotation.
			Eventually(func(g Gomega) {
				g.Expect(k8sClient.Get(context.TODO(), client.ObjectKeyFromObject(resHub), resHub)).To(Succeed())
				anno := resHub.GetAnnotations()
				g.Expect(anno).To(HaveLen(0))
			}).Should(Succeed())
		})

		// +crdbumper:scaffold:spoketest="lus.LustreFileSystem"
	})

	// +crdbumper:scaffold:webhooksuitetest
})
