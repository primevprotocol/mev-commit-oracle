package keysigner

import (
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type KeySigner interface {
	GetAddress() common.Address
	GetAuth(chainID *big.Int) (*bind.TransactOpts, error)
}

type privateKeySigner struct {
	privKey *ecdsa.PrivateKey
}

func NewPrivateKeySigner(privKey *ecdsa.PrivateKey) *privateKeySigner {
	return &privateKeySigner{
		privKey: privKey,
	}
}

func (pks *privateKeySigner) GetAddress() common.Address {
	return crypto.PubkeyToAddress(pks.privKey.PublicKey)
}

func (pks *privateKeySigner) GetAuth(chainID *big.Int) (*bind.TransactOpts, error) {
	return bind.NewKeyedTransactorWithChainID(pks.privKey, chainID)
}

type keystoreSigner struct {
	keystore *keystore.KeyStore
	password string
	account  accounts.Account
}

func NewKeystoreSigner(keystore *keystore.KeyStore, password string, account accounts.Account) *keystoreSigner {
	return &keystoreSigner{
		keystore: keystore,
		password: password,
		account:  account,
	}
}

func (kss *keystoreSigner) GetAddress() common.Address {
	return kss.account.Address
}

func (kss *keystoreSigner) GetAuth(chainID *big.Int) (*bind.TransactOpts, error) {
	if err := kss.keystore.Unlock(kss.account, kss.password); err != nil {
		return nil, err
	}

	return bind.NewKeyStoreTransactorWithChainID(kss.keystore, kss.account, chainID)
}
