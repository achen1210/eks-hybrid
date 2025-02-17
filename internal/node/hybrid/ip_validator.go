package hybrid

import (
	"fmt"
	"net"
	"strings"

	apimachinerynet "k8s.io/apimachinery/pkg/util/net"
	nodeutil "k8s.io/component-helpers/node/util"
	k8snet "k8s.io/utils/net"

	"github.com/aws/eks-hybrid/internal/aws/eks"
)

const (
	nodeIPFlag           = "--node-ip="
	hostnameOverrideFlag = "--hostname-override="
)

func containsIP(cidr string, ip net.IP) (bool, error) {
	_, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return false, err
	}

	return ipnet.Contains(ip), nil
}

func isIPInClusterNetworks(ip net.IP, remoteNetworkConfig *eks.RemoteNetworkConfig) (bool, error) {
	if ip.To4() == nil || remoteNetworkConfig == nil {
		return false, fmt.Errorf("error: ip is invalid or remoteNetworkConfig is nil")
	}

	for _, network := range remoteNetworkConfig.RemoteNodeNetworks {
		for _, cidr := range network.CIDRs {
			if cidr == nil {
				continue
			}

			if inNetwork, err := containsIP(*cidr, ip); err != nil {
				return false, fmt.Errorf("error checking IP in CIDR %s: %w", *cidr, err)
			} else if inNetwork {
				return true, nil
			}
		}
	}

	return false, nil
}

func validateIP(ipAddr net.IP, hnp *HybridNodeProvider) error {
	if validIP, err := isIPInClusterNetworks(ipAddr, hnp.remoteNetworkConfig); err != nil {
		return err
	} else if !validIP {
		cidrs := getClusterCIDRs(hnp.remoteNetworkConfig)

		return fmt.Errorf(
			"node IP %s is not in any of the remote network CIDR blocks: %s; "+
				"use .spec.kubelet.flags field in config-source yaml to set node-ip to an IP within one of these CIDR blocks"+
				"(e.g. --node-ip=10.0.0.1) "+
				"or use --skip ip-validation",
			ipAddr, cidrs,
		)
	}
	return nil
}

func getClusterCIDRs(remoteNetworkConfig *eks.RemoteNetworkConfig) []string {
	var cidrs []string
	for _, network := range remoteNetworkConfig.RemoteNodeNetworks {
		for _, cidr := range network.CIDRs {
			if cidr != nil {
				cidrs = append(cidrs, *cidr)
			}
		}
	}
	return cidrs
}

func extractFlagValue(kubeletArgs []string, flag string) (string, error) {
	var flagValue string

	// pick last instance of the flag if it exists
	for _, s := range kubeletArgs {
		if strings.HasPrefix(s, flag) {
			flagValue = strings.TrimPrefix(s, flag)
		}
	}

	return flagValue, nil
}

func extractNodeIPFromFlags(kubeletArgs []string) (net.IP, error) {
	ipStr, err := extractFlagValue(kubeletArgs, nodeIPFlag)
	if err != nil {
		return nil, err
	}

	if ipStr != "" {
		ip := net.ParseIP(ipStr)
		if ip == nil {
			return nil, fmt.Errorf("invalid ip %s in --node-ip flag. only 1 IPv4 address is allowed", ipStr)
		} else if ip.To4() == nil {
			return nil, fmt.Errorf("invalid IPv6 address %s in --node-ip flag. only IPv4 is supported", ipStr)
		}
		return ip, nil
	}

	//--node-ip flag not set
	return nil, nil
}

func extractHostName(kubeletArgs []string) (string, error) {
	hostnameOverride, err := extractFlagValue(kubeletArgs, hostnameOverrideFlag)
	if err != nil {
		return "", fmt.Errorf("failed to extract hostname override: %w", err)
	}

	hostname, err := nodeutil.GetHostname(hostnameOverride) // returns error if it cannot resolve to a non-empty hostname
	if err != nil {
		return "", fmt.Errorf("failed to get hostname: %w", err)
	}

	return hostname, nil
}

// Validate given node IP belongs to the current host.
//
// validateNodeIP adapts the unexported 'validateNodeIP' function from kubelet.
// Source: https://github.com/kubernetes/kubernetes/blob/master/pkg/kubelet/kubelet_node_status.go#L796
func validateNodeIP(nodeIP net.IP) error {
	// Honor IP limitations set in setNodeStatus()
	if nodeIP.To4() == nil && nodeIP.To16() == nil {
		return fmt.Errorf("nodeIP must be a valid IP address")
	}
	if nodeIP.IsLoopback() {
		return fmt.Errorf("nodeIP can't be loopback address")
	}
	if nodeIP.IsMulticast() {
		return fmt.Errorf("nodeIP can't be a multicast address")
	}
	if nodeIP.IsLinkLocalUnicast() {
		return fmt.Errorf("nodeIP can't be a link-local unicast address")
	}
	if nodeIP.IsUnspecified() {
		return fmt.Errorf("nodeIP can't be an all zeros address")
	}

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return err
	}
	for _, addr := range addrs {
		var ip net.IP
		switch v := addr.(type) {
		case *net.IPNet:
			ip = v.IP
		case *net.IPAddr:
			ip = v.IP
		}
		if ip != nil && ip.Equal(nodeIP) {
			return nil
		}
	}
	return fmt.Errorf("node IP: %q not found in the host's network interfaces", nodeIP.String())
}

// getNodeIP determines the node's IP address based on kubelet configuration and system information.
func getNodeIP(kubeletArgs []string, IAMNodeName string) (net.IP, error) {
	// Follows algorithm used by kubelet to assign nodeIP
	// Implementation adapted for hybrid nodes
	// 1) Use nodeIP if set (and not "0.0.0.0"/"::")
	// 2) If the user has specified an IP to HostnameOverride, use it
	// 3) Lookup the IP from node name by DNS
	// 4) Try to get the IP from the network interface used as default gateway
	// Source: https://github.com/kubernetes/kubernetes/blob/master/pkg/kubelet/nodestatus/setters.go#L206

	nodeIP, err := extractNodeIPFromFlags(kubeletArgs)
	if err != nil {
		return nil, err
	}
	hostname, err := extractHostName(kubeletArgs)
	if err != nil {
		return nil, err
	}

	var ipAddr net.IP

	nodeIPSpecified := nodeIP != nil && nodeIP.To4() != nil && !nodeIP.IsUnspecified()

	if nodeIPSpecified {
		ipAddr = nodeIP
	} else if parsedAddr := k8snet.ParseIPSloppy(hostname); parsedAddr != nil {
		if parsedAddr.To4() == nil {
			return nil, fmt.Errorf("hostname address %s is not IPv4", parsedAddr)
		}
		ipAddr = parsedAddr
	} else {
		// If using SSM, the node name will be set at initialization to the SSM instance ID,
		// so it won't resolve to anything via DNS, hence we're only checking in the case of IAM-RA
		if IAMNodeName != "" {
			var addrs []net.IP

			addrs, _ = net.LookupIP(IAMNodeName)
			for _, addr := range addrs {
				if err = validateNodeIP(addr); addr.To4() != nil && err == nil {
					ipAddr = addr
					break
				}
			}
		}

		if ipAddr == nil {
			// current standard function for resolving bind address: https://github.com/kubernetes/kubernetes/blob/master/pkg/kubelet/kubelet_node_status.go#L768
			ipAddr, err = apimachinerynet.ResolveBindAddress(nodeIP)
		}

		if ipAddr == nil {
			// We tried everything we could, but the IP address wasn't fetchable; error out
			return nil, fmt.Errorf("couldn't get ip address of node: %w", err)
		}

	}

	return ipAddr, nil
}
