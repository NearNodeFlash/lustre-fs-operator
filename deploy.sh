#!/bin/bash

# Copyright 2023 Hewlett Packard Enterprise Development LP
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

set -e

# Deploy/undeploy controller to the K8s cluster specified in ~/.kube/config.

CMD=$1
KUSTOMIZE=$2
OVERLAY_DIR=$3

if [[ $CMD == 'deploy' ]]; then
    $KUSTOMIZE build $OVERLAY_DIR | kubectl apply -f -

    # Deploy the ServiceMonitor resource if its CRD is found. The CRD would
    # have been installed by a metrics service such as Prometheus.
    if kubectl get crd servicemonitors.monitoring.coreos.com > /dev/null 2>&1; then
        $KUSTOMIZE build config/prometheus | kubectl apply -f-
    fi
fi

if [[ $CMD == 'undeploy' ]]; then
    $KUSTOMIZE build config/prometheus | kubectl delete --ignore-not-found -f-
    $KUSTOMIZE build $OVERLAY_DIR | kubectl delete --ignore-not-found -f -
fi

exit 0
