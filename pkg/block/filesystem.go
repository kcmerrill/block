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

func (b *block) isDirIgnored(name, osPathname string) bool {
	for ignoreDir := range b.ignore {
		if string(ignoreDir[0]) == "/" && strings.HasPrefix(osPathname, ignoreDir) {
			return true
		}
		if strings.Contains(osPathname, ignoreDir) {
			return true
		}
	}

	// hidden folders
	if len(name) >= 3 && string(name[0]) == "." {
		return true
	}

	return false
}

func (b *block) filesystem(category, dir string, ignore map[string]bool) {
	// validate it exists
	if _, exists := os.Stat(dir); exists != nil {
		return
	}

	godirwalk.Walk(dir, &godirwalk.Options{
		Unsorted: true,
		Callback: func(osPathname string, de *godirwalk.Dirent) error {
			if de.IsDir() {
				if b.isDirIgnored(de.Name(), osPathname) {
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
