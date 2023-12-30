// Copyright 2014 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package types

import (
	"bytes"
	"container/heap"
	"errors"
	"github.com/DogeProtocol/dp/crypto"
	"github.com/DogeProtocol/dp/crypto/cryptobase"
	"io"
	"math/big"
	"runtime/debug"
	"sort"
	"sync/atomic"
	"time"

	"github.com/DogeProtocol/dp/common"
	"github.com/DogeProtocol/dp/rlp"
)

var (
	ErrInvalidSig           = errors.New("invalid transaction v, r, s values")
	ErrUnexpectedProtection = errors.New("transaction type does not supported EIP-155 protected signatures")
	ErrInvalidTxType        = errors.New("transaction type not valid in this context")
	ErrTxTypeNotSupported   = errors.New("transaction type not supported")
	ErrGasFeeCapTooLow      = errors.New("fee cap less than base fee")
	errEmptyTypedTx         = errors.New("empty typed transaction bytes")
)

// Transaction types.
const (
	DefaultFeeTxType = iota
)

// Transaction is an Ethereum transaction.
type Transaction struct {
	inner TxData    // Consensus contents of a transaction
	time  time.Time // Time first seen locally (spam avoidance)

	// caches
	hash atomic.Value
	size atomic.Value
	from atomic.Value
}

// NewTx creates a new transaction.
func NewTx(inner TxData) *Transaction {
	tx := new(Transaction)
	tx.setDecoded(inner.copy(), 0)
	return tx
}

// TxData is the underlying data of a transaction.
//
// This is implemented by DynamicFeeTx, LegacyTx and AccessListTx.
type TxData interface {
	txType() byte // returns the type ID
	copy() TxData // creates a deep copy and initializes all fields

	chainID() *big.Int
	accessList() AccessList
	data() []byte
	gas() uint64
	gasPrice() *big.Int
	maxGasTier() GasTier
	value() *big.Int
	nonce() uint64
	to() *common.Address
	remarks() []byte
	verifyFields() bool

	rawSignatureValues() (v, r, s *big.Int)
	setSignatureValues(chainID, v, r, s *big.Int)
}

// EncodeRLP implements rlp.Encoder
func (tx *Transaction) EncodeRLP(w io.Writer) error {
	// It's an EIP-2718 typed TX envelope.
	buf := encodeBufferPool.Get().(*bytes.Buffer)
	defer encodeBufferPool.Put(buf)
	buf.Reset()
	if err := tx.encodeTyped(buf); err != nil {
		return err
	}
	return rlp.Encode(w, buf.Bytes())
}

// encodeTyped writes the canonical encoding of a typed transaction to w.
func (tx *Transaction) encodeTyped(w *bytes.Buffer) error {
	w.WriteByte(tx.Type())
	return rlp.Encode(w, tx.inner)
}

// MarshalBinary returns the canonical encoding of the transaction.
// For legacy transactions, it returns the RLP encoding. For EIP-2718 typed
// transactions, it returns the type and payload.
func (tx *Transaction) MarshalBinary() ([]byte, error) {
	var buf bytes.Buffer
	err := tx.encodeTyped(&buf)
	return buf.Bytes(), err
}

// DecodeRLP implements rlp.Decoder
func (tx *Transaction) DecodeRLP(s *rlp.Stream) error {
	kind, _, err := s.Kind()
	switch {
	case err != nil:
		return err
	case kind == rlp.String:
		// It's an EIP-2718 typed TX envelope.
		var b []byte
		if b, err = s.Bytes(); err != nil {
			return err
		}
		inner, err := tx.decodeTyped(b)
		if err == nil {
			tx.setDecoded(inner, len(b))
		}
		return err
	default:
		return rlp.ErrExpectedList
	}
}

// UnmarshalBinary decodes the canonical encoding of transactions.
// It supports legacy RLP transactions and EIP2718 typed transactions.
func (tx *Transaction) UnmarshalBinary(b []byte) error {
	if len(b) > 0 && b[0] > 0x7f {
		// It's a legacy transaction.
		return errors.New("unsupported txn")
	}
	// It's an EIP2718 typed transaction envelope.
	inner, err := tx.decodeTyped(b)
	if err != nil {
		return err
	}
	tx.setDecoded(inner, len(b))
	return nil
}

// decodeTyped decodes a typed transaction from the canonical format.
func (tx *Transaction) decodeTyped(b []byte) (TxData, error) {
	if len(b) == 0 {
		return nil, errEmptyTypedTx
	}
	switch b[0] {

	case DefaultFeeTxType:
		var inner DefaultFeeTx
		err := rlp.DecodeBytes(b[1:], &inner)
		return &inner, err
	default:
		return nil, ErrTxTypeNotSupported
	}
}

// setDecoded sets the inner transaction and size after decoding.
func (tx *Transaction) setDecoded(inner TxData, size int) {
	tx.inner = inner
	//tx.time = time.Now()
	tx.time = time.Date(
		2009, 11, 17, 20, 34, 58, 651387237, time.UTC)
	if size > 0 {
		tx.size.Store(common.StorageSize(size))
	}
}

func sanityCheckSignature(digestHash []byte, v *big.Int, r *big.Int, s *big.Int, maybeProtected bool) error {
	if isProtectedV(v) && !maybeProtected {
		return ErrUnexpectedProtection
	}

	var plainV byte
	if isProtectedV(v) {
		chainID := deriveChainId(v).Uint64()
		plainV = byte(v.Uint64() - 35 - 2*chainID)
	} else if maybeProtected {
		// Only EIP-155 signatures can be optionally protected. Since
		// we determined this v value is not protected, it must be a
		// raw 27 or 28.
		plainV = byte(v.Uint64() - 27)
	} else {
		// If the signature is not optionally protected, we assume it
		// must already be equal to the recovery id.
		plainV = byte(v.Uint64())
	}
	if !cryptobase.SigAlg.ValidateSignatureValues(digestHash, plainV, r, s) {
		return ErrInvalidSig
	}

	return nil
}

func isProtectedV(V *big.Int) bool {
	if V.BitLen() <= 8 {
		v := V.Uint64()
		return v != 27 && v != 28 && v != 1 && v != 0
	}
	// anything not 27 or 28 is considered protected
	return true
}

// Protected says whether the transaction is replay-protected.
func (tx *Transaction) Protected() bool {
	return true
}

// Type returns the transaction type.
func (tx *Transaction) Type() uint8 {
	return tx.inner.txType()
}

// ChainId returns the EIP155 chain ID of the transaction. The return value will always be
// non-nil. For legacy transactions which are not replay-protected, the return value is
// zero.
func (tx *Transaction) ChainId() *big.Int {
	return tx.inner.chainID()
}

// Data returns the input data of the transaction.
func (tx *Transaction) Data() []byte { return tx.inner.data() }

// AccessList returns the access list of the transaction.
func (tx *Transaction) AccessList() AccessList { return tx.inner.accessList() }

// Gas returns the gas limit of the transaction.
func (tx *Transaction) Gas() uint64 { return tx.inner.gas() }

// GasPrice returns the gas price of the transaction.
func (tx *Transaction) GasPrice() *big.Int { return new(big.Int).Set(tx.inner.gasPrice()) }

func (tx *Transaction) MaxGasTier() *big.Int { return new(big.Int).Set(tx.inner.gasPrice()) }

// Value returns the ether amount of the transaction.
func (tx *Transaction) Value() *big.Int { return new(big.Int).Set(tx.inner.value()) }

// Nonce returns the sender account nonce of the transaction.
func (tx *Transaction) Nonce() uint64 { return tx.inner.nonce() }

func (tx *Transaction) Remarks() []byte { return tx.inner.remarks() }

func (tx *Transaction) VerifyFields() bool { return tx.inner.verifyFields() }

// To returns the recipient address of the transaction.
// For contract-creation transactions, To returns nil.
func (tx *Transaction) To() *common.Address {
	// Copy the pointed-to address.
	ito := tx.inner.to()
	if ito == nil {
		return nil
	}
	cpy := *ito
	return &cpy
}

// Cost returns gas * gasPrice + value.
func (tx *Transaction) Cost() *big.Int {
	total := new(big.Int).Mul(tx.GasPrice(), new(big.Int).SetUint64(tx.Gas()))
	total.Add(total, tx.Value())
	return total
}

// RawSignatureValues returns the V, R, S signature values of the transaction.
// The return values should not be modified by the caller.
func (tx *Transaction) RawSignatureValues() (v, r, s *big.Int) {
	return tx.inner.rawSignatureValues()
}

// Hash returns the transaction hash.
func (tx *Transaction) Hash() common.Hash {
	if hash := tx.hash.Load(); hash != nil {
		return hash.(common.Hash)
	}

	var h common.Hash
	h = prefixedRlpHash(tx.Type(), tx.inner)
	tx.hash.Store(h)
	return h
}

// Size returns the true RLP encoded storage size of the transaction, either by
// encoding and returning it, or returning a previously cached value.
func (tx *Transaction) Size() common.StorageSize {
	if size := tx.size.Load(); size != nil {
		return size.(common.StorageSize)
	}
	c := writeCounter(0)
	rlp.Encode(&c, &tx.inner)
	tx.size.Store(common.StorageSize(c))
	return common.StorageSize(c)
}

// WithSignature returns a new transaction with the given signature.
// This signature needs to be in the [R || S || V] format where V is 0 or 1.
func (tx *Transaction) WithSignature(signer Signer, sig []byte) (*Transaction, error) {
	r, s, v, err := signer.SignatureValues(tx, sig)
	if err != nil {
		return nil, err
	}
	cpy := tx.inner.copy()
	cpy.setSignatureValues(signer.ChainID(), v, r, s)
	t := time.Date(
		2009, 11, 17, 20, 34, 58, 651387237, time.UTC)

	copiedTxn := &Transaction{inner: cpy, time: t}
	_, err = Sender(signer, copiedTxn)
	if err != nil {
		return nil, err
	}
	return &Transaction{inner: cpy, time: t}, nil
}

func (tx *Transaction) Verify(digestHash []byte) bool {
	_, r, s := tx.RawSignatureValues()
	return cryptobase.SigAlg.ValidateSignatureValues(digestHash, 1, r, s)
}

// Transactions implements DerivableList for transactions.
type Transactions []*Transaction

func (s Transactions) Less(i, j int) bool { return s[i].Nonce() < s[j].Nonce() }
func (s Transactions) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

// Len returns the length of s.
func (s Transactions) Len() int { return len(s) }

// EncodeIndex encodes the i'th transaction to w. Note that this does not check for errors
// because we assume that *Transaction will only ever contain valid txns that were either
// constructed by decoding or via public API in this package.
func (s Transactions) EncodeIndex(i int, w *bytes.Buffer) {
	tx := s[i]
	tx.encodeTyped(w)
}

// TxDifference returns a new set which is the difference between a and b.
func TxDifference(a, b Transactions) Transactions {
	keep := make(Transactions, 0, len(a))

	remove := make(map[common.Hash]struct{})
	for _, tx := range b {
		remove[tx.Hash()] = struct{}{}
	}

	for _, tx := range a {
		if _, ok := remove[tx.Hash()]; !ok {
			keep = append(keep, tx)
		}
	}

	return keep
}

// TxByNonce implements the sort interface to allow sorting a list of transactions
// by their nonces. This is usually only useful for sorting transactions from a
// single account, otherwise a nonce comparison doesn't make much sense.
type TxByNonce Transactions

func (s TxByNonce) Len() int           { return len(s) }
func (s TxByNonce) Less(i, j int) bool { return s[i].Nonce() < s[j].Nonce() }
func (s TxByNonce) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

// WrapperTxn wraps a transaction with its gas price or effective miner gasTipCap
type WrapperTxn struct {
	tx         *Transaction
	sortPrefix []byte
}

// NewWrapperTxn creates a wrapped transaction, calculating the effective
// miner gasTipCap if a base fee is provided.
// Returns error in case of a negative effective miner gasTipCap.
func NewWrapperTxn(tx *Transaction, sortPrefix []byte) (*WrapperTxn, error) {
	return &WrapperTxn{
		tx:         tx,
		sortPrefix: sortPrefix,
	}, nil
}

// TxBySortPrefix implements both the sort and the heap interface, making it useful
// for all at once sorting as well as individually adding and removing elements.
type TxBySortPrefix []*WrapperTxn

func (s TxBySortPrefix) Len() int { return len(s) }
func (s TxBySortPrefix) Less(i, j int) bool {
	cmp := bytes.Compare(s[i].sortPrefix, s[j].sortPrefix) < 0
	return cmp
}
func (s TxBySortPrefix) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s *TxBySortPrefix) Push(x interface{}) {
	*s = append(*s, x.(*WrapperTxn))
}

func (s *TxBySortPrefix) Pop() interface{} {
	old := *s
	n := len(old)
	x := old[n-1]
	*s = old[0 : n-1]
	return x
}

// TransactionsByNonce represents a set of transactions supporting removing
// entire batches of transactions for non-executable accounts.
type TransactionsByNonce struct {
	txns             map[common.Address]Transactions // Per account nonce-sorted list of transactions
	heads            TxBySortPrefix                  // Next transaction for each unique account (heap)
	signer           Signer                          // Signer for the set of transactions
	parentHash       common.Hash
	orderedAddresses []common.Address
	addressIndex     int
	round            int
}

// NewTransactionsByNonce creates a transaction set that can retrieve transactions in a nonce-honouring way.
// Note, the input map is reowned so the caller should not interact any more with
// if after providing it to the constructor.
func NewTransactionsByNonce(signer Signer, txs map[common.Address]Transactions, parentHash common.Hash) *TransactionsByNonce {
	// Initialize a time based heap with the head transactions
	heads := make(TxBySortPrefix, 0, len(txs))
	for from, accTxs := range txs {
		for i := 0; i < len(accTxs); i++ {
			_, err := Sender(signer, accTxs[i])
			if err != nil {
				delete(txs, from)
				continue
			}
		}
		if len(accTxs) == 0 {
			delete(txs, from)
			continue
		}
		sort.Sort(accTxs)
		prevTxn := accTxs[0]
		for i := 1; i < len(accTxs); i++ {
			if accTxs[i].Nonce() != prevTxn.Nonce()+1 {
				accTxs = accTxs[:i]
				break
			}
			prevTxn = accTxs[i]
		}
		if len(accTxs) == 0 {
			delete(txs, from)
			continue
		}
		acc, _ := Sender(signer, accTxs[0])
		sortPrefix := crypto.Keccak256(parentHash.Bytes(), acc.Bytes())
		wrapped, err := NewWrapperTxn(accTxs[0], sortPrefix)
		// Remove transaction if sender doesn't match from, or if wrapping fails.
		if acc != from || err != nil {
			delete(txs, from)
			continue
		}
		heads = append(heads, wrapped)
		txs[from] = accTxs[0:]
	}
	heap.Init(&heads)

	// Assemble and return the transaction set
	output := &TransactionsByNonce{
		txns:       txs,
		heads:      heads,
		signer:     signer,
		parentHash: parentHash,
	}
	output.internalSort()
	output.ResetCursor()
	return output
}

func (t *TransactionsByNonce) GetList() []common.Hash {
	txnList := make([]common.Hash, 0, len(t.txns))

	for _, accTxs := range t.txns {
		for i := 0; i < len(accTxs); i++ {
			txnList = append(txnList, accTxs[i].Hash())
		}
	}

	return txnList
}

func (t *TransactionsByNonce) GetTotalCount() int {
	count := 0

	for _, accTxs := range t.txns {
		count = count + len(accTxs)
	}

	return count
}

func (t *TransactionsByNonce) GetMap() map[common.Address]Transactions {
	return t.txns
}

// Peek returns the next transaction
func (t *TransactionsByNonce) internalSort() {
	t.orderedAddresses = make([]common.Address, len(t.txns))
	txnIndex := 0
	for from, _ := range t.txns {
		t.orderedAddresses[txnIndex] = from
		txnIndex = txnIndex + 1
	}
	parentHashBytes := t.parentHash.Bytes()
	sort.SliceStable(t.orderedAddresses, func(i, j int) bool {
		sortPrefixI := crypto.Keccak256(parentHashBytes, t.orderedAddresses[i].Bytes())
		sortPrefixJ := crypto.Keccak256(parentHashBytes, t.orderedAddresses[j].Bytes())
		cmp := bytes.Compare(sortPrefixI, sortPrefixJ) < 0
		return cmp
	})
}

func (t *TransactionsByNonce) PeekCursor() *Transaction {
	if t.addressIndex < 0 || len(t.txns) == 0 {
		return nil
	}
	return t.txns[t.orderedAddresses[t.addressIndex]][t.round]
}

func (t *TransactionsByNonce) ResetCursor() {
	t.round = 0
	t.addressIndex = -1
}

func (t *TransactionsByNonce) NextCursor() bool {
	if t.addressIndex == -2 {
		return false
	}
	t.addressIndex = t.addressIndex + 1
	if t.addressIndex >= len(t.orderedAddresses) {
		t.addressIndex = 0
		t.round = t.round + 1
	}

	for i := t.addressIndex; i < len(t.orderedAddresses); i++ {
		if t.addressIndex < 0 {
			debug.PrintStack()
		}
		if t.txns[t.orderedAddresses[i]].Len() > t.round {
			t.addressIndex = i
			return true
		}
	}

	t.round = t.round + 1
	for i := 0; i < t.addressIndex; i++ {
		if t.txns[t.orderedAddresses[i]].Len() > t.round {
			t.addressIndex = i
			return true
		}
	}
	t.addressIndex = -2
	return false
}

// Peek returns the next transaction
func (t *TransactionsByNonce) Peek1() *Transaction {
	if len(t.heads) == 0 {
		return nil
	}
	return t.heads[0].tx
}

// Shift replaces the current best head with the next one from the same account.
func (t *TransactionsByNonce) Shift1() {
	acc, _ := Sender(t.signer, t.heads[0].tx)
	if txs, ok := t.txns[acc]; ok && len(txs) > 0 {
		sortPrefix := crypto.Keccak256(t.parentHash.Bytes(), acc.Bytes())
		if wrapped, err := NewWrapperTxn(txs[0], sortPrefix); err == nil {
			t.heads[0], t.txns[acc] = wrapped, txs[1:]
			heap.Fix(&t.heads, 0)
			return
		}
	}
	heap.Pop(&t.heads)
}

// Pop removes the best transaction, *not* replacing it with the next one from
// the same account. This should be used when a transaction cannot be executed
// and hence all subsequent ones should be discarded from the same account.
func (t *TransactionsByNonce) Pop1() {
	heap.Pop(&t.heads)
}

// Message is a fully derived transaction and implements core.Message
//
// NOTE: In a future PR this will be removed.
type Message struct {
	to         *common.Address
	from       common.Address
	nonce      uint64
	amount     *big.Int
	gasLimit   uint64
	gasPrice   *big.Int
	data       []byte
	accessList AccessList
	checkNonce bool
	remarks    []byte
}

func NewMessage(from common.Address, to *common.Address, nonce uint64, amount *big.Int, gasLimit uint64, gasPrice *big.Int, data []byte, accessList AccessList, checkNonce bool) Message {
	return Message{
		from:       from,
		to:         to,
		nonce:      nonce,
		amount:     amount,
		gasLimit:   gasLimit,
		gasPrice:   gasPrice,
		data:       data,
		accessList: accessList,
		checkNonce: checkNonce,
	}
}

// AsMessage returns the transaction as a core.Message.
func (tx *Transaction) AsMessage(s Signer) (Message, error) {
	msg := Message{
		nonce:      tx.Nonce(),
		gasLimit:   tx.Gas(),
		gasPrice:   new(big.Int).Set(tx.GasPrice()),
		to:         tx.To(),
		amount:     tx.Value(),
		data:       tx.Data(),
		accessList: tx.AccessList(),
		checkNonce: true,
		remarks:    tx.Remarks(),
	}
	var err error
	msg.from, err = Sender(s, tx)
	return msg, err
}

func (m Message) From() common.Address            { return m.from }
func (m Message) To() *common.Address             { return m.to }
func (m Message) GasPrice() *big.Int              { return m.gasPrice }
func (m Message) Value() *big.Int                 { return m.amount }
func (m Message) Gas() uint64                     { return m.gasLimit }
func (m Message) Nonce() uint64                   { return m.nonce }
func (m Message) Data() []byte                    { return m.data }
func (m Message) AccessList() AccessList          { return m.accessList }
func (m Message) CheckNonce() bool                { return m.checkNonce }
func (m Message) Remarks() []byte                 { return m.remarks }
func (m Message) OverrideGasPrice(price *big.Int) { m.gasPrice.Set(price) }
