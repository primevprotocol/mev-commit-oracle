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

// ClientMetaData contains all meta data concerning the Client contract.
var ClientMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_preConfContract\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string[]\",\"name\":\"txnList\",\"type\":\"string[]\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"blockBuilderName\",\"type\":\"string\"}],\"name\":\"BlockDataReceived\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"}],\"name\":\"BlockDataRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"commitmentHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"isSlash\",\"type\":\"bool\"}],\"name\":\"CommitmentProcessed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"builderName\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"builderAddress\",\"type\":\"address\"}],\"name\":\"addBuilderAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"name\":\"blockBuilderNameToAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"commitmentIndex\",\"type\":\"bytes32\"},{\"internalType\":\"bool\",\"name\":\"isSlash\",\"type\":\"bool\"}],\"name\":\"processCommitment\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string[]\",\"name\":\"txnList\",\"type\":\"string[]\"},{\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"blockBuilderName\",\"type\":\"string\"}],\"name\":\"receiveBlockData\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"}],\"name\":\"requestBlockData\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
}

// ClientABI is the input ABI used to generate the binding from.
// Deprecated: Use ClientMetaData.ABI instead.
var ClientABI = ClientMetaData.ABI

// Client is an auto generated Go binding around an Ethereum contract.
type Client struct {
	ClientCaller     // Read-only binding to the contract
	ClientTransactor // Write-only binding to the contract
	ClientFilterer   // Log filterer for contract events
}

// ClientCaller is an auto generated read-only Go binding around an Ethereum contract.
type ClientCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ClientTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ClientTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ClientFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ClientFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ClientSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ClientSession struct {
	Contract     *Client           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ClientCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ClientCallerSession struct {
	Contract *ClientCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// ClientTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ClientTransactorSession struct {
	Contract     *ClientTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ClientRaw is an auto generated low-level Go binding around an Ethereum contract.
type ClientRaw struct {
	Contract *Client // Generic contract binding to access the raw methods on
}

// ClientCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ClientCallerRaw struct {
	Contract *ClientCaller // Generic read-only contract binding to access the raw methods on
}

// ClientTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ClientTransactorRaw struct {
	Contract *ClientTransactor // Generic write-only contract binding to access the raw methods on
}

// NewClient creates a new instance of Client, bound to a specific deployed contract.
func NewClient(address common.Address, backend bind.ContractBackend) (*Client, error) {
	contract, err := bindClient(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Client{ClientCaller: ClientCaller{contract: contract}, ClientTransactor: ClientTransactor{contract: contract}, ClientFilterer: ClientFilterer{contract: contract}}, nil
}

// NewClientCaller creates a new read-only instance of Client, bound to a specific deployed contract.
func NewClientCaller(address common.Address, caller bind.ContractCaller) (*ClientCaller, error) {
	contract, err := bindClient(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ClientCaller{contract: contract}, nil
}

// NewClientTransactor creates a new write-only instance of Client, bound to a specific deployed contract.
func NewClientTransactor(address common.Address, transactor bind.ContractTransactor) (*ClientTransactor, error) {
	contract, err := bindClient(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ClientTransactor{contract: contract}, nil
}

// NewClientFilterer creates a new log filterer instance of Client, bound to a specific deployed contract.
func NewClientFilterer(address common.Address, filterer bind.ContractFilterer) (*ClientFilterer, error) {
	contract, err := bindClient(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ClientFilterer{contract: contract}, nil
}

// bindClient binds a generic wrapper to an already deployed contract.
func bindClient(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ClientMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Client *ClientRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Client.Contract.ClientCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Client *ClientRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Client.Contract.ClientTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Client *ClientRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Client.Contract.ClientTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Client *ClientCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Client.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Client *ClientTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Client.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Client *ClientTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Client.Contract.contract.Transact(opts, method, params...)
}

// BlockBuilderNameToAddress is a free data retrieval call binding the contract method 0xeebac3ac.
//
// Solidity: function blockBuilderNameToAddress(string ) view returns(address)
func (_Client *ClientCaller) BlockBuilderNameToAddress(opts *bind.CallOpts, arg0 string) (common.Address, error) {
	var out []interface{}
	err := _Client.contract.Call(opts, &out, "blockBuilderNameToAddress", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// BlockBuilderNameToAddress is a free data retrieval call binding the contract method 0xeebac3ac.
//
// Solidity: function blockBuilderNameToAddress(string ) view returns(address)
func (_Client *ClientSession) BlockBuilderNameToAddress(arg0 string) (common.Address, error) {
	return _Client.Contract.BlockBuilderNameToAddress(&_Client.CallOpts, arg0)
}

// BlockBuilderNameToAddress is a free data retrieval call binding the contract method 0xeebac3ac.
//
// Solidity: function blockBuilderNameToAddress(string ) view returns(address)
func (_Client *ClientCallerSession) BlockBuilderNameToAddress(arg0 string) (common.Address, error) {
	return _Client.Contract.BlockBuilderNameToAddress(&_Client.CallOpts, arg0)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Client *ClientCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Client.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Client *ClientSession) Owner() (common.Address, error) {
	return _Client.Contract.Owner(&_Client.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Client *ClientCallerSession) Owner() (common.Address, error) {
	return _Client.Contract.Owner(&_Client.CallOpts)
}

// AddBuilderAddress is a paid mutator transaction binding the contract method 0x0bd0a9e1.
//
// Solidity: function addBuilderAddress(string builderName, address builderAddress) returns()
func (_Client *ClientTransactor) AddBuilderAddress(opts *bind.TransactOpts, builderName string, builderAddress common.Address) (*types.Transaction, error) {
	return _Client.contract.Transact(opts, "addBuilderAddress", builderName, builderAddress)
}

// AddBuilderAddress is a paid mutator transaction binding the contract method 0x0bd0a9e1.
//
// Solidity: function addBuilderAddress(string builderName, address builderAddress) returns()
func (_Client *ClientSession) AddBuilderAddress(builderName string, builderAddress common.Address) (*types.Transaction, error) {
	return _Client.Contract.AddBuilderAddress(&_Client.TransactOpts, builderName, builderAddress)
}

// AddBuilderAddress is a paid mutator transaction binding the contract method 0x0bd0a9e1.
//
// Solidity: function addBuilderAddress(string builderName, address builderAddress) returns()
func (_Client *ClientTransactorSession) AddBuilderAddress(builderName string, builderAddress common.Address) (*types.Transaction, error) {
	return _Client.Contract.AddBuilderAddress(&_Client.TransactOpts, builderName, builderAddress)
}

// ProcessCommitment is a paid mutator transaction binding the contract method 0x09f750a1.
//
// Solidity: function processCommitment(bytes32 commitmentIndex, bool isSlash) returns()
func (_Client *ClientTransactor) ProcessCommitment(opts *bind.TransactOpts, commitmentIndex [32]byte, isSlash bool) (*types.Transaction, error) {
	return _Client.contract.Transact(opts, "processCommitment", commitmentIndex, isSlash)
}

// ProcessCommitment is a paid mutator transaction binding the contract method 0x09f750a1.
//
// Solidity: function processCommitment(bytes32 commitmentIndex, bool isSlash) returns()
func (_Client *ClientSession) ProcessCommitment(commitmentIndex [32]byte, isSlash bool) (*types.Transaction, error) {
	return _Client.Contract.ProcessCommitment(&_Client.TransactOpts, commitmentIndex, isSlash)
}

// ProcessCommitment is a paid mutator transaction binding the contract method 0x09f750a1.
//
// Solidity: function processCommitment(bytes32 commitmentIndex, bool isSlash) returns()
func (_Client *ClientTransactorSession) ProcessCommitment(commitmentIndex [32]byte, isSlash bool) (*types.Transaction, error) {
	return _Client.Contract.ProcessCommitment(&_Client.TransactOpts, commitmentIndex, isSlash)
}

// ReceiveBlockData is a paid mutator transaction binding the contract method 0x0ec508dd.
//
// Solidity: function receiveBlockData(string[] txnList, uint256 blockNumber, string blockBuilderName) returns()
func (_Client *ClientTransactor) ReceiveBlockData(opts *bind.TransactOpts, txnList []string, blockNumber *big.Int, blockBuilderName string) (*types.Transaction, error) {
	return _Client.contract.Transact(opts, "receiveBlockData", txnList, blockNumber, blockBuilderName)
}

// ReceiveBlockData is a paid mutator transaction binding the contract method 0x0ec508dd.
//
// Solidity: function receiveBlockData(string[] txnList, uint256 blockNumber, string blockBuilderName) returns()
func (_Client *ClientSession) ReceiveBlockData(txnList []string, blockNumber *big.Int, blockBuilderName string) (*types.Transaction, error) {
	return _Client.Contract.ReceiveBlockData(&_Client.TransactOpts, txnList, blockNumber, blockBuilderName)
}

// ReceiveBlockData is a paid mutator transaction binding the contract method 0x0ec508dd.
//
// Solidity: function receiveBlockData(string[] txnList, uint256 blockNumber, string blockBuilderName) returns()
func (_Client *ClientTransactorSession) ReceiveBlockData(txnList []string, blockNumber *big.Int, blockBuilderName string) (*types.Transaction, error) {
	return _Client.Contract.ReceiveBlockData(&_Client.TransactOpts, txnList, blockNumber, blockBuilderName)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Client *ClientTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Client.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Client *ClientSession) RenounceOwnership() (*types.Transaction, error) {
	return _Client.Contract.RenounceOwnership(&_Client.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Client *ClientTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Client.Contract.RenounceOwnership(&_Client.TransactOpts)
}

// RequestBlockData is a paid mutator transaction binding the contract method 0x9943545e.
//
// Solidity: function requestBlockData(uint256 blockNumber) returns()
func (_Client *ClientTransactor) RequestBlockData(opts *bind.TransactOpts, blockNumber *big.Int) (*types.Transaction, error) {
	return _Client.contract.Transact(opts, "requestBlockData", blockNumber)
}

// RequestBlockData is a paid mutator transaction binding the contract method 0x9943545e.
//
// Solidity: function requestBlockData(uint256 blockNumber) returns()
func (_Client *ClientSession) RequestBlockData(blockNumber *big.Int) (*types.Transaction, error) {
	return _Client.Contract.RequestBlockData(&_Client.TransactOpts, blockNumber)
}

// RequestBlockData is a paid mutator transaction binding the contract method 0x9943545e.
//
// Solidity: function requestBlockData(uint256 blockNumber) returns()
func (_Client *ClientTransactorSession) RequestBlockData(blockNumber *big.Int) (*types.Transaction, error) {
	return _Client.Contract.RequestBlockData(&_Client.TransactOpts, blockNumber)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Client *ClientTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Client.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Client *ClientSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Client.Contract.TransferOwnership(&_Client.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Client *ClientTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Client.Contract.TransferOwnership(&_Client.TransactOpts, newOwner)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Client *ClientTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _Client.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Client *ClientSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _Client.Contract.Fallback(&_Client.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Client *ClientTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _Client.Contract.Fallback(&_Client.TransactOpts, calldata)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Client *ClientTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Client.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Client *ClientSession) Receive() (*types.Transaction, error) {
	return _Client.Contract.Receive(&_Client.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Client *ClientTransactorSession) Receive() (*types.Transaction, error) {
	return _Client.Contract.Receive(&_Client.TransactOpts)
}

// ClientBlockDataReceivedIterator is returned from FilterBlockDataReceived and is used to iterate over the raw logs and unpacked data for BlockDataReceived events raised by the Client contract.
type ClientBlockDataReceivedIterator struct {
	Event *ClientBlockDataReceived // Event containing the contract specifics and raw log

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
func (it *ClientBlockDataReceivedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ClientBlockDataReceived)
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
		it.Event = new(ClientBlockDataReceived)
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
func (it *ClientBlockDataReceivedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ClientBlockDataReceivedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ClientBlockDataReceived represents a BlockDataReceived event raised by the Client contract.
type ClientBlockDataReceived struct {
	TxnList          []string
	BlockNumber      *big.Int
	BlockBuilderName string
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterBlockDataReceived is a free log retrieval operation binding the contract event 0x7f6309d67e0de7438df797c33b4cd882df8e732c453296b2998bab55dd2ed005.
//
// Solidity: event BlockDataReceived(string[] txnList, uint256 blockNumber, string blockBuilderName)
func (_Client *ClientFilterer) FilterBlockDataReceived(opts *bind.FilterOpts) (*ClientBlockDataReceivedIterator, error) {

	logs, sub, err := _Client.contract.FilterLogs(opts, "BlockDataReceived")
	if err != nil {
		return nil, err
	}
	return &ClientBlockDataReceivedIterator{contract: _Client.contract, event: "BlockDataReceived", logs: logs, sub: sub}, nil
}

// WatchBlockDataReceived is a free log subscription operation binding the contract event 0x7f6309d67e0de7438df797c33b4cd882df8e732c453296b2998bab55dd2ed005.
//
// Solidity: event BlockDataReceived(string[] txnList, uint256 blockNumber, string blockBuilderName)
func (_Client *ClientFilterer) WatchBlockDataReceived(opts *bind.WatchOpts, sink chan<- *ClientBlockDataReceived) (event.Subscription, error) {

	logs, sub, err := _Client.contract.WatchLogs(opts, "BlockDataReceived")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ClientBlockDataReceived)
				if err := _Client.contract.UnpackLog(event, "BlockDataReceived", log); err != nil {
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
func (_Client *ClientFilterer) ParseBlockDataReceived(log types.Log) (*ClientBlockDataReceived, error) {
	event := new(ClientBlockDataReceived)
	if err := _Client.contract.UnpackLog(event, "BlockDataReceived", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ClientBlockDataRequestedIterator is returned from FilterBlockDataRequested and is used to iterate over the raw logs and unpacked data for BlockDataRequested events raised by the Client contract.
type ClientBlockDataRequestedIterator struct {
	Event *ClientBlockDataRequested // Event containing the contract specifics and raw log

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
func (it *ClientBlockDataRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ClientBlockDataRequested)
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
		it.Event = new(ClientBlockDataRequested)
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
func (it *ClientBlockDataRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ClientBlockDataRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ClientBlockDataRequested represents a BlockDataRequested event raised by the Client contract.
type ClientBlockDataRequested struct {
	BlockNumber *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterBlockDataRequested is a free log retrieval operation binding the contract event 0xa4e655157264f5c4f534fdcbf662f33a4bac8f9544f8554511e53e8745c7ea62.
//
// Solidity: event BlockDataRequested(uint256 blockNumber)
func (_Client *ClientFilterer) FilterBlockDataRequested(opts *bind.FilterOpts) (*ClientBlockDataRequestedIterator, error) {

	logs, sub, err := _Client.contract.FilterLogs(opts, "BlockDataRequested")
	if err != nil {
		return nil, err
	}
	return &ClientBlockDataRequestedIterator{contract: _Client.contract, event: "BlockDataRequested", logs: logs, sub: sub}, nil
}

// WatchBlockDataRequested is a free log subscription operation binding the contract event 0xa4e655157264f5c4f534fdcbf662f33a4bac8f9544f8554511e53e8745c7ea62.
//
// Solidity: event BlockDataRequested(uint256 blockNumber)
func (_Client *ClientFilterer) WatchBlockDataRequested(opts *bind.WatchOpts, sink chan<- *ClientBlockDataRequested) (event.Subscription, error) {

	logs, sub, err := _Client.contract.WatchLogs(opts, "BlockDataRequested")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ClientBlockDataRequested)
				if err := _Client.contract.UnpackLog(event, "BlockDataRequested", log); err != nil {
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
func (_Client *ClientFilterer) ParseBlockDataRequested(log types.Log) (*ClientBlockDataRequested, error) {
	event := new(ClientBlockDataRequested)
	if err := _Client.contract.UnpackLog(event, "BlockDataRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ClientCommitmentProcessedIterator is returned from FilterCommitmentProcessed and is used to iterate over the raw logs and unpacked data for CommitmentProcessed events raised by the Client contract.
type ClientCommitmentProcessedIterator struct {
	Event *ClientCommitmentProcessed // Event containing the contract specifics and raw log

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
func (it *ClientCommitmentProcessedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ClientCommitmentProcessed)
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
		it.Event = new(ClientCommitmentProcessed)
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
func (it *ClientCommitmentProcessedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ClientCommitmentProcessedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ClientCommitmentProcessed represents a CommitmentProcessed event raised by the Client contract.
type ClientCommitmentProcessed struct {
	CommitmentHash [32]byte
	IsSlash        bool
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterCommitmentProcessed is a free log retrieval operation binding the contract event 0xddc1768a3a762a04e5fd3abea8ae3b60e23bcf290f4a032280e6a726611d41f5.
//
// Solidity: event CommitmentProcessed(bytes32 commitmentHash, bool isSlash)
func (_Client *ClientFilterer) FilterCommitmentProcessed(opts *bind.FilterOpts) (*ClientCommitmentProcessedIterator, error) {

	logs, sub, err := _Client.contract.FilterLogs(opts, "CommitmentProcessed")
	if err != nil {
		return nil, err
	}
	return &ClientCommitmentProcessedIterator{contract: _Client.contract, event: "CommitmentProcessed", logs: logs, sub: sub}, nil
}

// WatchCommitmentProcessed is a free log subscription operation binding the contract event 0xddc1768a3a762a04e5fd3abea8ae3b60e23bcf290f4a032280e6a726611d41f5.
//
// Solidity: event CommitmentProcessed(bytes32 commitmentHash, bool isSlash)
func (_Client *ClientFilterer) WatchCommitmentProcessed(opts *bind.WatchOpts, sink chan<- *ClientCommitmentProcessed) (event.Subscription, error) {

	logs, sub, err := _Client.contract.WatchLogs(opts, "CommitmentProcessed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ClientCommitmentProcessed)
				if err := _Client.contract.UnpackLog(event, "CommitmentProcessed", log); err != nil {
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
func (_Client *ClientFilterer) ParseCommitmentProcessed(log types.Log) (*ClientCommitmentProcessed, error) {
	event := new(ClientCommitmentProcessed)
	if err := _Client.contract.UnpackLog(event, "CommitmentProcessed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ClientOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Client contract.
type ClientOwnershipTransferredIterator struct {
	Event *ClientOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *ClientOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ClientOwnershipTransferred)
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
		it.Event = new(ClientOwnershipTransferred)
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
func (it *ClientOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ClientOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ClientOwnershipTransferred represents a OwnershipTransferred event raised by the Client contract.
type ClientOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Client *ClientFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*ClientOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Client.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &ClientOwnershipTransferredIterator{contract: _Client.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Client *ClientFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ClientOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Client.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ClientOwnershipTransferred)
				if err := _Client.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_Client *ClientFilterer) ParseOwnershipTransferred(log types.Log) (*ClientOwnershipTransferred, error) {
	event := new(ClientOwnershipTransferred)
	if err := _Client.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
