package packagemanager

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/aws/eks-hybrid/internal/artifact"
	"github.com/aws/eks-hybrid/internal/containerd"
	"github.com/aws/eks-hybrid/internal/system"
	"github.com/aws/eks-hybrid/internal/util"
	"github.com/aws/eks-hybrid/internal/util/cmd"
)

// Package manager types and commands
const (
	aptPackageManager  = "apt"
	snapPackageManager = "snap"
	yumPackageManager  = "yum"

	// Snap-specific commands
	snapInstallVerb = "install"
	snapUpdateVerb  = "refresh"
	snapRemoveVerb  = "remove"

	// YUM utilities
	yumUtilsManager    = "yum-config-manager"
	yumUtilsManagerPkg = "yum-utils"

	// Repository URLs and paths
	centOsDockerRepo            = "https://download.docker.com/linux/centos/docker-ce.repo"
	ubuntuDockerRepo            = "https://download.docker.com/linux/ubuntu"
	ubuntuDockerGpgKey          = "https://download.docker.com/linux/ubuntu/gpg"
	ubuntuDockerGpgKeyPath      = "/etc/apt/keyrings/docker.asc"
	ubuntuDockerGpgKeyFilePerms = 0o755
	aptDockerRepoSourceFilePath = "/etc/apt/sources.list.d/docker.list"
	yumDockerRepoSourceFilePath = "/etc/yum.repos.d/docker-ce.repo"

	// Package names
	containerdDistroPkgName = "containerd"
	containerdDockerPkgName = "containerd.io"
	runcPkgName             = "runc"

	caCertsPkgName  = "ca-certificates"
	iptablesPkgName = "iptables"
	ssmPkgName      = "amazon-ssm-agent"
)

// Default containerd repository content
// TODO: Change to include download from closest region
const containerdRepoContent = `deb http://us-west-2.ec2.archive.ubuntu.com/ubuntu jammy-updates main
deb http://security.ubuntu.com/ubuntu jammy-security main
deb http://us-west-2.ec2.archive.ubuntu.com/ubuntu jammy main`

// Maps for package manager commands
var (
	packageManagerInstallCmd = map[string]string{
		aptPackageManager: "install",
		yumPackageManager: "install",
	}

	packageManagerUpdateCmd = map[string]string{
		aptPackageManager: "update",
		yumPackageManager: "update",
	}

	packageManagerDeleteCmd = map[string]string{
		aptPackageManager: "autoremove",
		yumPackageManager: "remove",
	}

	packageManagerMetadataRefreshCmd = map[string]string{
		aptPackageManager: "update",
		yumPackageManager: "makecache",
	}

	managerToDockerRepoMap = map[string]string{
		yumPackageManager: centOsDockerRepo,
		aptPackageManager: ubuntuDockerRepo,
	}
)

// aptDockerRepoConfig generates the apt repository configuration for Docker
var aptDockerRepoConfig = fmt.Sprintf("deb [arch=%s signed-by=%s] %s %s stable\n", runtime.GOARCH, ubuntuDockerGpgKeyPath,
	ubuntuDockerRepo, system.GetVersionCodeName())

// DistroPackageManager defines a new package manager using apt or yum
type DistroPackageManager struct {
	manager             string
	installVerb         string
	updateVerb          string
	deleteVerb          string
	refreshMetadataVerb string
	dockerRepo          string
	logger              *zap.Logger
}

func New(containerdSource containerd.SourceName, logger *zap.Logger) (*DistroPackageManager, error) {
	manager, err := getOsPackageManager()
	if err != nil {
		return nil, err
	}

	pm := &DistroPackageManager{
		manager:             manager,
		logger:              logger,
		installVerb:         packageManagerInstallCmd[manager],
		updateVerb:          packageManagerUpdateCmd[manager],
		deleteVerb:          packageManagerDeleteCmd[manager],
		refreshMetadataVerb: packageManagerMetadataRefreshCmd[manager],
	}
	if containerdSource == containerd.ContainerdSourceDocker {
		pm.dockerRepo = managerToDockerRepoMap[manager]
	}
	return pm, nil
}

// Configure configures the package manager.
func (pm *DistroPackageManager) Configure(ctx context.Context) error {
	// Add docker repos to the package manager
	if pm.dockerRepo != "" {
		return pm.configureDockerRepo(ctx)
	} else if pm.manager == aptPackageManager {
		pm.logger.Info("Updating package metadata for containerd installation")
		return pm.updateContainerdAptPackagesWithRetries(ctx)
	}
	return nil
}

// configureDockerRepo configures the appropriate Docker repositories based on package manager
func (pm *DistroPackageManager) configureDockerRepo(ctx context.Context) error {
	if pm.manager == yumPackageManager {
		return pm.configureYumPackageManagerWithDockerRepo(ctx)
	}
	if pm.manager == aptPackageManager {
		return pm.configureAptPackageManagerWithDockerRepo(ctx)
	}
	return nil
}

// configureYumPackageManagerWithDockerRepo configures yum package manager with docker repos
func (pm *DistroPackageManager) configureYumPackageManagerWithDockerRepo(ctx context.Context) error {
	// Check and remove runc if installed, as it conflicts with docker repo
	if _, errNotFound := exec.LookPath(runcPkgName); errNotFound == nil {
		pm.logger.Info("Removing runc to avoid package conflicts from docker repos...")
		if err := cmd.Retry(ctx, pm.runcPackage().UninstallCmd, 5*time.Second); err != nil {
			return errors.Wrapf(err, "failed to remove runc using package manager")
		}
	}

	// Sometimes install fails due to conflicts with other processes
	// updating packages, specially when automating at machine startup.
	// We assume errors are transient and just retry for a bit.
	if err := cmd.Retry(ctx, pm.yumUtilsPackage().InstallCmd, 5*time.Second); err != nil {
		return errors.Wrapf(err, "failed to install %s using package manager", yumUtilsManagerPkg)
	}

	// Get yumUtilsManager full path
	yumUtilsManagerPath, err := exec.LookPath(yumUtilsManager)
	if err != nil {
		return errors.Wrapf(err, "failed to locate yum utils manager in $PATH")
	}
	pm.logger.Info("Adding docker repo to package manager...")
	configureCmd := exec.Command(yumUtilsManagerPath, "--add-repo", centOsDockerRepo)
	out, err := configureCmd.CombinedOutput()
	if err != nil {
		return errors.Wrapf(err, "failed adding docker repo to package manager: %s", out)
	}

	return nil
}

// configureAptPackageManagerWithDockerRepo configures apt package manager with docker repos
func (pm *DistroPackageManager) configureAptPackageManagerWithDockerRepo(ctx context.Context) error {
	// Sometimes install fails due to conflicts with other processes
	// updating packages, specially when automating at machine startup.
	// We assume errors are transient and just retry for a bit.
	if err := cmd.Retry(ctx, pm.caCertsPackage().InstallCmd, 5*time.Second); err != nil {
		return errors.Wrapf(err, "failed running commands to configure package manager")
	}

	// Download docker gpg key and write it to file
	resp, err := http.Get(ubuntuDockerGpgKey)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := util.WriteFileWithDirFromReader(ubuntuDockerGpgKeyPath, resp.Body, ubuntuDockerGpgKeyFilePerms); err != nil {
		return err
	}

	// Add docker repo config for ubuntu-apt to apt sources
	if err := util.WriteFileWithDir(aptDockerRepoSourceFilePath, []byte(aptDockerRepoConfig), ubuntuDockerGpgKeyFilePerms); err != nil {
		return err
	}

	// Run update to pull docker repo's metadata
	pm.logger.Info("Updating packages to refresh docker repo metadata...")
	err = pm.updateDockerAptPackagesWithRetries(ctx)
	if err != nil {
		return errors.Wrapf(err, "failed running commands to configure package manager")
	}
	return nil
}

// uninstallDockerRepo uninstalls docker repos installed by package managers when containerd source is docker
func (pm *DistroPackageManager) uninstallDockerRepo() error {
	removeRepoFile := func(path, pkgType string) error {
		_, err := os.Stat(path)

		if os.IsNotExist(err) {
			return nil
		}
		if err != nil {
			return errors.Wrapf(err, "encountered error while trying to reach %s docker repo file at %s",
				pkgType, path)
		}

		if err := os.Remove(path); err != nil {
			return errors.Wrapf(err, "failed to remove %s docker repo from %s",
				pkgType, path)
		}

		return nil
	}

	switch pm.manager {
	case yumPackageManager:
		return removeRepoFile(yumDockerRepoSourceFilePath, yumPackageManager)
	case aptPackageManager:
		if err := os.Remove(ubuntuDockerGpgKeyPath); err != nil {
			if !os.IsNotExist(err) {
				return err
			}
		}

		return removeRepoFile(aptDockerRepoSourceFilePath, aptPackageManager)
	default:
		return nil
	}
}

// updateDockerAptPackagesCommand creates a command to update Docker packages
func (pm *DistroPackageManager) updateDockerAptPackagesCommand(ctx context.Context) *exec.Cmd {
	return exec.CommandContext(ctx, pm.manager, pm.updateVerb, "-y", "-o", fmt.Sprintf("Dir::Etc::sourcelist=\"%s\"", aptDockerRepoSourceFilePath))
}

// updateDockerAptPackagesWithRetries retries the update Docker packages command
func (pm *DistroPackageManager) updateDockerAptPackagesWithRetries(ctx context.Context) error {
	return cmd.Retry(ctx, pm.updateDockerAptPackagesCommand, 5*time.Second)
}

// updateContainerdAptPackagesWithRetries updates package info using only the containerd repos
func (pm *DistroPackageManager) updateContainerdAptPackagesWithRetries(ctx context.Context) error {
	// Create an in-memory temp file and use it to update Ubuntu containerd package metadata
	tmpFile, err := os.CreateTemp("", "containerd-sources.*.list")
	if err != nil {
		return errors.Wrap(err, "failed to create temporary containerd sources file")
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(containerdRepoContent); err != nil {
		return errors.Wrap(err, "failed to write to temporary containerd sources file")
	}
	if err := tmpFile.Close(); err != nil {
		return errors.Wrap(err, "failed to close temporary containerd sources file")
	}

	updateCmd := exec.CommandContext(ctx, pm.manager, pm.updateVerb, "-y",
		"-o", fmt.Sprintf("Dir::Etc::sourcelist=%s", tmpFile.Name()),
		"-o", "Dir::Etc::sourceparts=-",
		"-o", "APT::Get::List-Cleanup=0")

	return cmd.Retry(ctx, func(ctx context.Context) *exec.Cmd { return updateCmd }, 5*time.Second)
}

func (pm *DistroPackageManager) appendPackageVersion(packageName, version string) string {
	if version == "" {
		return packageName
	}
	switch pm.manager {
	case yumPackageManager:
		return fmt.Sprintf("%s-%s", packageName, version)
	case aptPackageManager:
		return fmt.Sprintf("%s=%s", packageName, version)
	default:
		return packageName
	}
}

// getContainerdPackageNameWithVersion gets the appropriate containerd package name with version
func (pm *DistroPackageManager) getContainerdPackageNameWithVersion(version string) string {
	containerdPkgName := containerdDistroPkgName
	if pm.dockerRepo != "" {
		containerdPkgName = containerdDockerPkgName
	}
	return pm.appendPackageVersion(containerdPkgName, version)
}

// RefreshMetadataCache refreshes the package managers metadata cache
func (pm *DistroPackageManager) RefreshMetadataCache(ctx context.Context) error {
	pm.logger.Info("Refreshing package metadata cache")
	return cmd.Retry(ctx, pm.refreshMetadataCacheCommand, 5*time.Second)
}

// refreshMetadataCacheCommand creates a command to refresh package metadata
func (pm *DistroPackageManager) refreshMetadataCacheCommand(ctx context.Context) *exec.Cmd {
	return exec.CommandContext(ctx, pm.manager, pm.refreshMetadataVerb)
}

// GetContainerd gets the Package
// Satisfies the containerd source interface
func (pm *DistroPackageManager) GetContainerd(version string) artifact.Package {
	packageName := pm.getContainerdPackageNameWithVersion(version)
	return artifact.NewPackageSource(
		artifact.NewCmd(pm.manager, pm.installVerb, packageName, "-y"),
		artifact.NewCmd(pm.manager, pm.deleteVerb, packageName, "-y"),
		artifact.NewCmd(pm.manager, pm.updateVerb, packageName, "-y"),
	)
}

// GetIptables satisfies the getiptables source interface
func (pm *DistroPackageManager) GetIptables() artifact.Package {
	return artifact.NewPackageSource(
		artifact.NewCmd(pm.manager, pm.installVerb, iptablesPkgName, "-y"),
		artifact.NewCmd(pm.manager, pm.deleteVerb, iptablesPkgName, "-y"),
		artifact.NewCmd(pm.manager, pm.updateVerb, iptablesPkgName, "-y"),
	)
}

// GetSSMPackage satisfies the getssmpackage source interface
func (pm *DistroPackageManager) GetSSMPackage() artifact.Package {
	// SSM is installed using snap package manager. If apt package manager
	// is detected, use snap to install/uninstall SSM.
	if pm.manager == aptPackageManager {
		return artifact.NewPackageSource(
			artifact.NewCmd(snapPackageManager, snapInstallVerb, ssmPkgName),
			artifact.NewCmd(snapPackageManager, snapRemoveVerb, ssmPkgName),
			artifact.NewCmd(snapPackageManager, snapUpdateVerb, ssmPkgName),
		)
	}
	return artifact.NewPackageSource(
		artifact.NewCmd(pm.manager, pm.installVerb, ssmPkgName, "-y"),
		artifact.NewCmd(pm.manager, pm.deleteVerb, ssmPkgName, "-y"),
		artifact.NewCmd(pm.manager, pm.updateVerb, ssmPkgName, "-y"),
	)
}

// caCertsPackage returns a Package object for CA certificates
func (pm *DistroPackageManager) caCertsPackage() artifact.Package {
	return artifact.NewPackageSource(
		artifact.NewCmd(pm.manager, pm.installVerb, caCertsPkgName, "-y"),
		artifact.NewCmd(pm.manager, pm.deleteVerb, caCertsPkgName, "-y"),
		artifact.NewCmd(pm.manager, pm.updateVerb, caCertsPkgName, "-y"),
	)
}

// yumUtilsPackage returns a Package object for yum-utils
func (pm *DistroPackageManager) yumUtilsPackage() artifact.Package {
	return artifact.NewPackageSource(
		artifact.NewCmd(pm.manager, pm.installVerb, yumUtilsManagerPkg, "-y"),
		artifact.NewCmd(pm.manager, pm.deleteVerb, yumUtilsManagerPkg, "-y"),
		artifact.NewCmd(pm.manager, pm.updateVerb, yumUtilsManagerPkg, "-y"),
	)
}

// runcPackage returns a Package object for runc
func (pm *DistroPackageManager) runcPackage() artifact.Package {
	return artifact.NewPackageSource(
		artifact.NewCmd(pm.manager, pm.installVerb, runcPkgName, "-y"),
		artifact.NewCmd(pm.manager, pm.deleteVerb, runcPkgName, "-y"),
		artifact.NewCmd(pm.manager, pm.updateVerb, runcPkgName, "-y"),
	)
}

// Cleanup cleans up any artifacts used by package manager during nodeadm install process
func (pm *DistroPackageManager) Cleanup() error {
	// Removes docker repos if installed by nodeadm ("Containerd: docker" was set in tracker file)
	if pm.dockerRepo != "" {
		pm.logger.Info("Cleaning up Docker repositories")
		if err := pm.uninstallDockerRepo(); err != nil {
			return err
		}
	}

	return nil
}

// getOsPackageManager determines which package manager is available on the system
func getOsPackageManager() (string, error) {
	supportedManagers := []string{yumPackageManager, aptPackageManager}
	for _, manager := range supportedManagers {
		if _, err := exec.LookPath(manager); err == nil {
			return manager, nil
		}
	}
	return "", errors.New("unsupported package manager encountered. Please run nodeadm from a supported os")
}
