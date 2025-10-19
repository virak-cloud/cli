package cli

import (
	"fmt"
	"strings"

	"github.com/oklog/ulid/v2"
)

// Values is an interface for accessing flag values.
// It is used by the validation rules to get the values of flags.
type Values interface {
	// GetString returns the string value of a flag.
	GetString(name string) string
	// GetBool returns the bool value of a flag.
	GetBool(name string) bool
	// Changed returns true if the flag was set by the user.
	Changed(name string) bool
}

// Rule is an interface for a validation rule.
type Rule interface{ Validate(v Values) error }

// RuleFunc is an adapter to allow the use of ordinary functions as validation rules.
type RuleFunc func(v Values) error

// Validate calls the underlying function.
func (f RuleFunc) Validate(v Values) error { return f(v) }

// Required returns a Rule that checks if a flag is present.
func Required(name string) Rule {
	return RuleFunc(func(v Values) error {
		// If the flag was explicitly provided by the user, consider it present.
		if v.Changed(name) {
			return nil
		}
		val := strings.TrimSpace(v.GetString(name))
		if val == "" || val == "-1" {
			return fmt.Errorf("--%s is required", name)
		}
		return nil
	})
}

// OneOf returns a Rule that checks if a flag's value is one of the allowed values.
func OneOf(name string, allowed ...string) Rule {
	return RuleFunc(func(v Values) error {
		val := v.GetString(name)
		if val == "" {
			return nil
		} // leave to Required if needed
		for _, a := range allowed {
			if val == a {
				return nil
			}
		}
		return fmt.Errorf("--%s must be one of: %s", name, strings.Join(allowed, ", "))
	})
}

// RequiredIf returns a Rule that checks if a flag is present if a predicate is true.
func RequiredIf(name string, predicate func(v Values) bool) Rule {
	return RuleFunc(func(v Values) error {
		if predicate(v) && !v.Changed(name) {
			return fmt.Errorf("--%s is required due to other flags", name)
		}
		return nil
	})
}

// MutuallyExclusive returns a Rule that checks if two flags are mutually exclusive.
func MutuallyExclusive(a, b string) Rule {
	return RuleFunc(func(v Values) error {
		if v.GetString(a) != "" && v.GetString(b) != "" {
			return fmt.Errorf("--%s and --%s are mutually exclusive", a, b)
		}
		return nil
	})
}

// ExactlyOne returns a Rule that checks if exactly one of a list of flags is present.
func ExactlyOne(names ...string) Rule {
	return RuleFunc(func(v Values) error {
		count := 0
		for _, n := range names {
			if v.GetString(n) != "" || v.GetBool(n) {
				count++
			}
		}
		if count != 1 {
			return fmt.Errorf("exactly one of (%s) is required", strings.Join(names, ", "))
		}
		return nil
	})
}

// IsUlid returns a Rule that checks if a flag's value is a valid ULID.
func IsUlid(name string) Rule {
	return RuleFunc(func(v Values) error {
		val := v.GetString(name)
		if val == "" {
			return nil
		} // leave to Required if needed
		if !isValidUlid(val) {
			return fmt.Errorf("--%s must be a valid ULID", name)
		}
		return nil
	})
}

// isValidUlid checks if a string is a valid ULID.
func isValidUlid(s string) bool {
	_, err := ulid.Parse(s)
	return err == nil
}

// MinLength returns a Rule that checks if a flag's value has a minimum length.
func MinLength(name string, min int) Rule {
	return RuleFunc(func(v Values) error {
		val := v.GetString(name)
		if val == "" {
			return nil // leave to Required if needed
		}
		if len(val) < min {
			return fmt.Errorf("--%s must be at least %d characters", name, min)
		}
		return nil
	})
}

// MaxLength returns a Rule that checks if a flag's value has a maximum length.
func MaxLength(name string, max int) Rule {
	return RuleFunc(func(v Values) error {
		val := v.GetString(name)
		if val == "" {
			return nil // leave to Required if needed
		}
		if len(val) > max {
			return fmt.Errorf("--%s must be at most %d characters", name, max)
		}
		return nil
	})
}
