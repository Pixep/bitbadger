package bitbadger

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
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

// GetBadgeType returns a BadgeType from a string, and an
// error if there is no corresponding BadgeType.
func GetBadgeType(badgeString string) (BadgeType, error) {
	badgeType := BadgeType(badgeString)
	if BadgeTypeValid(badgeType) {
		return badgeType, nil
	}

	return badgeType, errors.New("Invalid badge type '" + badgeString + "'." +
		"Badge type can be 'open-pr-count', 'avg-pr-time', or 'oldest-pr-time'")
}

// BadgeTypeValid returns true if the BadgeType provided is
// valid, false otherwise.
func BadgeTypeValid(badgeType BadgeType) bool {
	switch badgeType {
	case OpenPRCountType, AveragePRTimeType, OldestOpenPRTime:
		return true
	default:
		return false
	}
}

// GenerateBadgeInfo generates a badge from a type
// and pull request information.
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

func printDuration(duration time.Duration) string {
	days := int64(duration.Hours() / 24)
	hours := int64(math.Mod(duration.Hours(), 24))
	minutes := int64(math.Mod(duration.Minutes(), 60))
	seconds := int64(math.Mod(duration.Seconds(), 60))

	chunks := []struct {
		singularName string
		amount       int64
	}{
		{"day", days},
		{"hour", hours},
		{"minute", minutes},
		{"second", seconds},
	}

	parts := []string{}

	printedChunks := 0
	for _, chunk := range chunks {
		switch chunk.amount {
		case 0:
			continue
		case 1:
			parts = append(parts, fmt.Sprintf("%d %s", chunk.amount, chunk.singularName))
			printedChunks++
		default:
			parts = append(parts, fmt.Sprintf("%d %ss", chunk.amount, chunk.singularName))
			printedChunks++
		}

		if printedChunks >= 2 {
			break
		}
	}

	return strings.Join(parts, " ")
}
