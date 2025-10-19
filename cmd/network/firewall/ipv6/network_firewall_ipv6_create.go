package ipv6

import (
	"fmt"
	"log/slog"

	"github.com/virak-cloud/cli/internal/cli"
	"github.com/virak-cloud/cli/pkg/http"

	"github.com/spf13/cobra"
)

// NetworkFirewallIPv6CreateOptions holds the options for the `network firewall ipv6 create` command.
// These options are populated from command-line flags.
type NetworkFirewallIPv6CreateOptions struct {
	// ZoneID is the ID of the zone where the network is located.
	// This is optional if a default zone is set in the config.
	ZoneID        string `flag:"zoneId" usage:"Zone ID to use (optional if default.zoneId is set in config)"`
	// NetworkID is the ID of the network to add the firewall rule to.
	NetworkID     string `flag:"networkId"`
	// TrafficType is the type of traffic to match (Ingress or Egress).
	TrafficType   string `flag:"trafficType"`
	// ProtocolType is the protocol to match (TCP, UDP, or ICMP).
	ProtocolType  string `flag:"protocolType"`
	// IPSource is the source IP address or CIDR block to match.
	IPSource      string `flag:"ipSource"`
	// IPDestination is the destination IP address or CIDR block to match.
	IPDestination string `flag:"ipDestination"`
	// PortStart is the starting port number for TCP or UDP traffic.
	PortStart     int    `flag:"portStart"`
	// PortEnd is the ending port number for TCP or UDP traffic.
	PortEnd       int    `flag:"portEnd"`
	// ICMPCode is the ICMP code to match.
	ICMPCode      int    `flag:"icmpCode"`
	// ICMPType is the ICMP type to match.
	ICMPType      int    `flag:"icmpType"`
}

var firewallIPv6CreateOptions NetworkFirewallIPv6CreateOptions

// NetworkFirewallIPv6CreateCmd represents the `network firewall ipv6 create` command.
// It creates a new IPv6 firewall rule for a network.
var NetworkFirewallIPv6CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an IPv6 firewall rule for a network",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// Centralized login + zone handling
		if err := cli.Preflight(true)(cmd, args); err != nil {
			return err
		}
		// Declarative validation rules
		return cli.Validate(cmd,
			cli.Required("networkId"),
			cli.Required("trafficType"),
			cli.Required("protocolType"),
			cli.Required("ipSource"),
			cli.Required("ipDestination"),
			cli.OneOf("trafficType", "Ingress", "Egress"),
			cli.OneOf("protocolType", "TCP", "UDP", "ICMP"),
		)
	},
	RunE: func(cmd *cobra.Command, args []string) error {

		token := cli.TokenFromContext(cmd.Context())
		zoneId := cli.ZoneIDFromContext(cmd.Context())
		if err := cli.LoadFromCobraFlags(cmd, &firewallIPv6CreateOptions); err != nil {
			return err
		}

		httpClient := http.NewClient(token)
		body := map[string]interface{}{
			"traffic_type":   firewallIPv6CreateOptions.TrafficType,
			"protocol_type":  firewallIPv6CreateOptions.ProtocolType,
			"ip_source":      firewallIPv6CreateOptions.IPSource,
			"ip_destination": firewallIPv6CreateOptions.IPDestination,
			"port_start":     firewallIPv6CreateOptions.PortStart,
			"port_end":       firewallIPv6CreateOptions.PortEnd,
			"icmp_code":      firewallIPv6CreateOptions.ICMPCode,
			"icmp_type":      firewallIPv6CreateOptions.ICMPType,
		}
		resp, err := httpClient.CreateIPv6FirewallRule(zoneId, firewallIPv6CreateOptions.NetworkID, body)
		if err != nil {
			slog.Error("failed to create IPv6 firewall rule", "error", err)
			return fmt.Errorf("failed to create IPv6 firewall rule: %w", err)
		}
		if resp.Data.Success {
			fmt.Println("IPv6 firewall rule created successfully.")
		} else {
			fmt.Println("Failed to create IPv6 firewall rule.")
		}
		return nil
	},
}

// init registers the `network firewall ipv6 create` command with the parent `network firewall ipv6` command
// and binds the flags for the `NetworkFirewallIPv6CreateOptions` struct.
func init() {
	NetworkFirewallIPv6Cmd.AddCommand(NetworkFirewallIPv6CreateCmd)
	_ = cli.BindFlagsFromStruct(NetworkFirewallIPv6CreateCmd, &firewallIPv6CreateOptions)
}
