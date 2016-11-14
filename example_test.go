package fresh_test

import (
	"fmt"
	"net/http"

	"github.com/go-http-utils/fresh"
)

func ExampleFresh_IsFresh() {
	reqHeader, resHeader := make(http.Header), make(http.Header)

	reqHeader.Set("if-none-match", "foo")
	resHeader.Set("etag", "bar")

	fmt.Println(fresh.IsFresh(reqHeader, resHeader))
	// -> false

	reqHeader, resHeader = make(http.Header), make(http.Header)

	reqHeader.Set("if-modified-since", "Mon, 14 Nov 2016 22:05:49 GMT")
	resHeader.Set("last-modified", "Mon, 14 Nov 2016 22:05:47 GMT")

	fmt.Println(fresh.IsFresh(reqHeader, resHeader))
	// -> true
}
