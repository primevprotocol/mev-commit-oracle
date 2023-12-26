// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package preconf

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// PreConfCommitmentStorePreConfCommitment is an auto generated low-level Go binding around an user-defined struct.
type PreConfCommitmentStorePreConfCommitment struct {
	CommitmentUsed      bool
	Bidder              common.Address
	Commiter            common.Address
	Bid                 uint64
	BlockNumber         uint64
	BidHash             [32]byte
	TxnHash             string
	CommitmentHash      [32]byte
	BidSignature        []byte
	CommitmentSignature []byte
}

// PreconfMetaData contains all meta data concerning the Preconf contract.
var PreconfMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_providerRegistry\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_bidderRegistry\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_oracle\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_owner\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"fallback\",\"stateMutability\":\"payable\"},{\"type\":\"receive\",\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"DOMAIN_SEPARATOR_BID\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"DOMAIN_SEPARATOR_PRECONF\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"EIP712_COMMITMENT_TYPEHASH\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"EIP712_MESSAGE_TYPEHASH\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"_bytesToHexString\",\"inputs\":[{\"name\":\"_bytes\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"bidderRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIBidderRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"blockCommitments\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"commitmentCount\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"commitments\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"commitmentUsed\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"bidder\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"commiter\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"bid\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"blockNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"bidHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"txnHash\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"commitmentHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"bidSignature\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"commitmentSignature\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"commitmentsCount\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getBidHash\",\"inputs\":[{\"name\":\"_txnHash\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_bid\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"_blockNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCommitment\",\"inputs\":[{\"name\":\"commitmentIndex\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structPreConfCommitmentStore.PreConfCommitment\",\"components\":[{\"name\":\"commitmentUsed\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"bidder\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"commiter\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"bid\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"blockNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"bidHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"txnHash\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"commitmentHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"bidSignature\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"commitmentSignature\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCommitmentIndex\",\"inputs\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structPreConfCommitmentStore.PreConfCommitment\",\"components\":[{\"name\":\"commitmentUsed\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"bidder\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"commiter\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"bid\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"blockNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"bidHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"txnHash\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"commitmentHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"bidSignature\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"commitmentSignature\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"getCommitmentsByBlockNumber\",\"inputs\":[{\"name\":\"blockNumber\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCommitmentsByCommitter\",\"inputs\":[{\"name\":\"commiter\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPreConfHash\",\"inputs\":[{\"name\":\"_txnHash\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_bid\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"_blockNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"_bidHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"_bidSignature\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getTxnHashFromCommitment\",\"inputs\":[{\"name\":\"commitmentIndex\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"txnHash\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initateReward\",\"inputs\":[{\"name\":\"commitmentIndex\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"initiateSlash\",\"inputs\":[{\"name\":\"commitmentIndex\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"lastProcessedBlock\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"oracle\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"providerCommitments\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"providerRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIProviderRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"storeCommitment\",\"inputs\":[{\"name\":\"bid\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"blockNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"txnHash\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"bidSignature\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"commitmentSignature\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"commitmentIndex\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateBidderRegistry\",\"inputs\":[{\"name\":\"newBidderRegistry\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateOracle\",\"inputs\":[{\"name\":\"newOracle\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateProviderRegistry\",\"inputs\":[{\"name\":\"newProviderRegistry\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"verifyBid\",\"inputs\":[{\"name\":\"bid\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"blockNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"txnHash\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"bidSignature\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"messageDigest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"recoveredAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"stake\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyPreConfCommitment\",\"inputs\":[{\"name\":\"txnHash\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"bid\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"blockNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"bidHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"bidSignature\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"commitmentSignature\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"preConfHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"commiterAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SignatureVerified\",\"inputs\":[{\"name\":\"signer\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"txnHash\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"bid\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"blockNumber\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false}]",
}

// PreconfABI is the input ABI used to generate the binding from.
// Deprecated: Use PreconfMetaData.ABI instead.
var PreconfABI = PreconfMetaData.ABI

// Preconf is an auto generated Go binding around an Ethereum contract.
type Preconf struct {
	PreconfCaller     // Read-only binding to the contract
	PreconfTransactor // Write-only binding to the contract
	PreconfFilterer   // Log filterer for contract events
}

// PreconfCaller is an auto generated read-only Go binding around an Ethereum contract.
type PreconfCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PreconfTransactor is an auto generated write-only Go binding around an Ethereum contract.
type PreconfTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PreconfFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type PreconfFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PreconfSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type PreconfSession struct {
	Contract     *Preconf          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// PreconfCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type PreconfCallerSession struct {
	Contract *PreconfCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// PreconfTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type PreconfTransactorSession struct {
	Contract     *PreconfTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// PreconfRaw is an auto generated low-level Go binding around an Ethereum contract.
type PreconfRaw struct {
	Contract *Preconf // Generic contract binding to access the raw methods on
}

// PreconfCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type PreconfCallerRaw struct {
	Contract *PreconfCaller // Generic read-only contract binding to access the raw methods on
}

// PreconfTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type PreconfTransactorRaw struct {
	Contract *PreconfTransactor // Generic write-only contract binding to access the raw methods on
}

// NewPreconf creates a new instance of Preconf, bound to a specific deployed contract.
func NewPreconf(address common.Address, backend bind.ContractBackend) (*Preconf, error) {
	contract, err := bindPreconf(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Preconf{PreconfCaller: PreconfCaller{contract: contract}, PreconfTransactor: PreconfTransactor{contract: contract}, PreconfFilterer: PreconfFilterer{contract: contract}}, nil
}

// NewPreconfCaller creates a new read-only instance of Preconf, bound to a specific deployed contract.
func NewPreconfCaller(address common.Address, caller bind.ContractCaller) (*PreconfCaller, error) {
	contract, err := bindPreconf(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &PreconfCaller{contract: contract}, nil
}

// NewPreconfTransactor creates a new write-only instance of Preconf, bound to a specific deployed contract.
func NewPreconfTransactor(address common.Address, transactor bind.ContractTransactor) (*PreconfTransactor, error) {
	contract, err := bindPreconf(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &PreconfTransactor{contract: contract}, nil
}

// NewPreconfFilterer creates a new log filterer instance of Preconf, bound to a specific deployed contract.
func NewPreconfFilterer(address common.Address, filterer bind.ContractFilterer) (*PreconfFilterer, error) {
	contract, err := bindPreconf(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &PreconfFilterer{contract: contract}, nil
}

// bindPreconf binds a generic wrapper to an already deployed contract.
func bindPreconf(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := PreconfMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Preconf *PreconfRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Preconf.Contract.PreconfCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Preconf *PreconfRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Preconf.Contract.PreconfTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Preconf *PreconfRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Preconf.Contract.PreconfTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Preconf *PreconfCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Preconf.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Preconf *PreconfTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Preconf.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Preconf *PreconfTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Preconf.Contract.contract.Transact(opts, method, params...)
}

// DOMAINSEPARATORBID is a free data retrieval call binding the contract method 0x940b5765.
//
// Solidity: function DOMAIN_SEPARATOR_BID() view returns(bytes32)
func (_Preconf *PreconfCaller) DOMAINSEPARATORBID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Preconf.contract.Call(opts, &out, "DOMAIN_SEPARATOR_BID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DOMAINSEPARATORBID is a free data retrieval call binding the contract method 0x940b5765.
//
// Solidity: function DOMAIN_SEPARATOR_BID() view returns(bytes32)
func (_Preconf *PreconfSession) DOMAINSEPARATORBID() ([32]byte, error) {
	return _Preconf.Contract.DOMAINSEPARATORBID(&_Preconf.CallOpts)
}

// DOMAINSEPARATORBID is a free data retrieval call binding the contract method 0x940b5765.
//
// Solidity: function DOMAIN_SEPARATOR_BID() view returns(bytes32)
func (_Preconf *PreconfCallerSession) DOMAINSEPARATORBID() ([32]byte, error) {
	return _Preconf.Contract.DOMAINSEPARATORBID(&_Preconf.CallOpts)
}

// DOMAINSEPARATORPRECONF is a free data retrieval call binding the contract method 0xe5ae370f.
//
// Solidity: function DOMAIN_SEPARATOR_PRECONF() view returns(bytes32)
func (_Preconf *PreconfCaller) DOMAINSEPARATORPRECONF(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Preconf.contract.Call(opts, &out, "DOMAIN_SEPARATOR_PRECONF")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DOMAINSEPARATORPRECONF is a free data retrieval call binding the contract method 0xe5ae370f.
//
// Solidity: function DOMAIN_SEPARATOR_PRECONF() view returns(bytes32)
func (_Preconf *PreconfSession) DOMAINSEPARATORPRECONF() ([32]byte, error) {
	return _Preconf.Contract.DOMAINSEPARATORPRECONF(&_Preconf.CallOpts)
}

// DOMAINSEPARATORPRECONF is a free data retrieval call binding the contract method 0xe5ae370f.
//
// Solidity: function DOMAIN_SEPARATOR_PRECONF() view returns(bytes32)
func (_Preconf *PreconfCallerSession) DOMAINSEPARATORPRECONF() ([32]byte, error) {
	return _Preconf.Contract.DOMAINSEPARATORPRECONF(&_Preconf.CallOpts)
}

// EIP712COMMITMENTTYPEHASH is a free data retrieval call binding the contract method 0x10ce6471.
//
// Solidity: function EIP712_COMMITMENT_TYPEHASH() view returns(bytes32)
func (_Preconf *PreconfCaller) EIP712COMMITMENTTYPEHASH(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Preconf.contract.Call(opts, &out, "EIP712_COMMITMENT_TYPEHASH")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// EIP712COMMITMENTTYPEHASH is a free data retrieval call binding the contract method 0x10ce6471.
//
// Solidity: function EIP712_COMMITMENT_TYPEHASH() view returns(bytes32)
func (_Preconf *PreconfSession) EIP712COMMITMENTTYPEHASH() ([32]byte, error) {
	return _Preconf.Contract.EIP712COMMITMENTTYPEHASH(&_Preconf.CallOpts)
}

// EIP712COMMITMENTTYPEHASH is a free data retrieval call binding the contract method 0x10ce6471.
//
// Solidity: function EIP712_COMMITMENT_TYPEHASH() view returns(bytes32)
func (_Preconf *PreconfCallerSession) EIP712COMMITMENTTYPEHASH() ([32]byte, error) {
	return _Preconf.Contract.EIP712COMMITMENTTYPEHASH(&_Preconf.CallOpts)
}

// EIP712MESSAGETYPEHASH is a free data retrieval call binding the contract method 0xc24fb639.
//
// Solidity: function EIP712_MESSAGE_TYPEHASH() view returns(bytes32)
func (_Preconf *PreconfCaller) EIP712MESSAGETYPEHASH(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Preconf.contract.Call(opts, &out, "EIP712_MESSAGE_TYPEHASH")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// EIP712MESSAGETYPEHASH is a free data retrieval call binding the contract method 0xc24fb639.
//
// Solidity: function EIP712_MESSAGE_TYPEHASH() view returns(bytes32)
func (_Preconf *PreconfSession) EIP712MESSAGETYPEHASH() ([32]byte, error) {
	return _Preconf.Contract.EIP712MESSAGETYPEHASH(&_Preconf.CallOpts)
}

// EIP712MESSAGETYPEHASH is a free data retrieval call binding the contract method 0xc24fb639.
//
// Solidity: function EIP712_MESSAGE_TYPEHASH() view returns(bytes32)
func (_Preconf *PreconfCallerSession) EIP712MESSAGETYPEHASH() ([32]byte, error) {
	return _Preconf.Contract.EIP712MESSAGETYPEHASH(&_Preconf.CallOpts)
}

// BytesToHexString is a free data retrieval call binding the contract method 0xca64db2e.
//
// Solidity: function _bytesToHexString(bytes _bytes) pure returns(string)
func (_Preconf *PreconfCaller) BytesToHexString(opts *bind.CallOpts, _bytes []byte) (string, error) {
	var out []interface{}
	err := _Preconf.contract.Call(opts, &out, "_bytesToHexString", _bytes)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// BytesToHexString is a free data retrieval call binding the contract method 0xca64db2e.
//
// Solidity: function _bytesToHexString(bytes _bytes) pure returns(string)
func (_Preconf *PreconfSession) BytesToHexString(_bytes []byte) (string, error) {
	return _Preconf.Contract.BytesToHexString(&_Preconf.CallOpts, _bytes)
}

// BytesToHexString is a free data retrieval call binding the contract method 0xca64db2e.
//
// Solidity: function _bytesToHexString(bytes _bytes) pure returns(string)
func (_Preconf *PreconfCallerSession) BytesToHexString(_bytes []byte) (string, error) {
	return _Preconf.Contract.BytesToHexString(&_Preconf.CallOpts, _bytes)
}

// BidderRegistry is a free data retrieval call binding the contract method 0x909e54e2.
//
// Solidity: function bidderRegistry() view returns(address)
func (_Preconf *PreconfCaller) BidderRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Preconf.contract.Call(opts, &out, "bidderRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// BidderRegistry is a free data retrieval call binding the contract method 0x909e54e2.
//
// Solidity: function bidderRegistry() view returns(address)
func (_Preconf *PreconfSession) BidderRegistry() (common.Address, error) {
	return _Preconf.Contract.BidderRegistry(&_Preconf.CallOpts)
}

// BidderRegistry is a free data retrieval call binding the contract method 0x909e54e2.
//
// Solidity: function bidderRegistry() view returns(address)
func (_Preconf *PreconfCallerSession) BidderRegistry() (common.Address, error) {
	return _Preconf.Contract.BidderRegistry(&_Preconf.CallOpts)
}

// BlockCommitments is a free data retrieval call binding the contract method 0x159efb47.
//
// Solidity: function blockCommitments(uint256 , uint256 ) view returns(bytes32)
func (_Preconf *PreconfCaller) BlockCommitments(opts *bind.CallOpts, arg0 *big.Int, arg1 *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _Preconf.contract.Call(opts, &out, "blockCommitments", arg0, arg1)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// BlockCommitments is a free data retrieval call binding the contract method 0x159efb47.
//
// Solidity: function blockCommitments(uint256 , uint256 ) view returns(bytes32)
func (_Preconf *PreconfSession) BlockCommitments(arg0 *big.Int, arg1 *big.Int) ([32]byte, error) {
	return _Preconf.Contract.BlockCommitments(&_Preconf.CallOpts, arg0, arg1)
}

// BlockCommitments is a free data retrieval call binding the contract method 0x159efb47.
//
// Solidity: function blockCommitments(uint256 , uint256 ) view returns(bytes32)
func (_Preconf *PreconfCallerSession) BlockCommitments(arg0 *big.Int, arg1 *big.Int) ([32]byte, error) {
	return _Preconf.Contract.BlockCommitments(&_Preconf.CallOpts, arg0, arg1)
}

// CommitmentCount is a free data retrieval call binding the contract method 0xc44956d1.
//
// Solidity: function commitmentCount() view returns(uint256)
func (_Preconf *PreconfCaller) CommitmentCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Preconf.contract.Call(opts, &out, "commitmentCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CommitmentCount is a free data retrieval call binding the contract method 0xc44956d1.
//
// Solidity: function commitmentCount() view returns(uint256)
func (_Preconf *PreconfSession) CommitmentCount() (*big.Int, error) {
	return _Preconf.Contract.CommitmentCount(&_Preconf.CallOpts)
}

// CommitmentCount is a free data retrieval call binding the contract method 0xc44956d1.
//
// Solidity: function commitmentCount() view returns(uint256)
func (_Preconf *PreconfCallerSession) CommitmentCount() (*big.Int, error) {
	return _Preconf.Contract.CommitmentCount(&_Preconf.CallOpts)
}

// Commitments is a free data retrieval call binding the contract method 0x839df945.
//
// Solidity: function commitments(bytes32 ) view returns(bool commitmentUsed, address bidder, address commiter, uint64 bid, uint64 blockNumber, bytes32 bidHash, string txnHash, bytes32 commitmentHash, bytes bidSignature, bytes commitmentSignature)
func (_Preconf *PreconfCaller) Commitments(opts *bind.CallOpts, arg0 [32]byte) (struct {
	CommitmentUsed      bool
	Bidder              common.Address
	Commiter            common.Address
	Bid                 uint64
	BlockNumber         uint64
	BidHash             [32]byte
	TxnHash             string
	CommitmentHash      [32]byte
	BidSignature        []byte
	CommitmentSignature []byte
}, error) {
	var out []interface{}
	err := _Preconf.contract.Call(opts, &out, "commitments", arg0)

	outstruct := new(struct {
		CommitmentUsed      bool
		Bidder              common.Address
		Commiter            common.Address
		Bid                 uint64
		BlockNumber         uint64
		BidHash             [32]byte
		TxnHash             string
		CommitmentHash      [32]byte
		BidSignature        []byte
		CommitmentSignature []byte
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.CommitmentUsed = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.Bidder = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)
	outstruct.Commiter = *abi.ConvertType(out[2], new(common.Address)).(*common.Address)
	outstruct.Bid = *abi.ConvertType(out[3], new(uint64)).(*uint64)
	outstruct.BlockNumber = *abi.ConvertType(out[4], new(uint64)).(*uint64)
	outstruct.BidHash = *abi.ConvertType(out[5], new([32]byte)).(*[32]byte)
	outstruct.TxnHash = *abi.ConvertType(out[6], new(string)).(*string)
	outstruct.CommitmentHash = *abi.ConvertType(out[7], new([32]byte)).(*[32]byte)
	outstruct.BidSignature = *abi.ConvertType(out[8], new([]byte)).(*[]byte)
	outstruct.CommitmentSignature = *abi.ConvertType(out[9], new([]byte)).(*[]byte)

	return *outstruct, err

}

// Commitments is a free data retrieval call binding the contract method 0x839df945.
//
// Solidity: function commitments(bytes32 ) view returns(bool commitmentUsed, address bidder, address commiter, uint64 bid, uint64 blockNumber, bytes32 bidHash, string txnHash, bytes32 commitmentHash, bytes bidSignature, bytes commitmentSignature)
func (_Preconf *PreconfSession) Commitments(arg0 [32]byte) (struct {
	CommitmentUsed      bool
	Bidder              common.Address
	Commiter            common.Address
	Bid                 uint64
	BlockNumber         uint64
	BidHash             [32]byte
	TxnHash             string
	CommitmentHash      [32]byte
	BidSignature        []byte
	CommitmentSignature []byte
}, error) {
	return _Preconf.Contract.Commitments(&_Preconf.CallOpts, arg0)
}

// Commitments is a free data retrieval call binding the contract method 0x839df945.
//
// Solidity: function commitments(bytes32 ) view returns(bool commitmentUsed, address bidder, address commiter, uint64 bid, uint64 blockNumber, bytes32 bidHash, string txnHash, bytes32 commitmentHash, bytes bidSignature, bytes commitmentSignature)
func (_Preconf *PreconfCallerSession) Commitments(arg0 [32]byte) (struct {
	CommitmentUsed      bool
	Bidder              common.Address
	Commiter            common.Address
	Bid                 uint64
	BlockNumber         uint64
	BidHash             [32]byte
	TxnHash             string
	CommitmentHash      [32]byte
	BidSignature        []byte
	CommitmentSignature []byte
}, error) {
	return _Preconf.Contract.Commitments(&_Preconf.CallOpts, arg0)
}

// CommitmentsCount is a free data retrieval call binding the contract method 0x25f5cf21.
//
// Solidity: function commitmentsCount(address ) view returns(uint256)
func (_Preconf *PreconfCaller) CommitmentsCount(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Preconf.contract.Call(opts, &out, "commitmentsCount", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CommitmentsCount is a free data retrieval call binding the contract method 0x25f5cf21.
//
// Solidity: function commitmentsCount(address ) view returns(uint256)
func (_Preconf *PreconfSession) CommitmentsCount(arg0 common.Address) (*big.Int, error) {
	return _Preconf.Contract.CommitmentsCount(&_Preconf.CallOpts, arg0)
}

// CommitmentsCount is a free data retrieval call binding the contract method 0x25f5cf21.
//
// Solidity: function commitmentsCount(address ) view returns(uint256)
func (_Preconf *PreconfCallerSession) CommitmentsCount(arg0 common.Address) (*big.Int, error) {
	return _Preconf.Contract.CommitmentsCount(&_Preconf.CallOpts, arg0)
}

// GetBidHash is a free data retrieval call binding the contract method 0xbc1d6cdc.
//
// Solidity: function getBidHash(string _txnHash, uint64 _bid, uint64 _blockNumber) view returns(bytes32)
func (_Preconf *PreconfCaller) GetBidHash(opts *bind.CallOpts, _txnHash string, _bid uint64, _blockNumber uint64) ([32]byte, error) {
	var out []interface{}
	err := _Preconf.contract.Call(opts, &out, "getBidHash", _txnHash, _bid, _blockNumber)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetBidHash is a free data retrieval call binding the contract method 0xbc1d6cdc.
//
// Solidity: function getBidHash(string _txnHash, uint64 _bid, uint64 _blockNumber) view returns(bytes32)
func (_Preconf *PreconfSession) GetBidHash(_txnHash string, _bid uint64, _blockNumber uint64) ([32]byte, error) {
	return _Preconf.Contract.GetBidHash(&_Preconf.CallOpts, _txnHash, _bid, _blockNumber)
}

// GetBidHash is a free data retrieval call binding the contract method 0xbc1d6cdc.
//
// Solidity: function getBidHash(string _txnHash, uint64 _bid, uint64 _blockNumber) view returns(bytes32)
func (_Preconf *PreconfCallerSession) GetBidHash(_txnHash string, _bid uint64, _blockNumber uint64) ([32]byte, error) {
	return _Preconf.Contract.GetBidHash(&_Preconf.CallOpts, _txnHash, _bid, _blockNumber)
}

// GetCommitment is a free data retrieval call binding the contract method 0x7795820c.
//
// Solidity: function getCommitment(bytes32 commitmentIndex) view returns((bool,address,address,uint64,uint64,bytes32,string,bytes32,bytes,bytes))
func (_Preconf *PreconfCaller) GetCommitment(opts *bind.CallOpts, commitmentIndex [32]byte) (PreConfCommitmentStorePreConfCommitment, error) {
	var out []interface{}
	err := _Preconf.contract.Call(opts, &out, "getCommitment", commitmentIndex)

	if err != nil {
		return *new(PreConfCommitmentStorePreConfCommitment), err
	}

	out0 := *abi.ConvertType(out[0], new(PreConfCommitmentStorePreConfCommitment)).(*PreConfCommitmentStorePreConfCommitment)

	return out0, err

}

// GetCommitment is a free data retrieval call binding the contract method 0x7795820c.
//
// Solidity: function getCommitment(bytes32 commitmentIndex) view returns((bool,address,address,uint64,uint64,bytes32,string,bytes32,bytes,bytes))
func (_Preconf *PreconfSession) GetCommitment(commitmentIndex [32]byte) (PreConfCommitmentStorePreConfCommitment, error) {
	return _Preconf.Contract.GetCommitment(&_Preconf.CallOpts, commitmentIndex)
}

// GetCommitment is a free data retrieval call binding the contract method 0x7795820c.
//
// Solidity: function getCommitment(bytes32 commitmentIndex) view returns((bool,address,address,uint64,uint64,bytes32,string,bytes32,bytes,bytes))
func (_Preconf *PreconfCallerSession) GetCommitment(commitmentIndex [32]byte) (PreConfCommitmentStorePreConfCommitment, error) {
	return _Preconf.Contract.GetCommitment(&_Preconf.CallOpts, commitmentIndex)
}

// GetCommitmentIndex is a free data retrieval call binding the contract method 0x3f051705.
//
// Solidity: function getCommitmentIndex((bool,address,address,uint64,uint64,bytes32,string,bytes32,bytes,bytes) commitment) pure returns(bytes32)
func (_Preconf *PreconfCaller) GetCommitmentIndex(opts *bind.CallOpts, commitment PreConfCommitmentStorePreConfCommitment) ([32]byte, error) {
	var out []interface{}
	err := _Preconf.contract.Call(opts, &out, "getCommitmentIndex", commitment)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetCommitmentIndex is a free data retrieval call binding the contract method 0x3f051705.
//
// Solidity: function getCommitmentIndex((bool,address,address,uint64,uint64,bytes32,string,bytes32,bytes,bytes) commitment) pure returns(bytes32)
func (_Preconf *PreconfSession) GetCommitmentIndex(commitment PreConfCommitmentStorePreConfCommitment) ([32]byte, error) {
	return _Preconf.Contract.GetCommitmentIndex(&_Preconf.CallOpts, commitment)
}

// GetCommitmentIndex is a free data retrieval call binding the contract method 0x3f051705.
//
// Solidity: function getCommitmentIndex((bool,address,address,uint64,uint64,bytes32,string,bytes32,bytes,bytes) commitment) pure returns(bytes32)
func (_Preconf *PreconfCallerSession) GetCommitmentIndex(commitment PreConfCommitmentStorePreConfCommitment) ([32]byte, error) {
	return _Preconf.Contract.GetCommitmentIndex(&_Preconf.CallOpts, commitment)
}

// GetCommitmentsByBlockNumber is a free data retrieval call binding the contract method 0x82da12de.
//
// Solidity: function getCommitmentsByBlockNumber(uint256 blockNumber) view returns(bytes32[])
func (_Preconf *PreconfCaller) GetCommitmentsByBlockNumber(opts *bind.CallOpts, blockNumber *big.Int) ([][32]byte, error) {
	var out []interface{}
	err := _Preconf.contract.Call(opts, &out, "getCommitmentsByBlockNumber", blockNumber)

	if err != nil {
		return *new([][32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([][32]byte)).(*[][32]byte)

	return out0, err

}

// GetCommitmentsByBlockNumber is a free data retrieval call binding the contract method 0x82da12de.
//
// Solidity: function getCommitmentsByBlockNumber(uint256 blockNumber) view returns(bytes32[])
func (_Preconf *PreconfSession) GetCommitmentsByBlockNumber(blockNumber *big.Int) ([][32]byte, error) {
	return _Preconf.Contract.GetCommitmentsByBlockNumber(&_Preconf.CallOpts, blockNumber)
}

// GetCommitmentsByBlockNumber is a free data retrieval call binding the contract method 0x82da12de.
//
// Solidity: function getCommitmentsByBlockNumber(uint256 blockNumber) view returns(bytes32[])
func (_Preconf *PreconfCallerSession) GetCommitmentsByBlockNumber(blockNumber *big.Int) ([][32]byte, error) {
	return _Preconf.Contract.GetCommitmentsByBlockNumber(&_Preconf.CallOpts, blockNumber)
}

// GetCommitmentsByCommitter is a free data retrieval call binding the contract method 0xac8c8a0e.
//
// Solidity: function getCommitmentsByCommitter(address commiter) view returns(bytes32[])
func (_Preconf *PreconfCaller) GetCommitmentsByCommitter(opts *bind.CallOpts, commiter common.Address) ([][32]byte, error) {
	var out []interface{}
	err := _Preconf.contract.Call(opts, &out, "getCommitmentsByCommitter", commiter)

	if err != nil {
		return *new([][32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([][32]byte)).(*[][32]byte)

	return out0, err

}

// GetCommitmentsByCommitter is a free data retrieval call binding the contract method 0xac8c8a0e.
//
// Solidity: function getCommitmentsByCommitter(address commiter) view returns(bytes32[])
func (_Preconf *PreconfSession) GetCommitmentsByCommitter(commiter common.Address) ([][32]byte, error) {
	return _Preconf.Contract.GetCommitmentsByCommitter(&_Preconf.CallOpts, commiter)
}

// GetCommitmentsByCommitter is a free data retrieval call binding the contract method 0xac8c8a0e.
//
// Solidity: function getCommitmentsByCommitter(address commiter) view returns(bytes32[])
func (_Preconf *PreconfCallerSession) GetCommitmentsByCommitter(commiter common.Address) ([][32]byte, error) {
	return _Preconf.Contract.GetCommitmentsByCommitter(&_Preconf.CallOpts, commiter)
}

// GetPreConfHash is a free data retrieval call binding the contract method 0xe501d8dd.
//
// Solidity: function getPreConfHash(string _txnHash, uint64 _bid, uint64 _blockNumber, bytes32 _bidHash, string _bidSignature) view returns(bytes32)
func (_Preconf *PreconfCaller) GetPreConfHash(opts *bind.CallOpts, _txnHash string, _bid uint64, _blockNumber uint64, _bidHash [32]byte, _bidSignature string) ([32]byte, error) {
	var out []interface{}
	err := _Preconf.contract.Call(opts, &out, "getPreConfHash", _txnHash, _bid, _blockNumber, _bidHash, _bidSignature)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetPreConfHash is a free data retrieval call binding the contract method 0xe501d8dd.
//
// Solidity: function getPreConfHash(string _txnHash, uint64 _bid, uint64 _blockNumber, bytes32 _bidHash, string _bidSignature) view returns(bytes32)
func (_Preconf *PreconfSession) GetPreConfHash(_txnHash string, _bid uint64, _blockNumber uint64, _bidHash [32]byte, _bidSignature string) ([32]byte, error) {
	return _Preconf.Contract.GetPreConfHash(&_Preconf.CallOpts, _txnHash, _bid, _blockNumber, _bidHash, _bidSignature)
}

// GetPreConfHash is a free data retrieval call binding the contract method 0xe501d8dd.
//
// Solidity: function getPreConfHash(string _txnHash, uint64 _bid, uint64 _blockNumber, bytes32 _bidHash, string _bidSignature) view returns(bytes32)
func (_Preconf *PreconfCallerSession) GetPreConfHash(_txnHash string, _bid uint64, _blockNumber uint64, _bidHash [32]byte, _bidSignature string) ([32]byte, error) {
	return _Preconf.Contract.GetPreConfHash(&_Preconf.CallOpts, _txnHash, _bid, _blockNumber, _bidHash, _bidSignature)
}

// GetTxnHashFromCommitment is a free data retrieval call binding the contract method 0xfc4fbe32.
//
// Solidity: function getTxnHashFromCommitment(bytes32 commitmentIndex) view returns(string txnHash)
func (_Preconf *PreconfCaller) GetTxnHashFromCommitment(opts *bind.CallOpts, commitmentIndex [32]byte) (string, error) {
	var out []interface{}
	err := _Preconf.contract.Call(opts, &out, "getTxnHashFromCommitment", commitmentIndex)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// GetTxnHashFromCommitment is a free data retrieval call binding the contract method 0xfc4fbe32.
//
// Solidity: function getTxnHashFromCommitment(bytes32 commitmentIndex) view returns(string txnHash)
func (_Preconf *PreconfSession) GetTxnHashFromCommitment(commitmentIndex [32]byte) (string, error) {
	return _Preconf.Contract.GetTxnHashFromCommitment(&_Preconf.CallOpts, commitmentIndex)
}

// GetTxnHashFromCommitment is a free data retrieval call binding the contract method 0xfc4fbe32.
//
// Solidity: function getTxnHashFromCommitment(bytes32 commitmentIndex) view returns(string txnHash)
func (_Preconf *PreconfCallerSession) GetTxnHashFromCommitment(commitmentIndex [32]byte) (string, error) {
	return _Preconf.Contract.GetTxnHashFromCommitment(&_Preconf.CallOpts, commitmentIndex)
}

// LastProcessedBlock is a free data retrieval call binding the contract method 0x33de61d2.
//
// Solidity: function lastProcessedBlock() view returns(uint256)
func (_Preconf *PreconfCaller) LastProcessedBlock(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Preconf.contract.Call(opts, &out, "lastProcessedBlock")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LastProcessedBlock is a free data retrieval call binding the contract method 0x33de61d2.
//
// Solidity: function lastProcessedBlock() view returns(uint256)
func (_Preconf *PreconfSession) LastProcessedBlock() (*big.Int, error) {
	return _Preconf.Contract.LastProcessedBlock(&_Preconf.CallOpts)
}

// LastProcessedBlock is a free data retrieval call binding the contract method 0x33de61d2.
//
// Solidity: function lastProcessedBlock() view returns(uint256)
func (_Preconf *PreconfCallerSession) LastProcessedBlock() (*big.Int, error) {
	return _Preconf.Contract.LastProcessedBlock(&_Preconf.CallOpts)
}

// Oracle is a free data retrieval call binding the contract method 0x7dc0d1d0.
//
// Solidity: function oracle() view returns(address)
func (_Preconf *PreconfCaller) Oracle(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Preconf.contract.Call(opts, &out, "oracle")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Oracle is a free data retrieval call binding the contract method 0x7dc0d1d0.
//
// Solidity: function oracle() view returns(address)
func (_Preconf *PreconfSession) Oracle() (common.Address, error) {
	return _Preconf.Contract.Oracle(&_Preconf.CallOpts)
}

// Oracle is a free data retrieval call binding the contract method 0x7dc0d1d0.
//
// Solidity: function oracle() view returns(address)
func (_Preconf *PreconfCallerSession) Oracle() (common.Address, error) {
	return _Preconf.Contract.Oracle(&_Preconf.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Preconf *PreconfCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Preconf.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Preconf *PreconfSession) Owner() (common.Address, error) {
	return _Preconf.Contract.Owner(&_Preconf.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Preconf *PreconfCallerSession) Owner() (common.Address, error) {
	return _Preconf.Contract.Owner(&_Preconf.CallOpts)
}

// ProviderCommitments is a free data retrieval call binding the contract method 0x91b51cda.
//
// Solidity: function providerCommitments(address , uint256 ) view returns(bytes32)
func (_Preconf *PreconfCaller) ProviderCommitments(opts *bind.CallOpts, arg0 common.Address, arg1 *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _Preconf.contract.Call(opts, &out, "providerCommitments", arg0, arg1)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ProviderCommitments is a free data retrieval call binding the contract method 0x91b51cda.
//
// Solidity: function providerCommitments(address , uint256 ) view returns(bytes32)
func (_Preconf *PreconfSession) ProviderCommitments(arg0 common.Address, arg1 *big.Int) ([32]byte, error) {
	return _Preconf.Contract.ProviderCommitments(&_Preconf.CallOpts, arg0, arg1)
}

// ProviderCommitments is a free data retrieval call binding the contract method 0x91b51cda.
//
// Solidity: function providerCommitments(address , uint256 ) view returns(bytes32)
func (_Preconf *PreconfCallerSession) ProviderCommitments(arg0 common.Address, arg1 *big.Int) ([32]byte, error) {
	return _Preconf.Contract.ProviderCommitments(&_Preconf.CallOpts, arg0, arg1)
}

// ProviderRegistry is a free data retrieval call binding the contract method 0x545921d9.
//
// Solidity: function providerRegistry() view returns(address)
func (_Preconf *PreconfCaller) ProviderRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Preconf.contract.Call(opts, &out, "providerRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ProviderRegistry is a free data retrieval call binding the contract method 0x545921d9.
//
// Solidity: function providerRegistry() view returns(address)
func (_Preconf *PreconfSession) ProviderRegistry() (common.Address, error) {
	return _Preconf.Contract.ProviderRegistry(&_Preconf.CallOpts)
}

// ProviderRegistry is a free data retrieval call binding the contract method 0x545921d9.
//
// Solidity: function providerRegistry() view returns(address)
func (_Preconf *PreconfCallerSession) ProviderRegistry() (common.Address, error) {
	return _Preconf.Contract.ProviderRegistry(&_Preconf.CallOpts)
}

// VerifyBid is a free data retrieval call binding the contract method 0x5fd7f29b.
//
// Solidity: function verifyBid(uint64 bid, uint64 blockNumber, string txnHash, bytes bidSignature) view returns(bytes32 messageDigest, address recoveredAddress, uint256 stake)
func (_Preconf *PreconfCaller) VerifyBid(opts *bind.CallOpts, bid uint64, blockNumber uint64, txnHash string, bidSignature []byte) (struct {
	MessageDigest    [32]byte
	RecoveredAddress common.Address
	Stake            *big.Int
}, error) {
	var out []interface{}
	err := _Preconf.contract.Call(opts, &out, "verifyBid", bid, blockNumber, txnHash, bidSignature)

	outstruct := new(struct {
		MessageDigest    [32]byte
		RecoveredAddress common.Address
		Stake            *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.MessageDigest = *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	outstruct.RecoveredAddress = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)
	outstruct.Stake = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// VerifyBid is a free data retrieval call binding the contract method 0x5fd7f29b.
//
// Solidity: function verifyBid(uint64 bid, uint64 blockNumber, string txnHash, bytes bidSignature) view returns(bytes32 messageDigest, address recoveredAddress, uint256 stake)
func (_Preconf *PreconfSession) VerifyBid(bid uint64, blockNumber uint64, txnHash string, bidSignature []byte) (struct {
	MessageDigest    [32]byte
	RecoveredAddress common.Address
	Stake            *big.Int
}, error) {
	return _Preconf.Contract.VerifyBid(&_Preconf.CallOpts, bid, blockNumber, txnHash, bidSignature)
}

// VerifyBid is a free data retrieval call binding the contract method 0x5fd7f29b.
//
// Solidity: function verifyBid(uint64 bid, uint64 blockNumber, string txnHash, bytes bidSignature) view returns(bytes32 messageDigest, address recoveredAddress, uint256 stake)
func (_Preconf *PreconfCallerSession) VerifyBid(bid uint64, blockNumber uint64, txnHash string, bidSignature []byte) (struct {
	MessageDigest    [32]byte
	RecoveredAddress common.Address
	Stake            *big.Int
}, error) {
	return _Preconf.Contract.VerifyBid(&_Preconf.CallOpts, bid, blockNumber, txnHash, bidSignature)
}

// VerifyPreConfCommitment is a free data retrieval call binding the contract method 0x88e28d42.
//
// Solidity: function verifyPreConfCommitment(string txnHash, uint64 bid, uint64 blockNumber, bytes32 bidHash, bytes bidSignature, bytes commitmentSignature) view returns(bytes32 preConfHash, address commiterAddress)
func (_Preconf *PreconfCaller) VerifyPreConfCommitment(opts *bind.CallOpts, txnHash string, bid uint64, blockNumber uint64, bidHash [32]byte, bidSignature []byte, commitmentSignature []byte) (struct {
	PreConfHash     [32]byte
	CommiterAddress common.Address
}, error) {
	var out []interface{}
	err := _Preconf.contract.Call(opts, &out, "verifyPreConfCommitment", txnHash, bid, blockNumber, bidHash, bidSignature, commitmentSignature)

	outstruct := new(struct {
		PreConfHash     [32]byte
		CommiterAddress common.Address
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.PreConfHash = *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	outstruct.CommiterAddress = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)

	return *outstruct, err

}

// VerifyPreConfCommitment is a free data retrieval call binding the contract method 0x88e28d42.
//
// Solidity: function verifyPreConfCommitment(string txnHash, uint64 bid, uint64 blockNumber, bytes32 bidHash, bytes bidSignature, bytes commitmentSignature) view returns(bytes32 preConfHash, address commiterAddress)
func (_Preconf *PreconfSession) VerifyPreConfCommitment(txnHash string, bid uint64, blockNumber uint64, bidHash [32]byte, bidSignature []byte, commitmentSignature []byte) (struct {
	PreConfHash     [32]byte
	CommiterAddress common.Address
}, error) {
	return _Preconf.Contract.VerifyPreConfCommitment(&_Preconf.CallOpts, txnHash, bid, blockNumber, bidHash, bidSignature, commitmentSignature)
}

// VerifyPreConfCommitment is a free data retrieval call binding the contract method 0x88e28d42.
//
// Solidity: function verifyPreConfCommitment(string txnHash, uint64 bid, uint64 blockNumber, bytes32 bidHash, bytes bidSignature, bytes commitmentSignature) view returns(bytes32 preConfHash, address commiterAddress)
func (_Preconf *PreconfCallerSession) VerifyPreConfCommitment(txnHash string, bid uint64, blockNumber uint64, bidHash [32]byte, bidSignature []byte, commitmentSignature []byte) (struct {
	PreConfHash     [32]byte
	CommiterAddress common.Address
}, error) {
	return _Preconf.Contract.VerifyPreConfCommitment(&_Preconf.CallOpts, txnHash, bid, blockNumber, bidHash, bidSignature, commitmentSignature)
}

// InitateReward is a paid mutator transaction binding the contract method 0x78f7f847.
//
// Solidity: function initateReward(bytes32 commitmentIndex) returns()
func (_Preconf *PreconfTransactor) InitateReward(opts *bind.TransactOpts, commitmentIndex [32]byte) (*types.Transaction, error) {
	return _Preconf.contract.Transact(opts, "initateReward", commitmentIndex)
}

// InitateReward is a paid mutator transaction binding the contract method 0x78f7f847.
//
// Solidity: function initateReward(bytes32 commitmentIndex) returns()
func (_Preconf *PreconfSession) InitateReward(commitmentIndex [32]byte) (*types.Transaction, error) {
	return _Preconf.Contract.InitateReward(&_Preconf.TransactOpts, commitmentIndex)
}

// InitateReward is a paid mutator transaction binding the contract method 0x78f7f847.
//
// Solidity: function initateReward(bytes32 commitmentIndex) returns()
func (_Preconf *PreconfTransactorSession) InitateReward(commitmentIndex [32]byte) (*types.Transaction, error) {
	return _Preconf.Contract.InitateReward(&_Preconf.TransactOpts, commitmentIndex)
}

// InitiateSlash is a paid mutator transaction binding the contract method 0x6012a4f7.
//
// Solidity: function initiateSlash(bytes32 commitmentIndex) returns()
func (_Preconf *PreconfTransactor) InitiateSlash(opts *bind.TransactOpts, commitmentIndex [32]byte) (*types.Transaction, error) {
	return _Preconf.contract.Transact(opts, "initiateSlash", commitmentIndex)
}

// InitiateSlash is a paid mutator transaction binding the contract method 0x6012a4f7.
//
// Solidity: function initiateSlash(bytes32 commitmentIndex) returns()
func (_Preconf *PreconfSession) InitiateSlash(commitmentIndex [32]byte) (*types.Transaction, error) {
	return _Preconf.Contract.InitiateSlash(&_Preconf.TransactOpts, commitmentIndex)
}

// InitiateSlash is a paid mutator transaction binding the contract method 0x6012a4f7.
//
// Solidity: function initiateSlash(bytes32 commitmentIndex) returns()
func (_Preconf *PreconfTransactorSession) InitiateSlash(commitmentIndex [32]byte) (*types.Transaction, error) {
	return _Preconf.Contract.InitiateSlash(&_Preconf.TransactOpts, commitmentIndex)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Preconf *PreconfTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Preconf.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Preconf *PreconfSession) RenounceOwnership() (*types.Transaction, error) {
	return _Preconf.Contract.RenounceOwnership(&_Preconf.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Preconf *PreconfTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Preconf.Contract.RenounceOwnership(&_Preconf.TransactOpts)
}

// StoreCommitment is a paid mutator transaction binding the contract method 0x8c43bc7e.
//
// Solidity: function storeCommitment(uint64 bid, uint64 blockNumber, string txnHash, bytes bidSignature, bytes commitmentSignature) returns(bytes32 commitmentIndex)
func (_Preconf *PreconfTransactor) StoreCommitment(opts *bind.TransactOpts, bid uint64, blockNumber uint64, txnHash string, bidSignature []byte, commitmentSignature []byte) (*types.Transaction, error) {
	return _Preconf.contract.Transact(opts, "storeCommitment", bid, blockNumber, txnHash, bidSignature, commitmentSignature)
}

// StoreCommitment is a paid mutator transaction binding the contract method 0x8c43bc7e.
//
// Solidity: function storeCommitment(uint64 bid, uint64 blockNumber, string txnHash, bytes bidSignature, bytes commitmentSignature) returns(bytes32 commitmentIndex)
func (_Preconf *PreconfSession) StoreCommitment(bid uint64, blockNumber uint64, txnHash string, bidSignature []byte, commitmentSignature []byte) (*types.Transaction, error) {
	return _Preconf.Contract.StoreCommitment(&_Preconf.TransactOpts, bid, blockNumber, txnHash, bidSignature, commitmentSignature)
}

// StoreCommitment is a paid mutator transaction binding the contract method 0x8c43bc7e.
//
// Solidity: function storeCommitment(uint64 bid, uint64 blockNumber, string txnHash, bytes bidSignature, bytes commitmentSignature) returns(bytes32 commitmentIndex)
func (_Preconf *PreconfTransactorSession) StoreCommitment(bid uint64, blockNumber uint64, txnHash string, bidSignature []byte, commitmentSignature []byte) (*types.Transaction, error) {
	return _Preconf.Contract.StoreCommitment(&_Preconf.TransactOpts, bid, blockNumber, txnHash, bidSignature, commitmentSignature)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Preconf *PreconfTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Preconf.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Preconf *PreconfSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Preconf.Contract.TransferOwnership(&_Preconf.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Preconf *PreconfTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Preconf.Contract.TransferOwnership(&_Preconf.TransactOpts, newOwner)
}

// UpdateBidderRegistry is a paid mutator transaction binding the contract method 0x66544c41.
//
// Solidity: function updateBidderRegistry(address newBidderRegistry) returns()
func (_Preconf *PreconfTransactor) UpdateBidderRegistry(opts *bind.TransactOpts, newBidderRegistry common.Address) (*types.Transaction, error) {
	return _Preconf.contract.Transact(opts, "updateBidderRegistry", newBidderRegistry)
}

// UpdateBidderRegistry is a paid mutator transaction binding the contract method 0x66544c41.
//
// Solidity: function updateBidderRegistry(address newBidderRegistry) returns()
func (_Preconf *PreconfSession) UpdateBidderRegistry(newBidderRegistry common.Address) (*types.Transaction, error) {
	return _Preconf.Contract.UpdateBidderRegistry(&_Preconf.TransactOpts, newBidderRegistry)
}

// UpdateBidderRegistry is a paid mutator transaction binding the contract method 0x66544c41.
//
// Solidity: function updateBidderRegistry(address newBidderRegistry) returns()
func (_Preconf *PreconfTransactorSession) UpdateBidderRegistry(newBidderRegistry common.Address) (*types.Transaction, error) {
	return _Preconf.Contract.UpdateBidderRegistry(&_Preconf.TransactOpts, newBidderRegistry)
}

// UpdateOracle is a paid mutator transaction binding the contract method 0x1cb44dfc.
//
// Solidity: function updateOracle(address newOracle) returns()
func (_Preconf *PreconfTransactor) UpdateOracle(opts *bind.TransactOpts, newOracle common.Address) (*types.Transaction, error) {
	return _Preconf.contract.Transact(opts, "updateOracle", newOracle)
}

// UpdateOracle is a paid mutator transaction binding the contract method 0x1cb44dfc.
//
// Solidity: function updateOracle(address newOracle) returns()
func (_Preconf *PreconfSession) UpdateOracle(newOracle common.Address) (*types.Transaction, error) {
	return _Preconf.Contract.UpdateOracle(&_Preconf.TransactOpts, newOracle)
}

// UpdateOracle is a paid mutator transaction binding the contract method 0x1cb44dfc.
//
// Solidity: function updateOracle(address newOracle) returns()
func (_Preconf *PreconfTransactorSession) UpdateOracle(newOracle common.Address) (*types.Transaction, error) {
	return _Preconf.Contract.UpdateOracle(&_Preconf.TransactOpts, newOracle)
}

// UpdateProviderRegistry is a paid mutator transaction binding the contract method 0x92d2e3e7.
//
// Solidity: function updateProviderRegistry(address newProviderRegistry) returns()
func (_Preconf *PreconfTransactor) UpdateProviderRegistry(opts *bind.TransactOpts, newProviderRegistry common.Address) (*types.Transaction, error) {
	return _Preconf.contract.Transact(opts, "updateProviderRegistry", newProviderRegistry)
}

// UpdateProviderRegistry is a paid mutator transaction binding the contract method 0x92d2e3e7.
//
// Solidity: function updateProviderRegistry(address newProviderRegistry) returns()
func (_Preconf *PreconfSession) UpdateProviderRegistry(newProviderRegistry common.Address) (*types.Transaction, error) {
	return _Preconf.Contract.UpdateProviderRegistry(&_Preconf.TransactOpts, newProviderRegistry)
}

// UpdateProviderRegistry is a paid mutator transaction binding the contract method 0x92d2e3e7.
//
// Solidity: function updateProviderRegistry(address newProviderRegistry) returns()
func (_Preconf *PreconfTransactorSession) UpdateProviderRegistry(newProviderRegistry common.Address) (*types.Transaction, error) {
	return _Preconf.Contract.UpdateProviderRegistry(&_Preconf.TransactOpts, newProviderRegistry)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Preconf *PreconfTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _Preconf.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Preconf *PreconfSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _Preconf.Contract.Fallback(&_Preconf.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Preconf *PreconfTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _Preconf.Contract.Fallback(&_Preconf.TransactOpts, calldata)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Preconf *PreconfTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Preconf.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Preconf *PreconfSession) Receive() (*types.Transaction, error) {
	return _Preconf.Contract.Receive(&_Preconf.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Preconf *PreconfTransactorSession) Receive() (*types.Transaction, error) {
	return _Preconf.Contract.Receive(&_Preconf.TransactOpts)
}

// PreconfOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Preconf contract.
type PreconfOwnershipTransferredIterator struct {
	Event *PreconfOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PreconfOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PreconfOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PreconfOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PreconfOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PreconfOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PreconfOwnershipTransferred represents a OwnershipTransferred event raised by the Preconf contract.
type PreconfOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Preconf *PreconfFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*PreconfOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Preconf.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &PreconfOwnershipTransferredIterator{contract: _Preconf.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Preconf *PreconfFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *PreconfOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Preconf.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PreconfOwnershipTransferred)
				if err := _Preconf.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Preconf *PreconfFilterer) ParseOwnershipTransferred(log types.Log) (*PreconfOwnershipTransferred, error) {
	event := new(PreconfOwnershipTransferred)
	if err := _Preconf.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PreconfSignatureVerifiedIterator is returned from FilterSignatureVerified and is used to iterate over the raw logs and unpacked data for SignatureVerified events raised by the Preconf contract.
type PreconfSignatureVerifiedIterator struct {
	Event *PreconfSignatureVerified // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PreconfSignatureVerifiedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PreconfSignatureVerified)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PreconfSignatureVerified)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PreconfSignatureVerifiedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PreconfSignatureVerifiedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PreconfSignatureVerified represents a SignatureVerified event raised by the Preconf contract.
type PreconfSignatureVerified struct {
	Signer      common.Address
	TxnHash     string
	Bid         uint64
	BlockNumber uint64
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterSignatureVerified is a free log retrieval operation binding the contract event 0x48db0394d84b81a6f3cb6c61ea2dceff3cad797a9b889fe499fc051f08969c4d.
//
// Solidity: event SignatureVerified(address indexed signer, string txnHash, uint64 indexed bid, uint64 blockNumber)
func (_Preconf *PreconfFilterer) FilterSignatureVerified(opts *bind.FilterOpts, signer []common.Address, bid []uint64) (*PreconfSignatureVerifiedIterator, error) {

	var signerRule []interface{}
	for _, signerItem := range signer {
		signerRule = append(signerRule, signerItem)
	}

	var bidRule []interface{}
	for _, bidItem := range bid {
		bidRule = append(bidRule, bidItem)
	}

	logs, sub, err := _Preconf.contract.FilterLogs(opts, "SignatureVerified", signerRule, bidRule)
	if err != nil {
		return nil, err
	}
	return &PreconfSignatureVerifiedIterator{contract: _Preconf.contract, event: "SignatureVerified", logs: logs, sub: sub}, nil
}

// WatchSignatureVerified is a free log subscription operation binding the contract event 0x48db0394d84b81a6f3cb6c61ea2dceff3cad797a9b889fe499fc051f08969c4d.
//
// Solidity: event SignatureVerified(address indexed signer, string txnHash, uint64 indexed bid, uint64 blockNumber)
func (_Preconf *PreconfFilterer) WatchSignatureVerified(opts *bind.WatchOpts, sink chan<- *PreconfSignatureVerified, signer []common.Address, bid []uint64) (event.Subscription, error) {

	var signerRule []interface{}
	for _, signerItem := range signer {
		signerRule = append(signerRule, signerItem)
	}

	var bidRule []interface{}
	for _, bidItem := range bid {
		bidRule = append(bidRule, bidItem)
	}

	logs, sub, err := _Preconf.contract.WatchLogs(opts, "SignatureVerified", signerRule, bidRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PreconfSignatureVerified)
				if err := _Preconf.contract.UnpackLog(event, "SignatureVerified", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseSignatureVerified is a log parse operation binding the contract event 0x48db0394d84b81a6f3cb6c61ea2dceff3cad797a9b889fe499fc051f08969c4d.
//
// Solidity: event SignatureVerified(address indexed signer, string txnHash, uint64 indexed bid, uint64 blockNumber)
func (_Preconf *PreconfFilterer) ParseSignatureVerified(log types.Log) (*PreconfSignatureVerified, error) {
	event := new(PreconfSignatureVerified)
	if err := _Preconf.contract.UnpackLog(event, "SignatureVerified", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
