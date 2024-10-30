#!/usr/bin/env bash

set -o pipefail
set -o nounset
set -o errexit

if [ "$ENABLE_ACCELERATOR" != "nvidia" ]; then
  exit 0
fi

#Detect Isolated partitions
function is-isolated-partition() {
  PARTITION=$(imds /latest/meta-data/services/partition)
  NON_ISOLATED_PARTITIONS=("aws" "aws-cn" "aws-us-gov")
  for NON_ISOLATED_PARTITION in "${NON_ISOLATED_PARTITIONS[@]}"; do
    if [ "${NON_ISOLATED_PARTITION}" = "${PARTITION}" ]; then
      return 1
    fi
  done
  return 0
}

echo "Installing NVIDIA ${NVIDIA_DRIVER_MAJOR_VERSION} drivers..."

################################################################################
### Add repository #############################################################
################################################################################
# Determine the domain based on the region
if is-isolated-partition; then
  echo '[amzn2023-nvidia]
  name=Amazon Linux 2023 Nvidia repository
  mirrorlist=https://al2023-repos-$awsregion-de612dc2.s3.$awsregion.$awsdomain/nvidia/mirrors/$releasever/$basearch/mirror.list
  priority=20
  enabled=1
  repo_gpgcheck=0
  type=rpm
  gpgcheck=0
  gpgkey=file:///etc/pki/rpm-gpg/RPM-GPG-KEY-amazon-linux-2023' | sudo tee /etc/yum.repos.d/amzn2023-nvidia.repo

else
  if [[ $AWS_REGION == cn-* ]]; then
    DOMAIN="nvidia.cn"
  else
    DOMAIN="nvidia.com"
  fi

  sudo dnf config-manager --add-repo https://developer.download.${DOMAIN}/compute/cuda/repos/amzn2023/x86_64/cuda-amzn2023.repo
  sudo dnf config-manager --add-repo https://nvidia.github.io/libnvidia-container/stable/rpm/nvidia-container-toolkit.repo

  sudo sed -i 's/gpgcheck=0/gpgcheck=1/g' /etc/yum.repos.d/nvidia-container-toolkit.repo /etc/yum.repos.d/cuda-amzn2023.repo
fi

################################################################################
### Install drivers ############################################################
################################################################################
sudo mv ${WORKING_DIR}/gpu/gpu-ami-util /usr/bin/
sudo mv ${WORKING_DIR}/gpu/kmod-util /usr/bin/

sudo mkdir -p /etc/dkms
echo "MAKE[0]=\"'make' -j$(grep -c processor /proc/cpuinfo) module\"" | sudo tee /etc/dkms/nvidia.conf
sudo dnf -y install kernel-modules-extra.x86_64

function archive-open-kmods() {
  sudo dnf -y module install nvidia-driver:${NVIDIA_DRIVER_MAJOR_VERSION}-open
  # The DKMS package name differs between the RPM and the dkms.conf in the OSS kmod sources
  # TODO: can be removed if this is merged: https://github.com/NVIDIA/open-gpu-kernel-modules/pull/567
  sudo sed -i 's/PACKAGE_NAME="nvidia"/PACKAGE_NAME="nvidia-open"/g' /var/lib/dkms/nvidia-open/$(kmod-util module-version nvidia-open)/source/dkms.conf

  sudo kmod-util archive nvidia-open

  KMOD_MAJOR_VERSION=$(sudo kmod-util module-version nvidia-open | cut -d. -f1)
  SUPPORTED_DEVICE_FILE="${WORKING_DIR}/gpu/nvidia-open-supported-devices-${KMOD_MAJOR_VERSION}.txt"
  sudo mv "${SUPPORTED_DEVICE_FILE}" /etc/eks/

  sudo kmod-util remove nvidia-open

  sudo dnf -y module remove --all nvidia-driver
  sudo dnf -y module reset nvidia-driver
}

function archive-proprietary-kmod() {
  sudo dnf -y module install nvidia-driver:${NVIDIA_DRIVER_MAJOR_VERSION}-dkms
  sudo kmod-util archive nvidia
  sudo kmod-util remove nvidia
}

archive-open-kmods
archive-proprietary-kmod

################################################################################
### Prepare for nvidia init ####################################################
################################################################################

sudo mv ${WORKING_DIR}/gpu/nvidia-kmod-load.sh /etc/eks/
sudo mv ${WORKING_DIR}/gpu/set-nvidia-clocks.sh /etc/eks/
sudo mv ${WORKING_DIR}/gpu/nvidia-kmod-load.service /etc/systemd/system/nvidia-kmod-load.service
sudo mv ${WORKING_DIR}/gpu/set-nvidia-clocks.service /etc/systemd/system/set-nvidia-clocks.service
sudo systemctl daemon-reload
sudo systemctl enable nvidia-kmod-load.service
sudo systemctl enable set-nvidia-clocks.service

################################################################################
### Install other dependencies #################################################
################################################################################
sudo dnf -y install nvidia-fabric-manager nvidia-container-toolkit

sudo systemctl enable nvidia-fabricmanager
sudo systemctl enable nvidia-persistenced
