package schnorrmpc

import (
	"crypto/ecdsa"
	"math/big"
)

var bigZero = big.NewInt(0)

var bigOne = big.NewInt(1)

// Structure definition for polynomial
// only save self Polynomial
type Polynomial []big.Int

// polynomial commit
type PolynomialG []ecdsa.PublicKey
type PolynomialGSig []big.Int

// key: 	smIndex
// value: 	polyCommitG
type PolyGMap  		map[uint16]PolynomialG
type PolyGSigMap  	map[uint16]PolynomialGSig