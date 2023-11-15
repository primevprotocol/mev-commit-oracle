package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/primevprotocol/oracle/pkg/rollupclient"
)

func getAuth(privateKey *ecdsa.PrivateKey, chainID *big.Int, client *ethclient.Client) (opts *bind.TransactOpts) {
	// Set transaction opts
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		fmt.Errorf("Failed to construct auth: %v", auth)
	}
	fmt.Println(client.ChainID(context.Background()))

	// Set nonce (optional)
	nonce, err := client.PendingNonceAt(context.Background(), auth.From)
	if err != nil {
		fmt.Printf("Failed to get nonce: %v", err)
	}
	auth.Nonce = big.NewInt(int64(nonce))

	// Set gas price (optional)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		fmt.Printf("Failed to get gas price: %v", err)
	}
	auth.GasPrice = gasPrice

	// Set gas limit (you need to estimate or set a fixed value)
	auth.GasLimit = uint64(300000) // Example value

	return auth
}

func main() {
	fmt.Println("Hello")

	CONTRACT_ADDRESS := "0xCf7Ed3AccA5a467e9e704C703E8D87F634fB0Fc9"
	CLIENT_URL := "http://localhost:8545"

	client, err := ethclient.Dial(CLIENT_URL)
	if err != nil {
		fmt.Printf("Failed to connect to the Ethereum client: %v", err)
	}

	// tracer := chaintracer.NewIncrementingTracer(18293308)

	// for details, builder, err := tracer.RetrieveDetails(); err != nil; {

	// }

	rc, err := rollupclient.NewClient(common.HexToAddress(CONTRACT_ADDRESS), client)
	if err != nil {
		fmt.Println("error creating rollup client")
	}
	CHAIN_ID := big.NewInt(31337)

	privateKey, err := crypto.HexToECDSA("ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")
	if err != nil {
		fmt.Println("error creating private key")
	}

	txn, err := rc.AddBuilderAddress(getAuth(privateKey, CHAIN_ID, client), "k builder", common.HexToAddress("0x5FbDB2315678afecb367f032d93F642f64180aa3"))
	if err != nil {
		fmt.Println("error adding builder address")
	}
	fmt.Println(txn.Hash().String())

	txn2, err := rc.ReceiveBlockData(getAuth(privateKey, CHAIN_ID, client), []string{"txn1", "txn2", "txn3"}, big.NewInt(2000), "k builder")
	if err != nil {
		fmt.Printf("error on recieve block data %v", err)
	}
	fmt.Println(txn2.Hash().String())
}
