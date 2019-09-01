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
		newBadgeImage, err := GenerateBadge(*request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}

		CacheRequestResult(*request, badgeImage)
		badgeImage = newBadgeImage
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

func sendHTTPReponse(w http.ResponseWriter, badgeImage *BadgeImage) {
	w.Header().Set("Content-Type", "image/"+badgeImage.Extension)
	fmt.Fprintf(w, "%s", badgeImage.Data)
}
