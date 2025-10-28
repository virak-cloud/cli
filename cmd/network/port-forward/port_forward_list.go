package portforward

import (
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"

	"github.com/virak-cloud/cli/internal/cli"
	"github.com/virak-cloud/cli/internal/presenter"
	"github.com/virak-cloud/cli/pkg/http"
)

type portForwardListOptions struct {
	ZoneID    string `flag:"zoneId" desc:"Zone ID (optional if default.zoneId is set in config)"`
	NetworkID string `flag:"networkId" desc:"Network ID (required)"`
}

var portForwardListOpts portForwardListOptions

var NetworkPortForwardListCmd = &cobra.Command{
	Use:   "list",
	Short: "List port forwarding rules for a network",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := cli.Preflight(true)(cmd, args); err != nil {
			return err
		}
		return cli.Validate(cmd,
			cli.Required("networkId"),
			cli.IsUlid("networkId"),
		)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		token := cli.TokenFromContext(cmd.Context())
		zoneId := cli.ZoneIDFromContext(cmd.Context())

		if err := cli.LoadFromCobraFlags(cmd, &portForwardListOpts); err != nil {
			return err
		}

		httpClient := http.NewClient(token)
		resp, err := httpClient.ListPortForwards(zoneId, portForwardListOpts.NetworkID)
		if err != nil {
			slog.Error("failed to list port forwarding rules", "error", err)
			return fmt.Errorf("failed to list port forwarding rules: %w", err)
		}

		if len(resp.Data) == 0 {
			fmt.Println("No port forwarding rules found.")
			return nil
		}

		presenter.RenderPortForwardList(resp.Data)
		return nil
	},
}

func init() {
	NetworkPortForwardCmd.AddCommand(NetworkPortForwardListCmd)
	_ = cli.BindFlagsFromStruct(NetworkPortForwardListCmd, &portForwardListOpts)
}