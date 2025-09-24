package cli

import "github.com/spf13/cobra"

func Validate(cmd *cobra.Command, rules ...Rule) error {
	v := NewCobraValues(cmd)
	for _, r := range rules {
		if err := r.Validate(v); err != nil {
			return err
		}
	}
	return nil
}
