package cli

import (
	"fmt"
	"strings"

	"github.com/oklog/ulid/v2"
)

type Values interface {
	GetString(name string) string
	GetBool(name string) bool
	Changed(name string) bool
}

type Rule interface{ Validate(v Values) error }

type RuleFunc func(v Values) error

func (f RuleFunc) Validate(v Values) error { return f(v) }

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

func RequiredIf(name string, predicate func(v Values) bool) Rule {
	return RuleFunc(func(v Values) error {
		if predicate(v) && !v.Changed(name) {
			return fmt.Errorf("--%s is required due to other flags", name)
		}
		return nil
	})
}

func MutuallyExclusive(a, b string) Rule {
	return RuleFunc(func(v Values) error {
		if v.GetString(a) != "" && v.GetString(b) != "" {
			return fmt.Errorf("--%s and --%s are mutually exclusive", a, b)
		}
		return nil
	})
}

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

func isValidUlid(s string) bool {
	_, err := ulid.Parse(s)
	return err == nil
}

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
