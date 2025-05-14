package cache

import (
	"sync"
	"time"

	model "grab-bootcamp-be-team13-2025/internal/domain/models"
)

type rssCache struct {
	Data      []model.Outbreak
	ExpiresAt time.Time
}

var (
	cacheLock sync.Mutex
	cacheData rssCache
)

// SetCache sets the RSS cache with expiration
func SetCache(data []model.Outbreak, ttl time.Duration) {
	cacheLock.Lock()
	defer cacheLock.Unlock()
	cacheData = rssCache{
		Data:      data,
		ExpiresAt: time.Now().Add(ttl),
	}
}

// GetCache returns cached data if valid, else nil
func GetCache() ([]model.Outbreak, bool) {
	cacheLock.Lock()
	defer cacheLock.Unlock()
	if time.Now().Before(cacheData.ExpiresAt) {
		return cacheData.Data, true
	}
	return nil, false
}
