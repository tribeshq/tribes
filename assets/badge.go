package assets

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

//go:embed artifacts/Badge.json
var BadgeJson []byte

type BadgeArtifact struct {
	Bytecode string `json:"bytecode"`
}

func GetBadgeBytecode() ([]byte, error) {
	var artifact BadgeArtifact
	if err := json.Unmarshal(BadgeJson, &artifact); err != nil {
		return nil, fmt.Errorf("failed to parse embedded Badge.json: %w", err)
	}

	if artifact.Bytecode == "" {
		return nil, fmt.Errorf("bytecode not found in embedded Badge.json")
	}

	bytecode := common.Hex2Bytes(strings.TrimPrefix(artifact.Bytecode, "0x"))

	return bytecode, nil
}
