package accounts

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/republicprotocol/republic-go/blockchain/ethereum"
	"github.com/republicprotocol/republic-go/order"
)

// DarknodeRegistry is the dark node interface
type RenExAccounts struct {
	network      ethereum.Network
	context      context.Context
	conn         ethereum.Conn
	transactOpts *bind.TransactOpts
	callOpts     *bind.CallOpts
	binding      *TraderAccounts
	address      common.Address
}

// NewRenExAccounts returns a Dark node registrar
func NewRenExAccounts(ctx context.Context, conn ethereum.Conn, transactOpts *bind.TransactOpts, callOpts *bind.CallOpts) (RenExAccounts, error) {
	contract, err := NewTraderAccounts(common.HexToAddress(conn.Config.RenExAccountsAddress), bind.ContractBackend(conn.Client))
	if err != nil {
		return RenExAccounts{}, err
	}

	return RenExAccounts{
		network:      conn.Config.Network,
		context:      ctx,
		conn:         conn,
		transactOpts: transactOpts,
		callOpts:     callOpts,
		binding:      contract,
		address:      common.HexToAddress(conn.Config.RenExAccountsAddress),
	}, nil
}

// Register a new token in the RenEx
func (accounts *RenExAccounts) RegisterToken(tokenCode uint32, address string, decimals uint8) error {
	contractAddress := common.HexToAddress(address)
	tx, err := accounts.binding.RegisterToken(accounts.transactOpts, tokenCode, contractAddress, decimals)
	if err != nil {
		return err
	}
	_, err = accounts.conn.PatchedWaitMined(accounts.context, tx)
	return err
}

// Deregister a token from the RenEx
func (accounts *RenExAccounts) DeregisterToken(tokenCode uint32) error {
	tx, err := accounts.binding.DeregisterToken(accounts.transactOpts, tokenCode)
	if err != nil {
		return err
	}
	_, err = accounts.conn.PatchedWaitMined(accounts.context, tx)
	return err
}

// Deposit money into your RenEx account
func (accounts *RenExAccounts) Deposit(tokenCode uint32, value int) error {
	tx, err := accounts.binding.Deposit(accounts.transactOpts, tokenCode, big.NewInt(int64(value)))
	if err != nil {
		return err
	}
	_, err = accounts.conn.PatchedWaitMined(accounts.context, tx)
	return err
}

// Withdraw money from your RenEx account
func (accounts *RenExAccounts) Withdraw(tokenCode uint32, value int) error {
	tx, err := accounts.binding.Withdraw(accounts.transactOpts, tokenCode, big.NewInt(int64(value)))
	if err != nil {
		return err
	}
	_, err = accounts.conn.PatchedWaitMined(accounts.context, tx)
	return err
}

// GetBalance of given token of specific trader in RenEx account
func (accounts *RenExAccounts) Balance(trader string, tokenCode order.Token) (float64, error) {
	traderAddress := common.HexToAddress(trader)
	// Fixme : currently it will just return the value in the token;s minimum unit
	value, err := accounts.binding.GetBalance(accounts.callOpts, traderAddress, uint32(tokenCode))
	floatValue := new(big.Float).SetInt(value)
	if err != nil {
		return 0, err
	}
	float64Value, _ := floatValue.Float64()

	return float64Value, nil
}

// SubmitOrder to the RenEx accounts
func (accounts *RenExAccounts) SubmitOrder(ord order.Order) error {
	nonce := big.NewInt(int64(ord.Nonce))
	tx, err := accounts.binding.SubmitOrder(accounts.transactOpts, ord.ID, uint8(ord.Type), uint8(ord.Parity), uint64(ord.Expiry.Unix()), uint64(ord.Tokens), uint16(ord.Price.Co), uint16(ord.Price.Exp), uint16(ord.Volume.Co), uint16(ord.Volume.Exp), uint16(ord.MinimumVolume.Co), uint16(ord.MinimumVolume.Exp), nonce)
	if err != nil {
		return err
	}
	_, err = accounts.conn.PatchedWaitMined(accounts.context, tx)
	return err
}

// Submit order to the RenEx accounts
func (accounts *RenExAccounts) SubmitMatch(buy, sell order.ID) error {
	accounts.transactOpts.GasPrice = big.NewInt(1000000000)
	tx, err := accounts.binding.SubmitMatch(accounts.transactOpts, buy, sell)
	if err != nil {
		return err
	}
	_, err = accounts.conn.PatchedWaitMined(accounts.context, tx)
	return err
}

// Settle will settle the confirmed order pairs in the RenEx accounts
func (accounts *RenExAccounts) Settle(buy order.Order, sell order.Order) error {
	err := accounts.SubmitOrder(buy)
	if err != nil {
		return err
	}
	err = accounts.SubmitOrder(sell)
	if err != nil {
		return err
	}

	return accounts.SubmitMatch(buy.ID, sell.ID)
}
