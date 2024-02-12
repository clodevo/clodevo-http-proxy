package utils

import (
	"sync"
	"time"
)

type AuthCacheItem struct {
	TenantName string
	Expiry     time.Time
}

type AuthCache struct {
	items map[string]AuthCacheItem
	mutex *sync.Mutex
}

func NewAuthCache() *AuthCache {
	return &AuthCache{
		items: make(map[string]AuthCacheItem),
		mutex: &sync.Mutex{},
	}
}

// Check looks up the cache for the given tenantName and apiKey.
// It returns true if the credentials are valid and found in the cache, along with the tenant name.
// If the credentials are not found or expired, it returns false.
func (c *AuthCache) Check(tenantName, apiKey string) (bool, string) {
	cacheKey := tenantName + ":" + apiKey

	c.mutex.Lock()
	defer c.mutex.Unlock()

	if item, found := c.items[cacheKey]; found && item.Expiry.After(time.Now()) {
		// Cache hit and not expired
		return true, item.TenantName
	}

	// Not found in cache or expired
	return false, ""
}

// Update adds or updates the cache with the given tenantName and apiKey.
func (c *AuthCache) Update(tenantName, apiKey string) {
	cacheKey := tenantName + ":" + apiKey

	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.items[cacheKey] = AuthCacheItem{
		TenantName: tenantName,
		Expiry:     time.Now().Add(5 * time.Minute), // Adjust expiry time as needed
	}
}
