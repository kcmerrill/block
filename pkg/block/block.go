package block

import (
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Block ...
type Block struct {
	Query              string
	IgnoreIfContains   []string
	IgnoreIfStartsWith []string
	Timeout            time.Duration
	Action             string
	Debug              bool
	Boosted            map[string]int
	Overrides          map[string]string
	DefaultCommand     string

	debugMsgs     []string
	inventory     chan *Inventory
	completed     chan bool
	cwd           string
	queryRegEx    *regexp.Regexp
	queryRegExStr string
	maxInventory  *Inventory
	lock          *sync.Mutex
	debugLock     *sync.Mutex
}

func (b *Block) debugMsg(subject, msg string) {
	b.debugLock.Lock()
	b.debugMsgs = append(b.debugMsgs, subject+": "+msg)
	b.debugLock.Unlock()
}

func (b *Block) processInventory() {
	go b.FindInventory("/")

	count := 0
	for {
		var done bool
		select {
		case <-time.After(b.Timeout):
			done = true
		case done = <-b.completed:
		case inventory := <-b.inventory:
			count++
			go b.score(inventory)
		}
		if done {
			break
		}
	}

	b.debugMsg("#Scored", strconv.Itoa(count))
}

func (b *Block) act(inventory *Inventory) string {

	cmd := ""
	if inventory.Type == "directory" {
		cmd = "cd"
	}

	if inventory.Type != "directory" && inventory.Type != "file" {
		// we should use it
		cmd = inventory.Type
	}

	for startsWith, override := range b.Overrides {
		if strings.HasPrefix(inventory.ActionLowerCase, startsWith) {
			cmd = override
			break
		}
	}

	if b.Action != "" {
		cmd = b.Action
	}

	return cmd + " " + inventory.Action
}
