package relay

import "errors"

var (
	DIAL_READ_URL  = "http://127.0.0.1:8545" //"//./pipe/geth.ipc"
	DIAL_WRITE_URL = "http://127.0.0.1:8545" //"//./pipe/geth.ipc"
)

var (
	InfoTitleAccountDetails  = "Get account details"
	InfoTitleTransaction     = "Get Transaction"
	InfoTitleSendTransaction = "Send Transaction"
)

var (
	MsgDial               = "Dial"
	MsgAddress            = "Address"
	MsgBalance            = "Balance"
	MsgNonce              = "Nonce"
	MsgBlockNumber        = "Block number"
	MsgHash               = "Hash"
	MsgTransaction        = "Transaction"
	MsgTransactionReceipt = "Transaction receipt"
	MsgSend               = "Send"
	MsgRawRawTxHex        = "Raw tx hex"
	MsgRawTxData          = "Raw tx data"
)

var (
	ErrEmptyAddress   = errors.New("empty address")
	ErrInvalidAddress = errors.New("invalid address")
	ErrEmptyHash      = errors.New("empty hash")
	ErrInvalidHash    = errors.New("invalid hash")
	ErrEmptyRawTxHex  = errors.New("empty raw tx")

	ErrTransWithoutSign = errors.New("server returned transaction without signature")
)
