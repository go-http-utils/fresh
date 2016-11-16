package fresh

import (
	"net/http"
	"testing"
	"time"

	"github.com/go-http-utils/headers"
	"github.com/stretchr/testify/suite"
)

type FreshSuite struct {
	suite.Suite

	reqHeader http.Header
	resHeader http.Header
}

func (s *FreshSuite) SetupTest() {
	s.reqHeader = make(http.Header)
	s.resHeader = make(http.Header)
}

func (s FreshSuite) TestEtagEmpty() {
	s.False(IsFresh(s.reqHeader, s.resHeader))
}

func (s FreshSuite) TestEtagMatch() {
	s.reqHeader.Set(headers.IfNoneMatch, "foo")
	s.resHeader.Set(headers.ETag, "foo")

	s.True(IsFresh(s.reqHeader, s.resHeader))
}

func (s FreshSuite) TestEtagMismatch() {
	s.reqHeader.Set(headers.IfNoneMatch, "foo")
	s.resHeader.Set(headers.ETag, "bar")

	s.False(IsFresh(s.reqHeader, s.resHeader))
}

func (s FreshSuite) TestEtagMissing() {
	s.reqHeader.Set(headers.IfNoneMatch, "foo")

	s.False(IsFresh(s.reqHeader, s.resHeader))
}

func (s FreshSuite) TestWeakEtagMatch() {
	s.reqHeader.Set(headers.IfNoneMatch, `W/"foo"`)
	s.resHeader.Set(headers.ETag, `W/"foo"`)

	s.True(IsFresh(s.reqHeader, s.resHeader))
}

func (s FreshSuite) TestEtagStrongMatch() {
	s.reqHeader.Set(headers.IfNoneMatch, `W/"foo"`)
	s.resHeader.Set(headers.ETag, `"foo"`)

	s.True(IsFresh(s.reqHeader, s.resHeader))
}

func (s FreshSuite) TestStaleOnEtagWeakMatch() {
	s.reqHeader.Set(headers.IfNoneMatch, `"foo"`)
	s.resHeader.Set(headers.ETag, `W/"foo"`)

	s.False(IsFresh(s.reqHeader, s.resHeader))
}

func (s FreshSuite) TestEtagAsterisk() {
	s.reqHeader.Set(headers.IfNoneMatch, "*")
	s.resHeader.Set(headers.ETag, `"foo"`)

	s.True(IsFresh(s.reqHeader, s.resHeader))
}

func (s FreshSuite) TestModifiedFresh() {
	s.reqHeader.Set(headers.IfModifiedSince, getFormattedTime(4*time.Second))
	s.resHeader.Set(headers.LastModified, getFormattedTime(2*time.Second))

	s.True(IsFresh(s.reqHeader, s.resHeader))
}

func (s FreshSuite) TestModifiedStale() {
	s.reqHeader.Set(headers.IfModifiedSince, getFormattedTime(2*time.Second))
	s.resHeader.Set(headers.LastModified, getFormattedTime(4*time.Second))

	s.False(IsFresh(s.reqHeader, s.resHeader))
}

func (s FreshSuite) TestEmptyLastModified() {
	s.reqHeader.Set(headers.IfModifiedSince, getFormattedTime(4*time.Second))

	s.False(IsFresh(s.reqHeader, s.resHeader))
}

func (s FreshSuite) TestBoshAndModifiedFresh() {
	s.reqHeader.Set(headers.IfNoneMatch, "foo")
	s.reqHeader.Set(headers.IfModifiedSince, getFormattedTime(4*time.Second))

	s.resHeader.Set(headers.ETag, "bar")
	s.resHeader.Set(headers.LastModified, getFormattedTime(2*time.Second))

	s.True(IsFresh(s.reqHeader, s.resHeader))
}

func (s FreshSuite) TestBoshAndETagFresh() {
	s.reqHeader.Set(headers.IfNoneMatch, "foo")
	s.reqHeader.Set(headers.IfModifiedSince, getFormattedTime(2*time.Second))

	s.resHeader.Set(headers.ETag, "foo")
	s.resHeader.Set(headers.LastModified, getFormattedTime(4*time.Second))

	s.True(IsFresh(s.reqHeader, s.resHeader))
}

func (s FreshSuite) TestBoshFresh() {
	s.reqHeader.Set(headers.IfNoneMatch, "foo")
	s.reqHeader.Set(headers.IfModifiedSince, getFormattedTime(4*time.Second))

	s.resHeader.Set(headers.ETag, "foo")
	s.resHeader.Set(headers.LastModified, getFormattedTime(2*time.Second))

	s.True(IsFresh(s.reqHeader, s.resHeader))
}

func (s FreshSuite) TestBoshStale() {
	s.reqHeader.Set(headers.IfNoneMatch, "foo")
	s.reqHeader.Set(headers.IfModifiedSince, getFormattedTime(2*time.Second))

	s.resHeader.Set(headers.ETag, "bar")
	s.resHeader.Set(headers.LastModified, getFormattedTime(4*time.Second))

	s.False(IsFresh(s.reqHeader, s.resHeader))
}

func TestFresh(t *testing.T) {
	suite.Run(t, new(FreshSuite))
}

func getFormattedTime(d time.Duration) string {
	return time.Now().Add(d).Format(http.TimeFormat)
}
