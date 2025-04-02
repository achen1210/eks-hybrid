package containerd

import (
	"bytes"
	_ "embed"
	"text/template"

	"github.com/pelletier/go-toml/v2"
	"go.uber.org/zap"
	"golang.org/x/mod/semver"

	"github.com/aws/eks-hybrid/internal/api"
	"github.com/aws/eks-hybrid/internal/util"
)

const ContainerRuntimeEndpoint = "unix:///run/containerd/containerd.sock"

const (
	containerdConfigDir               = "/etc/containerd"
	containerdConfigFile              = "/etc/containerd/config.toml"
	containerdConfigImportDir         = "/etc/containerd/config.d"
	containerdKernelModulesConfigFile = "/etc/modules-load.d/containerd.conf"
	containerdConfigPerm              = 0o644
)

var (
	//go:embed config.template.toml
	containerdConfigTemplateData string
	containerdConfigTemplate     = template.Must(template.New(containerdConfigFile).Parse(containerdConfigTemplateData))

	//go:embed kernel-modules.conf
	containerdKernelModulesFileData string
)

type containerdTemplateVars struct {
	EnableCDI         bool
	SandboxImage      string
	RuntimeName       string
	RuntimeBinaryName string
}

func writeContainerdConfig(cfg *api.NodeConfig) error {
	if err := writeBaseRuntimeSpec(cfg); err != nil {
		return err
	}

	containerdConfig, err := generateContainerdConfig(cfg)
	if err != nil {
		return err
	}

	// because the logic in containerd's import merge decides to completely
	// overwrite entire sections, we want to implement this merging ourselves.
	// see: https://github.com/containerd/containerd/blob/a91b05d99ceac46329be06eb43f7ae10b89aad45/cmd/containerd/server/config/config.go#L407-L431
	if len(cfg.Spec.Containerd.Config) > 0 {
		containerdConfigMap, err := util.Merge(containerdConfig, []byte(cfg.Spec.Containerd.Config), toml.Marshal, toml.Unmarshal)
		if err != nil {
			return err
		}
		containerdConfig, err = toml.Marshal(containerdConfigMap)
		if err != nil {
			return err
		}
	}

	zap.L().Info("Writing containerd config to file..", zap.String("path", containerdConfigFile))
	return util.WriteFileWithDir(containerdConfigFile, containerdConfig, containerdConfigPerm)
}

func generateContainerdConfig(cfg *api.NodeConfig) ([]byte, error) {
	instanceOptions := applyInstanceTypeMixins(cfg.Status.Instance.Type)

	configVars := containerdTemplateVars{
		SandboxImage:      cfg.Status.Defaults.SandboxImage,
		RuntimeBinaryName: instanceOptions.RuntimeBinaryName,
		RuntimeName:       instanceOptions.RuntimeName,
		EnableCDI:         semver.Compare(cfg.Status.KubeletVersion, "v1.32.0") >= 0,
	}
	var buf bytes.Buffer
	if err := containerdConfigTemplate.Execute(&buf, configVars); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func writeContainerdKernelModulesConfig() error {
	return util.WriteFileWithDir(containerdKernelModulesConfigFile, []byte(containerdKernelModulesFileData), containerdConfigPerm)
}
