package instance

import (
	"fmt"

	"github.com/virak-cloud/cli/internal/cli"
	"github.com/virak-cloud/cli/pkg/http"

	"github.com/spf13/cobra"
)

type consoleOptions struct {
	ZoneID     string `flag:"zoneId" usage:"Zone ID to use (optional if default.zoneId is set in config)"`
	InstanceID string `flag:"instanceId" usage:"Instance ID"`
}

var consoleOpt consoleOptions

var instanceConsoleCmd = &cobra.Command{
	Use:   "console",
	Short: "Get console URL for an instance",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := cli.Preflight(true)(cmd, args); err != nil {
			return err
		}

		return cli.Validate(cmd,
			cli.Required("instanceId"),
		)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		token := cli.TokenFromContext(cmd.Context())
		zoneID := cli.ZoneIDFromContext(cmd.Context())

		if err := cli.LoadFromCobraFlags(cmd, &consoleOpt); err != nil {
			return err
		}

		httpClient := http.NewClient(token)

		resp, err := httpClient.GetInstanceConsole(zoneID, consoleOpt.InstanceID)
		if err != nil {
			return fmt.Errorf("could not get instance console: %w", err)
		}

		fmt.Printf("Console URL: %s\n", resp.Data.URL)
		return nil
	},
}

func init() {
	InstanceCmd.AddCommand(instanceConsoleCmd)
	_ = cli.BindFlagsFromStruct(instanceConsoleCmd, &consoleOpt)
}