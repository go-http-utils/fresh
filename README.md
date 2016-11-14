# fresh
[![Build Status](https://travis-ci.org/go-http-utils/fresh.svg?branch=master)](https://travis-ci.org/go-http-utils/fresh)
[![Coverage Status](https://coveralls.io/repos/github/go-http-utils/fresh/badge.svg?branch=master)](https://coveralls.io/github/go-http-utils/fresh?branch=master)

HTTP response freshness testing for Go

## Installation

```sh
go get -u github.com/go-http-utils/fresh
```

## Documentation

API documentation can be found here: https://godoc.org/github.com/go-http-utils/fresh

## Usage

```go
import (
  "net/http"

  "github.com/go-http-utils/fresh"
)
```

```go
reqHeader, resHeader := make(http.Header), make(http.Header)

reqHeader.Set("if-none-match", "foo")
resHeader.Set("etag", "bar")

fresh.IsFresh(reqHeader, resHeader)
// -> false
```

```go
reqHeader, resHeader := make(http.Header), make(http.Header)

reqHeader.Set("if-modified-since", "Mon, 14 Nov 2016 22:05:49 GMT")
resHeader.Set("last-modified", "Mon, 14 Nov 2016 22:05:47 GMT")

fresh.IsFresh(reqHeader, resHeader)
// -> true
```
