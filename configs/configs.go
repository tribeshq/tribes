// The config package manages the node configuration, which comes from environment variables.
// The sub-package generate specifies these environment variables.
package configs

import (
	"encoding/hex"
	"fmt"
	"log/slog"
	"net/url"
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
	URL            = *url.URL
	Duration       = time.Duration
	LogLevel       = slog.Level
	RedactedString = Redacted[string]
	RedactedUint   = Redacted[uint32]
	Address        = common.Address
)

// ------------------------------------------------------------------------------------------------
// Auth Kind
// ------------------------------------------------------------------------------------------------

type AuthKind uint8

const (
	AuthKindPrivateKeyVar AuthKind = iota
	AuthKindPrivateKeyFile
	AuthKindMnemonicVar
	AuthKindMnemonicFile
	AuthKindAWS
)

// ------------------------------------------------------------------------------------------------
// Parsing functions
// ------------------------------------------------------------------------------------------------

func ToUint64FromString(s string) (uint64, error) {
	value, err := strconv.ParseUint(s, 10, 64)
	return value, err
}

func ToStringFromString(s string) (string, error) {
	return s, nil
}

func ToDurationFromSeconds(s string) (time.Duration, error) {
	return time.ParseDuration(s + "s")
}

func ToLogLevelFromString(s string) (LogLevel, error) {
	var m = map[string]LogLevel{
		"debug": slog.LevelDebug,
		"info":  slog.LevelInfo,
		"warn":  slog.LevelWarn,
		"error": slog.LevelError,
	}
	if v, ok := m[s]; ok {
		return v, nil
	} else {
		var zeroValue LogLevel
		return zeroValue, fmt.Errorf("invalid log level '%s'", s)
	}
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

func ToAuthKindFromString(s string) (AuthKind, error) {
	var m = map[string]AuthKind{
		"private_key":      AuthKindPrivateKeyVar,
		"private_key_file": AuthKindPrivateKeyFile,
		"mnemonic":         AuthKindMnemonicVar,
		"mnemonic_file":    AuthKindMnemonicFile,
		"aws":              AuthKindAWS,
	}
	if v, ok := m[s]; ok {
		return v, nil
	} else {
		var zeroValue AuthKind
		return zeroValue, fmt.Errorf("invalid auth kind '%s'", s)
	}
}

func ToRedactedStringFromString(s string) (RedactedString, error) {
	return RedactedString{s}, nil
}

func ToRedactedUint32FromString(s string) (RedactedUint, error) {
	value, err := strconv.ParseUint(s, 10, 32)
	return RedactedUint{uint32(value)}, err
}

func ToURLFromString(s string) (URL, error) {
	result, err := url.Parse(s)
	if err != nil {
		return nil, fmt.Errorf("invalid URL [Redacted]")
	}
	return result, nil
}

// Aliases to be used by the generated functions.
var (
	toBool           = strconv.ParseBool
	toUint64         = ToUint64FromString
	toString         = ToStringFromString
	toLogLevel       = ToLogLevelFromString
	toAuthKind       = ToAuthKindFromString
	toRedactedString = ToRedactedStringFromString
	toRedactedUint   = ToRedactedUint32FromString
	toURL            = ToURLFromString
	toAddress        = ToAddressFromString
)

var (
	notDefinedbool           = func() bool { return false }
	notDefineduint64         = func() uint64 { return 0 }
	notDefinedstring         = func() string { return "" }
	notDefinedLogLevel       = func() slog.Level { return slog.LevelInfo }
	notDefinedAuthKind       = func() AuthKind { return AuthKindMnemonicVar }
	notDefinedRedactedString = func() RedactedString { return RedactedString{""} }
	notDefinedRedactedUint   = func() RedactedUint { return RedactedUint{0} }
	notDefinedURL            = func() URL { return &url.URL{} }
	notDefinedAddress        = func() Address { return common.Address{} }
)
