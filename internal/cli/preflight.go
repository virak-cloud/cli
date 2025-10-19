package cli

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// ctxKey is a custom type for context keys to avoid collisions.
type ctxKey string

const (
	// ctxTokenKey is the context key for the authentication token.
	ctxTokenKey  ctxKey = "token"
	// ctxZoneIDKey is the context key for the zone ID.
	ctxZoneIDKey ctxKey = "zoneId"
)

// TokenFromContext retrieves the authentication token from the context.
func TokenFromContext(ctx context.Context) string {
	v := ctx.Value(ctxTokenKey)
	s, _ := v.(string)
	return s
}

// ZoneIDFromContext retrieves the zone ID from the context.
func ZoneIDFromContext(ctx context.Context) string {
	v := ctx.Value(ctxZoneIDKey)
	s, _ := v.(string)
	return s
}

// Preflight returns a cobra.PositionalArgs function that performs pre-run checks.
// It ensures that the user is logged in, and if zoneRequired is true, it resolves the zone ID.
// The resolved token and zone ID are then added to the command's context.
func Preflight(zoneRequired bool) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, _ []string) error {
		token := viper.GetString("auth.token")
		if token == "" {
			slog.Error("not logged in")
			return fmt.Errorf("you must be logged in to use this command. Please run 'virak-cli login' first")
		}

		zoneId, _ := cmd.Flags().GetString("zoneId")

		if zoneRequired {
			// Always try to use the default zoneId from config if set
			if defaultZoneId := viper.GetString("default.zoneId"); defaultZoneId != "" {
				zoneId = defaultZoneId
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
