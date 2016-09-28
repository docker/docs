package idauth

import (
	"container/heap"
	"sync"
	"time"

	"github.com/docker/orca/enzi/jose"
)

// A KeyCache is used to store public keys used to verify identity tokens.
type KeyCache interface {
	Start()
	Stop()
	Get(keyID string) *jose.PublicKey
	Set(key *jose.PublicKey, ttl time.Duration)
}

type cachedKey struct {
	expiration time.Time
	key        *jose.PublicKey
}

// A keyCache is a composite data-structure: a map for looking up keys by ID
// in constant time and a min-heap for efficiently removing those entries
// which have expired.
type keyCache struct {
	sync.Mutex

	ttlTimer   *time.Timer
	done       chan struct{}
	cachedKeys []cachedKey
	lookup     map[string]int
}

var (
	_ heap.Interface = (*keyCache)(nil)
	_ KeyCache       = (*keyCache)(nil)
)

// NewInMemoryKeyCache creates a KeyCache which caches public keys in memory.
// Keys will automatically expire when they are supposed to. There is no limit
// to the number of keys which can be cached - as many as can fit in memory.
func NewInMemoryKeyCache() KeyCache {
	return &keyCache{
		ttlTimer:   time.NewTimer(0),
		done:       make(chan struct{}),
		cachedKeys: make([]cachedKey, 0, 1),
		lookup:     map[string]int{},
	}
}

func (kc *keyCache) Len() int {
	return len(kc.cachedKeys)
}

func (kc *keyCache) Less(i, j int) bool {
	return kc.cachedKeys[i].expiration.Before(kc.cachedKeys[j].expiration)
}

func (kc *keyCache) Swap(i, j int) {
	// Start by swapping the values.
	entryI, entryJ := kc.cachedKeys[j], kc.cachedKeys[i] // Note indexes.
	kc.cachedKeys[i], kc.cachedKeys[j] = entryI, entryJ

	// Set the new lookup indexes.
	kc.lookup[entryI.key.ID] = i
	kc.lookup[entryJ.key.ID] = j
}

func (kc *keyCache) Push(x interface{}) {
	currentLen := len(kc.cachedKeys)
	currentCap := cap(kc.cachedKeys)

	// Reallocate slice if necessary.
	if currentCap < currentLen+1 {
		copied := make([]cachedKey, currentLen, 2*currentCap)
		copy(copied, kc.cachedKeys)
		kc.cachedKeys = copied
	}

	entry := x.(cachedKey)
	kc.cachedKeys = append(kc.cachedKeys, entry)
	kc.lookup[entry.key.ID] = currentLen
}

func (kc *keyCache) Pop() interface{} {
	currentLen := len(kc.cachedKeys) - 1
	currentCap := cap(kc.cachedKeys)

	entry := (kc.cachedKeys)[currentLen]
	kc.cachedKeys = (kc.cachedKeys)[:currentLen]
	delete(kc.lookup, entry.key.ID)

	// Reallocate slice if necassary.
	if currentLen > 1 && currentLen < currentCap/4 {
		copied := make([]cachedKey, currentLen, currentCap/2)
		copy(copied, kc.cachedKeys)
		kc.cachedKeys = copied
	}

	return entry
}

func (kc *keyCache) Get(keyID string) *jose.PublicKey {
	kc.Lock()
	defer kc.Unlock()

	i, exists := kc.lookup[keyID]
	if !exists {
		return nil
	}

	return kc.cachedKeys[i].key
}

func (kc *keyCache) Set(key *jose.PublicKey, ttl time.Duration) {
	kc.Lock()
	defer kc.Unlock()

	var initialTimerExpiration time.Time
	if len(kc.cachedKeys) > 0 {
		initialTimerExpiration = kc.cachedKeys[0].expiration
	}

	if i, exists := kc.lookup[key.ID]; exists {
		// The key is already in the cache. Just update the expiration.
		kc.cachedKeys[i].expiration = time.Now().Add(ttl)
		heap.Fix(kc, i)
	} else {
		heap.Push(kc, cachedKey{
			expiration: time.Now().Add(ttl),
			key:        key,
		})
	}

	if initialTimerExpiration != kc.cachedKeys[0].expiration {
		kc.resetTimer()
	}
}

func (kc *keyCache) resetTimer() {
	if len(kc.cachedKeys) > 0 {
		// Reset the TTL Timer to the new most recent expiration.
		newTTL := kc.cachedKeys[0].expiration.Sub(time.Now())
		kc.ttlTimer.Reset(newTTL)
	}
}

func (kc *keyCache) Start() {
	go kc.timerLoop()
}

func (kc *keyCache) Stop() {
	close(kc.done)
	if kc.ttlTimer != nil {
		kc.ttlTimer.Stop()
	}
}

func (kc *keyCache) timerLoop() {
	for {
		select {
		case <-kc.done:
			return
		case <-kc.ttlTimer.C:
			kc.handleExpire()
		}
	}
}

func (kc *keyCache) handleExpire() {
	kc.Lock()
	defer kc.Unlock()

	if len(kc.cachedKeys) == 0 {
		return // Nothing to do.
	}

	expiration := kc.cachedKeys[0].expiration
	if expiration.Before(time.Now()) {
		heap.Remove(kc, 0)
	}

	kc.resetTimer()
}
