package chaintracer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

type TransactionHash []byte

type BlockDetails struct {
	Transactions []string `json:"transactions"`
	BlockNumber  string   `json:"number"`
}

type InfuraResponse struct {
	Jsonrpc string       `json:"jsonrpc"`
	ID      int          `json:"id"`
	Result  BlockDetails `json:"result"`
}

func NewIncrementingTracer(startingBlockNumber int64) Tracer {
	return &IncrementingTracer{
		BlockNumber: startingBlockNumber,
	}
}

type IncrementingTracer struct {
	BlockNumber int64
}

type Tracer interface {
	IncrementBlock() (NewBlockNumber int64)
	RetrieveDetails() (block *BlockDetails, BlockBuilder string, err error)
}

type Payload struct {
	Jsonrpc string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      int           `json:"id"`
}

func (it *IncrementingTracer) IncrementBlock() (NewBlockNumber int64) {
	it.BlockNumber += 1
	fmt.Println(it.BlockNumber)
	return it.BlockNumber
}

func (it *IncrementingTracer) RetrieveDetails() (block *BlockDetails, BlockBuilder string, err error) {
	retries := 0
	var blockData *BlockDetails
	// Retrieve Block Details from Infura
	for blockData == nil {
		blockData = InfuraData(it.BlockNumber)
		time.Sleep(time.Duration(retries*2) * time.Second)
		if retries > 5 {
			return nil, "", errors.New("Error: Could not retrieve block data")
		}
		retries += 1
	}
	retries = 0
	var builderName string
	for builderName == "" {
		builderName, err = PayloadsDe(it.BlockNumber)
		if err != nil {
			fmt.Println("Error: Could not retrieve block data")
		}
		time.Sleep(time.Duration(retries*2) * time.Second)
		if retries > 5 {
			return nil, "", errors.New("Error: Could not retrieve block data")
		}
		retries += 1
	}

	return blockData, builderName, nil
}

func PayloadsDe(blockNumber int64) (string, error) {
	url := fmt.Sprintf("https://api.payload.de/block_info?block=%d", blockNumber)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return "", err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return "", err
	}
	var payload struct {
		Builder string `json:"builder"`
	}
	json.Unmarshal(data, &payload)
	return payload.Builder, nil
}

func InfuraData(blockNumber int64) *BlockDetails {
	hex := strconv.FormatInt(blockNumber, 16)
	url := "https://mainnet.infura.io/v3/b8877b173a0543bea7dca82c313e7347"
	payload := Payload{
		Jsonrpc: "2.0",
		Method:  "eth_getBlockByNumber",
		Params:  []interface{}{"0x" + hex, false},
		ID:      1,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error marshalling payload:", err)
		return nil
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth("", "66f7989079d54ad6988e3b083da0723a")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return nil
	}
	defer resp.Body.Close()

	var infuraresp struct {
		Result BlockDetails `json:"result"`
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return nil
	}
	json.Unmarshal(data, &infuraresp)
	if infuraresp.Result.BlockNumber == "" {
		fmt.Println("Error: Block data is empty")
		return nil
	}
	return &infuraresp.Result
}
