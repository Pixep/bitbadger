package bitbadger

import (
	"errors"
	"time"
)

// BadgeRequest holds the information relative
// to a client badge generation request
type BadgeRequest struct {
	Username   string
	Repository string
	Type       BadgeType
}

// PullRequestsInfo holds the pull request data
// used to generate the badges.
type PullRequestsInfo struct {
	OpenCount       int
	OpenAverageTime time.Duration
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
