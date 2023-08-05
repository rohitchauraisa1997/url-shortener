package models

import "time"

type UrlAnalyticDetails struct {
	URL  string `json:"url"`
	Hits uint64 `json:"urlHits"`
	TTL  uint64 `json:"ttl"`
}

type ShortenedUrlAndDetail struct {
	ShortenedUrl  string             `json:"shortenedUrl"`
	UrlsAnalytics UrlAnalyticDetails `json:"urlsAnalytics"`
}

// helps for sorting the result in descending order of hits
// (hits refers to the amount of time shortenedUrl was used to access that particular URL).
type ShortenedUrlAndDetailsSlice []ShortenedUrlAndDetail

func (s ShortenedUrlAndDetailsSlice) Len() int {
	return len(s)
}

func (s ShortenedUrlAndDetailsSlice) Less(i, j int) bool {
	return s[i].UrlsAnalytics.Hits > s[j].UrlsAnalytics.Hits
}

func (s ShortenedUrlAndDetailsSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

type UrlHitsAndCreatedAtInfo struct {
	CreatedAt time.Time `json:"createdAt"`
	Urlhits   uint64    `json:"urlhits"`
}
