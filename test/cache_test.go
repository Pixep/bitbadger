package bitbadger

import (
	"bytes"
	"testing"
	"time"

	"github.com/Pixep/bitbadger/internal/bitbadger"
)

func TestSetTestPolicy(t *testing.T) {
	newPolicy := bitbadger.CachePolicy{
		ValidityDuration: 66 * time.Minute,
	}
	bitbadger.SetCachePolicy(newPolicy)

	currentPolicy := bitbadger.GetCachePolicy()

	if newPolicy != currentPolicy {
		t.Errorf("Set/GetCachePolicy: Cache policies differ")
	}
}

func TestRequestCaching(t *testing.T) {
	bitbadger.ClearCache()
	bitbadger.SetCachePolicy(bitbadger.CachePolicy{
		ValidityDuration: 2 * time.Minute,
		MaxCachedResults: 100,
	})

	request1 := bitbadger.BadgeRequest{
		Username:   "test",
		Repository: "repo",
		Type:       bitbadger.OpenPRType,
	}
	request2 := bitbadger.BadgeRequest{
		Username:   "test",
		Repository: "repo2",
		Type:       bitbadger.OpenPRType,
	}
	request1bis := bitbadger.BadgeRequest{
		Username:   "test",
		Repository: "repo",
		Type:       bitbadger.OpenPRType,
	}

	imageData := []byte("some-image")
	imageExtension := "jpeg"
	result := bitbadger.BadgeImage{
		Data:      imageData,
		Extension: imageExtension,
	}

	bitbadger.CacheRequestResult(request1, &result)

	if !bitbadger.RequestCached(request1) {
		t.Errorf("RequestCached: Request 1 should be cached")
	}
	if !bitbadger.RequestCached(request1bis) {
		t.Errorf("RequestCached: Request 1 bis should be cached")
	}
	if bitbadger.RequestCached(request2) {
		t.Errorf("RequestCached: Request 2 should not be cached")
	}

	cacheResult := bitbadger.GetCachedResult(request1)
	if cacheResult == nil {
		t.Errorf("GetCachedResult: Cache result should be valid")
	}
	if bitbadger.GetCachedResult(request1bis) == nil {
		t.Errorf("GetCachedResult: Cache result should be valid for request 1 bis")
	}
	if bitbadger.GetCachedResult(request2) != nil {
		t.Errorf("GetCachedResult: Other request should not be cached")
	}

	if bytes.Compare(cacheResult.Data, imageData) != 0 {
		t.Errorf("GetCachedResult: Invalid image data")
	}

	if cacheResult.Extension != imageExtension {
		t.Errorf("GetCachedResult: Invalid image extension")
	}
}

func TestCacheValidityDuration(t *testing.T) {
	cacheValidityDuration := 2 * time.Second

	bitbadger.ClearCache()
	bitbadger.SetCachePolicy(bitbadger.CachePolicy{
		ValidityDuration: cacheValidityDuration,
		MaxCachedResults: 100,
	})

	request := bitbadger.BadgeRequest{
		Username:   "test",
		Repository: "repo",
		Type:       bitbadger.OpenPRType,
	}

	bitbadger.CacheRequestResult(request, &bitbadger.BadgeImage{})

	// Wait for cache to become invalid
	time.Sleep(cacheValidityDuration + 1*time.Second)

	if bitbadger.RequestCached(request) {
		t.Errorf("RequestCached: Cached request should not be valid")
	}

	if bitbadger.GetCachedResult(request) != nil {
		t.Errorf("GetCachedResult: Cached result should not be valid")
	}
}

func TestCacheMaxCount(t *testing.T) {
	bitbadger.ClearCache()
	bitbadger.SetCachePolicy(bitbadger.CachePolicy{
		ValidityDuration: 10 * time.Minute,
		MaxCachedResults: 2,
	})

	request1 := bitbadger.BadgeRequest{
		Username: "request1",
	}
	request2 := bitbadger.BadgeRequest{
		Username: "request2",
	}
	request3 := bitbadger.BadgeRequest{
		Username: "request3",
	}

	bitbadger.CacheRequestResult(request1, &bitbadger.BadgeImage{})
	bitbadger.CacheRequestResult(request2, &bitbadger.BadgeImage{})

	if !bitbadger.RequestCached(request1) {
		t.Errorf("RequestCached: Request1 should be cached")
	}
	if !bitbadger.RequestCached(request2) {
		t.Errorf("RequestCached: Request2 should be cached")
	}
	if bitbadger.RequestCached(request3) {
		t.Errorf("RequestCached: Request3 should be not cached")
	}

	bitbadger.CacheRequestResult(request3, &bitbadger.BadgeImage{})

	if bitbadger.RequestCached(request1) {
		t.Errorf("RequestCached: Request1 should not be cached anymore")
	}
	if !bitbadger.RequestCached(request2) {
		t.Errorf("RequestCached: Request2 should be cached")
	}
	if !bitbadger.RequestCached(request3) {
		t.Errorf("RequestCached: Request3 should be cached")
	}

	bitbadger.CacheRequestResult(request2, &bitbadger.BadgeImage{})

	if bitbadger.RequestCached(request1) {
		t.Errorf("RequestCached: Request1 should not be cached anymore")
	}
	if !bitbadger.RequestCached(request2) {
		t.Errorf("RequestCached: Request2 should be cached")
	}
	if !bitbadger.RequestCached(request3) {
		t.Errorf("RequestCached: Request3 should be cached")
	}

	bitbadger.CacheRequestResult(request1, &bitbadger.BadgeImage{})

	if !bitbadger.RequestCached(request1) {
		t.Errorf("RequestCached: Request1 should be cached")
	}
	if !bitbadger.RequestCached(request2) {
		t.Errorf("RequestCached: Request2 should be cached")
	}
	if bitbadger.RequestCached(request3) {
		t.Errorf("RequestCached: Request3 should not be cached anymore")
	}
}
