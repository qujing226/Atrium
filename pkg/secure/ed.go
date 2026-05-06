package secure

import (
	"crypto/ed25519"
	"crypto/rand"
	"errors"
	"fmt"
)

// SignKeyPair 封装 Ed25519 签名密钥对
type SignKeyPair struct {
	pk ed25519.PublicKey
	sk ed25519.PrivateKey
}

// NewSignKeyPair 生成新的签名密钥对
func NewSignKeyPair() (*SignKeyPair, error) {
	pk, sk, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to generate ed25519 key: %w", err)
	}
	return &SignKeyPair{pk: pk, sk: sk}, nil
}

// LoadSignKeyFromBytes 加载现有的密钥
func LoadSignKeyFromBytes(pkBytes, skBytes []byte) (*SignKeyPair, error) {
	kp := &SignKeyPair{}

	if len(pkBytes) > 0 {
		if len(pkBytes) != ed25519.PublicKeySize {
			return nil, errors.New("invalid public key size")
		}
		kp.pk = make([]byte, ed25519.PublicKeySize)
		copy(kp.pk, pkBytes)
	}

	if len(skBytes) > 0 {
		if len(skBytes) != ed25519.PrivateKeySize {
			return nil, errors.New("invalid private key size")
		}
		kp.sk = make([]byte, ed25519.PrivateKeySize)
		copy(kp.sk, skBytes)
	}

	return kp, nil
}

// Export 导出密钥字节
func (kp *SignKeyPair) Export() ([]byte, []byte) {
	// 复制一份以防外部修改
	pkCopy := make([]byte, len(kp.pk))
	copy(pkCopy, kp.pk)
	
	var skCopy []byte
	if kp.sk != nil {
		skCopy = make([]byte, len(kp.sk))
		copy(skCopy, kp.sk)
	}
	
	return pkCopy, skCopy
}

// Sign 对消息进行签名
func (kp *SignKeyPair) Sign(message []byte) ([]byte, error) {
	if kp.sk == nil {
		return nil, errors.New("private key not available")
	}
	return ed25519.Sign(kp.sk, message), nil
}

// Verify 验证签名
// 这是一个静态方法或者是由于它只需要公钥，也可以作为 method
func Verify(pubKey []byte, message, signature []byte) bool {
	if len(pubKey) != ed25519.PublicKeySize {
		return false
	}
	return ed25519.Verify(pubKey, message, signature)
}
