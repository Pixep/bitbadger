package bitbadger

import (
	"time"
)

// CachePolicy holds information used to check if a cache entry is still valid.
type CachePolicy struct {
	ValidityDuration time.Duration
	MaxCachedResults int
}

// CacheEntry holds the request, result and last time refreshed.
type CacheEntry struct {
	Request     BadgeRequest
	ImageResult *BadgeImage
	RefreshTime time.Time
}

var cachePolicy CachePolicy
var cacheMap map[BadgeRequest]CacheEntry

func init() {
	cachePolicy = CachePolicy{
		ValidityDuration: 10 * time.Minute,
		MaxCachedResults: 100,
	}

	ClearCache()
}

// ClearCache clears the full cache content.
func ClearCache() {
	cacheMap = make(map[BadgeRequest]CacheEntry)
}

// SetCachePolicy sets the global cache policy.
func SetCachePolicy(policy CachePolicy) {
	cachePolicy = policy
}

// GetCachePolicy returns the current global test policy.
func GetCachePolicy() CachePolicy {
	return cachePolicy
}

// CacheRequestResult caches the result to a request and sets it as refreshed
// "Now()".
func CacheRequestResult(request BadgeRequest, image *BadgeImage) {
	// Don't cache request if cache is disabled
	if cachePolicy.ValidityDuration == 0 {
		return
	}

	cacheMap[request] = CacheEntry{
		Request:     request,
		ImageResult: image,
		RefreshTime: time.Now(),
	}

	cleanupCache()
}

func cleanupCache() {
	if cachePolicy.MaxCachedResults <= 0 && len(cacheMap) > 0 {
		ClearCache()
	}

	// Find and remove oldest entry, until below threshold.
	// Will definitely need a better implementation :)
	for len(cacheMap) > cachePolicy.MaxCachedResults {
		var oldestRequest BadgeRequest
		var oldestRefreshTime time.Time
		init := false
		for request, entry := range cacheMap {
			if !init || oldestRefreshTime.After(entry.RefreshTime) {
				init = true
				oldestRefreshTime = entry.RefreshTime
				oldestRequest = request
			}
		}

		delete(cacheMap, oldestRequest)
	}
}

// cacheEntryValid returns true if the cache entry is valid.
func cacheEntryValid(entry *CacheEntry) bool {
	if entry == nil {
		return false
	}

	return time.Now().Sub(entry.RefreshTime) < cachePolicy.ValidityDuration
}

// RequestCached returns true if the request is cached and valid.
func RequestCached(request BadgeRequest) bool {
	entry, cached := cacheMap[request]
	return cached && cacheEntryValid(&entry)
}

// GetCachedResult returns the cached result for the request or nil if the
// request is not cached, or if the cached result is not valid.
func GetCachedResult(request BadgeRequest) *BadgeImage {
	entry, cached := cacheMap[request]
	if !cached || !cacheEntryValid(&entry) {
		return nil
	}

	return entry.ImageResult
}
