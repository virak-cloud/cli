package cli

import "github.com/spf13/cobra"

// CobraValues is an adapter for a cobra.Command that implements the Values interface.
type CobraValues struct{ cmd *cobra.Command }

// NewCobraValues returns a new CobraValues adapter for the given command.
func NewCobraValues(cmd *cobra.Command) CobraValues { return CobraValues{cmd: cmd} }

// GetString returns the string value of a flag.
func (c CobraValues) GetString(name string) string { v, _ := c.cmd.Flags().GetString(name); return v }

// GetBool returns the bool value of a flag.
func (c CobraValues) GetBool(name string) bool     { v, _ := c.cmd.Flags().GetBool(name); return v }

// Changed returns true if the flag was set by the user.
func (c CobraValues) Changed(name string) bool     { return c.cmd.Flags().Changed(name) }
