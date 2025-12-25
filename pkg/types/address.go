package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

type Address common.Address

func (a *Address) Scan(value any) error {
	var hex string
	switch v := value.(type) {
	case string:
		hex = v
	case []byte:
		hex = string(v)
	default:
		return fmt.Errorf("unsupported type for address scan: %T", value)
	}
	if !common.IsHexAddress(hex) {
		return fmt.Errorf("invalid hex address: %s", hex)
	}
	*a = Address(common.HexToAddress(hex))
	return nil
}

func (a Address) Value() (driver.Value, error) {
	return a.Hex(), nil
}

func (a Address) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.Hex())
}

func (a *Address) UnmarshalJSON(data []byte) error {
	var hex string
	if err := json.Unmarshal(data, &hex); err != nil {
		return fmt.Errorf("failed to unmarshal address: %w", err)
	}
	if !common.IsHexAddress(hex) {
		return fmt.Errorf("invalid hex address: %s", hex)
	}
	*a = Address(common.HexToAddress(hex))
	return nil
}

func (a Address) String() string {
	return a.Hex()
}

func (a Address) Hex() string {
	return common.Address(a).Hex()
}

func (a Address) IsZero() bool {
	return a == Address{}
}

func HexToAddress(hex string) Address {
	return Address(common.HexToAddress(hex))
}
