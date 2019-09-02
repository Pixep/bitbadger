package bitbadger

import (
	"testing"
	"time"

	"github.com/Pixep/bitbadger/internal/bitbadger"
)

func TestGetBadgeType(t *testing.T) {
	cases := []struct {
		in       string
		expected bitbadger.BadgeType
	}{
		{string(bitbadger.OpenPRCountType), bitbadger.OpenPRCountType},
		{string(bitbadger.OpenPRAverageAgeType), bitbadger.OpenPRAverageAgeType},
		{string(bitbadger.OldestOpenPRAge), bitbadger.OldestOpenPRAge},
		{string(bitbadger.AveragePRMergeTime), bitbadger.AveragePRMergeTime},
	}

	for _, c := range cases {
		badge, err := bitbadger.GetBadgeType(c.in)
		if badge != c.expected {
			t.Errorf("Incorrect BadgeType '%s' from string %s", badge, c.in)
		}
		if err != nil {
			t.Errorf("Should not generate an error")
		}
	}

	_, err := bitbadger.GetBadgeType("unknown")
	if err == nil {
		t.Errorf("Should generate an error")
	}
}

func TestGenerateBadgeInfo(t *testing.T) {
	cases := []struct {
		inType          bitbadger.BadgeType
		inInfo          bitbadger.PullRequestsInfo
		expectedLabel   string
		expectedMessage string
		expectedColor   string
	}{
		{bitbadger.OpenPRCountType, bitbadger.PullRequestsInfo{OpenCount: 999},
			"Open PRs", "999", "red"},
		{bitbadger.OpenPRAverageAgeType, bitbadger.PullRequestsInfo{
			OpenAverageTime: 5 * time.Minute},
			"Avg. current PRs age", "5 mins", "green"},
		{bitbadger.OldestOpenPRAge, bitbadger.PullRequestsInfo{
			OldestOpenPR: 5 * time.Minute},
			"Oldest PR age", "5 mins", "green"},
		{bitbadger.AveragePRMergeTime, bitbadger.PullRequestsInfo{
			AveragePRMergeTime: 5 * time.Minute},
			"Avg. PR merge time", "5 mins", "green"},
	}

	for _, c := range cases {
		badgeInfo, _ := bitbadger.GenerateBadgeInfo(c.inType, c.inInfo)

		if badgeInfo.Label != c.expectedLabel {
			t.Errorf("Incorrect label for OpenPRCountType: %s", badgeInfo.Label)
		}
		if badgeInfo.Message != c.expectedMessage {
			t.Errorf("Incorrect message for OpenPRCountType: %s", badgeInfo.Message)
		}
		if badgeInfo.Color != c.expectedColor {
			t.Errorf("Incorrect color for OpenPRCountType: %s", badgeInfo.Color)
		}
	}

	_, err := bitbadger.GenerateBadgeInfo("invalid", bitbadger.PullRequestsInfo{})
	if err == nil {
		t.Errorf("Should generate an error")
	}
}
