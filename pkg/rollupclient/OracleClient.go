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

// RollupclientMetaData contains all meta data concerning the Rollupclient contract.
var RollupclientMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_preConfContract\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string[]\",\"name\":\"txnList\",\"type\":\"string[]\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"blockBuilderName\",\"type\":\"string\"}],\"name\":\"BlockDataReceived\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"}],\"name\":\"BlockDataRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"commitmentHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"isSlash\",\"type\":\"bool\"}],\"name\":\"CommitmentProcessed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"builderName\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"builderAddress\",\"type\":\"address\"}],\"name\":\"addBuilderAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"name\":\"blockBuilderNameToAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getNextRequestedBlockNumber\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"nextRequestedBlockNumber\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"commitmentIndex\",\"type\":\"bytes32\"},{\"internalType\":\"bool\",\"name\":\"isSlash\",\"type\":\"bool\"}],\"name\":\"processCommitment\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string[]\",\"name\":\"txnList\",\"type\":\"string[]\"},{\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"blockBuilderName\",\"type\":\"string\"}],\"name\":\"receiveBlockData\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"}],\"name\":\"requestBlockData\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
}

// RollupclientABI is the input ABI used to generate the binding from.
// Deprecated: Use RollupclientMetaData.ABI instead.
var RollupclientABI = RollupclientMetaData.ABI

// Rollupclient is an auto generated Go binding around an Ethereum contract.
type Rollupclient struct {
	RollupclientCaller     // Read-only binding to the contract
	RollupclientTransactor // Write-only binding to the contract
	RollupclientFilterer   // Log filterer for contract events
}

// RollupclientCaller is an auto generated read-only Go binding around an Ethereum contract.
type RollupclientCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RollupclientTransactor is an auto generated write-only Go binding around an Ethereum contract.
type RollupclientTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RollupclientFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type RollupclientFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RollupclientSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type RollupclientSession struct {
	Contract     *Rollupclient     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// RollupclientCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type RollupclientCallerSession struct {
	Contract *RollupclientCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// RollupclientTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type RollupclientTransactorSession struct {
	Contract     *RollupclientTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// RollupclientRaw is an auto generated low-level Go binding around an Ethereum contract.
type RollupclientRaw struct {
	Contract *Rollupclient // Generic contract binding to access the raw methods on
}

// RollupclientCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type RollupclientCallerRaw struct {
	Contract *RollupclientCaller // Generic read-only contract binding to access the raw methods on
}

// RollupclientTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type RollupclientTransactorRaw struct {
	Contract *RollupclientTransactor // Generic write-only contract binding to access the raw methods on
}

// NewRollupclient creates a new instance of Rollupclient, bound to a specific deployed contract.
func NewRollupclient(address common.Address, backend bind.ContractBackend) (*Rollupclient, error) {
	contract, err := bindRollupclient(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Rollupclient{RollupclientCaller: RollupclientCaller{contract: contract}, RollupclientTransactor: RollupclientTransactor{contract: contract}, RollupclientFilterer: RollupclientFilterer{contract: contract}}, nil
}

// NewRollupclientCaller creates a new read-only instance of Rollupclient, bound to a specific deployed contract.
func NewRollupclientCaller(address common.Address, caller bind.ContractCaller) (*RollupclientCaller, error) {
	contract, err := bindRollupclient(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &RollupclientCaller{contract: contract}, nil
}

// NewRollupclientTransactor creates a new write-only instance of Rollupclient, bound to a specific deployed contract.
func NewRollupclientTransactor(address common.Address, transactor bind.ContractTransactor) (*RollupclientTransactor, error) {
	contract, err := bindRollupclient(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &RollupclientTransactor{contract: contract}, nil
}

// NewRollupclientFilterer creates a new log filterer instance of Rollupclient, bound to a specific deployed contract.
func NewRollupclientFilterer(address common.Address, filterer bind.ContractFilterer) (*RollupclientFilterer, error) {
	contract, err := bindRollupclient(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &RollupclientFilterer{contract: contract}, nil
}

// bindRollupclient binds a generic wrapper to an already deployed contract.
func bindRollupclient(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := RollupclientMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Rollupclient *RollupclientRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Rollupclient.Contract.RollupclientCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Rollupclient *RollupclientRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Rollupclient.Contract.RollupclientTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Rollupclient *RollupclientRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Rollupclient.Contract.RollupclientTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Rollupclient *RollupclientCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Rollupclient.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Rollupclient *RollupclientTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Rollupclient.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Rollupclient *RollupclientTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Rollupclient.Contract.contract.Transact(opts, method, params...)
}

// BlockBuilderNameToAddress is a free data retrieval call binding the contract method 0xeebac3ac.
//
// Solidity: function blockBuilderNameToAddress(string ) view returns(address)
func (_Rollupclient *RollupclientCaller) BlockBuilderNameToAddress(opts *bind.CallOpts, arg0 string) (common.Address, error) {
	var out []interface{}
	err := _Rollupclient.contract.Call(opts, &out, "blockBuilderNameToAddress", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// BlockBuilderNameToAddress is a free data retrieval call binding the contract method 0xeebac3ac.
//
// Solidity: function blockBuilderNameToAddress(string ) view returns(address)
func (_Rollupclient *RollupclientSession) BlockBuilderNameToAddress(arg0 string) (common.Address, error) {
	return _Rollupclient.Contract.BlockBuilderNameToAddress(&_Rollupclient.CallOpts, arg0)
}

// BlockBuilderNameToAddress is a free data retrieval call binding the contract method 0xeebac3ac.
//
// Solidity: function blockBuilderNameToAddress(string ) view returns(address)
func (_Rollupclient *RollupclientCallerSession) BlockBuilderNameToAddress(arg0 string) (common.Address, error) {
	return _Rollupclient.Contract.BlockBuilderNameToAddress(&_Rollupclient.CallOpts, arg0)
}

// GetNextRequestedBlockNumber is a free data retrieval call binding the contract method 0xfce2a502.
//
// Solidity: function getNextRequestedBlockNumber() view returns(uint256)
func (_Rollupclient *RollupclientCaller) GetNextRequestedBlockNumber(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Rollupclient.contract.Call(opts, &out, "getNextRequestedBlockNumber")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetNextRequestedBlockNumber is a free data retrieval call binding the contract method 0xfce2a502.
//
// Solidity: function getNextRequestedBlockNumber() view returns(uint256)
func (_Rollupclient *RollupclientSession) GetNextRequestedBlockNumber() (*big.Int, error) {
	return _Rollupclient.Contract.GetNextRequestedBlockNumber(&_Rollupclient.CallOpts)
}

// GetNextRequestedBlockNumber is a free data retrieval call binding the contract method 0xfce2a502.
//
// Solidity: function getNextRequestedBlockNumber() view returns(uint256)
func (_Rollupclient *RollupclientCallerSession) GetNextRequestedBlockNumber() (*big.Int, error) {
	return _Rollupclient.Contract.GetNextRequestedBlockNumber(&_Rollupclient.CallOpts)
}

// NextRequestedBlockNumber is a free data retrieval call binding the contract method 0xc512c561.
//
// Solidity: function nextRequestedBlockNumber() view returns(uint256)
func (_Rollupclient *RollupclientCaller) NextRequestedBlockNumber(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Rollupclient.contract.Call(opts, &out, "nextRequestedBlockNumber")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// NextRequestedBlockNumber is a free data retrieval call binding the contract method 0xc512c561.
//
// Solidity: function nextRequestedBlockNumber() view returns(uint256)
func (_Rollupclient *RollupclientSession) NextRequestedBlockNumber() (*big.Int, error) {
	return _Rollupclient.Contract.NextRequestedBlockNumber(&_Rollupclient.CallOpts)
}

// NextRequestedBlockNumber is a free data retrieval call binding the contract method 0xc512c561.
//
// Solidity: function nextRequestedBlockNumber() view returns(uint256)
func (_Rollupclient *RollupclientCallerSession) NextRequestedBlockNumber() (*big.Int, error) {
	return _Rollupclient.Contract.NextRequestedBlockNumber(&_Rollupclient.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Rollupclient *RollupclientCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Rollupclient.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Rollupclient *RollupclientSession) Owner() (common.Address, error) {
	return _Rollupclient.Contract.Owner(&_Rollupclient.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Rollupclient *RollupclientCallerSession) Owner() (common.Address, error) {
	return _Rollupclient.Contract.Owner(&_Rollupclient.CallOpts)
}

// AddBuilderAddress is a paid mutator transaction binding the contract method 0x0bd0a9e1.
//
// Solidity: function addBuilderAddress(string builderName, address builderAddress) returns()
func (_Rollupclient *RollupclientTransactor) AddBuilderAddress(opts *bind.TransactOpts, builderName string, builderAddress common.Address) (*types.Transaction, error) {
	return _Rollupclient.contract.Transact(opts, "addBuilderAddress", builderName, builderAddress)
}

// AddBuilderAddress is a paid mutator transaction binding the contract method 0x0bd0a9e1.
//
// Solidity: function addBuilderAddress(string builderName, address builderAddress) returns()
func (_Rollupclient *RollupclientSession) AddBuilderAddress(builderName string, builderAddress common.Address) (*types.Transaction, error) {
	return _Rollupclient.Contract.AddBuilderAddress(&_Rollupclient.TransactOpts, builderName, builderAddress)
}

// AddBuilderAddress is a paid mutator transaction binding the contract method 0x0bd0a9e1.
//
// Solidity: function addBuilderAddress(string builderName, address builderAddress) returns()
func (_Rollupclient *RollupclientTransactorSession) AddBuilderAddress(builderName string, builderAddress common.Address) (*types.Transaction, error) {
	return _Rollupclient.Contract.AddBuilderAddress(&_Rollupclient.TransactOpts, builderName, builderAddress)
}

// ProcessCommitment is a paid mutator transaction binding the contract method 0x09f750a1.
//
// Solidity: function processCommitment(bytes32 commitmentIndex, bool isSlash) returns()
func (_Rollupclient *RollupclientTransactor) ProcessCommitment(opts *bind.TransactOpts, commitmentIndex [32]byte, isSlash bool) (*types.Transaction, error) {
	return _Rollupclient.contract.Transact(opts, "processCommitment", commitmentIndex, isSlash)
}

// ProcessCommitment is a paid mutator transaction binding the contract method 0x09f750a1.
//
// Solidity: function processCommitment(bytes32 commitmentIndex, bool isSlash) returns()
func (_Rollupclient *RollupclientSession) ProcessCommitment(commitmentIndex [32]byte, isSlash bool) (*types.Transaction, error) {
	return _Rollupclient.Contract.ProcessCommitment(&_Rollupclient.TransactOpts, commitmentIndex, isSlash)
}

// ProcessCommitment is a paid mutator transaction binding the contract method 0x09f750a1.
//
// Solidity: function processCommitment(bytes32 commitmentIndex, bool isSlash) returns()
func (_Rollupclient *RollupclientTransactorSession) ProcessCommitment(commitmentIndex [32]byte, isSlash bool) (*types.Transaction, error) {
	return _Rollupclient.Contract.ProcessCommitment(&_Rollupclient.TransactOpts, commitmentIndex, isSlash)
}

// ReceiveBlockData is a paid mutator transaction binding the contract method 0x0ec508dd.
//
// Solidity: function receiveBlockData(string[] txnList, uint256 blockNumber, string blockBuilderName) returns()
func (_Rollupclient *RollupclientTransactor) ReceiveBlockData(opts *bind.TransactOpts, txnList []string, blockNumber *big.Int, blockBuilderName string) (*types.Transaction, error) {
	return _Rollupclient.contract.Transact(opts, "receiveBlockData", txnList, blockNumber, blockBuilderName)
}

// ReceiveBlockData is a paid mutator transaction binding the contract method 0x0ec508dd.
//
// Solidity: function receiveBlockData(string[] txnList, uint256 blockNumber, string blockBuilderName) returns()
func (_Rollupclient *RollupclientSession) ReceiveBlockData(txnList []string, blockNumber *big.Int, blockBuilderName string) (*types.Transaction, error) {
	return _Rollupclient.Contract.ReceiveBlockData(&_Rollupclient.TransactOpts, txnList, blockNumber, blockBuilderName)
}

// ReceiveBlockData is a paid mutator transaction binding the contract method 0x0ec508dd.
//
// Solidity: function receiveBlockData(string[] txnList, uint256 blockNumber, string blockBuilderName) returns()
func (_Rollupclient *RollupclientTransactorSession) ReceiveBlockData(txnList []string, blockNumber *big.Int, blockBuilderName string) (*types.Transaction, error) {
	return _Rollupclient.Contract.ReceiveBlockData(&_Rollupclient.TransactOpts, txnList, blockNumber, blockBuilderName)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Rollupclient *RollupclientTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Rollupclient.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Rollupclient *RollupclientSession) RenounceOwnership() (*types.Transaction, error) {
	return _Rollupclient.Contract.RenounceOwnership(&_Rollupclient.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Rollupclient *RollupclientTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Rollupclient.Contract.RenounceOwnership(&_Rollupclient.TransactOpts)
}

// RequestBlockData is a paid mutator transaction binding the contract method 0x9943545e.
//
// Solidity: function requestBlockData(uint256 blockNumber) returns()
func (_Rollupclient *RollupclientTransactor) RequestBlockData(opts *bind.TransactOpts, blockNumber *big.Int) (*types.Transaction, error) {
	return _Rollupclient.contract.Transact(opts, "requestBlockData", blockNumber)
}

// RequestBlockData is a paid mutator transaction binding the contract method 0x9943545e.
//
// Solidity: function requestBlockData(uint256 blockNumber) returns()
func (_Rollupclient *RollupclientSession) RequestBlockData(blockNumber *big.Int) (*types.Transaction, error) {
	return _Rollupclient.Contract.RequestBlockData(&_Rollupclient.TransactOpts, blockNumber)
}

// RequestBlockData is a paid mutator transaction binding the contract method 0x9943545e.
//
// Solidity: function requestBlockData(uint256 blockNumber) returns()
func (_Rollupclient *RollupclientTransactorSession) RequestBlockData(blockNumber *big.Int) (*types.Transaction, error) {
	return _Rollupclient.Contract.RequestBlockData(&_Rollupclient.TransactOpts, blockNumber)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Rollupclient *RollupclientTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Rollupclient.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Rollupclient *RollupclientSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Rollupclient.Contract.TransferOwnership(&_Rollupclient.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Rollupclient *RollupclientTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Rollupclient.Contract.TransferOwnership(&_Rollupclient.TransactOpts, newOwner)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Rollupclient *RollupclientTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _Rollupclient.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Rollupclient *RollupclientSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _Rollupclient.Contract.Fallback(&_Rollupclient.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Rollupclient *RollupclientTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _Rollupclient.Contract.Fallback(&_Rollupclient.TransactOpts, calldata)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Rollupclient *RollupclientTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Rollupclient.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Rollupclient *RollupclientSession) Receive() (*types.Transaction, error) {
	return _Rollupclient.Contract.Receive(&_Rollupclient.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Rollupclient *RollupclientTransactorSession) Receive() (*types.Transaction, error) {
	return _Rollupclient.Contract.Receive(&_Rollupclient.TransactOpts)
}

// RollupclientBlockDataReceivedIterator is returned from FilterBlockDataReceived and is used to iterate over the raw logs and unpacked data for BlockDataReceived events raised by the Rollupclient contract.
type RollupclientBlockDataReceivedIterator struct {
	Event *RollupclientBlockDataReceived // Event containing the contract specifics and raw log

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
func (it *RollupclientBlockDataReceivedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RollupclientBlockDataReceived)
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
		it.Event = new(RollupclientBlockDataReceived)
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
func (it *RollupclientBlockDataReceivedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RollupclientBlockDataReceivedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RollupclientBlockDataReceived represents a BlockDataReceived event raised by the Rollupclient contract.
type RollupclientBlockDataReceived struct {
	TxnList          []string
	BlockNumber      *big.Int
	BlockBuilderName string
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterBlockDataReceived is a free log retrieval operation binding the contract event 0x7f6309d67e0de7438df797c33b4cd882df8e732c453296b2998bab55dd2ed005.
//
// Solidity: event BlockDataReceived(string[] txnList, uint256 blockNumber, string blockBuilderName)
func (_Rollupclient *RollupclientFilterer) FilterBlockDataReceived(opts *bind.FilterOpts) (*RollupclientBlockDataReceivedIterator, error) {

	logs, sub, err := _Rollupclient.contract.FilterLogs(opts, "BlockDataReceived")
	if err != nil {
		return nil, err
	}
	return &RollupclientBlockDataReceivedIterator{contract: _Rollupclient.contract, event: "BlockDataReceived", logs: logs, sub: sub}, nil
}

// WatchBlockDataReceived is a free log subscription operation binding the contract event 0x7f6309d67e0de7438df797c33b4cd882df8e732c453296b2998bab55dd2ed005.
//
// Solidity: event BlockDataReceived(string[] txnList, uint256 blockNumber, string blockBuilderName)
func (_Rollupclient *RollupclientFilterer) WatchBlockDataReceived(opts *bind.WatchOpts, sink chan<- *RollupclientBlockDataReceived) (event.Subscription, error) {

	logs, sub, err := _Rollupclient.contract.WatchLogs(opts, "BlockDataReceived")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RollupclientBlockDataReceived)
				if err := _Rollupclient.contract.UnpackLog(event, "BlockDataReceived", log); err != nil {
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

// ParseBlockDataReceived is a log parse operation binding the contract event 0x7f6309d67e0de7438df797c33b4cd882df8e732c453296b2998bab55dd2ed005.
//
// Solidity: event BlockDataReceived(string[] txnList, uint256 blockNumber, string blockBuilderName)
func (_Rollupclient *RollupclientFilterer) ParseBlockDataReceived(log types.Log) (*RollupclientBlockDataReceived, error) {
	event := new(RollupclientBlockDataReceived)
	if err := _Rollupclient.contract.UnpackLog(event, "BlockDataReceived", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RollupclientBlockDataRequestedIterator is returned from FilterBlockDataRequested and is used to iterate over the raw logs and unpacked data for BlockDataRequested events raised by the Rollupclient contract.
type RollupclientBlockDataRequestedIterator struct {
	Event *RollupclientBlockDataRequested // Event containing the contract specifics and raw log

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
func (it *RollupclientBlockDataRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RollupclientBlockDataRequested)
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
		it.Event = new(RollupclientBlockDataRequested)
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
func (it *RollupclientBlockDataRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RollupclientBlockDataRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RollupclientBlockDataRequested represents a BlockDataRequested event raised by the Rollupclient contract.
type RollupclientBlockDataRequested struct {
	BlockNumber *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterBlockDataRequested is a free log retrieval operation binding the contract event 0xa4e655157264f5c4f534fdcbf662f33a4bac8f9544f8554511e53e8745c7ea62.
//
// Solidity: event BlockDataRequested(uint256 blockNumber)
func (_Rollupclient *RollupclientFilterer) FilterBlockDataRequested(opts *bind.FilterOpts) (*RollupclientBlockDataRequestedIterator, error) {

	logs, sub, err := _Rollupclient.contract.FilterLogs(opts, "BlockDataRequested")
	if err != nil {
		return nil, err
	}
	return &RollupclientBlockDataRequestedIterator{contract: _Rollupclient.contract, event: "BlockDataRequested", logs: logs, sub: sub}, nil
}

// WatchBlockDataRequested is a free log subscription operation binding the contract event 0xa4e655157264f5c4f534fdcbf662f33a4bac8f9544f8554511e53e8745c7ea62.
//
// Solidity: event BlockDataRequested(uint256 blockNumber)
func (_Rollupclient *RollupclientFilterer) WatchBlockDataRequested(opts *bind.WatchOpts, sink chan<- *RollupclientBlockDataRequested) (event.Subscription, error) {

	logs, sub, err := _Rollupclient.contract.WatchLogs(opts, "BlockDataRequested")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RollupclientBlockDataRequested)
				if err := _Rollupclient.contract.UnpackLog(event, "BlockDataRequested", log); err != nil {
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

// ParseBlockDataRequested is a log parse operation binding the contract event 0xa4e655157264f5c4f534fdcbf662f33a4bac8f9544f8554511e53e8745c7ea62.
//
// Solidity: event BlockDataRequested(uint256 blockNumber)
func (_Rollupclient *RollupclientFilterer) ParseBlockDataRequested(log types.Log) (*RollupclientBlockDataRequested, error) {
	event := new(RollupclientBlockDataRequested)
	if err := _Rollupclient.contract.UnpackLog(event, "BlockDataRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RollupclientCommitmentProcessedIterator is returned from FilterCommitmentProcessed and is used to iterate over the raw logs and unpacked data for CommitmentProcessed events raised by the Rollupclient contract.
type RollupclientCommitmentProcessedIterator struct {
	Event *RollupclientCommitmentProcessed // Event containing the contract specifics and raw log

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
func (it *RollupclientCommitmentProcessedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RollupclientCommitmentProcessed)
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
		it.Event = new(RollupclientCommitmentProcessed)
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
func (it *RollupclientCommitmentProcessedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RollupclientCommitmentProcessedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RollupclientCommitmentProcessed represents a CommitmentProcessed event raised by the Rollupclient contract.
type RollupclientCommitmentProcessed struct {
	CommitmentHash [32]byte
	IsSlash        bool
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterCommitmentProcessed is a free log retrieval operation binding the contract event 0xddc1768a3a762a04e5fd3abea8ae3b60e23bcf290f4a032280e6a726611d41f5.
//
// Solidity: event CommitmentProcessed(bytes32 commitmentHash, bool isSlash)
func (_Rollupclient *RollupclientFilterer) FilterCommitmentProcessed(opts *bind.FilterOpts) (*RollupclientCommitmentProcessedIterator, error) {

	logs, sub, err := _Rollupclient.contract.FilterLogs(opts, "CommitmentProcessed")
	if err != nil {
		return nil, err
	}
	return &RollupclientCommitmentProcessedIterator{contract: _Rollupclient.contract, event: "CommitmentProcessed", logs: logs, sub: sub}, nil
}

// WatchCommitmentProcessed is a free log subscription operation binding the contract event 0xddc1768a3a762a04e5fd3abea8ae3b60e23bcf290f4a032280e6a726611d41f5.
//
// Solidity: event CommitmentProcessed(bytes32 commitmentHash, bool isSlash)
func (_Rollupclient *RollupclientFilterer) WatchCommitmentProcessed(opts *bind.WatchOpts, sink chan<- *RollupclientCommitmentProcessed) (event.Subscription, error) {

	logs, sub, err := _Rollupclient.contract.WatchLogs(opts, "CommitmentProcessed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RollupclientCommitmentProcessed)
				if err := _Rollupclient.contract.UnpackLog(event, "CommitmentProcessed", log); err != nil {
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
func (_Rollupclient *RollupclientFilterer) ParseCommitmentProcessed(log types.Log) (*RollupclientCommitmentProcessed, error) {
	event := new(RollupclientCommitmentProcessed)
	if err := _Rollupclient.contract.UnpackLog(event, "CommitmentProcessed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RollupclientOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Rollupclient contract.
type RollupclientOwnershipTransferredIterator struct {
	Event *RollupclientOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *RollupclientOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RollupclientOwnershipTransferred)
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
		it.Event = new(RollupclientOwnershipTransferred)
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
func (it *RollupclientOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RollupclientOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RollupclientOwnershipTransferred represents a OwnershipTransferred event raised by the Rollupclient contract.
type RollupclientOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Rollupclient *RollupclientFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*RollupclientOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Rollupclient.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &RollupclientOwnershipTransferredIterator{contract: _Rollupclient.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Rollupclient *RollupclientFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *RollupclientOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Rollupclient.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RollupclientOwnershipTransferred)
				if err := _Rollupclient.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_Rollupclient *RollupclientFilterer) ParseOwnershipTransferred(log types.Log) (*RollupclientOwnershipTransferred, error) {
	event := new(RollupclientOwnershipTransferred)
	if err := _Rollupclient.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
