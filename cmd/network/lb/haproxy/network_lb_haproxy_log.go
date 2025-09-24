package haproxy

import (
	"fmt"
	"log/slog"
	"virak-cli/internal/cli"
	"virak-cli/pkg/http"

	"github.com/spf13/cobra"
)

type lbHaproxyLogOptions struct {
	DefaultZone bool   `flag:"default-zone" usage:"Use default.zoneId from config"`
	ZoneID      string `flag:"zoneId" usage:"Zone ID to use"`
	NetworkID   string `flag:"networkId" usage:"Network ID for the load balancer"`
}

var lbHaproxyLogOpts lbHaproxyLogOptions

var NetworkLbHaproxyLogCmd = &cobra.Command{
	Use:   "log",
	Short: "Get HAProxy logs for the network",
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
		if err := cli.LoadFromCobraFlags(cmd, &lbHaproxyLogOpts); err != nil {
			return err
		}
		httpClient := http.NewClient(token)
		resp, err := httpClient.GetHaproxyLog(zoneId, lbHaproxyLogOpts.NetworkID)
		if err != nil {
			slog.Error("failed to get HAProxy log report", "error", err)
			return fmt.Errorf("error: %w", err)
		}
		if len(resp.Data) == 0 {
			fmt.Println("No HAProxy logs found.")
			return nil
		}
		fmt.Println("HAProxy Logs:")
		for i, log := range resp.Data {
			fmt.Printf("%d: %v\n", i+1, log)
		}
		return nil
	},
}

func init() {
	NetworkLbHaproxyCmd.AddCommand(NetworkLbHaproxyLogCmd)
	_ = cli.BindFlagsFromStruct(NetworkLbHaproxyLogCmd, &lbHaproxyLogOpts)
}
