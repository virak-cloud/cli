package ipv4

import (
	"fmt"
	"log/slog"

	"github.com/virak-cloud/cli/pkg/http"

	"github.com/virak-cloud/cli/internal/cli"

	"github.com/spf13/cobra"
)

// firewallIPv4CreateOptions holds the options for the `network firewall ipv4 create` command.
// These options are populated from command-line flags.
type firewallIPv4CreateOptions struct {
	// ZoneID is the ID of the zone where the network is located.
	// This is optional if a default zone is set in the config.
	ZoneID        string `flag:"zoneId" usage:"Zone ID to use (optional if default.zoneId is set in config)"`
	// NetworkID is the ID of the network to add the firewall rule to.
	NetworkID     string `flag:"networkId" desc:"Network ID (required)"`
	// TrafficType is the type of traffic to match (Ingress or Egress).
	TrafficType   string `flag:"trafficType" desc:"Traffic type (Ingress/Egress) [required]"`
	// ProtocolType is the protocol to match (TCP, UDP, or ICMP).
	ProtocolType  string `flag:"protocolType" desc:"Protocol type (TCP/UDP/ICMP) [required]"`
	// PublicIpId is the ID of the public IP address to associate with the rule.
	PublicIpId    string `flag:"publicIpId" desc:"Public IP ID (optional)"`
	// IPSource is the source IP address or CIDR block to match.
	IPSource      string `flag:"ipSource" desc:"Source IP [required]"`
	// IPDestination is the destination IP address or CIDR block to match.
	IPDestination string `flag:"ipDestination" desc:"Destination IP [required]"`
	// PortStart is the starting port number for TCP or UDP traffic.
	PortStart     int    `flag:"portStart" desc:"Start port (optional)"`
	// PortEnd is the ending port number for TCP or UDP traffic.
	PortEnd       int    `flag:"portEnd" desc:"End port (optional)"`
	// ICMPCode is the ICMP code to match.
	ICMPCode      int    `flag:"icmpCode" desc:"ICMP code (optional)"`
	// ICMPType is the ICMP type to match.
	ICMPType      int    `flag:"icmpType" desc:"ICMP type (optional)"`
}

var firewallIPv4CreateOpts firewallIPv4CreateOptions

// NetworkFirewallIPv4CreateCmd represents the `network firewall ipv4 create` command.
// It creates a new IPv4 firewall rule for a network.
var NetworkFirewallIPv4CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an IPv4 firewall rule for a network",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := cli.Preflight(true)(cmd, args); err != nil {
			return err
		}
		return cli.Validate(cmd,
			cli.Required("networkId"),
			cli.IsUlid("networkId"),
			cli.Required("trafficType"),
			cli.OneOf("trafficType", "Ingress", "Egress"),
			cli.Required("protocolType"),
			cli.OneOf("protocolType", "TCP", "UDP", "ICMP"),
			cli.Required("ipSource"),
			cli.Required("ipDestination"),
			cli.RequiredIf("publicIpId", func(v cli.Values) bool { return v.GetString("trafficType") == "Ingress" }),
			cli.RequiredIf("icmpCode", func(v cli.Values) bool { return v.GetString("protocolType") == "ICMP" }),
			cli.RequiredIf("icmpType", func(v cli.Values) bool { return v.GetString("protocolType") == "ICMP" }),
		)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		token := cli.TokenFromContext(cmd.Context())
		zoneId := cli.ZoneIDFromContext(cmd.Context())

		if err := cli.LoadFromCobraFlags(cmd, &firewallIPv4CreateOpts); err != nil {
			return err
		}

		body := map[string]interface{}{
			"traffic_type":   firewallIPv4CreateOpts.TrafficType,
			"protocol_type":  firewallIPv4CreateOpts.ProtocolType,
			"public_ip_id":   firewallIPv4CreateOpts.PublicIpId,
			"ip_source":      firewallIPv4CreateOpts.IPSource,
			"ip_destination": firewallIPv4CreateOpts.IPDestination,
		}

		if firewallIPv4CreateOpts.ProtocolType == "ICMP" {
			body["icmp_code"] = firewallIPv4CreateOpts.ICMPCode
			body["icmp_type"] = firewallIPv4CreateOpts.ICMPType
		} else {
			body["port_start"] = firewallIPv4CreateOpts.PortStart
			body["port_end"] = firewallIPv4CreateOpts.PortEnd
		}

		httpClient := http.NewClient(token)
		resp, err := httpClient.CreateIPv4FirewallRule(zoneId, firewallIPv4CreateOpts.NetworkID, body)
		if err != nil {
			slog.Error("failed to create IPv4 firewall rule", "error", err)
			return fmt.Errorf("failed to create IPv4 firewall rule: %w", err)
		}

		if resp.Data.Success {
			fmt.Println("IPv4 firewall rule created successfully.")
		} else {
			fmt.Println("Failed to create IPv4 firewall rule.")
		}
		return nil
	},
}

// init registers the `network firewall ipv4 create` command with the parent `network firewall ipv4` command
// and binds the flags for the `firewallIPv4CreateOptions` struct.
func init() {
	NetworkFirewallIPv4Cmd.AddCommand(NetworkFirewallIPv4CreateCmd)
	_ = cli.BindFlagsFromStruct(NetworkFirewallIPv4CreateCmd, &firewallIPv4CreateOpts)
}
