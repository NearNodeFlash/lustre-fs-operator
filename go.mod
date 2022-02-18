module github.hpe.com/hpe/hpc-rabsw-lustre-fs-operator

go 1.16

require (
	github.com/go-logr/logr v0.4.0
	github.com/onsi/ginkgo v1.16.4
	github.com/onsi/gomega v1.15.0
	github.hpe.com/hpe/hpc-rabsw-lustre-csi-driver v0.0.0-20220217210743-a86c23a50c7a
	k8s.io/api v0.22.1
	k8s.io/apimachinery v0.22.1
	k8s.io/client-go v0.22.1
	sigs.k8s.io/controller-runtime v0.10.0
)
