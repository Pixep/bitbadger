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

func TestCacheRequestResult(t *testing.T) {
	cacheValidityDuration := 2 * time.Second
	newPolicy := bitbadger.CachePolicy{
		ValidityDuration: cacheValidityDuration,
	}
	bitbadger.SetCachePolicy(newPolicy)

	request := bitbadger.BadgeRequest{
		Username:   "test",
		Repository: "repo",
		Type:       bitbadger.OpenPRType,
	}
	otherRequest := bitbadger.BadgeRequest{
		Username:   "test",
		Repository: "repo2",
		Type:       bitbadger.OpenPRType,
	}
	imageData := []byte("some-image")
	imageExtension := "jpeg"
	result := bitbadger.BadgeImage{
		Data:      imageData,
		Extension: imageExtension,
	}

	bitbadger.CacheRequestResult(request, &result)

	if !bitbadger.RequestCached(request) {
		t.Errorf("RequestCached: Request should be cached")
	}

	if bitbadger.RequestCached(otherRequest) {
		t.Errorf("RequestCached: The other request should not be cached")
	}

	cacheResult := bitbadger.GetCachedResult(request)
	if cacheResult == nil {
		t.Errorf("GetCachedResult: Cache result should be valid")
	}
	if bitbadger.GetCachedResult(otherRequest) != nil {
		t.Errorf("GetCachedResult: Other request should not be cached")
	}

	if bytes.Compare(cacheResult.Data, imageData) != 0 {
		t.Errorf("GetCachedResult: Invalid image data")
	}

	if cacheResult.Extension != imageExtension {
		t.Errorf("GetCachedResult: Invalid image extension")
	}

	// Wait for cache to become invalid
	time.Sleep(cacheValidityDuration + 1*time.Second)

	if bitbadger.RequestCached(request) {
		t.Errorf("RequestCached: Cached request should not be valid")
	}

	if bitbadger.GetCachedResult(request) != nil {
		t.Errorf("GetCachedResult: Cached result should not be valid")
	}
}
