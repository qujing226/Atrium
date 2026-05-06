package blockchain

import (
	"bytes"
	"errors"
	"sync"
	"time"
)

// Chain defines the interface for interacting with the DID registry
type Chain interface {
	RegisterDidDoc(did string, doc []byte)
	ResolveDidDoc(did string) ([]byte, error)
}

// SimulatedChain implements Chain with an in-memory map and simulated latency
// 用于模拟“慢速”的区块链网络
type SimulatedChain struct {
	didDocs map[string][]byte
	mu      sync.RWMutex
	latency time.Duration
}

func NewSimulatedChain(latency time.Duration) *SimulatedChain {
	return &SimulatedChain{
		didDocs: make(map[string][]byte),
		latency: latency,
	}
}

func (m *SimulatedChain) RegisterDidDoc(did string, doc []byte) {
	m.mu.Lock()
	defer m.mu.Unlock()
	// Simulate consensus time
	time.Sleep(m.latency)
	m.didDocs[did] = doc
}

func (m *SimulatedChain) ResolveDidDoc(did string) ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	// Simulate network/query latency
	if m.latency > 0 {
		time.Sleep(m.latency)
	}
	
	doc, ok := m.didDocs[did]
	if !ok {
		return nil, errors.New("did doc not found on chain")
	}
	// Return a copy to prevent modification
	out := make([]byte, len(doc))
	copy(out, doc)
	return out, nil
}

// VerificationCallback 当后台验证发现缓存过期/欺诈时触发
// 参数: did, 缓存的旧文档, 链上的新文档
type VerificationCallback func(did string, cachedDoc, freshDoc []byte)

// OptimisticCache 实现“乐观验证缓存”策略
type OptimisticCache struct {
	chain      Chain
	cache      map[string][]byte
	mu         sync.RWMutex
	onMismatch VerificationCallback
}

func NewOptimisticCache(chain Chain, onMismatch VerificationCallback) *OptimisticCache {
	return &OptimisticCache{
		chain:      chain,
		cache:      make(map[string][]byte),
		onMismatch: onMismatch,
	}
}

func (c *OptimisticCache) RegisterDidDoc(did string, doc []byte) {
	// Write-through: 先写链，再写缓存
	c.chain.RegisterDidDoc(did, doc)
	
	c.mu.Lock()
	c.cache[did] = doc
	c.mu.Unlock()
}

func (c *OptimisticCache) Resolve(did string) ([]byte, error) {
	c.mu.RLock()
	cachedDoc, ok := c.cache[did]
	c.mu.RUnlock()

	if ok {
		// Fast Path: 立即返回缓存数据 (0-RTT)
		// 并在后台启动“乐观验证”
		go c.verifyInBackground(did, cachedDoc)
		return cachedDoc, nil
	}

	// Slow Path: 缓存未命中，不得不查链
	doc, err := c.chain.ResolveDidDoc(did)
	if err != nil {
		return nil, err
	}

	// 更新缓存
	c.mu.Lock()
	c.cache[did] = doc
	c.mu.Unlock()

	return doc, nil
}

func (c *OptimisticCache) verifyInBackground(did string, cachedDoc []byte) {
	// 这是一个耗时操作，在后台运行
	freshDoc, err := c.chain.ResolveDidDoc(did)
	if err != nil {
		// 链不可达？暂时保留缓存，或者根据策略决定是否报警
		return 
	}

	if !bytes.Equal(cachedDoc, freshDoc) {
		// Critical: 发现缓存不一致（DID文档已更新或被撤销）
		// 1. 立即更新缓存
		c.mu.Lock()
		c.cache[did] = freshDoc
		c.mu.Unlock()

		// 2. 触发回调（通知上层应用断开连接或重新握手）
		if c.onMismatch != nil {
			c.onMismatch(did, cachedDoc, freshDoc)
		}
	}
}
