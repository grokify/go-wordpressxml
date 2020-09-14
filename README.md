WordPress XML Parser
====================

[![Go Report Card][goreport-svg]][goreport-link]
[![Docs][docs-godoc-svg]][docs-godoc-link]
[![License][license-svg]][license-link]

## Overview

The `go-wordpressxml` package provides WordPress XML parser.

## Documentation

Documentation is provided using godoc and available on [GoDoc.org](https://godoc.org/github.com/grokify/go0wordpressxml).

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
	wp := wordpressxml.NewWordpressXml()
	err := wp.ReadXml("myblog.wordpress.2016-08-13.xml")
	if err != nil {
		panic(err)
	}
	wp.WriteMetaCsv("articles.csv")
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

Please report issues and feature requests on [Github](https://github.com/grokify/go-wordpressxml).

## License

WordPress XML Parser is available under the MIT license. See [LICENSE](LICENSE) for details.

WordPress XML Parser &copy; 2016 by John Wang

 [build-status-svg]: https://api.travis-ci.org/grokify/go-wordpressxml.svg?branch=master
 [build-status-link]: https://travis-ci.org/grokify/go-wordpressxml
 [goreport-svg]: https://goreportcard.com/badge/github.com/grokify/go-wordpressxml
 [goreport-link]: https://goreportcard.com/report/github.com/grokify/go-wordpressxml
 [docs-godoc-svg]: https://img.shields.io/badge/docs-godoc-blue.svg
 [docs-godoc-link]: https://godoc.org/github.com/grokify/go-wordpressxml
 [license-svg]: https://img.shields.io/badge/license-MIT-blue.svg
 [license-link]: https://github.com/grokify/go-wordpressxml/blob/master/LICENSE
