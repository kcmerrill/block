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
	Type            string
	Dir             string
	Action          string
	ActionLowerCase string
	ActionShortened string
	Score           int
	Scoring         []string
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
				Type:            "file",
				ActionShortened: de.Name(),
				Action:          osPathname,
				ActionLowerCase: strings.ToLower(osPathname),
				Dir:             osPathname,
			}

			if de.IsDir() {
				for _, contains := range b.IgnoreIfContains {
					if strings.Contains(i.ActionLowerCase, contains) {
						// we should skip it
						if b.Debug {
							fmt.Println(i.Action, fmt.Sprintf("skipping, contains '%s'", contains))
						}
						return filepath.SkipDir
					}
				}

				for _, startsWith := range b.IgnoreIfStartsWith {
					if strings.HasPrefix(i.ActionLowerCase, startsWith) {
						// we should skip it
						if b.Debug {
							fmt.Println(i.Action, fmt.Sprintf("skipping, starts with '%s'", startsWith))
						}
						return filepath.SkipDir
					}
				}

				i.Type = "directory"
			} else {
				if b.Action == "cd" {
					// we shouldn't emit anything, the user is looking for a directory, this is a file ...
					return nil
				}
			}

			// process the bit of inventory
			b.inventory <- i
			return nil
		},
	})
	b.completed <- true
}
