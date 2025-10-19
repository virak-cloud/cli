package cli

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// FlagSpec defines the specification for a command-line flag.
type FlagSpec struct {
	// Name is the name of the flag.
	Name  string
	// Usage is the help text for the flag.
	Usage string
	// Def is the default value for the flag.
	Def   string
}

// BindFlagsFromStruct declares flags on a cobra command based on struct tags.
// It supports the following tags:
// - `flag`: the name of the flag.
// - `usage`: the help text for the flag.
// - `default`: the default value for the flag.
func BindFlagsFromStruct(cmd *cobra.Command, opts any) error {
	t := reflect.TypeOf(opts)
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return fmt.Errorf("opts must be a struct or *struct")
	}

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		name := f.Tag.Get("flag")
		if name == "" {
			continue
		}
		usage := f.Tag.Get("usage")
		def := f.Tag.Get("default")
		switch f.Type.Kind() {
		case reflect.String:
			cmd.Flags().String(name, def, usage)
		case reflect.Bool:
			// bool has no default string, interpret def == "true"
			cmd.Flags().Bool(name, def == "true", usage)
		case reflect.Int:
			defInt, _ := strconv.Atoi(def)
			if def == "" {
				defInt = 0
			}
			cmd.Flags().Int(name, defInt, usage)
		case reflect.Slice:
			if f.Type.Elem().Kind() == reflect.String {
				defSlice := []string{}
				if def != "" {
					defSlice = strings.Split(def, ",")
				}
				cmd.Flags().StringSlice(name, defSlice, usage)
			}
		default:
			panic("unhandled default case")
		}
	}
	return nil
}

// LoadFromViper decodes viper keyspace into the provided struct.
// It assumes that the flags have already been bound to viper.
func LoadFromViper(opts any) error {
	dec, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{TagName: "mapstructure", Result: opts})
	if err != nil {
		return err
	}
	return dec.Decode(viper.AllSettings())
}

// LoadFromCobraFlags reads the values of flags defined on a cobra command
// and loads them into the fields of the provided struct, based on the `flag` tag.
// It supports string, bool, int, and []string field kinds.
func LoadFromCobraFlags(cmd *cobra.Command, opts any) error {
	v := reflect.ValueOf(opts)
	if v.Kind() != reflect.Pointer || v.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("opts must be a pointer to struct")
	}
	v = v.Elem()
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		flagName := field.Tag.Get("flag")
		if flagName == "" {
			continue
		}
		fv := v.Field(i)
		if !fv.CanSet() {
			continue
		}

		var (
			err error
		)
		switch field.Type.Kind() {
		case reflect.String:
			var val string
			val, err = cmd.Flags().GetString(flagName)
			if err == nil {
				fv.SetString(val)
			}
		case reflect.Bool:
			var val bool
			val, err = cmd.Flags().GetBool(flagName)
			if err == nil {
				fv.SetBool(val)
			}
		case reflect.Int:
			var val int
			val, err = cmd.Flags().GetInt(flagName)
			if err == nil {
				fv.SetInt(int64(val))
			}
		case reflect.Slice:
			if field.Type.Elem().Kind() == reflect.String {
				var val []string
				val, err = cmd.Flags().GetStringSlice(flagName)
				if err == nil {
					fv.Set(reflect.ValueOf(val))
				}
			}
		default:
			return fmt.Errorf("unsupported field kind %s for flag %q", field.Type.Kind(), flagName)
		}
		if err != nil {
			return fmt.Errorf("reading flag %q: %w", flagName, err)
		}
	}
	return nil
}
