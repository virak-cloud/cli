package zone

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/virak-cloud/cli/internal/cli"
	"github.com/virak-cloud/cli/pkg/http"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "A list of all available zones",
	Long:  `Get a list of all available zones from Virak API.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := cli.Preflight(false)(cmd, args); err != nil {
			return err
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		token := cli.TokenFromContext(cmd.Context())
		httpClient := http.NewClient(token)
		zones, err := httpClient.GetZoneList()
		if err != nil {
			slog.Error("failed to get zone list", "error", err)
			return fmt.Errorf("failed to get zone list: %w", err)
		}
		slog.Info("successfully retrieved zone list", "count", len(zones.Data))
		for i, zone := range zones.Data {
			fmt.Printf("[%d] Zone Name: %s, ID: %s \n", i+1, zone.Name, zone.ID)
		}

		// Ask user if they want to set a default zone
		fmt.Print("Do you want to set a default zone? (y/N): ")
		reader := bufio.NewReader(os.Stdin)
		answer, _ := reader.ReadString('\n')
		answer = strings.TrimSpace(answer)
		if answer == "y" || answer == "Y" {
			fmt.Print("Enter the number of the default zone: ")
			zoneInput, _ := reader.ReadString('\n')
			zoneInput = strings.TrimSpace(zoneInput)
			var zoneNumber int
			_, err := fmt.Sscanf(zoneInput, "%d", &zoneNumber)
			if err != nil {
				slog.Error("failed to read zone number", "error", err)
				fmt.Println("Invalid input.")
				return nil
			}
			if zoneNumber < 1 || zoneNumber > len(zones.Data) {
				slog.Error("invalid zone number", "zoneNumber", zoneNumber)
				fmt.Println("Invalid zone number.")
				return nil
			}
			defaultZoneID := zones.Data[zoneNumber-1].ID
			defaultZoneName := zones.Data[zoneNumber-1].Name
			err = cli.SetDefaultZone(defaultZoneID, defaultZoneName)
			if err != nil {
				slog.Error("failed to save default zone to config", "error", err)
				fmt.Println("Failed to save default zone to config:", err)
			} else {
				slog.Info("default zone set", "zoneName", defaultZoneName, "zoneId", defaultZoneID)
				fmt.Println("Default zone set to:", defaultZoneName)
			}
		}
		return nil
	},
}

func init() {
	ZoneCmd.AddCommand(listCmd)
}
