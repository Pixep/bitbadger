package bitbadger

import (
	"errors"
	log "github.com/Sirupsen/logrus"
)

// BadgeInfo holds information required to generate a badge image.
type BadgeInfo struct {
	Label   string
	Message string
	Color   string
}

// BadgeImage holds the badge image data and extension.
type BadgeImage struct {
	Data      []byte
	Extension string
}

// GenerateBadge generates a badge from a BadgeRequest.
func GenerateBadge(request BadgeRequest) (*BadgeImage, error) {
	prInfo, err := RetrieveBBPullRequestInfo(request)
	if err != nil {
		log.Error("Error while retrieving badge info: ", err)
		return nil, errors.New("Error while getting pull request info from the upstream server")
	}

	badge, err := GenerateBadgeInfo(request.Type, prInfo)
	if err != nil {
		log.Error("Failed to generate badge: ", err)
		return nil, errors.New("Failed to generate badge")
	}

	badgeImage, err := DownloadBadge(badge)
	if err != nil {
		log.Error("Error downloading badge: ", err)
		return nil, errors.New("Failed to download badge")
	}

	return badgeImage, nil
}
