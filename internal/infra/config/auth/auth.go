package auth

import (
	"fmt"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"

	. "github.com/tribeshq/tribes/internal/infra/config"
	"github.com/tribeshq/tribes/pkg/ethutil"
)

func GetTransactOpts(chainId *big.Int) (*bind.TransactOpts, error) {
	authKind, err := GetTribesAuthKind()
	if err != nil {
		return nil, err
	}
	switch authKind {
	case AuthKindMnemonicVar:
		mnemonic, err := GetTribesAuthMnemonic()
		if err != nil {
			return nil, err
		}
		accountIndex, err := GetTribesAuthMnemonicAccountIndex()
		if err != nil {
			return nil, err
		}
		privateKey, err := ethutil.MnemonicToPrivateKey(mnemonic.Value, accountIndex.Value)
		if err != nil {
			return nil, err
		}
		return bind.NewKeyedTransactorWithChainID(privateKey, chainId)
	case AuthKindMnemonicFile:
		mnemonicFile, err := GetTribesAuthMnemonicFile()
		if err != nil {
			return nil, err
		}
		mnemonic, err := os.ReadFile(mnemonicFile)
		if err != nil {
			return nil, err
		}
		accountIndex, err := GetTribesAuthMnemonicAccountIndex()
		if err != nil {
			return nil, err
		}
		privateKey, err := ethutil.MnemonicToPrivateKey(string(mnemonic), accountIndex.Value)
		if err != nil {
			return nil, err
		}
		return bind.NewKeyedTransactorWithChainID(privateKey, chainId)
	case AuthKindPrivateKeyVar:
		privateKey, err := GetTribesAuthPrivateKey()
		if err != nil {
			return nil, err
		}
		key, err := crypto.HexToECDSA(privateKey.Value)
		if err != nil {
			return nil, err
		}
		return bind.NewKeyedTransactorWithChainID(key, chainId)
	case AuthKindPrivateKeyFile:
		privateKeyFile, err := GetTribesAuthPrivateKeyFile()
		if err != nil {
			return nil, err
		}
		privateKey, err := os.ReadFile(privateKeyFile)
		if err != nil {
			return nil, err
		}
		key, err := crypto.HexToECDSA(string(privateKey))
		if err != nil {
			return nil, err
		}
		return bind.NewKeyedTransactorWithChainID(key, chainId)
	default:
		return nil, fmt.Errorf("no valid authentication method found")
	}
}
