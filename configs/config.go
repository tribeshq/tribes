// The config package manages the node configuration, which comes from environment variables.
// The sub-package generate specifies these environment variables.
package configs

import (
	"encoding/hex"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

// Redacted is a wrapper that redacts a given field from the logs.
type Redacted[T any] struct {
	Value T
}

func (r Redacted[T]) String() string {
	return "[REDACTED]"
}

type (
	Address  = common.Address
	Duration = time.Duration
)

// ------------------------------------------------------------------------------------------------
// Parsing functions
// ------------------------------------------------------------------------------------------------

func ToUint64FromString(s string) (uint64, error) {
	value, err := strconv.ParseUint(s, 10, 64)
	return value, err
}

func ToUint64FromDecimalOrHexString(s string) (uint64, error) {
	if len(s) >= 2 && (strings.HasPrefix(s, "0x") || strings.HasPrefix(s, "0X")) {
		return strconv.ParseUint(s[2:], 16, 64)
	}
	return ToUint64FromString(s)
}

func ToStringFromString(s string) (string, error) {
	return s, nil
}

func ToDurationFromSeconds(s string) (Duration, error) {
	return time.ParseDuration(s + "s")
}

func ToAddressFromString(s string) (Address, error) {
	if len(s) < 3 || (!strings.HasPrefix(s, "0x") && !strings.HasPrefix(s, "0X")) {
		return Address{}, fmt.Errorf("invalid address '%s'", s)
	}
	s = s[2:]
	b, err := hex.DecodeString(s)
	if err != nil {
		return Address{}, err
	}
	return common.BytesToAddress(b), nil
}

func ToApplicationNameFromString(s string) (string, error) {
	if s == "" {
		return "", fmt.Errorf("application name cannot be empty")
	}
	validNamePattern := regexp.MustCompile(`^[a-z0-9_-]+$`)
	if !validNamePattern.MatchString(s) {
		return "", fmt.Errorf("invalid application name '%s': must contain only lowercase letters, numbers, underscores, and hyphens", s)
	}
	return s, nil
}

// Aliases to be used by the generated functions.
var (
	toBool     = strconv.ParseBool
	toUint64   = ToUint64FromString
	toString   = ToStringFromString
	toDuration = ToDurationFromSeconds
	toAddress  = ToAddressFromString
)

var (
	notDefinedbool     = func() bool { return false }
	notDefineduint64   = func() uint64 { return 0 }
	notDefinedstring   = func() string { return "" }
	notDefinedDuration = func() Duration { return 0 }
	notDefinedAddress  = func() Address { return common.Address{} }
)
