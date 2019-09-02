package bitbadger

import (
	"testing"
	"time"
)

func TestGetBadgeType(t *testing.T) {
	cases := []struct {
		in       string
		expected BadgeType
	}{
		{string(OpenPRCountType), OpenPRCountType},
		{string(OpenPRAverageAgeType), OpenPRAverageAgeType},
		{string(OldestOpenPRAge), OldestOpenPRAge},
		{string(AveragePRMergeTime), AveragePRMergeTime},
	}

	for _, c := range cases {
		badge, err := GetBadgeType(c.in)
		if badge != c.expected {
			t.Errorf("Incorrect BadgeType '%s' from string %s", badge, c.in)
		}
		if err != nil {
			t.Errorf("Should not generate an error")
		}
	}

	_, err := GetBadgeType("unknown")
	if err == nil {
		t.Errorf("Should generate an error")
	}
}

func TestGenerateBadgeInfo(t *testing.T) {
	cases := []struct {
		inType          BadgeType
		inInfo          PullRequestsInfo
		expectedLabel   string
		expectedMessage string
		expectedColor   string
	}{
		{OpenPRCountType, PullRequestsInfo{OpenCount: 999},
			"Open PRs", "999", "red"},
		{OpenPRAverageAgeType, PullRequestsInfo{
			OpenAverageTime: 5 * time.Minute},
			"Avg. current PRs age", "5 mins", "green"},
		{OldestOpenPRAge, PullRequestsInfo{
			OldestOpenPR: 5 * time.Minute},
			"Oldest PR age", "5 mins", "green"},
		{AveragePRMergeTime, PullRequestsInfo{
			AveragePRMergeTime: 5 * time.Minute},
			"Avg. PR merge time", "5 mins", "green"},
	}

	for _, c := range cases {
		badgeInfo, _ := GenerateBadgeInfo(c.inType, c.inInfo)

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

	_, err := GenerateBadgeInfo("invalid", PullRequestsInfo{})
	if err == nil {
		t.Errorf("Should generate an error")
	}
}
