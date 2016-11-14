package fresh

import (
	"net/http"
	"testing"
	"time"

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
	s.reqHeader.Set(HeaderIfNoneMatch, "foo")
	s.resHeader.Set(HeaderETag, "foo")

	s.True(IsFresh(s.reqHeader, s.resHeader))
}

func (s FreshSuite) TestEtagMismatch() {
	s.reqHeader.Set(HeaderIfNoneMatch, "foo")
	s.resHeader.Set(HeaderETag, "bar")

	s.False(IsFresh(s.reqHeader, s.resHeader))
}

func (s FreshSuite) TestEtagMissing() {
	s.reqHeader.Set(HeaderIfNoneMatch, "foo")

	s.False(IsFresh(s.reqHeader, s.resHeader))
}

func (s FreshSuite) TestWeakEtagMatch() {
	s.reqHeader.Set(HeaderIfNoneMatch, `W/"foo"`)
	s.resHeader.Set(HeaderETag, `W/"foo"`)

	s.True(IsFresh(s.reqHeader, s.resHeader))
}

func (s FreshSuite) TestEtagStrongMatch() {
	s.reqHeader.Set(HeaderIfNoneMatch, `W/"foo"`)
	s.resHeader.Set(HeaderETag, `"foo"`)

	s.True(IsFresh(s.reqHeader, s.resHeader))
}

func (s FreshSuite) TestStaleOnEtagWeakMatch() {
	s.reqHeader.Set(HeaderIfNoneMatch, `"foo"`)
	s.resHeader.Set(HeaderETag, `W/"foo"`)

	s.False(IsFresh(s.reqHeader, s.resHeader))
}

func (s FreshSuite) TestEtagAsterisk() {
	s.reqHeader.Set(HeaderIfNoneMatch, "*")
	s.resHeader.Set(HeaderETag, `"foo"`)

	s.True(IsFresh(s.reqHeader, s.resHeader))
}

func (s FreshSuite) TestModifiedFresh() {
	s.reqHeader.Set(HeaderIfModifiedSince, getFormattedTime(4*time.Second))
	s.resHeader.Set(HeaderLastModified, getFormattedTime(2*time.Second))

	s.True(IsFresh(s.reqHeader, s.resHeader))
}

func (s FreshSuite) TestModifiedStale() {
	s.reqHeader.Set(HeaderIfModifiedSince, getFormattedTime(2*time.Second))
	s.resHeader.Set(HeaderLastModified, getFormattedTime(4*time.Second))

	s.False(IsFresh(s.reqHeader, s.resHeader))
}

func (s FreshSuite) TestEmptyLastModified() {
	s.reqHeader.Set(HeaderIfModifiedSince, getFormattedTime(4*time.Second))

	s.False(IsFresh(s.reqHeader, s.resHeader))
}

func (s FreshSuite) TestBoshAndModifiedFresh() {
	s.reqHeader.Set(HeaderIfNoneMatch, "foo")
	s.reqHeader.Set(HeaderIfModifiedSince, getFormattedTime(4*time.Second))

	s.resHeader.Set(HeaderETag, "bar")
	s.resHeader.Set(HeaderLastModified, getFormattedTime(2*time.Second))

	s.True(IsFresh(s.reqHeader, s.resHeader))
}

func (s FreshSuite) TestBoshAndETagFresh() {
	s.reqHeader.Set(HeaderIfNoneMatch, "foo")
	s.reqHeader.Set(HeaderIfModifiedSince, getFormattedTime(2*time.Second))

	s.resHeader.Set(HeaderETag, "foo")
	s.resHeader.Set(HeaderLastModified, getFormattedTime(4*time.Second))

	s.True(IsFresh(s.reqHeader, s.resHeader))
}

func (s FreshSuite) TestBoshFresh() {
	s.reqHeader.Set(HeaderIfNoneMatch, "foo")
	s.reqHeader.Set(HeaderIfModifiedSince, getFormattedTime(4*time.Second))

	s.resHeader.Set(HeaderETag, "foo")
	s.resHeader.Set(HeaderLastModified, getFormattedTime(2*time.Second))

	s.True(IsFresh(s.reqHeader, s.resHeader))
}

func (s FreshSuite) TestBoshStale() {
	s.reqHeader.Set(HeaderIfNoneMatch, "foo")
	s.reqHeader.Set(HeaderIfModifiedSince, getFormattedTime(2*time.Second))

	s.resHeader.Set(HeaderETag, "bar")
	s.resHeader.Set(HeaderLastModified, getFormattedTime(4*time.Second))

	s.False(IsFresh(s.reqHeader, s.resHeader))
}

func TestFresh(t *testing.T) {
	suite.Run(t, new(FreshSuite))
}

func getFormattedTime(d time.Duration) string {
	return time.Now().Add(d).Format(http.TimeFormat)
}
