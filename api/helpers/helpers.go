package helpers

import (
	"os"
	"sort"
	"strings"

	"github.com/rohitchauraisa1997/url-shortener/models"
)

func RemoveDomainError(url string) bool {
	if url == os.Getenv("DOMAIN") {
		return false
	}
	newURL := strings.Replace(url, "http://", "", 1)
	newURL = strings.Replace(newURL, "https://", "", 1)
	newURL = strings.Replace(newURL, "www.", "", 1)

	return newURL == os.Getenv("DOMAIN")
}

func EnforceHTTP(url string) string {
	if url[:4] != "http" {
		return "http://" + url
	}
	return url
}

func SortResponse(resp map[string]models.UrlAnalyticDetails) models.ShortenedUrlAndDetailsSlice {

	// Convert map to a slice for sorting
	var items models.ShortenedUrlAndDetailsSlice
	for key, value := range resp {
		items = append(items, models.ShortenedUrlAndDetail{ShortenedUrl: key, UrlsAnalytics: value})
	}

	// Sort the slice in descending order based on Hits
	sort.Sort(items)

	return items
}
