package block

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/karrick/godirwalk"
)

func (b *block) dirs(category, dir string) {
	// validate it exists
	if _, exists := os.Stat(dir); exists != nil {
		return
	}

	// We only need top level of directories for this callout
	dirs, _ := ioutil.ReadDir(dir)
	for _, d := range dirs {
		b.score(category, dir, dir+d.Name(), 0)
	}
}

func (b *block) filesystem(category, dir string, boost int) {
	// validate it exists
	if _, exists := os.Stat(dir); exists != nil {
		return
	}

	// these boots are made for walking ...
	godirwalk.Walk(dir, &godirwalk.Options{
		Unsorted: true,
		Callback: func(osPathname string, de *godirwalk.Dirent) error {
			if de.IsDir() {
				if isDirIgnored(de.Name(), osPathname, b.ignore) {
					return filepath.SkipDir
				}
				// score it
				b.score("cd", dir, osPathname, boost)
			} else {
				// score it
				b.score(category, dir, osPathname, boost)
			}
			return nil
		},
	})
}
