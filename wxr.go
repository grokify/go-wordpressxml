// wordpressxml provides WordPress XML parser with metadata
package wordpressxml

import (
	"io/ioutil"

	"github.com/frankbille/go-wxr-import"
)

func ReadFileWXR(filename string) (wxr.Wxr, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return wxr.Wxr{}, err
	}
	return wxr.ParseWxr(bytes), nil
}
