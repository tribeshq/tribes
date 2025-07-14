package deploy

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func ComputeCreate2AddressFromJSON(jsonPath string, bytecodeKey string, factory common.Address, salt common.Hash) (common.Address, error) {
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to read file: %w", err)
	}

	var obj map[string]interface{}
	if err := json.Unmarshal(data, &obj); err != nil {
		return common.Address{}, fmt.Errorf("failed to unmarshal json: %w", err)
	}

	val, ok := obj[bytecodeKey]
	if !ok {
		return common.Address{}, fmt.Errorf("key %s not found in json", bytecodeKey)
	}
	bytecodeHex, ok := val.(string)
	if !ok {
		return common.Address{}, fmt.Errorf("value at key %s is not a string", bytecodeKey)
	}

	bytecodeHex = strings.TrimPrefix(bytecodeHex, "0x")
	bytecode, err := hex.DecodeString(bytecodeHex)
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to decode bytecode hex: %w", err)
	}

	inithash := crypto.Keccak256(bytecode)
	address := crypto.CreateAddress2(factory, salt, inithash)
	return address, nil
}
