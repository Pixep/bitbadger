package bitbadger

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"
)

type serverError struct {
	Message         string
	HTTPErrorStatus int
}

var client = &http.Client{}

// ServeWithHTTP starts the HTTP bitbadger server on the specificed port
func ServeWithHTTP(port int) error {
	http.HandleFunc("/", handleHTTPRequest)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))
	return nil
}

// ServeWithHTTPS starts the HTTPS bitbadger server on the specificed port
func ServeWithHTTPS(port int, certFile, keyFile string) error {
	http.HandleFunc("/", handleHTTPRequest)
	log.Fatal(http.ListenAndServeTLS(":"+strconv.Itoa(port), certFile, keyFile, nil))
	return nil
}

func handleHTTPRequest(w http.ResponseWriter, r *http.Request) {
	request, httpError := parseHTTPRequest(r)
	if httpError != nil {
		http.Error(w, httpError.Message, httpError.HTTPErrorStatus)
		return
	}

	log.Info("Creating badge for ", request.Username, "/", request.Repository, "/", request.Type)

	badgeImage := GetCachedResult(*request)
	if badgeImage == nil {
		badgeImage, httpError = generateNewBadge(*request)
		if httpError != nil {
			http.Error(w, httpError.Message, httpError.HTTPErrorStatus)
			return
		}

		CacheRequestResult(*request, badgeImage)
	}

	sendHTTPReponse(w, badgeImage)
}

func parseHTTPRequest(r *http.Request) (*BadgeRequest, *serverError) {
	paths := strings.Split(r.URL.Path, "/")
	paths = paths[1:]

	if len(paths) < 3 {
		log.Warn("Invalid request: ", r.URL)
		errorMessage := "Requires a request of the form: '<username>/<repository-slug>/<type>'"
		return nil, &serverError{
			Message:         errorMessage,
			HTTPErrorStatus: http.StatusBadRequest,
		}
	}

	badgeType, err := GetBadgeType(strings.TrimSuffix(paths[2], ".svg"))
	if err != nil {
		log.Warn("Invalid request: ", r.URL)
		return nil, &serverError{
			Message:         err.Error(),
			HTTPErrorStatus: http.StatusBadRequest,
		}
	}

	return &BadgeRequest{
		Username:   paths[0],
		Repository: paths[1],
		Type:       badgeType,
	}, nil
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

func sendHTTPReponse(w http.ResponseWriter, badgeImage *BadgeImage) {
	w.Header().Set("Content-Type", "image/"+badgeImage.Extension)
	fmt.Fprintf(w, "%s", badgeImage.Data)
}
