package chaintracer_test

import (
	"context"
	"log"
	"math/big"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/primevprotocol/oracle/pkg/chaintracer"
	"github.com/primevprotocol/oracle/pkg/rollupclient"
)

func TestDataPull(t *testing.T) {

	client, err := ethclient.Dial("http://localhost:8545")
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	tracer := chaintracer.NewIncrementingTracer(18293308)
	_, builder, err := tracer.RetrieveDetails()

	if !reflect.DeepEqual("titanbuilder", builder) {
		t.Error("winning builder is not titanbuilder for block 18293308")
	}
	if err != nil {
		t.Error("error retrieving block details")
	}
	rc, err := rollupclient.NewClient(common.HexToAddress("0x5FbDB2315678afecb367f032d93F642f64180aa3"), client)
	if err != nil {
		t.Error("error creating rollup client")
	}
	var CHAIN_ID int64 = 31337

	privateKey, err := crypto.HexToECDSA("ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")
	if err != nil {
		t.Error("error creating private key")
	}
	// Set transaction opts
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(CHAIN_ID))
	if err != nil {
		t.Error("error creating transaction opts")
	}
	t.Log(client.ChainID(context.Background()))

	// Set nonce (optional)
	nonce, err := client.PendingNonceAt(context.Background(), auth.From)
	if err != nil {
		log.Fatalf("Failed to get nonce: %v", err)
	}
	auth.Nonce = big.NewInt(int64(nonce))

	// Set gas price (optional)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatalf("Failed to get gas price: %v", err)
	}
	auth.GasPrice = gasPrice

	// Set gas limit (you need to estimate or set a fixed value)
	auth.GasLimit = uint64(300000) // Example value

	txn, err := rc.AddBuilderAddress(auth, "k builder", common.HexToAddress("0x5FbDB2315678afecb367f032d93F642f64180aa3"))
	if err != nil {
		t.Error("error adding builder address")
	}
	t.Log(txn.Hash().String())

}
