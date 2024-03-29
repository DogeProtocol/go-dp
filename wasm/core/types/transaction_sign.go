package types

import (
	"errors"
	"github.com/DogeProtocol/dp/common"
	"math/big"
)

var ErrInvalidChainId = errors.New("invalid chain id for signer")

// sigCache is used to cache the derived sender and contains
// the signer used to derive it.
type sigCache struct {
	signer Signer
	from   common.Address
}

// Signer encapsulates transaction signature handling. The name of this type is slightly
// misleading because Signers don't actually sign, they're just for validating and
// processing of signatures.
//
// Note that this interface is not a stable API and may change at any time to accommodate
// new protocol rules.
type Signer interface {

	// SignatureValues returns the raw R, S, V values corresponding to the
	// given signature.
	SignatureValues(tx *Transaction, sig []byte) (r, s, v *big.Int, err error)
	ChainID() *big.Int

	// Hash returns 'signature hash', i.e. the transaction hash that is signed by the
	// private key. This hash does not uniquely identify the transaction.
	Hash(tx *Transaction) (common.Hash, error)

	// Equal returns true if the given signer is the same as the receiver.
	Equal(Signer) bool
}

type londonSigner struct{ chainId *big.Int }

// NewLondonSigner returns a signer that accepts
// - EIP-1559 dynamic fee transactions
// - EIP-2930 access list transactions,
// - EIP-155 replay protected transactions, and
// - legacy Homestead transactions.
func NewLondonSigner(chainId *big.Int) Signer {
	return &londonSigner{
		chainId: chainId,
	}
}

/*
	func NewLondonSignerDefaultChain() Signer {
		return &londonSigner{
			chainId: big.NewInt(DEFAULT_CHAIN_ID),
		}
	}
*/
func (s londonSigner) ChainID() *big.Int {
	return s.chainId
}

func (s londonSigner) Equal(s2 Signer) bool {
	x, ok := s2.(londonSigner)
	return ok && x.chainId.Cmp(s.chainId) == 0
}

func (s londonSigner) SignatureValues(tx *Transaction, sig []byte) (R, S, V *big.Int, err error) {
	txdata1, ok1 := tx.inner.(*DefaultFeeTx)
	if ok1 {
		// Check that chain ID of tx matches the signer. We also accept ID zero here,
		// because it indicates that the chain ID was not specified in the tx.
		if txdata1.ChainID.Sign() != 0 && txdata1.ChainID.Cmp(s.chainId) != 0 {
			return nil, nil, nil, ErrInvalidChainId
		}

		R, S, V, err = decodeSignature(sig)
		if err != nil {
			return nil, nil, nil, err
		}

		V = big.NewInt(1)
		return R, S, V, nil
	}

	return nil, nil, nil, errors.New("signature error")
}

// Hash returns the hash to be signed by the sender.
// It does not uniquely identify the transaction.
func (s londonSigner) Hash(tx *Transaction) (common.Hash, error) {
	if tx.VerifyFields() == false {
		return common.ZERO_HASH, errors.New("txn field verify failed")
	}
	if s.chainId == nil || tx.ChainId() == nil {
		return common.ZERO_HASH, errors.New("chain id is nil")
	}
	if s.chainId.Cmp(tx.ChainId()) != 0 {
		return common.ZERO_HASH, errors.New("signing failed, chainId mismatch")
	}
	return prefixedRlpHash(
		tx.Type(),
		[]interface{}{
			s.chainId,
			tx.Nonce(),
			tx.To(),
			tx.Gas(),
			tx.MaxGasTier(),
			tx.Value(),
			tx.Data(),
			tx.AccessList(),
			tx.Remarks(),
		}), nil
}

func decodeSignature(sig []byte) (r, s, v *big.Int, err error) {
	signature, publicKey, err := common.ExtractTwoParts(sig)
	if err != nil {
		return nil, nil, nil, err
	}

	r = new(big.Int).SetBytes(publicKey)
	s = new(big.Int).SetBytes(signature)
	v = new(big.Int).SetBytes([]byte{1 + 27})

	return r, s, v, nil
}
