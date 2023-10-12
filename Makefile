# Copyright 2021-2023 Hewlett Packard Enterprise Development LP
# Other additional copyright holders may be indicated within.
#
# The entirety of this work is licensed under the Apache License,
# Version 2.0 (the "License"); you may not use this file except
# in compliance with the License.
#
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Default container tool to use.
#   To use podman:
#   $ DOCKER=podman make docker-build
DOCKER ?= docker

# VERSION defines the project version for the bundle.
# Update this value when you upgrade the version of your project.
# To re-generate a bundle for another specific version without changing the standard setup, you can:
# - use the VERSION as arg of the bundle target (e.g make bundle VERSION=0.0.2)
# - use environment variables to overwrite this value (e.g export VERSION=0.0.2)
# NOTE: git-version-gen will generate a value for VERSION, unless you override it.

# CHANNELS define the bundle channels used in the bundle.
# Add a new line here if you would like to change its default config. (E.g CHANNELS = "candidate,fast,stable")
# To re-generate a bundle for other specific channels without changing the standard setup, you can:
# - use the CHANNELS as arg of the bundle target (e.g make bundle CHANNELS=candidate,fast,stable)
# - use environment variables to overwrite this value (e.g export CHANNELS="candidate,fast,stable")
ifneq ($(origin CHANNELS), undefined)
BUNDLE_CHANNELS := --channels=$(CHANNELS)
endif

# DEFAULT_CHANNEL defines the default channel used in the bundle.
# Add a new line here if you would like to change its default config. (E.g DEFAULT_CHANNEL = "stable")
# To re-generate a bundle for any other default channel without changing the default setup, you can:
# - use the DEFAULT_CHANNEL as arg of the bundle target (e.g make bundle DEFAULT_CHANNEL=stable)
# - use environment variables to overwrite this value (e.g export DEFAULT_CHANNEL="stable")
ifneq ($(origin DEFAULT_CHANNEL), undefined)
BUNDLE_DEFAULT_CHANNEL := --default-channel=$(DEFAULT_CHANNEL)
endif
BUNDLE_METADATA_OPTS ?= $(BUNDLE_CHANNELS) $(BUNDLE_DEFAULT_CHANNEL)

# IMAGE_TAG_BASE defines the docker.io namespace and part of the image name for remote images.
# This variable is used to construct full image tags for bundle and catalog images.
#
# For example, running 'make bundle-build bundle-push catalog-build catalog-push' will build and push both
# cray.hpe.com/lustre-fs-operator-bundle:$VERSION and cray.hpe.com/lustre-fs-operator-catalog:$VERSION.
IMAGE_TAG_BASE ?= ghcr.io/nearnodeflash/lustre-fs-operator

# BUNDLE_IMG defines the image:tag used for the bundle.
# You can use it as an arg. (E.g make bundle-build BUNDLE_IMG=<some-registry>/<project-name-bundle>:<tag>)

# Image URL to use all building/pushing image targets

# ENVTEST_K8S_VERSION refers to the version of kubebuilder assets to be downloaded by envtest binary.
ENVTEST_K8S_VERSION = 1.28.0

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

# Setting SHELL to bash allows bash commands to be executed by recipes.
# This is a requirement for 'setup-envtest.sh' in the test target.
# Options are set to exit when a recipe line exits non-zero or a piped command fails.
SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

all: build

##@ General

# The help target prints out all targets with their descriptions organized
# beneath their categories. The categories are represented by '##@' and the
# target descriptions by '##'. The awk commands is responsible for reading the
# entire set of makefiles included in this invocation, looking for lines of the
# file as xyz: ## something, and then pretty-format the target and help. Then,
# if there's a line with ##@ something, that gets pretty-printed as a category.
# More info on the usage of ANSI control characters for terminal formatting:
# https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_parameters
# More info on the awk command:
# http://linuxcommand.org/lc3_adv_awk.php

help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

manifests: controller-gen ## Generate WebhookConfiguration, ClusterRole and CustomResourceDefinition objects.
	$(CONTROLLER_GEN) rbac:roleName=manager-role crd webhook paths="./..." output:crd:artifacts:config=config/crd/bases

generate: controller-gen ## Generate code containing DeepCopy, DeepCopyInto, and DeepCopyObject method implementations.
	$(CONTROLLER_GEN) object:headerFile="hack/boilerplate.go.txt" paths="./..."

fmt: ## Run go fmt against code.
	go fmt ./...

vet: ## Run go vet against code.
	go vet ./...

# Explicitly specifying directories to test here to avoid running tests in .dws-operator for the time being.
# ./controllers/...
# ./api/...
# Below is a list of ginkgo test flags that may be used to generate different test patterns.
# Specifying 'count=1' is the idiomatic way to disable test caching
#   -ginkgo.debug
#         If set, ginkgo will emit node output to files when running in parallel.
#   -ginkgo.dryRun
#         If set, ginkgo will walk the test hierarchy without actually running anything.  Best paired with -v.
#   -ginkgo.failFast
#         If set, ginkgo will stop running a test suite after a failure occurs.
#   -ginkgo.failOnPending
#         If set, ginkgo will mark the test suite as failed if any specs are pending.
#   -ginkgo.flakeAttempts int
#         Make up to this many attempts to run each spec. Please note that if any of the attempts succeed, the suite will not be failed. But any failures will still be recorded. (default 1)
#   -ginkgo.focus value
#         If set, ginkgo will only run specs that match this regular expression. Can be specified multiple times, values are ORed.
#   -ginkgo.noColor
#         If set, suppress color output in default reporter.
#   -ginkgo.noisyPendings
#         If set, default reporter will shout about pending tests. (default true)
#   -ginkgo.noisySkippings
#         If set, default reporter will shout about skipping tests. (default true)
#   -ginkgo.parallel.node int
#         This worker node's (one-indexed) node number.  For running specs in parallel. (default 1)
#   -ginkgo.parallel.streamhost string
#         The address for the server that the running nodes should stream data to.
#   -ginkgo.parallel.synchost string
#         The address for the server that will synchronize the running nodes.
#   -ginkgo.parallel.total int
#         The total number of worker nodes.  For running specs in parallel. (default 1)
#   -ginkgo.progress
#         If set, ginkgo will emit progress information as each spec runs to the GinkgoWriter.
#   -ginkgo.randomizeAllSpecs
#         If set, ginkgo will randomize all specs together.  By default, ginkgo only randomizes the top level Describe, Context and When groups.
#   -ginkgo.regexScansFilePath
#         If set, ginkgo regex matching also will look at the file path (code location).
#   -ginkgo.reportFile string
#         Override the default reporter output file path.
#   -ginkgo.reportPassed
#         If set, default reporter prints out captured output of passed tests.
#   -ginkgo.seed int
#         The seed used to randomize the spec suite. (default 1635181882)
#   -ginkgo.skip value
#         If set, ginkgo will only run specs that do not match this regular expression. Can be specified multiple times, values are ORed.
#   -ginkgo.skipMeasurements
#         If set, ginkgo will skip any measurement specs.
#   -ginkgo.slowSpecThreshold float
#         (in seconds) Specs that take longer to run than this threshold are flagged as slow by the default reporter. (default 5)
#   -ginkgo.succinct
#         If set, default reporter prints out a very succinct report
#   -ginkgo.trace
#         If set, default reporter prints out the full stack trace when a failure occurs
#   -ginkgo.v
#         If set, default reporter print out all specs as they begin.
#

container-unit-test: VERSION ?= $(shell cat .version)
container-unit-test: .version ## Build docker image with the manager and execute unit tests.
	${DOCKER} build -f Dockerfile --label $(IMAGE_TAG_BASE)-$@:$(VERSION)-$@ -t $(IMAGE_TAG_BASE)-$@:$(VERSION) --target testing .
	${DOCKER} run --rm -t --name $@-lustre-fs-operator $(IMAGE_TAG_BASE)-$@:$(VERSION)

TESTDIR ?= ./controllers/... ./api/...
test: manifests generate fmt vet envtest ## Run tests.
	KUBEBUILDER_ASSETS="$(shell $(ENVTEST) use $(ENVTEST_K8S_VERSION) -p path --bin-dir $(LOCALBIN))" go test $(TESTDIR) -coverprofile cover.out -args -ginkgo.v

##@ Build

build: generate fmt vet ## Build manager binary.
	go build -o bin/manager main.go

run: manifests generate fmt vet ## Run a controller from your host.
	go run ./main.go

docker-build: VERSION ?= $(shell cat .version)
docker-build: .version ## Build docker image with the manager.
	${DOCKER} build -t $(IMAGE_TAG_BASE):$(VERSION) .

edit-image: VERSION ?= $(shell cat .version)
edit-image: .version kustomize ## Replace manager/kustomization.yaml image with name "controller" -> ghcr tagged container reference
	cd config/manager && $(KUSTOMIZE) edit set image controller=$(IMAGE_TAG_BASE):$(VERSION)

docker-push: VERSION ?= $(shell cat .version)
docker-push: .version ## Push docker image with the manager.
	${DOCKER} push $(IMAGE_TAG_BASE):$(VERSION)

kind-push: VERSION ?= $(shell cat .version)
kind-push: .version
	kind load docker-image $(IMAGE_TAG_BASE):$(VERSION)

##@ Deployment

install: manifests kustomize ## Install CRDs into the K8s cluster specified in ~/.kube/config.
	$(KUSTOMIZE) build config/crd | kubectl apply -f -

uninstall: manifests kustomize ## Uninstall CRDs from the K8s cluster specified in ~/.kube/config.
	$(KUSTOMIZE) build config/crd | kubectl delete --ignore-not-found -f -

deploy: kustomize edit-image ## Deploy controller to the K8s cluster specified in ~/.kube/config.
	$(KUSTOMIZE) build config/default | kubectl apply -f -

undeploy: kustomize edit-image ## Undeploy controller from the K8s cluster specified in ~/.kube/config.
	$(KUSTOMIZE) build config/default | kubectl delete --ignore-not-found -f -

installer: kustomize edit-image

# Let .version be phony so that a git update to the workarea can be reflected
# in it each time it's needed.
.PHONY: .version
.version: ## Uses the git-version-gen script to generate a tag version
	./git-version-gen --fallback `git rev-parse HEAD` > .version

clean:
	rm -f .version

## Location to install dependencies to
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

.PHONY: clean-bin
clean-bin:
	if [[ -d $(LOCALBIN) ]]; then \
	  chmod -R u+w $(LOCALBIN) && rm -rf $(LOCALBIN); \
	fi

## Tool Binaries
GO_INSTALL := ./github/cluster-api/scripts/go_install.sh
KUSTOMIZE ?= $(LOCALBIN)/kustomize
CONTROLLER_GEN ?= $(LOCALBIN)/controller-gen
ENVTEST ?= $(LOCALBIN)/setup-envtest

CONVERSION_GEN_BIN := conversion-gen
CONVERSION_GEN := $(LOCALBIN)/$(CONVERSION_GEN_BIN)
CONVERSION_GEN_OUTPUT_BASE := --output-base=.
CONVERSION_GEN_PKG := k8s.io/code-generator/cmd/conversion-gen

CONVERSION_VERIFIER_BIN := conversion-verifier
CONVERSION_VERIFIER := $(LOCALBIN)/$(CONVERSION_VERIFIER_BIN)
CONVERSION_VERIFIER_PKG := sigs.k8s.io/cluster-api/hack/tools/conversion-verifier

## Tool Versions
KUSTOMIZE_VERSION ?= v5.1.1
CONTROLLER_TOOLS_VERSION ?= v0.12.0
CONVERSION_GEN_VER := v0.28.2

# Can be "latest", but cannot be a tag, such as "v1.3.3".  However, it will
# work with the short-form git commit rev that has been tagged.
#CONVERSION_VERIFIER_VER := 09030092b # v1.3.3
CONVERSION_VERIFIER_VER := 3290c5a # v1.5.2

KUSTOMIZE_INSTALL_SCRIPT ?= "https://raw.githubusercontent.com/kubernetes-sigs/kustomize/master/hack/install_kustomize.sh"
.PHONY: kustomize
kustomize: $(LOCALBIN) ## Download kustomize locally if necessary.
	if [[ ! -s $(LOCALBIN)/kustomize || ! $$($(LOCALBIN)/kustomize version) =~ $(KUSTOMIZE_VERSION) ]]; then \
	  rm -f $(LOCALBIN)/kustomize && \
	  { curl -s $(KUSTOMIZE_INSTALL_SCRIPT) | bash -s -- $(subst v,,$(KUSTOMIZE_VERSION)) $(LOCALBIN); }; \
	fi

.PHONY: controller-gen
controller-gen: $(LOCALBIN) ## Download controller-gen locally if necessary.
	if [[ ! -s $(LOCALBIN)/controller-gen || $$($(LOCALBIN)/controller-gen --version | awk '{print $$2}') != $(CONTROLLER_TOOLS_VERSION) ]]; then \
	  rm -f $(LOCALBIN)/controller-gen && GOBIN=$(LOCALBIN) go install sigs.k8s.io/controller-tools/cmd/controller-gen@$(CONTROLLER_TOOLS_VERSION); \
	fi

.PHONY: $(CONVERSION_GEN_BIN)
$(CONVERSION_GEN_BIN): $(CONVERSION_GEN) ## Build a local copy of conversion-gen.

## We are forcing a rebuild of conversion-gen via PHONY so that we're always using an up-to-date version.
## We can't use a versioned name for the binary, because that would be reflected in generated files.
.PHONY: $(CONVERSION_GEN)
$(CONVERSION_GEN): $(LOCALBIN) # Build conversion-gen from tools folder.
	GOBIN=$(LOCALBIN) $(GO_INSTALL) $(CONVERSION_GEN_PKG) $(CONVERSION_GEN_BIN) $(CONVERSION_GEN_VER)

.PHONY: generate-go-conversions
# The SRC_DIRS value is a comma-separated list of paths to old versions.
# The --input-dirs value is a single path item; specify multiple --input-dirs
# parameters if you have multiple old versions.
generate-go-conversions: $(CONVERSION_GEN) ## Generate conversions go code
	$(MAKE) clean-generated-conversions SRC_DIRS="./api/v1alpha1"
	$(CONVERSION_GEN) \
		--input-dirs=./api/v1alpha1 \
		--build-tag=ignore_autogenerated_core \
		--output-file-base=zz_generated.conversion $(CONVERSION_GEN_OUTPUT_BASE) \
		--go-header-file=./hack/boilerplate.go.txt

.PHONY: clean-generated-conversions
clean-generated-conversions: ## Remove files generated by conversion-gen from the mentioned dirs
	(IFS=','; for i in $(SRC_DIRS); do find $$i -type f -name 'zz_generated.conversion*' -exec rm -f {} \;; done)

## We are forcing a rebuild of conversion-verifier via PHONY so that we're always using an up-to-date version.
.PHONY: $(CONVERSION_VERIFIER)
$(CONVERSION_VERIFIER): $(LOCALBIN) # Build conversion-verifier from tools folder.
	GOBIN=$(LOCALBIN) $(GO_INSTALL) $(CONVERSION_VERIFIER_PKG) $(CONVERSION_VERIFIER_BIN) $(CONVERSION_VERIFIER_VER)

.PHONY: $(CONVERSION_VERIFIER_BIN)
$(CONVERSION_VERIFIER_BIN): $(CONVERSION_VERIFIER) ## Build a local copy of conversion-verifier.

## -------------
## verify
## -------------

ALL_VERIFY_CHECKS = gen conversions

.PHONY: verify
verify: $(addprefix verify-,$(ALL_VERIFY_CHECKS)) ## Run all verify-* targets

.PHONY: verify-gen
verify-gen: generate manifests generate-go-conversions ## Verify go generated files are up to date
	@if !(git diff --quiet HEAD); then \
		git diff; \
		echo "generated files are out of date, run make generate"; exit 1; \
	fi

.PHONY: verify-conversions
verify-conversions: $(CONVERSION_VERIFIER)  ## Verifies expected API conversion are in place
	$(CONVERSION_VERIFIER)

.PHONY: envtest
envtest: $(ENVTEST) ## Download envtest-setup locally if necessary.
$(ENVTEST): $(LOCALBIN)
	test -s $(LOCALBIN)/setup-envtest || GOBIN=$(LOCALBIN) go install sigs.k8s.io/controller-runtime/tools/setup-envtest@latest

.PHONY: bundle
bundle: manifests kustomize edit-image ## Generate bundle manifests and metadata, then validate generated files.
	operator-sdk generate kustomize manifests -q
	cd config/manager && $(KUSTOMIZE) edit set image controller=$(IMG)
	$(KUSTOMIZE) build config/manifests | operator-sdk generate bundle -q --overwrite --version $(VERSION) $(BUNDLE_METADATA_OPTS)
	operator-sdk bundle validate ./bundle

.PHONY: bundle-build
bundle-build: VERSION ?= $(shell cat .version)
bundle-build: BUNDLE_IMG ?= $(IMAGE_TAG_BASE)-bundle:v$(VERSION)
bundle-build: .version ## Build the bundle image.
	${DOCKER} build -f bundle.Dockerfile -t $(BUNDLE_IMG) .

.PHONY: bundle-push
bundle-push: VERSION ?= $(shell cat .version)
bundle-push: BUNDLE_IMG ?= $(IMAGE_TAG_BASE)-bundle:v$(VERSION)
bundle-push: .version ## Push the bundle image.
	$(MAKE) docker-push IMG=$(BUNDLE_IMG)

.PHONY: opm
OPM = ./bin/opm
opm: ## Download opm locally if necessary.
ifeq (,$(wildcard $(OPM)))
ifeq (,$(shell which opm 2>/dev/null))
	@{ \
	set -e ;\
	mkdir -p $(dir $(OPM)) ;\
	OS=$(shell go env GOOS) && ARCH=$(shell go env GOARCH) && \
	curl -sSLo $(OPM) https://github.com/operator-framework/operator-registry/releases/download/v1.15.1/$${OS}-$${ARCH}-opm ;\
	chmod +x $(OPM) ;\
	}
else
OPM = $(shell which opm)
endif
endif

# A comma-separated list of bundle images (e.g. make catalog-build BUNDLE_IMGS=example.com/operator-bundle:v0.1.0,example.com/operator-bundle:v0.2.0).
# These images MUST exist in a registry and be pull-able.
BUNDLE_IMGS ?= $(BUNDLE_IMG)

# The image tag given to the resulting catalog image (e.g. make catalog-build CATALOG_IMG=example.com/operator-catalog:v0.2.0).
CATALOG_IMG ?= $(IMAGE_TAG_BASE)-catalog:v$(VERSION)

# Set CATALOG_BASE_IMG to an existing catalog image tag to add $BUNDLE_IMGS to that image.
ifneq ($(origin CATALOG_BASE_IMG), undefined)
FROM_INDEX_OPT := --from-index $(CATALOG_BASE_IMG)
endif

# Build a catalog image by adding bundle images to an empty catalog using the operator package manager tool, 'opm'.
# This recipe invokes 'opm' in 'semver' bundle add mode. For more information on add modes, see:
# https://github.com/operator-framework/community-operators/blob/7f1438c/docs/packaging-operator.md#updating-your-existing-operator
.PHONY: catalog-build
catalog-build: opm ## Build a catalog image.
	$(OPM) index add --container-tool ${DOCKER} --mode semver --tag $(CATALOG_IMG) --bundles $(BUNDLE_IMGS) $(FROM_INDEX_OPT)

# Push the catalog image.
.PHONY: catalog-push
catalog-push: ## Push a catalog image.
	$(MAKE) docker-push IMG=$(CATALOG_IMG)
