package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func init() {
	SchemeBuilder.Register(&NodeConfig{}, &NodeConfigList{})
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Cluster
// +kubebuilder:storageversion

// NodeConfig is the primary configuration object for `nodeadm`.
type NodeConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              NodeConfigSpec `json:"spec,omitempty"`
}

// +kubebuilder:object:root=true

type NodeConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NodeConfig `json:"items"`
}

type NodeConfigSpec struct {
	Cluster    ClusterDetails    `json:"cluster,omitempty"`
	Containerd ContainerdOptions `json:"containerd,omitempty"`
	Instance   InstanceOptions   `json:"instance,omitempty"`
	Kubelet    KubeletOptions    `json:"kubelet,omitempty"`
	Hybrid     *HybridOptions    `json:"hybrid,omitempty"`
	// FeatureGates holds key-value pairs to enable or disable application features.
	FeatureGates map[Feature]bool `json:"featureGates,omitempty"`
}

// ClusterDetails contains the coordinates of your EKS cluster.
// These details can be found using the [DescribeCluster API](https://docs.aws.amazon.com/eks/latest/APIReference/API_DescribeCluster.html).
type ClusterDetails struct {
	// Name is the name of your EKS cluster
	Name string `json:"name,omitempty"`

	// Region is an AWS region (e.g. us-east-1) used to retrieve regional artifacts
	// as well as region where EKS cluster lives.
	Region string `json:"region,omitempty"`

	// APIServerEndpoint is the URL of your EKS cluster's kube-apiserver.
	APIServerEndpoint string `json:"apiServerEndpoint,omitempty"`

	// CertificateAuthority is a base64-encoded string of your cluster's certificate authority chain.
	CertificateAuthority []byte `json:"certificateAuthority,omitempty"`

	// CIDR is your cluster's service CIDR block. This value is used to infer your cluster's DNS address.
	CIDR string `json:"cidr,omitempty"`

	// EnableOutpost determines how your node is configured when running on an AWS Outpost.
	EnableOutpost *bool `json:"enableOutpost,omitempty"`

	// ID is an identifier for your cluster; this is only used when your node is running on an AWS Outpost.
	ID string `json:"id,omitempty"`
}

// KubeletOptions are additional parameters passed to `kubelet`.
type KubeletOptions struct {
	// Config is a [`KubeletConfiguration`](https://kubernetes.io/docs/reference/config-api/kubelet-config.v1beta1/)
	// that will be merged with the defaults.
	Config map[string]runtime.RawExtension `json:"config,omitempty"`

	// Flags are [command-line `kubelet` arguments](https://kubernetes.io/docs/reference/command-line-tools-reference/kubelet/).
	// that will be appended to the defaults.
	Flags []string `json:"flags,omitempty"`
}

// ContainerdOptions are additional parameters passed to `containerd`.
type ContainerdOptions struct {
	// Config is an inline [`containerd` configuration TOML](https://github.com/containerd/containerd/blob/main/docs/man/containerd-config.toml.5.md)
	// that will be merged with the defaults.
	Config string `json:"config,omitempty"`

	// BaseRuntimeSpec is the OCI runtime specification upon which all containers will be based.
	// The provided spec will be merged with the default spec; so that a partial spec may be provided.
	// For more information, see: https://github.com/opencontainers/runtime-spec
	BaseRuntimeSpec map[string]runtime.RawExtension `json:"baseRuntimeSpec,omitempty"`
}

// InstanceOptions determines how the node's operating system and devices are configured.
type InstanceOptions struct {
	LocalStorage LocalStorageOptions `json:"localStorage,omitempty"`
}

// LocalStorageOptions control how [EC2 instance stores](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/InstanceStorage.html)
// are used when available.
type LocalStorageOptions struct {
	Strategy LocalStorageStrategy `json:"strategy,omitempty"`
}

// LocalStorageStrategy specifies how to handle an instance's local storage devices.
// +kubebuilder:validation:Enum={RAID0, RAID10, Mount}
type LocalStorageStrategy string

const (
	// LocalStorageRAID0 will create a single raid0 volume from any local disks
	LocalStorageRAID0 LocalStorageStrategy = "RAID0"

	// LocalStorageRAID10 will create a single raid10 volume from any local disks. Minimum of 4.
	LocalStorageRAID10 LocalStorageStrategy = "RAID10"

	// LocalStorageMount will mount each local disk individually
	LocalStorageMount LocalStorageStrategy = "Mount"
)

// HybridOptions defines the options specific to hybrid node enrollment.
type HybridOptions struct {
	// EnableCredentialsFile enables a shared credentials file on the host at /eks-hybrid/.aws/credentials
	// For SSM, this means that nodeadm will create a symlink from `/root/.aws/credentials` to `/eks-hybrid/.aws/credentials`.
	// For IAM Roles Anywhere, this means that nodeadm will set up a systemd service to write and refresh the credentials to `/eks-hybrid/.aws/credentials`.
	EnableCredentialsFile bool `json:"enableCredentialsFile,omitempty"`

	// IAMRolesAnywhere includes IAM Roles Anywhere specific configuration and is mutually exclusive
	// with SSM.
	IAMRolesAnywhere *IAMRolesAnywhere `json:"iamRolesAnywhere,omitempty"`

	// SSM includes Systems Manager specific configuration and is mutually exclusive with
	// IAMRolesAnywhere.
	SSM *SSM `json:"ssm,omitempty"`
}

// IsHybridNode returns true when the nc.Hybrid configuration is non-nil.
func (nc NodeConfig) IsHybridNode() bool {
	return nc.Spec.Hybrid != nil
}

// IsOutpostNode returns true when Output configuration is non-nil.
func (nc NodeConfig) IsOutpostNode() bool {
	enabled := nc.Spec.Cluster.EnableOutpost
	return enabled != nil && *enabled
}

// IAMRolesAnywhere defines IAM Roles Anywhere specific configuration.
type IAMRolesAnywhere struct {
	// NodeName is the name the node will adopt.
	NodeName string `json:"nodeName,omitempty"`

	// TrustAnchorARN is the ARN of the trust anchor.
	TrustAnchorARN string `json:"trustAnchorArn,omitempty"`

	// ProfileARN is the ARN of the profile linked with the Hybrid IAM Role.
	ProfileARN string `json:"profileArn,omitempty"`

	// RoleARN is the role to IAM roles anywhere gets authorized as to get temporary credentials.
	RoleARN string `json:"roleArn,omitempty"`

	// AwsConfigPath is the path where the Aws config is stored for hybrid nodes.
	// This field is only used to init phase
	// +optional
	AwsConfigPath string `json:"awsConfigPath,omitempty"`

	// CertificatePath is the location on disk for the certificate used to authenticate with AWS.
	// +optional
	CertificatePath string `json:"certificatePath,omitempty"`

	// PrivateKeyPath is the location on disk for the certificate's private key.
	// +optional
	PrivateKeyPath string `json:"privateKeyPath,omitempty"`
}

// SSM defines Systems Manager specific configuration.
// ActivationCode and ActivationID are generated on the aws console or cli during hybrid activations.
// During activation an IAM role is chosen for the SSM agent to assume. This is not overridable from the agent.
type SSM struct {
	// ActivationCode is the token generated when creating an SSM activation.
	ActivationCode string `json:"activationCode,omitempty"`

	// ActivationToken is the ID generated when creating an SSM activation.
	ActivationID string `json:"activationId,omitempty"`
}

// Feature specifies which feature gate should be toggled
// +kubebuilder:validation:Enum={InstanceIdNodeName}
type Feature string

const (
	// InstanceIdNodeName will use EC2 instance ID as node name
	InstanceIdNodeName Feature = "InstanceIdNodeName"
)
