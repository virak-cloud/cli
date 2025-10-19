package cli

import "github.com/spf13/cobra"

// Validate runs a set of validation rules against the command's flags.
func Validate(cmd *cobra.Command, rules ...Rule) error {
	v := NewCobraValues(cmd)
	for _, r := range rules {
		if err := r.Validate(v); err != nil {
			return err
		}
	}
	return nil
}
