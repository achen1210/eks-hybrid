package kubelet

import (
	"fmt"
	"strings"

	"github.com/aws/eks-hybrid/internal/util"
)

const (
	kubeletEnvironmentFilePath = "/etc/eks/kubelet/environment"
	KubeletArgsEnvironmentName = "NODEADM_KUBELET_ARGS"
)

// Write environment variables needed for kubelet runtime. This should be the
// last method called on the kubelet object so that environment side effects of
// other methods are properly recorded
func (k *Kubelet) writeKubeletEnvironment() error {
	// transform Kubelet flags into a single string and write them to the
	// Kubelet environment variable
	var kubeletFlags []string
	for flag, value := range k.flags {
		kubeletFlags = append(kubeletFlags, fmt.Sprintf("--%s=%s", flag, value))
	}
	// append user-provided flags at the end to give them precedence
	kubeletFlags = append(kubeletFlags, k.nodeConfig.Spec.Kubelet.Flags...)
	// expose these flags via an environment variable scoped to nodeadm
	k.environment[KubeletArgsEnvironmentName] = strings.Join(kubeletFlags, " ")
	// write additional environment variables
	var kubeletEnvironment []string
	for eKey, eValue := range k.environment {
		kubeletEnvironment = append(kubeletEnvironment, fmt.Sprintf(`%s="%s"`, eKey, eValue))
	}
	return util.WriteFileWithDir(kubeletEnvironmentFilePath, []byte(strings.Join(kubeletEnvironment, "\n")), kubeletConfigPerm)
}

// Add values to the environment variables map in a terse manner
func (k *Kubelet) setEnv(envName, envArg string) {
	k.environment[envName] = envArg
}

func (k *Kubelet) GetEnv(envName string) string {
	return k.environment[envName]
}
