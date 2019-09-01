package bitbadger

import (
	"errors"
	"fmt"
	"time"
)

// BadgeRequest holds the information relative
// to a client badge generation request
type BadgeRequest struct {
	Username   string
	Repository string
	Type       BadgeType
}

func compare(a, b BadgeRequest) bool {
	if a.Username != b.Username {
		return false
	}

	if a.Repository != b.Repository {
		return false
	}

	if a.Type != b.Type {
		return false
	}

	return true
}

// PullRequestsInfo holds the pull request data
// used to generate the badges.
type PullRequestsInfo struct {
	OpenCount       int
	OldestOpenPR    time.Duration
	OpenAverageTime time.Duration
}

func (info PullRequestsInfo) String() string {
	return fmt.Sprintf("%d", info.OpenCount)
}

type RepositoryType int

const (
	BitBucketCloud RepositoryType = iota
)

// RetrievePullRequestInfo retrieves information relative to pull requests
// from a specific repository
func RetrievePullRequestInfo(repoType RepositoryType, request BadgeRequest) (PullRequestsInfo, error) {
	switch repoType {
	case BitBucketCloud:
		return RetrieveBBPullRequestInfo(request)
	default:
		return PullRequestsInfo{}, errors.New("Invalid repository type")
	}
}
