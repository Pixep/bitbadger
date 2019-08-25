package bitbadger

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	log "github.com/Sirupsen/logrus"
)

// BadgeInfo holds information required
// to generate a badge image
type BadgeInfo struct {
	Label   string
	Message string
	Color   string
}

// BadgeImage holds the badge image data and extension
type BadgeImage struct {
	Data      []byte
	Extension string
}

// Label, message and color are '-' separated
func generateBadgeURL(badge BadgeInfo) string {
	badgetInfoURL := fmt.Sprintf("%s-%s-%s", badge.Label, badge.Message, badge.Color)
	// Use ReplaceAll to have "%20" in place of spaces, as Golang encode uses "+" instead
	return "https://img.shields.io/badge/" + strings.ReplaceAll(badgetInfoURL, " ", "%20")
}

// DownloadBadge downloads and returns a badge image
// from "img.shields.io", using badgeInfo.
func DownloadBadge(badgeInfo BadgeInfo) (*BadgeImage, error) {
	// Get the data
	badgeURL := generateBadgeURL(badgeInfo)
	log.Debug("Badge URL = ", badgeURL)

	resp, err := http.Get(badgeURL)
	if err != nil {
		log.Error("Error while retrieving badge at '", badgeURL, "': ", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("Failed to read badge image response: ", err)
		return nil, err
	}

	image := &BadgeImage{
		Data:      body,
		Extension: "svg",
	}
	return image, nil
}
