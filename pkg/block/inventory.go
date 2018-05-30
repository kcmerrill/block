package block

import (
	"fmt"
	"os"

	"github.com/karrick/godirwalk"
)

// Inventory ...
func Inventory(dir string, process chan string, finished chan bool, ignoreIfContains []string, ignoreIfStartsWith []string) {
	// validate it exists
	if _, exists := os.Stat(dir); exists != nil {
		return
	}

	// these boots are made for walking ...
	godirwalk.Walk(dir, &godirwalk.Options{
		Unsorted: true,
		Callback: func(osPathname string, de *godirwalk.Dirent) error {
			if de.IsDir() {
				fmt.Println(de.Name(), "is a directory")

			} else {
				fmt.Println(de.Name(), "is a file")
			}
			return nil
		},
	})
	finished <- true
}
