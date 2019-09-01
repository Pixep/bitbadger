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
}

type bbPullRequestsReponse struct {
	PullRequestsCount int             `json:"size"`
	PullRequests      []bbPullRequest `json:"values"`
	Page              int             `json:"page"`
	PageLength        int             `json:"pagelen"`
	NextPageURL       string          `json:"next"`
	PreviousPageURL   string          `json:"previous"`
}

// RetrieveBBPullRequestInfo retrieves information relative to pull requests
// from BitBucket Cloud
func RetrieveBBPullRequestInfo(request BadgeRequest) (PullRequestsInfo, error) {
	sourceServerRequest := "https://api.bitbucket.org/2.0/repositories/"
	sourceServerRequest += request.Username + "/" + request.Repository
	sourceServerRequest += "/pullrequests?state=OPEN"

	req, err := http.NewRequest("GET", sourceServerRequest, nil)
	if err != nil {
		return PullRequestsInfo{}, err
	}

	req.SetBasicAuth(config.Username, config.Password)

	resp, err := client.Do(req)
	if err != nil {
		log.Error("Get request failed: ", err)
		return PullRequestsInfo{}, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		log.Error("Non-200 response:\n", body)
		return PullRequestsInfo{}, err
	}

	log.Debug("BitBucket response:")
	log.Debug(string(body))

	var response bbPullRequestsReponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Error("Answer decoding failed for:")
		log.Error(body)
		return PullRequestsInfo{}, err
	}

	openPRTotalTime := time.Duration(0)
	oldestOpenPRTime := time.Duration(0)

	now := time.Now()
	prsWithValidTime := 0
	for _, pullRequest := range response.PullRequests {
		createdOnTime, err := time.Parse(time.RFC3339, pullRequest.CreatedOn)
		log.Debug(createdOnTime)
		if err != nil {
			log.Error("Failed to parse time:", pullRequest.CreatedOn)
		} else {
			openTime := now.Sub(createdOnTime)

			if openTime > oldestOpenPRTime {
				oldestOpenPRTime = openTime
			}

			openPRTotalTime += openTime
			prsWithValidTime++
		}
	}

	openPRAverageTime := time.Duration(0)
	if prsWithValidTime > 0 {
		openPRAverageTime = time.Duration(openPRTotalTime.Minutes()/float64(prsWithValidTime)) * time.Minute
	}

	prInfo := PullRequestsInfo{
		OpenCount:       response.PullRequestsCount,
		OldestOpenPR:    oldestOpenPRTime,
		OpenAverageTime: openPRAverageTime,
	}
	return prInfo, nil
}
