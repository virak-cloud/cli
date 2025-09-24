package cli

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type ctxKey string

const (
	ctxTokenKey  ctxKey = "token"
	ctxZoneIDKey ctxKey = "zoneId"
)

func TokenFromContext(ctx context.Context) string {
	v := ctx.Value(ctxTokenKey)
	s, _ := v.(string)
	return s
}

func ZoneIDFromContext(ctx context.Context) string {
	v := ctx.Value(ctxZoneIDKey)
	s, _ := v.(string)
	return s
}

// Preflight returns a PersistentPreRunE-compatible function that ensures login and, if zoneRequired, resolves zoneId
// via --default-zone or --zoneId flags with consistent error messages.
func Preflight(zoneRequired bool) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, _ []string) error {
		token := viper.GetString("auth.token")
		if token == "" {
			slog.Error("not logged in")
			return fmt.Errorf("you must be logged in to use this command. Please run 'virak-cli login' first")
		}

		defaultZoneFlag, _ := cmd.Flags().GetBool("default-zone")
		zoneId, _ := cmd.Flags().GetString("zoneId")

		if zoneRequired {
			if defaultZoneFlag {
				zoneId = viper.GetString("default.zoneId")
				if zoneId == "" {
					slog.Error("--default-zone flag used but no default.zoneId found in config")
					return fmt.Errorf("--default-zone flag used but no default.zoneId found in config")
				}
			} else if zoneId == "" {
				slog.Error("--zoneId flag required when --default-zone not set")
				return fmt.Errorf("--zoneId flag required when --default-zone not set")
			}
		}

		ctx := context.WithValue(cmd.Context(), ctxTokenKey, token)
		if zoneId != "" {
			ctx = context.WithValue(ctx, ctxZoneIDKey, zoneId)
		}
		cmd.SetContext(ctx)
		return nil
	}
}
