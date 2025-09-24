package internal

import (
	"slices"

	"github.com/oklog/ulid/v2"
)

// IsValidRecordType checks if the given record type is valid.
func IsValidRecordType(recordType string) bool {
	// A, AAAA, CNAME, MX, TXT, NS, SOA, SRV, etc.
	// Based on the provided JSON, we can see "A" is a valid type.
	// We can extend this list with other common types.
	validTypes := []string{"A", "AAAA", "CNAME", "MX", "TXT", "NS", "SOA", "SRV"}
	return slices.Contains(validTypes, recordType)
}

// IsValidULID checks if the given string is a valid ULID.
func IsValidULID(s string) bool {
	_, err := ulid.Parse(s)
	return err == nil
}
