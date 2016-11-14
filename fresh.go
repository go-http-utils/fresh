package fresh

import (
	"net/http"
	"strings"
	"time"
)

const (
	// Version is this package's verison
	Version = "0.0.1"

	timeForm = "Mon, 14 Nov 2016 07:53:39 GMT"
)

// IsFresh check whether cache can be used in this HTTP request
func IsFresh(reqHeader http.Header, resHeader http.Header) bool {
	isEtagMatched, isModifiedMatched := false, false

	ifModifiedSince := reqHeader.Get("if-modified-since")
	ifNoneMatch := reqHeader.Get("if-none-match")
	cacheControl := reqHeader.Get("cache-control")

	etag := reqHeader.Get("etag")
	lastModified := reqHeader.Get("last-modified")

	if ifModifiedSince == "" && ifNoneMatch == "" {
		return false
	}

	if strings.Contains(cacheControl, "no-cache") {
		return false
	}

	if etag != "" {
		isEtagMatched = checkEtagMatch(trimTags(strings.Split(ifNoneMatch, ",")), etag)
	}

	if lastModified != "" && ifModifiedSince != "" {
		isModifiedMatched = checkModifedMatch(lastModified, ifModifiedSince)
	}

	return isEtagMatched || isModifiedMatched
}

func trimTags(tags []string) []string {
	trimedTags := make([]string, len(tags))

	for i, tag := range tags {
		trimedTags[i] = strings.TrimSpace(tag)
	}

	return trimedTags
}

func checkEtagMatch(etagsToMatch []string, etag string) bool {
	for _, etagToMatch := range etagsToMatch {
		if etagToMatch == "*" || etagToMatch == etag || etagToMatch == "W/"+etag {
			return true
		}
	}

	return false
}

func checkModifedMatch(lastModified, ifModifiedSince string) bool {
	if lm, err := time.Parse(lastModified, timeForm); err != nil {
		if ims, err := time.Parse(ifModifiedSince, timeForm); err != nil {
			return lm.Before(ims)
		}
	}

	return false
}
