package containerd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	containerdConfigDumpFragment = `
[plugins]

  [plugins."io.containerd.gc.v1.scheduler"]
    deletion_threshold = 0
    mutation_threshold = 100
    pause_threshold = 0.02
    schedule_delay = "0s"
    startup_delay = "100ms"

  [plugins."io.containerd.grpc.v1.cri"]
    cdi_spec_dirs = ["/etc/cdi", "/var/run/cdi"]
    device_ownership_from_security_context = false
    disable_apparmor = false
    disable_cgroup = false
    disable_hugetlb_controller = true
    disable_proc_mount = false
    disable_tcp_service = true
    drain_exec_sync_io_timeout = "0s"
    enable_cdi = false
    enable_selinux = false
    enable_tls_streaming = false
    enable_unprivileged_icmp = false
    enable_unprivileged_ports = false
    ignore_image_defined_volumes = false
    image_pull_progress_timeout = "1m0s"
    max_concurrent_downloads = 3
    max_container_log_line_size = 16384
    netns_mounts_under_state_dir = false
    restrict_oom_score_adj = false
    sandbox_image = "registry.k8s.io/pause:3.8"
    selinux_category_range = 1024
    stats_collect_period = 10
    stream_idle_timeout = "4h0m0s"
    stream_server_address = "127.0.0.1"
    stream_server_port = "0"
    systemd_cgroup = false
    tolerate_missing_hugetlb_controller = true
    unset_seccomp_profile = ""
`

	// Test data for parseConfigVersion function
	containerdConfigV2 = `version = 2
[plugins."io.containerd.grpc.v1.cri"]
    sandbox_image = "registry.k8s.io/pause:3.8"`

	containerdConfigV3 = `version = 3
[plugins."io.containerd.grpc.v1.cri"]
    sandbox = "registry.k8s.io/pause:3.8"`

	containerdConfigNoVersion = `[plugins."io.containerd.grpc.v1.cri"]
    sandbox_image = "registry.k8s.io/pause:3.8"`

	containerdConfigInvalidVersion = `version = invalid
[plugins."io.containerd.grpc.v1.cri"]
    sandbox_image = "registry.k8s.io/pause:3.8"`
)

func TestSandboxImageV2Regex(t *testing.T) {
	matches := containerdSandboxImageV2Regex.FindStringSubmatch(containerdConfigDumpFragment)
	if matches == nil {
		t.Errorf("sandbox image could not be found in containerd config")
	}
	sandboxImage := matches[1]
	assert.Equal(t, sandboxImage, "registry.k8s.io/pause:3.8")
}

func TestParseConfigVersion(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    int
		expectError bool
	}{
		{
			name:        "valid v2 config",
			input:       containerdConfigV2,
			expected:    2,
			expectError: false,
		},
		{
			name:        "valid v3 config",
			input:       containerdConfigV3,
			expected:    3,
			expectError: false,
		},
		{
			name:        "config without version",
			input:       containerdConfigNoVersion,
			expected:    0,
			expectError: true,
		},
		{
			name:        "config with invalid version",
			input:       containerdConfigInvalidVersion,
			expected:    0,
			expectError: true,
		},
		{
			name:        "empty input",
			input:       "",
			expected:    0,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseConfigVersion([]byte(tt.input))

			if tt.expectError {
				assert.Error(t, err, "expected error for test case: %s", tt.name)
				assert.Equal(t, tt.expected, result, "expected result should be 0 for error cases")
			} else {
				assert.NoError(t, err, "unexpected error for test case: %s", tt.name)
				assert.Equal(t, tt.expected, result, "incorrect version parsed for test case: %s", tt.name)
			}
		})
	}
}
