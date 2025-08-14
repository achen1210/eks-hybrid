#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

source /helpers.sh

mock::aws
wait::dbus-ready

mock::kubelet 1.21.0
nodeadm init --skip run,install-validation,k8s-authentication-validation --config-source file://config.yaml
assert::file-not-contains /etc/kubernetes/kubelet/config.json '"kubeAPIQPS"'
assert::file-not-contains /etc/kubernetes/kubelet/config.json '"kubeAPIBurst"'

mock::kubelet 1.22.0-eks-5e0fdde
nodeadm init --skip run,install-validation,k8s-authentication-validation --config-source file://config.yaml
assert::file-contains /etc/kubernetes/kubelet/config.json '"kubeAPIQPS": 10'
assert::file-contains /etc/kubernetes/kubelet/config.json '"kubeAPIBurst": 20'

mock::kubelet 1.22.0
nodeadm init --skip run,install-validation,k8s-authentication-validation --config-source file://config.yaml
assert::file-contains /etc/kubernetes/kubelet/config.json '"kubeAPIQPS": 10'
assert::file-contains /etc/kubernetes/kubelet/config.json '"kubeAPIBurst": 20'

mock::kubelet 1.26.0-eks-5e0fdde
nodeadm init --skip run,install-validation,k8s-authentication-validation --config-source file://config.yaml
assert::file-contains /etc/kubernetes/kubelet/config.json '"kubeAPIQPS": 10'
assert::file-contains /etc/kubernetes/kubelet/config.json '"kubeAPIBurst": 20'

mock::kubelet 1.26.0
nodeadm init --skip run,install-validation,k8s-authentication-validation --config-source file://config.yaml
assert::file-contains /etc/kubernetes/kubelet/config.json '"kubeAPIQPS": 10'
assert::file-contains /etc/kubernetes/kubelet/config.json '"kubeAPIBurst": 20'

mock::kubelet 1.28.0-eks-5e0fdde
nodeadm init --skip run,install-validation,k8s-authentication-validation --config-source file://config.yaml
assert::file-not-contains /etc/kubernetes/kubelet/config.json '"kubeAPIQPS"'
assert::file-not-contains /etc/kubernetes/kubelet/config.json '"kubeAPIBurst"'
