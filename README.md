WordPress XML Parser
====================

[![Codacy Badge](https://api.codacy.com/project/badge/Grade/bc2a0f5a10f6459f8fba86f58f368553)](https://www.codacy.com/app/grokify/wordpress-xml-go?utm_source=github.com&utm_medium=referral&utm_content=grokify/wordpress-xml-go&utm_campaign=badger)

[![Docs][docs-godoc-svg]][docs-godoc-link]
[![License][license-svg]][license-link]

## Overview

The `wordpress-xml-go` package provides WordPress XML parser.

## Documentation

Documentation is provided using godoc and available on [GoDoc.org](https://godoc.org/github.com/grokify/wordpress-xml-go).

## Installation

Installing any of the packages will install the entire library. For example:

```bash
$ go get github.com/grokify/wordpress-xml-go
```

## Usage

```go
import (
	"github.com/grokify/wordpress-xml-go"
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

1. Fork it ( http://github.com/grokify/wordpress-xml-go/fork )
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new Pull Request

Please report issues and feature requests on [Github](https://github.com/grokify/wordpress-xml-go).

## License

WordPress XML Parser is available under the MIT license. See [LICENSE](LICENSE) for details.

WordPress XML Parser &copy; 2016 by John Wang

 [build-status-svg]: https://api.travis-ci.org/grokify/wordpress-xml-go.svg?branch=master
 [build-status-link]: https://travis-ci.org/grokify/wordpress-xml-go
 [docs-godoc-svg]: https://img.shields.io/badge/docs-godoc-blue.svg
 [docs-godoc-link]: https://godoc.org/github.com/grokify/wordpress-xml-go
 [license-svg]: https://img.shields.io/badge/license-MIT-blue.svg
 [license-link]: https://github.com/grokify/wordpress-xml-go/blob/master/LICENSE