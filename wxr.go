// wordpressxml provides WordPress XML parser with metadata
package wordpressxml

import (
	"os"

	wxr "github.com/frankbille/go-wxr-import"
)

func ReadFileWXR(filename string) (wxr.Wxr, error) {
	if bytes, err := os.ReadFile(filename); err != nil {
		return wxr.Wxr{}, err
	} else {
		return wxr.ParseWxr(bytes), nil
	}
}
