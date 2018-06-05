package block

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/mitchellh/go-homedir"
)

// New ...
func New(b *Block) {
	cwd, _ := os.Getwd()
	homeDir, _ := homedir.Dir()

	b.queryRegExStr = strings.Join(strings.Split(b.Query, ""), ".*?") // fuzzy matching
	b.cwd = strings.ToLower(cwd)
	b.inventory = make(chan *Inventory)
	b.completed = make(chan bool)
	b.lock = &sync.Mutex{}
	b.debugLock = &sync.Mutex{}
	b.maxInventory = &Inventory{
		Type:     "echo",
		FileName: "[BLOCK] No results found ... Try broadening your search",
		Score:    1,
		Scoring:  []string{"+1 Default"},
	}

	b.queryRegEx = regexp.MustCompile(b.queryRegExStr)
	b.boost = map[string]int{
		strings.ToLower(homeDir) + "/block/": 2,
	}
	b.override = map[string]string{
		strings.ToLower(homeDir) + "/block": "bash",
	}

	b.config()

	b.debugMsg("Query", b.Query)
	b.debugMsg("Fuzzy", b.queryRegExStr)
	b.debugMsg("Dir", b.cwd)

	b.processInventory()

	b.debugMsg("Found", b.maxInventory.FileName)
	b.debugMsg("Reasons", strings.Join(b.maxInventory.Scoring, ", "))

	fmt.Println(strings.Join(b.debugMsgs, "\n"))

	// lets do something!
	fmt.Println(b.act(b.maxInventory))
}
