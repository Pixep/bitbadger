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
	OpenPRType        BadgeType = "open-pr-count"
	AveragePRTimeType BadgeType = "avg-pr-time"
)

// GenerateBadgeInfo generates a badge from a type
// and pull request information
func GenerateBadgeInfo(badgeType BadgeType, prInfo PullRequestsInfo) (BadgeInfo, error) {
	switch badgeType {
	case OpenPRType:
		return generateOpenPRCountBadge(prInfo), nil
	case AveragePRTimeType:
		return generateAveragePRTimeBadge(prInfo), nil
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

func generateAveragePRTimeBadge(prInfo PullRequestsInfo) (badge BadgeInfo) {
	avgOpenTime := prInfo.OpenAverageTime
	badge = BadgeInfo{
		Label:   "Avg PR time",
		Message: printDuration(avgOpenTime),
	}

	switch {
	case avgOpenTime < 24*time.Hour:
		badge.Color = "green"
	case avgOpenTime < 48*time.Hour:
		badge.Color = "yellowgreen"
	case avgOpenTime < 72*time.Hour:
		badge.Color = "yellow"
	case avgOpenTime < 96*time.Hour:
		badge.Color = "orange"
	default:
		badge.Color = "red"
	}

	return
}
