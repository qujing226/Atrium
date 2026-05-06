package model

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/mr-tron/base58"
)

type DIDDocument struct {
	ID                 string               `json:"id"`
	Version            int                  `json:"version"`
	Revoked            bool                 `json:"revoked"`
	Created            time.Time            `json:"created"`
	Updated            time.Time            `json:"updated"`
	VerificationMethod []VerificationMethod `json:"verificationMethod"`
	Authentication     []string             `json:"authentication"`
	Proof              *Proof               `json:"proof,omitempty"`
}

type VerificationMethod struct {
	ID         string `json:"id"`
	Type       string `json:"type"` // "Ed25519VerificationKey2020" 或 "Kyber768PublicKey"
	Controller string `json:"controller,omitempty"`

	PublicKeyBase58 string `json:"publicKeyBase58,omitempty"` // for Ed25519
	PublicKeyBase64 string `json:"publicKeyBase64,omitempty"` // for Kyber
}

type Proof struct {
	Type               string    `json:"type"`
	Created            time.Time `json:"created"`
	ProofPurpose       string    `json:"proofPurpose"`
	VerificationMethod string    `json:"verificationMethod"`
	SignatureValue     string    `json:"signatureValue"` // Base58 编码的签名
}

func (doc *DIDDocument) GetKyberPubKey() ([]byte, error) {
	for _, vm := range doc.VerificationMethod {
		if vm.Type == "Kyber768PublicKey" {
			return base64.StdEncoding.DecodeString(vm.PublicKeyBase64)
		}
	}
	return nil, fmt.Errorf("kyber key not found in DID document")
}

func (doc *DIDDocument) GetSigningPubKey() ([]byte, error) {
	for _, vm := range doc.VerificationMethod {
		if vm.Type == "Ed25519VerificationKey2020" {
			return base58.Decode(vm.PublicKeyBase58)
		}
	}
	return nil, fmt.Errorf("signing key not found in DID document")
}
