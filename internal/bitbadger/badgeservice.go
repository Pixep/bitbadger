package bitbadger

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	log "github.com/Sirupsen/logrus"
)

func generateBadgeURL(badge BadgeInfo) string {
	// Label, message and color are '-' separate in shields.io format.
	badgetInfoURL := fmt.Sprintf("%s-%s-%s", badge.Label, badge.Message, badge.Color)
	// Use ReplaceAll to have "%20" in place of spaces, as Golang encode uses "+" instead
	return "https://img.shields.io/badge/" + strings.ReplaceAll(badgetInfoURL, " ", "%20")
}

// DownloadBadge downloads and returns a badge image from "img.shields.io",
// using badgeInfo.
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
		Extension: "svg+xml",
	}
	return image, nil
}
