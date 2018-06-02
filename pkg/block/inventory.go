package block

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/karrick/godirwalk"
)

// Inventory ...
type Inventory struct {
	Type                   string
	Dir                    string
	File                   string
	FileNameWithoutBaseDir string
	FileName               string
	FileNameLowerCase      string
	Score                  int
	Scoring                []string
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
				for _, contains := range b.IgnoreIfContains {
					if strings.Contains(i.FileNameLowerCase, contains) {
						// we should skip it
						if b.Debug {
							fmt.Println(i.FileName, fmt.Sprintf("skipping, contains '%s'", contains))
						}
						return filepath.SkipDir
					}
				}

				for _, startsWith := range b.IgnoreIfStartsWith {
					if strings.HasPrefix(i.FileNameLowerCase, startsWith) {
						// we should skip it
						if b.Debug {
							fmt.Println(i.FileName, fmt.Sprintf("skipping, starts with '%s'", startsWith))
						}
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
