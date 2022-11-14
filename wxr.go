// wordpressxml provides WordPress XML parser with metadata
package wordpressxml

import (
	"os"

	"github.com/frankbille/go-wxr-import"
)

func ReadFileWXR(filename string) (wxr.Wxr, error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return wxr.Wxr{}, err
	}
	return wxr.ParseWxr(bytes), nil
}
