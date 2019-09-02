package bitbadger

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
)

type bbPullRequest struct {
	Title     string `json:"title"`
	ID        int    `json:"id"`
	CreatedOn string `json:"created_on"`
	UpdatedOn string `json:"updated_on"`
}

type bbPullRequestsReponse struct {
	PullRequestsCount int             `json:"size"`
	PullRequests      []bbPullRequest `json:"values"`
	Page              int             `json:"page"`
	PageLength        int             `json:"pagelen"`
	NextPageURL       string          `json:"next"`
	PreviousPageURL   string          `json:"previous"`
}

type openPRInfo struct {
	OpenCount       int
	OldestOpenPR    time.Duration
	OpenAverageTime time.Duration
}

type mergedPRInfo struct {
	AveragePRMergeTime time.Duration
}

// RetrieveBBPullRequestInfo retrieves information relative to pull requests
// from BitBucket Cloud.
func RetrieveBBPullRequestInfo(request BadgeRequest) (PullRequestsInfo, error) {
	openPRInfo, err := retrieveBBOpenPRInfo(request)
	if err != nil {
		return PullRequestsInfo{}, err
	}

	mergedPRInfo, err := retrieveBBMergedPRInfo(request)
	if err != nil {
		return PullRequestsInfo{}, err
	}

	return PullRequestsInfo{
		OpenCount:          openPRInfo.OpenCount,
		OldestOpenPR:       openPRInfo.OldestOpenPR,
		OpenAverageTime:    openPRInfo.OpenAverageTime,
		AveragePRMergeTime: mergedPRInfo.AveragePRMergeTime,
	}, nil
}

func queryBB(request BadgeRequest, endpoint string) ([]byte, error) {
	sourceServerRequest := "https://api.bitbucket.org/2.0/repositories/"
	sourceServerRequest += request.Username + "/" + request.Repository
	sourceServerRequest += endpoint

	req, err := http.NewRequest("GET", sourceServerRequest, nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(config.Username, config.Password)

	resp, err := client.Do(req)
	if err != nil {
		log.Error("Get request failed: ", err)
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		log.Error("Non-200 response:\n", body)
		return nil, err
	}

	log.Debug("BitBucket response:")
	log.Debug(string(body))

	return body, nil
}

func retrieveBBOpenPRInfo(request BadgeRequest) (openPRInfo, error) {
	body, err := queryBB(request, "/pullrequests?state=OPEN")
	if err != nil {
		return openPRInfo{}, err
	}

	var response bbPullRequestsReponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Error("Answer decoding failed for:")
		log.Error(body)
		return openPRInfo{}, err
	}

	openPRTotalTime := time.Duration(0)
	oldestOpenPRAge := time.Duration(0)

	now := time.Now()
	prsWithValidTime := 0
	for _, pullRequest := range response.PullRequests {
		createdOnTime, err := time.Parse(time.RFC3339, pullRequest.CreatedOn)
		if err != nil {
			log.Error("Failed to parse time:", pullRequest.CreatedOn)
		} else {
			openTime := now.Sub(createdOnTime)

			if openTime > oldestOpenPRAge {
				oldestOpenPRAge = openTime
			}

			openPRTotalTime += openTime
			prsWithValidTime++
		}
	}

	openPRAverageTime := time.Duration(0)
	if prsWithValidTime > 0 {
		openPRAverageTime = time.Duration(
			openPRTotalTime.Minutes()/float64(prsWithValidTime)) * time.Minute
	}

	return openPRInfo{
		OpenCount:       response.PullRequestsCount,
		OldestOpenPR:    oldestOpenPRAge,
		OpenAverageTime: openPRAverageTime,
	}, nil
}

func retrieveBBMergedPRInfo(request BadgeRequest) (mergedPRInfo, error) {
	body, err := queryBB(request, "/pullrequests?state=MERGED")
	if err != nil {
		return mergedPRInfo{}, err
	}

	var response bbPullRequestsReponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Error("Answer decoding failed for:")
		log.Error(body)
		return mergedPRInfo{}, err
	}

	mergedPRTotalTime := time.Duration(0)
	mergedPRConsidered := 0

	for _, pullRequest := range response.PullRequests {
		// TODO: Use the activity feed of each PR instead. This value will be
		// incorrect if the PR is updated after it has been merged.
		createdOnTime, err := time.Parse(time.RFC3339, pullRequest.CreatedOn)
		updatedOnTime, err := time.Parse(time.RFC3339, pullRequest.UpdatedOn)
		if err != nil {
			log.Error("Failed to parse time:", pullRequest.CreatedOn)
		} else {
			openTime := updatedOnTime.Sub(createdOnTime)
			mergedPRTotalTime += openTime
			mergedPRConsidered++
		}
	}

	averagePRMergeTime := time.Duration(0)
	if mergedPRConsidered > 0 {
		averagePRMergeTime = time.Duration(
			mergedPRTotalTime.Minutes()/float64(mergedPRConsidered)) * time.Minute
	}

	return mergedPRInfo{
		AveragePRMergeTime: averagePRMergeTime,
	}, nil
}
