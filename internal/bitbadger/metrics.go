package bitbadger

import (
	"errors"
	"strconv"
	"time"
)

// BadgeType represents the type of badge / metric to
// show
type BadgeType string

const (
	OpenPRCountType   BadgeType = "open-pr-count"
	AveragePRTimeType BadgeType = "avg-pr-time"
	OldestOpenPRTime  BadgeType = "oldest-pr-time"
)

// GenerateBadgeInfo generates a badge from a type
// and pull request information
func GenerateBadgeInfo(badgeType BadgeType, prInfo PullRequestsInfo) (BadgeInfo, error) {
	switch badgeType {
	case OpenPRCountType:
		return generateOpenPRCountBadge(prInfo), nil
	case AveragePRTimeType:
		return generateAveragePRTimeBadge(prInfo), nil
	case OldestOpenPRTime:
		return generateOldestOpenPRTimeBadge(prInfo), nil
	default:
		return BadgeInfo{}, errors.New("Invalid badge type")
	}
}

func generateOpenPRCountBadge(prInfo PullRequestsInfo) (badge BadgeInfo) {
	badge = BadgeInfo{
		Label:   "Open PRs",
		Message: strconv.Itoa(prInfo.OpenCount),
	}

	switch {
	case prInfo.OpenCount <= 3:
		badge.Color = "green"
	case prInfo.OpenCount <= 5:
		badge.Color = "yellowgreen"
	case prInfo.OpenCount <= 7:
		badge.Color = "yellow"
	case prInfo.OpenCount <= 9:
		badge.Color = "orange"
	default:
		badge.Color = "red"
	}

	return
}

func generateAveragePRTimeBadge(prInfo PullRequestsInfo) BadgeInfo {
	return BadgeInfo{
		Label:   "Avg PR time",
		Message: printDuration(prInfo.OpenAverageTime),
		Color:   prOpenTimeColor(prInfo.OpenAverageTime),
	}
}

func generateOldestOpenPRTimeBadge(prInfo PullRequestsInfo) BadgeInfo {
	return BadgeInfo{
		Label:   "Oldest open PR",
		Message: printDuration(prInfo.OldestOpenPR),
		Color:   prOpenTimeColor(prInfo.OldestOpenPR),
	}
}

func prOpenTimeColor(openTime time.Duration) string {
	switch {
	case openTime < 24*time.Hour:
		return "green"
	case openTime < 48*time.Hour:
		return "yellowgreen"
	case openTime < 72*time.Hour:
		return "yellow"
	case openTime < 96*time.Hour:
		return "orange"
	default:
		return "red"
	}
}
