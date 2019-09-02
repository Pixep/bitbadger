package bitbadger

import (
	"errors"
	"fmt"
	"time"
)

// BadgeRequest holds the information relative to a client badge generation
// request.
type BadgeRequest struct {
	Username   string
	Repository string
	Type       BadgeType
}

// PullRequestsInfo holds the pull request data used to generate the badges.
type PullRequestsInfo struct {
	OpenCount          int
	OldestOpenPR       time.Duration
	OpenAverageTime    time.Duration
	AveragePRMergeTime time.Duration
}

func (info PullRequestsInfo) String() string {
	return fmt.Sprintf("%d", info.OpenCount)
}

// RepositoryType holds the type of repository service targetted.
type RepositoryType int

const (
	// BitBucketCloud represents BitBucket Cloud service.
	BitBucketCloud RepositoryType = iota
)

// RetrievePullRequestInfo retrieves information relative to pull requests
// from a specific repository.
func RetrievePullRequestInfo(repoType RepositoryType, request BadgeRequest) (PullRequestsInfo, error) {
	switch repoType {
	case BitBucketCloud:
		return RetrieveBBPullRequestInfo(request)
	default:
		return PullRequestsInfo{}, errors.New("Invalid repository type")
	}
}
