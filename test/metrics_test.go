package bitbadger

import (
	"testing"
	"time"

	"github.com/Pixep/bitbadger/internal/bitbadger"
)

func TestGenerateBadgeInfo(t *testing.T) {
	badgeInfo, _ := bitbadger.GenerateBadgeInfo(bitbadger.OpenPRCountType, bitbadger.PullRequestsInfo{
		OpenCount: 999,
	})

	if badgeInfo.Label != "Open PRs" {
		t.Errorf("Incorrect label for OpenPRCountType: %s", badgeInfo.Label)
	}
	if badgeInfo.Message != "999" {
		t.Errorf("Incorrect message for OpenPRCountType: %s", badgeInfo.Message)
	}
	if badgeInfo.Color != "red" {
		t.Errorf("Incorrect color for OpenPRCountType: %s", badgeInfo.Color)
	}

	badgeInfo, _ = bitbadger.GenerateBadgeInfo(bitbadger.AveragePRTimeType, bitbadger.PullRequestsInfo{
		OpenAverageTime: 5 * time.Minute,
	})

	if badgeInfo.Label != "Avg PR time" {
		t.Errorf("Incorrect label for OpenPRCountType: %s", badgeInfo.Label)
	}
	if badgeInfo.Message != "5 minutes" {
		t.Errorf("Incorrect message for OpenPRCountType: %s", badgeInfo.Message)
	}
	if badgeInfo.Color != "green" {
		t.Errorf("Incorrect color for OpenPRCountType: %s", badgeInfo.Color)
	}

	badgeInfo, _ = bitbadger.GenerateBadgeInfo(bitbadger.OldestOpenPRTime, bitbadger.PullRequestsInfo{
		OldestOpenPR: 5 * time.Minute,
	})

	if badgeInfo.Label != "Oldest open PR" {
		t.Errorf("Incorrect label for OpenPRCountType: %s", badgeInfo.Label)
	}
	if badgeInfo.Message != "5 minutes" {
		t.Errorf("Incorrect message for OpenPRCountType: %s", badgeInfo.Message)
	}
	if badgeInfo.Color != "green" {
		t.Errorf("Incorrect color for OpenPRCountType: %s", badgeInfo.Color)
	}

	_, err := bitbadger.GenerateBadgeInfo("invalid", bitbadger.PullRequestsInfo{})
	if err == nil {
		t.Errorf("Should generate an error")
	}
}
