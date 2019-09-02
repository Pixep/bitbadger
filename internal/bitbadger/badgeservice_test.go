package bitbadger

import (
	"testing"
)

func TestBadgeURL(t *testing.T) {
	badgeURL := generateBadgeURL(BadgeInfo{
		Label:   "My label",
		Message: "My message",
		Color:   "somecolor",
	})
	if badgeURL != "https://img.shields.io/badge/My%20label-My%20message-somecolor" {
		t.Errorf("generateBadgeURL: Invalid badge URL generated %s", badgeURL)
	}
}
