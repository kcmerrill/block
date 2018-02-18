package block

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

type block struct {
	flow          flow
	query         string
	queryRegEx    *regexp.Regexp
	queryRegExStr string
	ignore        map[string]bool
	scoring       map[string]int
	lock          sync.Mutex
	rootDir       string
	searchDir     string
	homeDir       string
	blockDir      string
	category      string
}

// Search something
func Search(cmd, category, query string) {
	b := block{
		query:         strings.ToLower(query),
		queryRegExStr: strings.Join(strings.Split(query, ""), ".*?"), // fuzzy matching
		category:      category,
		ignore: map[string]bool{
			// could get large  .... but I think it's worth it(for now)
			"/Library":      true,
			"/dev":          true,
			"/cores":        true,
			"/var":          true,
			"/Network":      true,
			"/System":       true,
			"/Volumes":      true,
			"/etc":          true,
			"/net":          true,
			"/private":      true,
			"/usr":          true,
			"/sbin":         true,
			"/Applications": true,
			".git/":         true,
			"vendor/":       true,
			"Downloads/":    true,
		},
		scoring: map[string]int{
			"cmd": 2,
		},
		lock:     sync.Mutex{},
		homeDir:  os.Getenv("HOME"),
		blockDir: os.Getenv("HOME") + "/block/",
		flow: flow{
			category: "echo",
			name:     "Unable to find anything.",
		},
	}

	b.queryRegEx = regexp.MustCompile(b.queryRegExStr)
	b.rootDir, _ = os.Getwd()
	b.rootDir += "/"
	b.homeDir = os.Getenv("HOME")
	b.searchDir = b.rootDir

	// figure out searchdir
	for {
		dir := filepath.Dir(b.searchDir)
		if dir != "/" {
			b.searchDir = dir
			continue
		}
		break
	}

	fmt.Println("query: ", b.query)
	fmt.Println("queryRegEx: ", b.queryRegExStr)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		b.filesystem("/bin/bash -c", b.blockDir, b.ignore)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		b.filesystem("open", b.rootDir, b.ignore)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		b.filesystem("open", b.searchDir, b.ignore)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		b.filesystem("open", "/", b.ignore)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		b.dirs("open", "/Applications/")
		wg.Done()
	}()

	wg.Wait()

	if category != "" {
		b.flow.category = category
	}

	// winner?
	fmt.Println(b.flow.category, b.flow.name)
}

func (b *block) score(category, name string) {
	// ghetto analytics yall
	// :D
	score := 0
	origName := name
	if b.category == "cd" && category != "cd" {
		// no need to do anything if the user doesn't want anything but directories
		return
	}

	if name == b.rootDir {
		// no need to do anything with the current directory
		return
	}

	// lets strip off the root directory if it exists
	if category != "cmd" {
		name = strings.ToLower(strings.Replace(name, b.rootDir, "", 1))
	}

	if strings.Contains(name, b.query) {
		// exact matches should get a boost
		score += 4
	} else {
		if b.queryRegEx.Match([]byte(name)) {
			score++
		}
	}

	if score == 0 {
		// no need to go on ... drop it on the floor
		return
	}

	modifier, exists := b.scoring[category]
	if exists {
		score += modifier
	}

	// do we want a directory?
	if b.category == category {
		score++
	}

	// boost if it ends with what we wanted
	if strings.HasSuffix(name, b.query) {
		score++
	}

	// score equal? which one is shorter?
	if score == b.flow.score {
		if len(name) < len(b.flow.name) {
			score++
		}
	}

	// we have a winner?
	if score > b.flow.score {
		b.lock.Lock()
		b.flow = flow{score: score, category: category, name: name}
		b.lock.Unlock()
		fmt.Println(score, category, origName)
	}
}
