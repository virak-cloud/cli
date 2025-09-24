package haproxy

import (
	"fmt"
	"log/slog"
	"virak-cli/internal/cli"
	"virak-cli/pkg/http"

	"github.com/spf13/cobra"
)

type lbHaproxyLiveOptions struct {
	DefaultZone bool   `flag:"default-zone" usage:"Use default.zoneId from config"`
	ZoneID      string `flag:"zoneId" usage:"Zone ID to use"`
	NetworkID   string `flag:"networkId" usage:"Network ID for the load balancer"`
}

var lbHaproxyLiveOpts lbHaproxyLiveOptions

var NetworkLbHaproxyLiveCmd = &cobra.Command{
	Use:   "live",
	Short: "Get live HAProxy statistics for the network",
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
		zoneId := cli.ZoneIDFromContext(cmd.Context())
		if err := cli.LoadFromCobraFlags(cmd, &lbHaproxyLiveOpts); err != nil {
			return err
		}
		httpClient := http.NewClient(token)
		resp, err := httpClient.GetHaproxyLive(zoneId, lbHaproxyLiveOpts.NetworkID)
		if err != nil {
			slog.Error("failed to get HAProxy live report", "error", err)
			return fmt.Errorf("error: %w", err)
		}
		fmt.Printf("Updated At: %d\n", resp.Data.UpdatedAt)
		if len(resp.Data.Rules) == 0 {
			fmt.Println("No HAProxy rules found.")
			return nil
		}
		fmt.Println("ID\tName\tAlgorithm\tPublicPort\tPrivatePort\tStatus")
		for _, rule := range resp.Data.Rules {
			fmt.Printf("%s\t%s\t%s\t%d\t%d\t%s\n", rule.ID, rule.Name, rule.Algorithm, rule.PublicPort, rule.PrivatePort, rule.Status)
		}
		return nil
	},
}

func init() {
	NetworkLbHaproxyCmd.AddCommand(NetworkLbHaproxyLiveCmd)
	_ = cli.BindFlagsFromStruct(NetworkLbHaproxyLiveCmd, &lbHaproxyLiveOpts)
}
