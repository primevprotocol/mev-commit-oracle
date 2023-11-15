package main

import (
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/primevprotocol/oracle/pkg/chaintracer"
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
	auth.GasLimit = uint64(30000000) // Example value

	return auth
}

type dummyTracer struct {
	blockNumberCurrent int64
}

func (d *dummyTracer) IncrementBlock() int64 {
	d.blockNumberCurrent += 1
	return d.blockNumberCurrent
}

func (d *dummyTracer) RetrieveDetails() (block *chaintracer.BlockDetails, BlockBuilder string, err error) {
	block = &chaintracer.BlockDetails{
		BlockNumber:  strconv.FormatInt(d.blockNumberCurrent, 10),
		Transactions: []string{},
	}

	for i := 0; i < 200; i++ {
		randomInt, err := rand.Int(rand.Reader, big.NewInt(1000))
		if err != nil {
			panic(err)
		}
		randomBytes := crypto.Keccak256(randomInt.Bytes())
		block.Transactions = append(block.Transactions, hex.EncodeToString(randomBytes))
	}

	sleepDuration, _ := rand.Int(rand.Reader, big.NewInt(12))

	time.Sleep(time.Duration(sleepDuration.Int64()) * time.Second)
	return block, "k builder", nil
}

func main() {
	fmt.Println("Hello")

	CONTRACT_ADDRESS := "0xCf7Ed3AccA5a467e9e704C703E8D87F634fB0Fc9"
	CLIENT_URL := "http://localhost:8545"

	client, err := ethclient.Dial(CLIENT_URL)
	if err != nil {
		fmt.Printf("Failed to connect to the Ethereum client: %v", err)
	}

	privateKey, err := crypto.HexToECDSA("ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")
	if err != nil {
		fmt.Println("error creating private key")
	}

	rc, err := rollupclient.NewClient(common.HexToAddress(CONTRACT_ADDRESS), client)
	if err != nil {
		fmt.Println("error creating rollup client")
	}
	CHAIN_ID := big.NewInt(31337)
	txn, err := rc.AddBuilderAddress(getAuth(privateKey, CHAIN_ID, client), "k builder", common.HexToAddress("0x5FbDB2315678afecb367f032d93F642f64180aa3"))
	if err != nil {
		fmt.Println("error adding builder address")
	}
	fmt.Println(txn.Hash().String())

	// tracer := chaintracer.NewIncrementingTracer(18293308)
	// TODO(@ckartik): Remove dummy tracer, only mean't to be used offline for testing
	tracer := dummyTracer{10}
	for {
		blockNumber := tracer.IncrementBlock()
		details, builder, err := tracer.RetrieveDetails()
		if err != nil {
			panic(err)
		}
		blockDataTxn, err := rc.ReceiveBlockData(getAuth(privateKey, CHAIN_ID, client), details.Transactions, big.NewInt(blockNumber), builder)
		if err != nil {

			fmt.Printf("error on recieve block data %v", err)
		}
		fmt.Println(blockDataTxn.Hash().String())
	}
}
