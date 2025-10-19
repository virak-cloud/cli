package publicip

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	"github.com/virak-cloud/cli/internal/cli"
	"github.com/virak-cloud/cli/pkg/http"
)

// listOptions holds the options for the `network public-ip list` command.
// These options are populated from command-line flags.
type listOptions struct {
	// ZoneID is the ID of the zone where the network is located.
	// This is optional if a default zone is set in the config.
	ZoneID    string `flag:"zoneId" usage:"Zone ID to use (optional if default.zoneId is set in config)"`
	// NetworkID is the ID of the network to list the public IPs for.
	NetworkID string `flag:"networkId" usage:"Network ID to associate the public IP with"`
}

var listOpts listOptions

// NetworkPublicIPListCmd represents the `network public-ip list` command.
// It lists all public IPs for a network.
var NetworkPublicIPListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all public IPs for a network",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := cli.Preflight(true)(cmd, args); err != nil {
			return err
		}
		return cli.Validate(cmd,
			cli.Required("networkId"),
		)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		token := cli.TokenFromContext(cmd.Context())
		zoneID := cli.ZoneIDFromContext(cmd.Context())
		if err := cli.LoadFromCobraFlags(cmd, &listOpts); err != nil {
			return err
		}

		client := http.NewClient(token)
		resp, err := client.ListNetworkPublicIps(zoneID, listOpts.NetworkID)
		if err != nil {
			slog.Error("failed to list public IPs", "error", err)
			return fmt.Errorf("error: %w", err)
		}
		if len(resp.Data) == 0 {
			fmt.Println("No public IPs found for this network.")
			return nil
		}
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "IP Address", "Is Source NAT", "Static NAT Enabled", "Created At"})
		for _, ip := range resp.Data {
			table.Append([]string{
				ip.ID,
				ip.IpAddress,
				fmt.Sprintf("%v", ip.IsSourceNat),
				fmt.Sprintf("%v", ip.StaticNatEnable),
				fmt.Sprintf("%v", ip.CreatedAt),
			})
		}
		table.Render()
		return nil
	},
}

// init registers the `network public-ip list` command with the parent `network public-ip` command
// and binds the flags for the `listOptions` struct.
func init() {
	NetworkPublicIPCmd.AddCommand(NetworkPublicIPListCmd)
	_ = cli.BindFlagsFromStruct(NetworkPublicIPListCmd, &listOpts)
}
