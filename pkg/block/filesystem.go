package block

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

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
		b.score(category, dir+d.Name())
	}
}

func (b *block) filesystem(category, dir string, ignore map[string]bool) {
	// validate it exists
	if _, exists := os.Stat(dir); exists != nil {
		return
	}

	godirwalk.Walk(dir, &godirwalk.Options{
		Callback: func(osPathname string, de *godirwalk.Dirent) error {
			if de.IsDir() {
				for ignoreDir := range ignore {
					if string(ignoreDir[0]) == "/" && strings.HasPrefix(osPathname, ignoreDir) {
						return filepath.SkipDir
					}
					if strings.Contains(osPathname, ignoreDir) {
						return filepath.SkipDir
					}
				}

				// hidden folders
				if len(de.Name()) >= 3 && string(de.Name()[0]) == "." {
					return filepath.SkipDir
				}

				// score it
				b.score("cd", osPathname)
			} else {
				// score it
				b.score(category, osPathname)
			}
			return nil
		},
	})
}
