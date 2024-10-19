package relay

import "errors"

var (
	InfoTitleLatestBlockDetails = "Get latest block details"
	InfoTitleAccountDetails     = "Get account details"
	InfoTitleTransaction        = "Get Transaction"
	InfoTitleSendTransaction    = "Send Transaction"
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
	MsgTimeDuration       = "Time Duration"
	MsgStatus             = "Status"
	MsgError              = "Error"
)

var (
	ErrEmptyAddress   = errors.New("empty address")
	ErrInvalidAddress = errors.New("invalid address")
	ErrEmptyHash      = errors.New("empty hash")
	ErrInvalidHash    = errors.New("invalid hash")
	ErrEmptyRawTxHex  = errors.New("empty raw tx")
)
