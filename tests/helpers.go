package tests

import (
	"encoding/base64"
	"net/url"
)

func getURL(urlStr string) string {
	aURL, _ := url.Parse(urlStr)
	return aURL.Scheme + "://" + aURL.Host + aURL.Path
}

func checkQuery(urlStr, key, val string) bool {
	aURL, _ := url.Parse(urlStr)
	q := aURL.Query()
	if q.Get(key) == val {
		return true
	}

	return false
}

func existQuery(urlStr, key string) bool {
	aURL, _ := url.Parse(urlStr)
	q := aURL.Query()
	if q.Get(key) != "" {
		return true
	}

	return false
}

func getQuery(urlStr, key string) string {
	aURL, _ := url.Parse(urlStr)
	q := aURL.Query()
	return q.Get(key)
}

func auth(clientID, clientSecret string) string {
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(clientID+":"+clientSecret))
}
