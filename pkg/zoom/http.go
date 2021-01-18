package zoom

import (
	"log"
	"net/http"
)

// make sure these are consistent
// TODO: make configurable
const USER_AGENT string = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36"
const USER_AGENT_SHORTHAND string = "Chrome87" // todo: figure out zooms algorithm for determining this

var INITIAL_HEADERS = map[string]string{
	"pragma":                    "no-cache",
	"cache-control":             "no-cache",
	"upgrade-insecure-requests": "1",
	"user-agent":                USER_AGENT,
	"accept":                    "application/json, text/plain, */*",
	"sec-fetch-site":            "none",
	"sec-fetch-mode":            "navigate",
	"sec-fetch-user":            "?1",
	"sec-fetch-dest":            "document",
	"accept-language":           "en-US,en;q=0.9",
}

func httpGet(httpClient *http.Client, link string, headers map[string]string) (*http.Response, error) {
	log.Printf("Requesting %s [GET]", link)

	request, err := http.NewRequest("GET", link, nil)
	if err != nil {
		return nil, err
	}

	for header, headerValue := range headers {
		request.Header.Set(header, headerValue)
	}

	return httpClient.Do(request)
}
