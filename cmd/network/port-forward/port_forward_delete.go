package portforward

import (
	"fmt"
	"log/slog"

	"github.com/virak-cloud/cli/internal/cli"
	"github.com/virak-cloud/cli/pkg/http"

	"github.com/spf13/cobra"
)

type portForwardDeleteOptions struct {
	ZoneID string `flag:"zoneId" desc:"Zone ID (optional if default.zoneId is set in config)"`
	ID     string `flag:"id" desc:"Port forwarding rule ID (required)"`
}

var portForwardDeleteOpts portForwardDeleteOptions

var NetworkPortForwardDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a port forwarding rule",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := cli.Preflight(true)(cmd, args); err != nil {
			return err
		}
		return cli.Validate(cmd,
			cli.Required("id"),
			cli.IsUlid("id"),
		)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		token := cli.TokenFromContext(cmd.Context())
		zoneId := cli.ZoneIDFromContext(cmd.Context())

		if err := cli.LoadFromCobraFlags(cmd, &portForwardDeleteOpts); err != nil {
			return err
		}

		httpClient := http.NewClient(token)
		resp, err := httpClient.DeletePortForward(zoneId, portForwardDeleteOpts.ID)
		if err != nil {
			slog.Error("failed to delete port forwarding rule", "error", err)
			return fmt.Errorf("failed to delete port forwarding rule: %w", err)
		}

		if resp.Data.Success {
			fmt.Println("Port forwarding rule deleted successfully.")
		} else {
			fmt.Println("Failed to delete port forwarding rule.")
		}
		return nil
	},
}

func init() {
	NetworkPortForwardCmd.AddCommand(NetworkPortForwardDeleteCmd)
	_ = cli.BindFlagsFromStruct(NetworkPortForwardDeleteCmd, &portForwardDeleteOpts)
}