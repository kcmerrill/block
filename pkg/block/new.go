package block

import (
	"os"
	"regexp"
	"strings"
	"time"
)

// New ...
func New(query, dir string, timeout time.Duration, ignoreIfContains, ignoreIfStartsWith []string) {
	cwd, _ := os.Getwd()
	b := &Block{
		inventory:          make(chan *Inventory),
		completed:          make(chan bool),
		ignoreIfContains:   ignoreIfContains,
		ignoreIfStartsWith: ignoreIfStartsWith,
		cwd:                strings.ToLower(cwd),
		timeout:            timeout,
		query:              query,
		queryRegExStr:      strings.Join(strings.Split(query, ""), ".*?"), // fuzzy matching
	}

	b.queryRegEx = regexp.MustCompile(b.queryRegExStr)
	go b.FindInventory("/")

	count := 0
	for {
		var done bool
		select {
		case <-time.After(b.timeout):
			done = true
		case done = <-b.completed:
		case <-b.inventory:
			count++
		}
		if done {
			break
		}
	}
}
