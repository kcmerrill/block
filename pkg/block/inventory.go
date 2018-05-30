package block

import (
	"os"

	"github.com/karrick/godirwalk"
)

func filesystem(dir string) {
	// validate it exists
	if _, exists := os.Stat(dir); exists != nil {
		return
	}

	// these boots are made for walking ...
	godirwalk.Walk(dir, &godirwalk.Options{
		Unsorted: true,
		Callback: func(osPathname string, de *godirwalk.Dirent) error {
			if de.IsDir() {
			} else {
			}
			return nil
		},
	})
}
