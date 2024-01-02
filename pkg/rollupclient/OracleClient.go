// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package rollupclient

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

// OracleClientMetaData contains all meta data concerning the OracleClient contract.
var OracleClientMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_preConfContract\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_nextRequestedBlockNumber\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_owner\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"fallback\",\"stateMutability\":\"payable\"},{\"type\":\"receive\",\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"addBuilderAddress\",\"inputs\":[{\"name\":\"builderName\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"builderAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"blockBuilderNameToAddress\",\"inputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getBuilder\",\"inputs\":[{\"name\":\"builderNameGrafiti\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getNextRequestedBlockNumber\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"moveToNextBlock\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"nextRequestedBlockNumber\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"processBuilderCommitmentForBlockNumber\",\"inputs\":[{\"name\":\"commitmentIndex\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"blockNumber\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"blockBuilderName\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"isSlash\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setNextBlock\",\"inputs\":[{\"name\":\"newBlockNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"CommitmentProcessed\",\"inputs\":[{\"name\":\"commitmentHash\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"isSlash\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false}]",
}

// OracleClientABI is the input ABI used to generate the binding from.
// Deprecated: Use OracleClientMetaData.ABI instead.
var OracleClientABI = OracleClientMetaData.ABI

// OracleClient is an auto generated Go binding around an Ethereum contract.
type OracleClient struct {
	OracleClientCaller     // Read-only binding to the contract
	OracleClientTransactor // Write-only binding to the contract
	OracleClientFilterer   // Log filterer for contract events
}

// OracleClientCaller is an auto generated read-only Go binding around an Ethereum contract.
type OracleClientCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OracleClientTransactor is an auto generated write-only Go binding around an Ethereum contract.
type OracleClientTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OracleClientFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type OracleClientFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OracleClientSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type OracleClientSession struct {
	Contract     *OracleClient     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// OracleClientCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type OracleClientCallerSession struct {
	Contract *OracleClientCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// OracleClientTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type OracleClientTransactorSession struct {
	Contract     *OracleClientTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// OracleClientRaw is an auto generated low-level Go binding around an Ethereum contract.
type OracleClientRaw struct {
	Contract *OracleClient // Generic contract binding to access the raw methods on
}

// OracleClientCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type OracleClientCallerRaw struct {
	Contract *OracleClientCaller // Generic read-only contract binding to access the raw methods on
}

// OracleClientTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type OracleClientTransactorRaw struct {
	Contract *OracleClientTransactor // Generic write-only contract binding to access the raw methods on
}

// NewOracleClient creates a new instance of OracleClient, bound to a specific deployed contract.
func NewOracleClient(address common.Address, backend bind.ContractBackend) (*OracleClient, error) {
	contract, err := bindOracleClient(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &OracleClient{OracleClientCaller: OracleClientCaller{contract: contract}, OracleClientTransactor: OracleClientTransactor{contract: contract}, OracleClientFilterer: OracleClientFilterer{contract: contract}}, nil
}

// NewOracleClientCaller creates a new read-only instance of OracleClient, bound to a specific deployed contract.
func NewOracleClientCaller(address common.Address, caller bind.ContractCaller) (*OracleClientCaller, error) {
	contract, err := bindOracleClient(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OracleClientCaller{contract: contract}, nil
}

// NewOracleClientTransactor creates a new write-only instance of OracleClient, bound to a specific deployed contract.
func NewOracleClientTransactor(address common.Address, transactor bind.ContractTransactor) (*OracleClientTransactor, error) {
	contract, err := bindOracleClient(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OracleClientTransactor{contract: contract}, nil
}

// NewOracleClientFilterer creates a new log filterer instance of OracleClient, bound to a specific deployed contract.
func NewOracleClientFilterer(address common.Address, filterer bind.ContractFilterer) (*OracleClientFilterer, error) {
	contract, err := bindOracleClient(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OracleClientFilterer{contract: contract}, nil
}

// bindOracleClient binds a generic wrapper to an already deployed contract.
func bindOracleClient(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := OracleClientMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OracleClient *OracleClientRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OracleClient.Contract.OracleClientCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OracleClient *OracleClientRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OracleClient.Contract.OracleClientTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OracleClient *OracleClientRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OracleClient.Contract.OracleClientTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OracleClient *OracleClientCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OracleClient.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OracleClient *OracleClientTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OracleClient.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OracleClient *OracleClientTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OracleClient.Contract.contract.Transact(opts, method, params...)
}

// BlockBuilderNameToAddress is a free data retrieval call binding the contract method 0xeebac3ac.
//
// Solidity: function blockBuilderNameToAddress(string ) view returns(address)
func (_OracleClient *OracleClientCaller) BlockBuilderNameToAddress(opts *bind.CallOpts, arg0 string) (common.Address, error) {
	var out []interface{}
	err := _OracleClient.contract.Call(opts, &out, "blockBuilderNameToAddress", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// BlockBuilderNameToAddress is a free data retrieval call binding the contract method 0xeebac3ac.
//
// Solidity: function blockBuilderNameToAddress(string ) view returns(address)
func (_OracleClient *OracleClientSession) BlockBuilderNameToAddress(arg0 string) (common.Address, error) {
	return _OracleClient.Contract.BlockBuilderNameToAddress(&_OracleClient.CallOpts, arg0)
}

// BlockBuilderNameToAddress is a free data retrieval call binding the contract method 0xeebac3ac.
//
// Solidity: function blockBuilderNameToAddress(string ) view returns(address)
func (_OracleClient *OracleClientCallerSession) BlockBuilderNameToAddress(arg0 string) (common.Address, error) {
	return _OracleClient.Contract.BlockBuilderNameToAddress(&_OracleClient.CallOpts, arg0)
}

// GetBuilder is a free data retrieval call binding the contract method 0x237ba8fb.
//
// Solidity: function getBuilder(string builderNameGrafiti) view returns(address)
func (_OracleClient *OracleClientCaller) GetBuilder(opts *bind.CallOpts, builderNameGrafiti string) (common.Address, error) {
	var out []interface{}
	err := _OracleClient.contract.Call(opts, &out, "getBuilder", builderNameGrafiti)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetBuilder is a free data retrieval call binding the contract method 0x237ba8fb.
//
// Solidity: function getBuilder(string builderNameGrafiti) view returns(address)
func (_OracleClient *OracleClientSession) GetBuilder(builderNameGrafiti string) (common.Address, error) {
	return _OracleClient.Contract.GetBuilder(&_OracleClient.CallOpts, builderNameGrafiti)
}

// GetBuilder is a free data retrieval call binding the contract method 0x237ba8fb.
//
// Solidity: function getBuilder(string builderNameGrafiti) view returns(address)
func (_OracleClient *OracleClientCallerSession) GetBuilder(builderNameGrafiti string) (common.Address, error) {
	return _OracleClient.Contract.GetBuilder(&_OracleClient.CallOpts, builderNameGrafiti)
}

// GetNextRequestedBlockNumber is a free data retrieval call binding the contract method 0xfce2a502.
//
// Solidity: function getNextRequestedBlockNumber() view returns(uint256)
func (_OracleClient *OracleClientCaller) GetNextRequestedBlockNumber(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _OracleClient.contract.Call(opts, &out, "getNextRequestedBlockNumber")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetNextRequestedBlockNumber is a free data retrieval call binding the contract method 0xfce2a502.
//
// Solidity: function getNextRequestedBlockNumber() view returns(uint256)
func (_OracleClient *OracleClientSession) GetNextRequestedBlockNumber() (*big.Int, error) {
	return _OracleClient.Contract.GetNextRequestedBlockNumber(&_OracleClient.CallOpts)
}

// GetNextRequestedBlockNumber is a free data retrieval call binding the contract method 0xfce2a502.
//
// Solidity: function getNextRequestedBlockNumber() view returns(uint256)
func (_OracleClient *OracleClientCallerSession) GetNextRequestedBlockNumber() (*big.Int, error) {
	return _OracleClient.Contract.GetNextRequestedBlockNumber(&_OracleClient.CallOpts)
}

// NextRequestedBlockNumber is a free data retrieval call binding the contract method 0xc512c561.
//
// Solidity: function nextRequestedBlockNumber() view returns(uint256)
func (_OracleClient *OracleClientCaller) NextRequestedBlockNumber(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _OracleClient.contract.Call(opts, &out, "nextRequestedBlockNumber")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// NextRequestedBlockNumber is a free data retrieval call binding the contract method 0xc512c561.
//
// Solidity: function nextRequestedBlockNumber() view returns(uint256)
func (_OracleClient *OracleClientSession) NextRequestedBlockNumber() (*big.Int, error) {
	return _OracleClient.Contract.NextRequestedBlockNumber(&_OracleClient.CallOpts)
}

// NextRequestedBlockNumber is a free data retrieval call binding the contract method 0xc512c561.
//
// Solidity: function nextRequestedBlockNumber() view returns(uint256)
func (_OracleClient *OracleClientCallerSession) NextRequestedBlockNumber() (*big.Int, error) {
	return _OracleClient.Contract.NextRequestedBlockNumber(&_OracleClient.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_OracleClient *OracleClientCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _OracleClient.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_OracleClient *OracleClientSession) Owner() (common.Address, error) {
	return _OracleClient.Contract.Owner(&_OracleClient.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_OracleClient *OracleClientCallerSession) Owner() (common.Address, error) {
	return _OracleClient.Contract.Owner(&_OracleClient.CallOpts)
}

// AddBuilderAddress is a paid mutator transaction binding the contract method 0x0bd0a9e1.
//
// Solidity: function addBuilderAddress(string builderName, address builderAddress) returns()
func (_OracleClient *OracleClientTransactor) AddBuilderAddress(opts *bind.TransactOpts, builderName string, builderAddress common.Address) (*types.Transaction, error) {
	return _OracleClient.contract.Transact(opts, "addBuilderAddress", builderName, builderAddress)
}

// AddBuilderAddress is a paid mutator transaction binding the contract method 0x0bd0a9e1.
//
// Solidity: function addBuilderAddress(string builderName, address builderAddress) returns()
func (_OracleClient *OracleClientSession) AddBuilderAddress(builderName string, builderAddress common.Address) (*types.Transaction, error) {
	return _OracleClient.Contract.AddBuilderAddress(&_OracleClient.TransactOpts, builderName, builderAddress)
}

// AddBuilderAddress is a paid mutator transaction binding the contract method 0x0bd0a9e1.
//
// Solidity: function addBuilderAddress(string builderName, address builderAddress) returns()
func (_OracleClient *OracleClientTransactorSession) AddBuilderAddress(builderName string, builderAddress common.Address) (*types.Transaction, error) {
	return _OracleClient.Contract.AddBuilderAddress(&_OracleClient.TransactOpts, builderName, builderAddress)
}

// MoveToNextBlock is a paid mutator transaction binding the contract method 0x32289b72.
//
// Solidity: function moveToNextBlock() returns()
func (_OracleClient *OracleClientTransactor) MoveToNextBlock(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OracleClient.contract.Transact(opts, "moveToNextBlock")
}

// MoveToNextBlock is a paid mutator transaction binding the contract method 0x32289b72.
//
// Solidity: function moveToNextBlock() returns()
func (_OracleClient *OracleClientSession) MoveToNextBlock() (*types.Transaction, error) {
	return _OracleClient.Contract.MoveToNextBlock(&_OracleClient.TransactOpts)
}

// MoveToNextBlock is a paid mutator transaction binding the contract method 0x32289b72.
//
// Solidity: function moveToNextBlock() returns()
func (_OracleClient *OracleClientTransactorSession) MoveToNextBlock() (*types.Transaction, error) {
	return _OracleClient.Contract.MoveToNextBlock(&_OracleClient.TransactOpts)
}

// ProcessBuilderCommitmentForBlockNumber is a paid mutator transaction binding the contract method 0x04a24484.
//
// Solidity: function processBuilderCommitmentForBlockNumber(bytes32 commitmentIndex, uint256 blockNumber, string blockBuilderName, bool isSlash) returns()
func (_OracleClient *OracleClientTransactor) ProcessBuilderCommitmentForBlockNumber(opts *bind.TransactOpts, commitmentIndex [32]byte, blockNumber *big.Int, blockBuilderName string, isSlash bool) (*types.Transaction, error) {
	return _OracleClient.contract.Transact(opts, "processBuilderCommitmentForBlockNumber", commitmentIndex, blockNumber, blockBuilderName, isSlash)
}

// ProcessBuilderCommitmentForBlockNumber is a paid mutator transaction binding the contract method 0x04a24484.
//
// Solidity: function processBuilderCommitmentForBlockNumber(bytes32 commitmentIndex, uint256 blockNumber, string blockBuilderName, bool isSlash) returns()
func (_OracleClient *OracleClientSession) ProcessBuilderCommitmentForBlockNumber(commitmentIndex [32]byte, blockNumber *big.Int, blockBuilderName string, isSlash bool) (*types.Transaction, error) {
	return _OracleClient.Contract.ProcessBuilderCommitmentForBlockNumber(&_OracleClient.TransactOpts, commitmentIndex, blockNumber, blockBuilderName, isSlash)
}

// ProcessBuilderCommitmentForBlockNumber is a paid mutator transaction binding the contract method 0x04a24484.
//
// Solidity: function processBuilderCommitmentForBlockNumber(bytes32 commitmentIndex, uint256 blockNumber, string blockBuilderName, bool isSlash) returns()
func (_OracleClient *OracleClientTransactorSession) ProcessBuilderCommitmentForBlockNumber(commitmentIndex [32]byte, blockNumber *big.Int, blockBuilderName string, isSlash bool) (*types.Transaction, error) {
	return _OracleClient.Contract.ProcessBuilderCommitmentForBlockNumber(&_OracleClient.TransactOpts, commitmentIndex, blockNumber, blockBuilderName, isSlash)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_OracleClient *OracleClientTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OracleClient.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_OracleClient *OracleClientSession) RenounceOwnership() (*types.Transaction, error) {
	return _OracleClient.Contract.RenounceOwnership(&_OracleClient.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_OracleClient *OracleClientTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _OracleClient.Contract.RenounceOwnership(&_OracleClient.TransactOpts)
}

// SetNextBlock is a paid mutator transaction binding the contract method 0x19072add.
//
// Solidity: function setNextBlock(uint64 newBlockNumber) returns()
func (_OracleClient *OracleClientTransactor) SetNextBlock(opts *bind.TransactOpts, newBlockNumber uint64) (*types.Transaction, error) {
	return _OracleClient.contract.Transact(opts, "setNextBlock", newBlockNumber)
}

// SetNextBlock is a paid mutator transaction binding the contract method 0x19072add.
//
// Solidity: function setNextBlock(uint64 newBlockNumber) returns()
func (_OracleClient *OracleClientSession) SetNextBlock(newBlockNumber uint64) (*types.Transaction, error) {
	return _OracleClient.Contract.SetNextBlock(&_OracleClient.TransactOpts, newBlockNumber)
}

// SetNextBlock is a paid mutator transaction binding the contract method 0x19072add.
//
// Solidity: function setNextBlock(uint64 newBlockNumber) returns()
func (_OracleClient *OracleClientTransactorSession) SetNextBlock(newBlockNumber uint64) (*types.Transaction, error) {
	return _OracleClient.Contract.SetNextBlock(&_OracleClient.TransactOpts, newBlockNumber)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_OracleClient *OracleClientTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _OracleClient.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_OracleClient *OracleClientSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _OracleClient.Contract.TransferOwnership(&_OracleClient.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_OracleClient *OracleClientTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _OracleClient.Contract.TransferOwnership(&_OracleClient.TransactOpts, newOwner)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_OracleClient *OracleClientTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _OracleClient.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_OracleClient *OracleClientSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _OracleClient.Contract.Fallback(&_OracleClient.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_OracleClient *OracleClientTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _OracleClient.Contract.Fallback(&_OracleClient.TransactOpts, calldata)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_OracleClient *OracleClientTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OracleClient.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_OracleClient *OracleClientSession) Receive() (*types.Transaction, error) {
	return _OracleClient.Contract.Receive(&_OracleClient.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_OracleClient *OracleClientTransactorSession) Receive() (*types.Transaction, error) {
	return _OracleClient.Contract.Receive(&_OracleClient.TransactOpts)
}

// OracleClientCommitmentProcessedIterator is returned from FilterCommitmentProcessed and is used to iterate over the raw logs and unpacked data for CommitmentProcessed events raised by the OracleClient contract.
type OracleClientCommitmentProcessedIterator struct {
	Event *OracleClientCommitmentProcessed // Event containing the contract specifics and raw log

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
func (it *OracleClientCommitmentProcessedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OracleClientCommitmentProcessed)
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
		it.Event = new(OracleClientCommitmentProcessed)
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
func (it *OracleClientCommitmentProcessedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OracleClientCommitmentProcessedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OracleClientCommitmentProcessed represents a CommitmentProcessed event raised by the OracleClient contract.
type OracleClientCommitmentProcessed struct {
	CommitmentHash [32]byte
	IsSlash        bool
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterCommitmentProcessed is a free log retrieval operation binding the contract event 0xddc1768a3a762a04e5fd3abea8ae3b60e23bcf290f4a032280e6a726611d41f5.
//
// Solidity: event CommitmentProcessed(bytes32 commitmentHash, bool isSlash)
func (_OracleClient *OracleClientFilterer) FilterCommitmentProcessed(opts *bind.FilterOpts) (*OracleClientCommitmentProcessedIterator, error) {

	logs, sub, err := _OracleClient.contract.FilterLogs(opts, "CommitmentProcessed")
	if err != nil {
		return nil, err
	}
	return &OracleClientCommitmentProcessedIterator{contract: _OracleClient.contract, event: "CommitmentProcessed", logs: logs, sub: sub}, nil
}

// WatchCommitmentProcessed is a free log subscription operation binding the contract event 0xddc1768a3a762a04e5fd3abea8ae3b60e23bcf290f4a032280e6a726611d41f5.
//
// Solidity: event CommitmentProcessed(bytes32 commitmentHash, bool isSlash)
func (_OracleClient *OracleClientFilterer) WatchCommitmentProcessed(opts *bind.WatchOpts, sink chan<- *OracleClientCommitmentProcessed) (event.Subscription, error) {

	logs, sub, err := _OracleClient.contract.WatchLogs(opts, "CommitmentProcessed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OracleClientCommitmentProcessed)
				if err := _OracleClient.contract.UnpackLog(event, "CommitmentProcessed", log); err != nil {
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

// ParseCommitmentProcessed is a log parse operation binding the contract event 0xddc1768a3a762a04e5fd3abea8ae3b60e23bcf290f4a032280e6a726611d41f5.
//
// Solidity: event CommitmentProcessed(bytes32 commitmentHash, bool isSlash)
func (_OracleClient *OracleClientFilterer) ParseCommitmentProcessed(log types.Log) (*OracleClientCommitmentProcessed, error) {
	event := new(OracleClientCommitmentProcessed)
	if err := _OracleClient.contract.UnpackLog(event, "CommitmentProcessed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OracleClientOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the OracleClient contract.
type OracleClientOwnershipTransferredIterator struct {
	Event *OracleClientOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *OracleClientOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OracleClientOwnershipTransferred)
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
		it.Event = new(OracleClientOwnershipTransferred)
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
func (it *OracleClientOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OracleClientOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OracleClientOwnershipTransferred represents a OwnershipTransferred event raised by the OracleClient contract.
type OracleClientOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_OracleClient *OracleClientFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*OracleClientOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _OracleClient.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &OracleClientOwnershipTransferredIterator{contract: _OracleClient.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_OracleClient *OracleClientFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *OracleClientOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _OracleClient.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OracleClientOwnershipTransferred)
				if err := _OracleClient.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_OracleClient *OracleClientFilterer) ParseOwnershipTransferred(log types.Log) (*OracleClientOwnershipTransferred, error) {
	event := new(OracleClientOwnershipTransferred)
	if err := _OracleClient.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
