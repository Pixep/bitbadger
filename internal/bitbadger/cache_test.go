package bitbadger

import (
	"bytes"
	"testing"
	"time"
)

func TestSetTestPolicy(t *testing.T) {
	newPolicy := CachePolicy{
		ValidityDuration: 66 * time.Minute,
	}
	SetCachePolicy(newPolicy)

	currentPolicy := GetCachePolicy()

	if newPolicy != currentPolicy {
		t.Errorf("Set/GetCachePolicy: Cache policies differ")
	}
}

func TestRequestCaching(t *testing.T) {
	ClearCache()
	SetCachePolicy(CachePolicy{
		ValidityDuration: 2 * time.Minute,
		MaxCachedResults: 100,
	})

	request1 := BadgeRequest{
		Username:   "test",
		Repository: "repo",
		Type:       OpenPRCountType,
	}
	request2 := BadgeRequest{
		Username:   "test",
		Repository: "repo2",
		Type:       OpenPRCountType,
	}
	request1bis := BadgeRequest{
		Username:   "test",
		Repository: "repo",
		Type:       OpenPRCountType,
	}

	imageData := []byte("some-image")
	imageExtension := "jpeg"
	result := BadgeImage{
		Data:      imageData,
		Extension: imageExtension,
	}

	CacheRequestResult(request1, &result)

	if !RequestCached(request1) {
		t.Errorf("RequestCached: Request 1 should be cached")
	}
	if !RequestCached(request1bis) {
		t.Errorf("RequestCached: Request 1 bis should be cached")
	}
	if RequestCached(request2) {
		t.Errorf("RequestCached: Request 2 should not be cached")
	}

	cacheResult := GetCachedResult(request1)
	if cacheResult == nil {
		t.Errorf("GetCachedResult: Cache result should be valid")
	}
	if GetCachedResult(request1bis) == nil {
		t.Errorf("GetCachedResult: Cache result should be valid for request 1 bis")
	}
	if GetCachedResult(request2) != nil {
		t.Errorf("GetCachedResult: Other request should not be cached")
	}

	if !bytes.Equal(cacheResult.Data, imageData) {
		t.Errorf("GetCachedResult: Invalid image data")
	}

	if cacheResult.Extension != imageExtension {
		t.Errorf("GetCachedResult: Invalid image extension")
	}
}

func TestCacheValidityDuration(t *testing.T) {
	cacheValidityDuration := 2 * time.Second

	ClearCache()
	SetCachePolicy(CachePolicy{
		ValidityDuration: cacheValidityDuration,
		MaxCachedResults: 100,
	})

	request := BadgeRequest{
		Username:   "test",
		Repository: "repo",
		Type:       OpenPRCountType,
	}

	CacheRequestResult(request, &BadgeImage{})

	// Wait for cache to become invalid
	time.Sleep(cacheValidityDuration + 1*time.Second)

	if RequestCached(request) {
		t.Errorf("RequestCached: Cached request should not be valid")
	}

	if GetCachedResult(request) != nil {
		t.Errorf("GetCachedResult: Cached result should not be valid")
	}
}

func TestCacheMaxCount(t *testing.T) {
	ClearCache()
	SetCachePolicy(CachePolicy{
		ValidityDuration: 10 * time.Minute,
		MaxCachedResults: 2,
	})

	request1 := BadgeRequest{
		Username: "request1",
	}
	request2 := BadgeRequest{
		Username: "request2",
	}
	request3 := BadgeRequest{
		Username: "request3",
	}

	CacheRequestResult(request1, &BadgeImage{})
	CacheRequestResult(request2, &BadgeImage{})

	if !RequestCached(request1) {
		t.Errorf("RequestCached: Request1 should be cached")
	}
	if !RequestCached(request2) {
		t.Errorf("RequestCached: Request2 should be cached")
	}
	if RequestCached(request3) {
		t.Errorf("RequestCached: Request3 should be not cached")
	}

	CacheRequestResult(request3, &BadgeImage{})

	if RequestCached(request1) {
		t.Errorf("RequestCached: Request1 should not be cached anymore")
	}
	if !RequestCached(request2) {
		t.Errorf("RequestCached: Request2 should be cached")
	}
	if !RequestCached(request3) {
		t.Errorf("RequestCached: Request3 should be cached")
	}

	CacheRequestResult(request2, &BadgeImage{})

	if RequestCached(request1) {
		t.Errorf("RequestCached: Request1 should not be cached anymore")
	}
	if !RequestCached(request2) {
		t.Errorf("RequestCached: Request2 should be cached")
	}
	if !RequestCached(request3) {
		t.Errorf("RequestCached: Request3 should be cached")
	}

	CacheRequestResult(request1, &BadgeImage{})

	if !RequestCached(request1) {
		t.Errorf("RequestCached: Request1 should be cached")
	}
	if !RequestCached(request2) {
		t.Errorf("RequestCached: Request2 should be cached")
	}
	if RequestCached(request3) {
		t.Errorf("RequestCached: Request3 should not be cached anymore")
	}
}
