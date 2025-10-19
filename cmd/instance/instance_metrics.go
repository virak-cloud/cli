package instance

import (
	"fmt"
	"github.com/virak-cloud/cli/internal/cli"
	"github.com/virak-cloud/cli/pkg/http"
	"log/slog"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// metricsOptions holds the options for the `instance metrics` command.
// These options are populated from command-line flags.
type metricsOptions struct {
	// ZoneID is the ID of the zone where the instance is located.
	// This is optional if a default zone is set in the config.
	ZoneID     string   `flag:"zoneId" usage:"Zone ID to use (optional if default.zoneId is set in config)"`
	// InstanceID is the ID of the instance to get metrics for.
	InstanceID string   `flag:"instanceId" usage:"Instance ID"`
	// Metrics is a list of metrics to query.
	Metrics    []string `flag:"metrics" usage:"Metrics to query"`
	// Time is the time window for the metrics.
	Time       int      `flag:"time" default:"1" usage:"Time window"`
	// Aggregator is the aggregator to use for the metrics.
	Aggregator string   `flag:"aggregator" default:"mean" usage:"Aggregator"`
}

var metricsOpt metricsOptions

// instanceMetricsCmd represents the `instance metrics` command.
// It retrieves performance metrics for a specified instance.
var instanceMetricsCmd = &cobra.Command{
	Use:   "metrics",
	Short: "Get instance performance metrics",
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

		if err := cli.LoadFromCobraFlags(cmd, &metricsOpt); err != nil {
			return err
		}

		metrics := metricsOpt.Metrics
		if len(metrics) == 0 {
			metrics = []string{"memoryusedkbs", "cpuused"}
		}

		httpClient := http.NewClient(token)
		resp, err := httpClient.GetInstanceMetrics(zoneID, metricsOpt.InstanceID, metrics, metricsOpt.Time, metricsOpt.Aggregator)
		if err != nil {
			slog.Error("failed to get instance metrics", "error", err, "zoneID", zoneID, "instanceID", metricsOpt.InstanceID)
			return fmt.Errorf("failed to get instance metrics: %w", err)
		}
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Metric", "Time", "Value"})
		for _, col := range resp.Data {
			for _, val := range col.Values {
				table.Append([]string{col.Column, val.Time, fmt.Sprintf("%v", val.Value)})
			}
		}
		table.Render()
		return nil
	},
}

// init registers the `instance metrics` command with the parent `instance` command
// and binds the flags for the `metricsOptions` struct.
func init() {
	InstanceCmd.AddCommand(instanceMetricsCmd)
	_ = cli.BindFlagsFromStruct(instanceMetricsCmd, &metricsOpt)
}
