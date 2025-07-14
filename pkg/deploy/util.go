package deploy

import (
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	"github.com/tidwall/gjson"
)

func GetBytecodeFromJSON(jsonPath string, bytecodeKey string) ([]byte, error) {
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	result := gjson.GetBytes(data, bytecodeKey)
	if !result.Exists() {
		return nil, fmt.Errorf("key %s not found in json", bytecodeKey)
	}

	bytecodeHex := result.String()
	bytecodeHex = strings.TrimPrefix(bytecodeHex, "0x")
	bytecode, err := hex.DecodeString(bytecodeHex)
	if err != nil {
		return nil, fmt.Errorf("failed to decode bytecode hex: %w", err)
	}

	return bytecode, nil
}
