package chaintracer

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
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

func NewIncrementingTracer(startingBlockNumber int64, rateLimit time.Duration) Tracer {
	return &IncrementingTracer{
		BlockNumber: startingBlockNumber,
		RateLimit:   rateLimit,
	}
}

type IncrementingTracer struct {
	BlockNumber int64
	RateLimit   time.Duration
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
	log.Info().Int64("BlockNumber", it.BlockNumber).Msg("Incremented Block Number")
	return it.BlockNumber
}

func (it *IncrementingTracer) RetrieveDetails() (block *BlockDetails, BlockBuilder string, err error) {
	log.Info().Str("RateLimit", it.RateLimit.String()).Msg("Waiting for Rate Limit")
	time.Sleep(it.RateLimit)
	retries := 0
	var blockData *BlockDetails
	log.Debug().Msg("Starting Retreival of Block Details")
	// Retrieve Block Details from Infura
	for blockData == nil {
		log.Debug().Msg("fetching infura data")
		blockData = InfuraData(it.BlockNumber)
		time.Sleep(time.Duration(retries*2) * time.Second)
		if retries > 5 {
			log.Error().Msg("Error: Could not retrieve block data")
			return nil, "", errors.New("Error: Could not retrieve block data")
		}
		retries += 1
	}
	retries = 0
	var builderName string
	for builderName == "" {
		log.Debug().Msg("fetching Payloadsde data")
		builderName, err = PayloadsDe(it.BlockNumber)
		if err != nil {
			log.Error().Err(err).Msg("Error: Could not retrieve block data")
		}
		time.Sleep(time.Duration(retries*2) * time.Second)
		if retries > 5 {
			return nil, "", errors.New("Error: Could not retrieve block data")
		}
		retries += 1
	}
	log.Info().
		Int64("BlockNumber", it.BlockNumber).
		Str("Builder", builderName).
		Int("Txns Received", len(blockData.Transactions)).
		Msg("Finished Retreival of Block Details")

	return blockData, builderName, nil
}

func PayloadsDe(blockNumber int64) (string, error) {
	url := "https://api.payload.de/block_info?block=" + strconv.FormatInt(blockNumber, 10)
	resp, err := http.Get(url)
	log.Debug().Msg("payloadsde data recieved")
	if err != nil {
		log.Error().Err(err).Str("URL", url).Msg("Error sending request")
		return "", err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error().Err(err).Msg("Error reading response")
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
		log.Error().Err(err).Msg("Error marshalling payload")
		return nil
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		log.Error().Err(err).Msg("Error creating request")
		return nil
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth("", "66f7989079d54ad6988e3b083da0723a")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("Error sending request")
		return nil
	}
	defer resp.Body.Close()

	var infuraresp struct {
		Result BlockDetails `json:"result"`
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error().Err(err).Msg("Error reading response")
		return nil
	}
	json.Unmarshal(data, &infuraresp)
	if infuraresp.Result.BlockNumber == "" {
		log.Error().Msg("Error: Block data is empty")
		return nil
	}
	return &infuraresp.Result
}
