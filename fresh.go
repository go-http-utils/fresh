package fresh

import (
	"net/http"
	"strings"
	"time"

	"github.com/go-http-utils/headers"
)

// Version is this package's verison
const Version = "0.3.0"

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
	if lm, err := time.Parse(http.TimeFormat, lastModified); err == nil {
		if ims, err := time.Parse(http.TimeFormat, ifModifiedSince); err == nil {
			return lm.Before(ims)
		}
	}

	return false
}

func checkUnmodifedMatch(lastModified, ifUnmodifiedSince string) bool {
	if lm, err := time.Parse(http.TimeFormat, lastModified); err == nil {
		if ius, err := time.Parse(http.TimeFormat, ifUnmodifiedSince); err == nil {
			return lm.After(ius)
		}
	}

	return false
}
