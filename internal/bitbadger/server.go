package bitbadger

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
)

var client = &http.Client{}

// ServeWithHTTP starts the HTTP bitbadger server on the specificed port
func ServeWithHTTP(port int) error {
	http.HandleFunc("/", httpHandler)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))
	return nil
}

// ServeWithHTTPS starts the HTTPS bitbadger server on the specificed port
func ServeWithHTTPS(port int, certFile, keyFile string) error {
	http.HandleFunc("/", httpHandler)
	log.Fatal(http.ListenAndServeTLS(":"+strconv.Itoa(port), certFile, keyFile, nil))
	return nil
}

func httpHandler(w http.ResponseWriter, r *http.Request) {
	paths := strings.Split(r.URL.Path, "/")
	paths = paths[1:]

	if len(paths) < 3 {
		errorMessage := "Requires a request of the form: '<username>/<repository-slug>/<type>'\n"
		errorMessage += "where <type> can be 'openpr', or 'avgprtime'"
		log.Warn("Invalid request with: ", r.URL)
		http.Error(w, errorMessage, http.StatusBadRequest)
		return
	}

	badgeType := BadgeType(strings.TrimSuffix(paths[2], ".svg"))
	switch badgeType {
	case OpenPRCountType, AveragePRTimeType:
	default:
		errorMessage := "Invalid badge type '" + paths[2] + "'\n"
		errorMessage += "Type can be can be 'open-pr-count', or 'avg-pr-time'"
		log.Warn("Invalid request with: ", r.URL)
		http.Error(w, errorMessage, http.StatusBadRequest)
		return
	}
	request := BadgeRequest{
		Username:   paths[0],
		Repository: paths[1],
		Type:       badgeType,
	}

	log.Info("Creating badge for ", request.Username, "/", request.Repository, "/", request.Type)

	badgeImage := GetCachedResult(request)
	if badgeImage == nil {
		newBadgeImage, httpError := generateNewBadge(request)
		if httpError != nil {
			http.Error(w, httpError.Message, httpError.HTTPErrorStatus)
			return
		}

		CacheRequestResult(request, newBadgeImage)
		badgeImage = newBadgeImage
	}

	w.Header().Set("Content-Type", "image/"+badgeImage.Extension)
	fmt.Fprintf(w, "%s", badgeImage.Data)
}

func (info PullRequestsInfo) String() string {
	return fmt.Sprintf("%d", info.OpenCount)
}

type serverError struct {
	Message         string
	HTTPErrorStatus int
}

func generateNewBadge(request BadgeRequest) (*BadgeImage, *serverError) {
	prInfo, err := RetrieveBBPullRequestInfo(request)
	if err != nil {
		log.Error("Error while retrieving badge info: ", err)
		return nil, &serverError{
			Message:         "Error while getting pull request info from the upstream server",
			HTTPErrorStatus: http.StatusBadGateway,
		}
	}

	badge, err := GenerateBadgeInfo(request.Type, prInfo)
	if err != nil {
		log.Error("Failed to generate badge: ", err)
		return nil, &serverError{
			Message:         "Failed to generate badge",
			HTTPErrorStatus: http.StatusInternalServerError,
		}
	}

	badgeImage, err := DownloadBadge(badge)
	if err != nil {
		log.Error("Error downloading badge: ", err)
		return nil, &serverError{
			Message:         "Failed to download badge",
			HTTPErrorStatus: http.StatusBadGateway,
		}
	}

	return badgeImage, nil
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
