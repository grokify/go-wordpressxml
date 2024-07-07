WordPress XML Parser
====================

[![Build Status][build-status-svg]][build-status-url]
[![Go Report Card][goreport-svg]][goreport-url]
[![Docs][docs-godoc-svg]][docs-godoc-url]
[![License][license-svg]][license-url]

## Overview

The `go-wordpressxml` package provides a WordPress WXR (WordPress eXtended RSS) XML parser.

## Documentation

Documentation is provided using godoc and available on [GoDoc.org](https://godoc.org/github.com/grokify/go-wordpressxml).

## Installation

Installing any of the packages will install the entire library. For example:

```bash
$ go get github.com/grokify/go-wordpressxml
```

## Usage

```go
import (
	"github.com/grokify/go-wordpressxml"
)

func main() {
	wp := wordpressxml.NewWordPressXML()
	err := wp.ReadFile("myblog.wordpress.2016-08-13.xml")
	if err != nil {
		panic(err)
	}
	wp.WriteMetaCSV("articles.csv")
}
```

## Notes

Since WordPress uses `content:encoded` and `excerpt:encoded`, Go's XML built-in parser treats both of these as the field `encoded` in different namespaces. This parser retrieves these fields as an array of `encoded` and then moves the data into the `Content` property.

## Contributing

Features, Issues, and Pull Requests are always welcome.

To contribute:

1. Fork it ( http://github.com/grokify/go-wordpressxml/fork )
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new Pull Request

Please report issues and feature requests on [GitHub](https://github.com/grokify/go-wordpressxml).

 [build-status-svg]: https://github.com/grokify/go-wordpressxml/workflows/test/badge.svg
 [build-status-url]: https://github.com/grokify/go-wordpressxml/actions
 [build-status-svg]: https://api.travis-ci.org/grokify/go-wordpressxml.svg?branch=master
 [build-status-url]: https://travis-ci.org/grokify/go-wordpressxml
 [goreport-svg]: https://goreportcard.com/badge/github.com/grokify/go-wordpressxml
 [goreport-url]: https://goreportcard.com/report/github.com/grokify/go-wordpressxml
 [docs-godoc-svg]: https://pkg.go.dev/badge/github.com/grokify/go-wordpressxml
 [docs-godoc-url]: https://pkg.go.dev/github.com/grokify/go-wordpressxml
 [license-svg]: https://img.shields.io/badge/license-MIT-blue.svg
 [license-url]: https://github.com/grokify/go-wordpressxml/blob/master/LICENSE
