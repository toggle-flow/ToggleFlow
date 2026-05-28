package api

import (
	"sync"
	"time"

	"toggleflow/internal/db"
)

// flagCache is a two-level in-memory store: projectID → environmentID → JSON bytes.
// Busted explicitly after any write that changes flag data for an environment.
type flagCache struct {
	mu    sync.RWMutex
	store map[int64]map[int64][]byte
}

func newFlagCache() *flagCache {
	return &flagCache{store: make(map[int64]map[int64][]byte)}
}

func (c *flagCache) get(projectID, envID int64) []byte {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if envs, ok := c.store[projectID]; ok {
		return envs[envID]
	}
	return nil
}

func (c *flagCache) set(projectID, envID int64, data []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.store[projectID] == nil {
		c.store[projectID] = make(map[int64][]byte)
	}
	c.store[projectID][envID] = data
}

// bust removes the cache for one environment — used when a single env changes
// (flag toggle, rules update).
func (c *flagCache) bust(projectID, envID int64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if envs, ok := c.store[projectID]; ok {
		delete(envs, envID)
	}
}

// bustProject removes all cached entries for a project — used when a change
// affects every environment (flag renamed, flag deleted, variations changed).
func (c *flagCache) bustProject(projectID int64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.store, projectID)
}

// sdkKeyEntry caches a validated SDK key → environment mapping.
// TTL is capped at the key's own expiry so we never serve an expired key from cache.
type sdkKeyEntry struct {
	env       *db.Environment
	expiresAt time.Time
}

// sdkKeyCache maps hashed SDK key → cached environment, with a 5-minute TTL.
type sdkKeyCache struct {
	mu    sync.RWMutex
	store map[string]sdkKeyEntry
}

func newSDKKeyCache() *sdkKeyCache {
	return &sdkKeyCache{store: make(map[string]sdkKeyEntry)}
}

func (c *sdkKeyCache) get(keyHash string) *db.Environment {
	c.mu.RLock()
	entry, ok := c.store[keyHash]
	c.mu.RUnlock()
	if !ok || time.Now().After(entry.expiresAt) {
		return nil
	}
	return entry.env
}

func (c *sdkKeyCache) set(keyHash string, env *db.Environment, keyExpiry *time.Time) {
	ttl := time.Now().Add(5 * time.Minute)
	// Never cache past the key's own expiry
	if keyExpiry != nil && keyExpiry.Before(ttl) {
		ttl = *keyExpiry
	}
	c.mu.Lock()
	c.store[keyHash] = sdkKeyEntry{env: env, expiresAt: ttl}
	c.mu.Unlock()
}
