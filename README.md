WordPress XML Parser
====================

[![Go CI][go-ci-svg]][go-ci-url]
[![Go Lint][go-lint-svg]][go-lint-url]
[![Go SAST][go-sast-svg]][go-sast-url]
[![Go Report Card][goreport-svg]][goreport-url]
[![Docs][docs-godoc-svg]][docs-godoc-url]
[![Visualization][viz-svg]][viz-url]
[![License][license-svg]][license-url]

 [go-ci-svg]: https://github.com/grokify/go-wordpressxml/actions/workflows/go-ci.yaml/badge.svg?branch=main
 [go-ci-url]: https://github.com/grokify/go-wordpressxml/actions/workflows/go-ci.yaml
 [go-lint-svg]: https://github.com/grokify/go-wordpressxml/actions/workflows/go-lint.yaml/badge.svg?branch=main
 [go-lint-url]: https://github.com/grokify/go-wordpressxml/actions/workflows/go-lint.yaml
 [go-sast-svg]: https://github.com/grokify/go-wordpressxml/actions/workflows/go-sast-codeql.yaml/badge.svg?branch=main
 [go-sast-url]: https://github.com/grokify/go-wordpressxml/actions/workflows/go-sast-codeql.yaml
 [goreport-svg]: https://goreportcard.com/badge/github.com/grokify/go-wordpressxml
 [goreport-url]: https://goreportcard.com/report/github.com/grokify/go-wordpressxml
 [docs-godoc-svg]: https://pkg.go.dev/badge/github.com/grokify/go-wordpressxml
 [docs-godoc-url]: https://pkg.go.dev/github.com/grokify/go-wordpressxml
 [viz-svg]: https://img.shields.io/badge/visualization-Go-blue.svg
 [viz-url]: https://mango-dune-07a8b7110.1.azurestaticapps.net/?repo=grokify%2Fgo-wordpressxml
 [loc-svg]: https://tokei.rs/b1/github/grokify/go-wordpressxml
 [repo-url]: https://github.com/grokify/go-wordpressxml
 [license-svg]: https://img.shields.io/badge/license-MIT-blue.svg
 [license-url]: https://github.com/grokify/go-wordpressxml/blob/main/LICENSE

## Overview

The `go-wordpressxml` package provides a WordPress WXR (WordPress eXtended RSS) XML parser.

## Features

- Parse WordPress XML export files
- Extract post metadata (authors, dates, categories, comments)
- Convert WordPress exports to [Hugo](https://gohugo.io/) static site format
- Export article metadata to CSV
- Export articles to HTML

## Documentation

Documentation is available on [pkg.go.dev](https://pkg.go.dev/github.com/grokify/go-wordpressxml).

## Installation

Installing any of the packages will install the entire library. For example:

```bash
$ go get github.com/grokify/go-wordpressxml
```

## Usage

### Parse WordPress XML and Export to CSV

```go
import (
	"github.com/grokify/go-wordpressxml"
)

func main() {
	wp := wordpressxml.NewWordPressXML()
	err := wp.ReadFile("myblog.wordpress.xml")
	if err != nil {
		panic(err)
	}
	err = wp.WriteMetaCSV("articles.csv")
	if err != nil {
		panic(err)
	}
}
```

### Convert to Hugo Format

```go
import (
	"github.com/grokify/go-wordpressxml"
	"github.com/grokify/go-wordpressxml/hugo"
)

func main() {
	wp := wordpressxml.NewWordPressXML()
	err := wp.ReadFile("myblog.wordpress.xml")
	if err != nil {
		panic(err)
	}

	converter := hugo.WxrConverter{}
	posts := converter.ConvertPosts(wp.Channel.Items)
	// Use posts to generate Hugo content files
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
