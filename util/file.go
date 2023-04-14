package util

import (
	"path/filepath"

	"github.com/katallaxie/autobot/res"
)

// ReadFile returns the bytes of a file searched in the path and beyond it
func ReadFile(path string) (bytes []byte) {
	bytes, err := res.Files.ReadFile(filepath.Clean(path))
	if err != nil {
		panic(err)
	}

	return bytes
}
