#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

source /helpers.sh

mock::aws
mock::kubelet 1.28.0
wait::dbus-ready

for config in config.*; do
  nodeadm init --skip run,install-validation --config-source file://${config}
  assert::json-files-equal /etc/kubernetes/kubelet/config.json expected-kubelet-config.json
done
