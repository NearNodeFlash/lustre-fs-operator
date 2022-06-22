module github.com/NearNodeFlash/lustre-fs-operator

go 1.16

replace github.com/HewlettPackard/lustre-csi-driver => ../lustre-csi-driver

require (
	github.com/HewlettPackard/lustre-csi-driver v0.0.0-20220516192757-17ac28565db5
	github.com/go-logr/logr v0.4.0
	github.com/onsi/ginkgo/v2 v2.1.4
	github.com/onsi/gomega v1.19.0
	k8s.io/api v0.22.1
	k8s.io/apimachinery v0.22.1
	k8s.io/client-go v0.22.1
	sigs.k8s.io/controller-runtime v0.10.0
)
