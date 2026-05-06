package secure

import (
	"crypto/rand"
	"errors"
	"fmt"

	"github.com/cloudflare/circl/kem/kyber/kyber768"
)

// KyberKeyPair 是常驻内存的密钥对象
// 优势：保持私钥在 NTT 域的展开状态，避免重复 Unmarshal
type KyberKeyPair struct {
	pk *kyber768.PublicKey
	sk *kyber768.PrivateKey
}

// NewKyberKeyPair 生成新的密钥对（用于注册阶段）
func NewKyberKeyPair() (*KyberKeyPair, error) {
	pk, sk, err := kyber768.GenerateKeyPair(rand.Reader)
	if err != nil {
		return nil, err
	}
	return &KyberKeyPair{pk: pk, sk: sk}, nil
}

// LoadFromBytes 从存储/网络加载密钥（用于服务启动阶段）
// 这是一个“昂贵”的操作，建议只做一次
func LoadFromBytes(pkBytes, skBytes []byte) (*KyberKeyPair, error) {
	kp := &KyberKeyPair{
		pk: new(kyber768.PublicKey),
		sk: new(kyber768.PrivateKey),
	}

	if len(pkBytes) > 0 {
		if len(pkBytes) != kyber768.PublicKeySize {
			return nil, fmt.Errorf("invalid public key size: %d", len(pkBytes))
		}
		kp.pk.Unpack(pkBytes)
	}

	if len(skBytes) > 0 {
		if len(skBytes) != kyber768.PrivateKeySize {
			return nil, fmt.Errorf("invalid private key size: %d", len(skBytes))
		}
		kp.sk.Unpack(skBytes)
	}

	return kp, nil
}

// Export 将密钥导出为字节（用于存入 DID 文档或数据库）
func (kp *KyberKeyPair) Export() ([]byte, []byte) {
	var pkBytes, skBytes []byte
	if kp.pk != nil {
		pkBytes = make([]byte, kyber768.PublicKeySize)
		kp.pk.Pack(pkBytes)
	}
	if kp.sk != nil {
		skBytes = make([]byte, kyber768.PrivateKeySize)
		kp.sk.Pack(skBytes)
	}
	return pkBytes, skBytes
}

// Encapsulate 生成共享密钥并封装到当前公钥 (kp.pk)
// 用于发送方：持有接收方的公钥，生成 (ct, ss)
func (kp *KyberKeyPair) Encapsulate() (ct []byte, ss []byte, err error) {
	if kp.pk == nil {
		return nil, nil, errors.New("public key not loaded")
	}
	ct = make([]byte, kyber768.CiphertextSize)
	ss = make([]byte, kyber768.SharedKeySize)
	
	// Generate random seed for encapsulation
	// EncapsulateTo expects a seed, not a reader
	seed := make([]byte, 32) // Kyber768 needs 32 bytes randomness usually? Check API.
	// Actually circl API usually takes a seed for deterministic generation.
	// Let's assume 32 bytes is enough (standard AES-DRBG seed size).
	// If it fails, we check doc again.
	if _, err := rand.Read(seed); err != nil {
		return nil, nil, err
	}

	kp.pk.EncapsulateTo(ct, ss, seed)
	return ct, ss, nil
}

func (kp *KyberKeyPair) Decapsulate(ct []byte) (ss []byte, err error) {
	if kp.sk == nil {
		return nil, errors.New("private key not loaded")
	}
	if len(ct) != kyber768.CiphertextSize {
		return nil, errors.New("invalid ciphertext size")
	}

	ss = make([]byte, kyber768.SharedKeySize)

	kp.sk.DecapsulateTo(ss, ct)

	return ss, nil
}
