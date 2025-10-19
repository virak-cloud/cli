package cli

import (
	"fmt"

	"github.com/spf13/viper"
)

// SetDefaultZone sets the default zone ID and name in the viper configuration
// and writes the configuration to disk.
func SetDefaultZone(zoneID, zoneName string) error {
	viper.Set("default.zoneId", zoneID)
	viper.Set("default.zoneName", zoneName)
	if err := viper.SafeWriteConfig(); err != nil {
		// If config file already exists, fallback to WriteConfig
		if err := viper.WriteConfig(); err != nil {
			return fmt.Errorf("failed to write config: %w", err)
		}
	}
	return nil
}
