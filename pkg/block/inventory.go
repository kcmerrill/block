package block

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/karrick/godirwalk"
)

// Inventory ...
type Inventory struct {
	Type              string
	Dir               string
	File              string
	FileName          string
	FileNameLowerCase string
}

// FindInventory ...
func (b *Block) FindInventory(dir string) {
	// validate it exists
	if _, exists := os.Stat(dir); exists != nil {
		return
	}

	// these boots are made for walking ...
	godirwalk.Walk(dir, &godirwalk.Options{
		Unsorted: true,
		ErrorCallback: func(osPathname string, err error) godirwalk.ErrorAction {
			return godirwalk.SkipNode
		},
		Callback: func(osPathname string, de *godirwalk.Dirent) error {
			i := &Inventory{
				Type:              "file",
				File:              de.Name(),
				FileName:          osPathname,
				FileNameLowerCase: strings.ToLower(osPathname),
				Dir:               osPathname,
			}

			if de.IsDir() {
				for _, contains := range b.ignoreIfContains {
					if strings.Contains(i.FileNameLowerCase, contains) {
						// we should skip it
						return filepath.SkipDir
					}
				}

				for _, startsWith := range b.ignoreIfStartsWith {
					if strings.HasPrefix(i.FileNameLowerCase, startsWith) {
						// we should skip it
						return filepath.SkipDir
					}
				}

				i.Type = "directory"
			}

			// process the bit of inventory
			b.inventory <- i
			return nil
		},
	})
	b.completed <- true
}
