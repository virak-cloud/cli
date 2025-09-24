package cli

import "github.com/spf13/cobra"

type CobraValues struct{ cmd *cobra.Command }

func NewCobraValues(cmd *cobra.Command) CobraValues { return CobraValues{cmd: cmd} }

func (c CobraValues) GetString(name string) string { v, _ := c.cmd.Flags().GetString(name); return v }
func (c CobraValues) GetBool(name string) bool     { v, _ := c.cmd.Flags().GetBool(name); return v }
func (c CobraValues) Changed(name string) bool     { return c.cmd.Flags().Changed(name) }
