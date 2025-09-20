// The config package manages the node configuration, which comes from environment variables.
// The sub-package generate specifies these environment variables.
package configs

import (
	"encoding/hex"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

type (
	URL      = *url.URL
	Duration = time.Duration
	Address  = common.Address
)

// ------------------------------------------------------------------------------------------------
// Parsing functions
// ------------------------------------------------------------------------------------------------

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

func ToURLFromString(s string) (URL, error) {
	result, err := url.Parse(s)
	if err != nil {
		return nil, fmt.Errorf("invalid URL [Redacted]")
	}
	return result, nil
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
	toURL      = ToURLFromString
	toDuration = ToDurationFromSeconds
	toAddress  = ToAddressFromString
)

var (
	notDefinedDuration = func() Duration { return 0 }
	notDefinedURL      = func() URL { return &url.URL{} }
	notDefinedAddress  = func() Address { return common.Address{} }
)
