package fresh

import (
	"net/http"
	"strings"
	"time"

	"github.com/go-http-utils/headers"
)

// Version is this package's verison
const Version = "0.3.1"

// IsFresh check whether cache can be used in this HTTP request
func IsFresh(reqHeader http.Header, resHeader http.Header) bool {
	isEtagMatched, isModifiedMatched := false, false

	ifModifiedSince := reqHeader.Get(headers.IfModifiedSince)
	ifUnmodifiedSince := reqHeader.Get(headers.IfUnmodifiedSince)
	ifNoneMatch := reqHeader.Get(headers.IfNoneMatch)
	cacheControl := reqHeader.Get(headers.CacheControl)

	etag := resHeader.Get(headers.ETag)
	lastModified := resHeader.Get(headers.LastModified)

	if ifModifiedSince == "" && ifUnmodifiedSince == "" && ifNoneMatch == "" {
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

	if lastModified != "" && ifUnmodifiedSince != "" {
		isModifiedMatched = checkUnmodifedMatch(lastModified, ifUnmodifiedSince)
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
	if lm, ims, ok := parseTimePairs(lastModified, ifModifiedSince); ok == true {
		return lm.Before(ims)
	}

	return false
}

func checkUnmodifedMatch(lastModified, ifUnmodifiedSince string) bool {
	if lm, ius, ok := parseTimePairs(lastModified, ifUnmodifiedSince); ok == true {
		return lm.After(ius)
	}

	return false
}

func parseTimePairs(s1, s2 string) (t1 time.Time, t2 time.Time, ok bool) {
	if t1, err := time.Parse(http.TimeFormat, s1); err == nil {
		if t2, err := time.Parse(http.TimeFormat, s2); err == nil {
			return t1, t2, true
		}
	}

	return t1, t2, false
}
