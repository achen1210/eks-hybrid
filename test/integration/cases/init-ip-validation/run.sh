#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

source /helpers.sh

# Set up mocks
mock::aws
mock::eks-cluster 172.16.0.0/24
echo "1.2.3.4 mock-hybrid-node" | sudo tee -a /etc/hosts #would cause kubelet to set node ip to 1.2.3.4 via DNS
wait::dbus-ready

# Install nodeadm
nodeadm install 1.30 --credential-provider iam-ra

#Should fail as 1.2.3.4 not in 172.16.0.0/24
exit_code=0
STDERR=$(nodeadm init --config-source file://config.yaml 2>&1) || exit_code=$?
if [ $exit_code -ne 0 ]; then
    assert::is-substring "$STDERR" "node IP 1.2.3.4 is not in any of the remote network CIDR blocks: [172.16.0.0/24]"
else
    echo "nodeadm init should have failed with: node IP 1.2.3.4 is not in any of the remote network CIDR blocks: [172.16.0.0/24]"
    exit 1
fi

