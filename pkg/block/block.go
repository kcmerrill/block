package block

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
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

	// this is seriously dumb ... lets refactor yo
	wg.Add(1)
	go func() {
		// the crappiest plugin system ever, but we can figure this out later
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

	dirs, _ := ioutil.ReadDir("/")
	for _, d := range dirs {
		dirPath := "/" + d.Name()
		if b.isDirIgnored(d.Name(), dirPath) {
			continue
		}
		wg.Add(1)
		go func(dir string) {
			b.filesystem("open", dir, b.ignore)
			wg.Done()
		}(dirPath)
	}

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
	fmt.Println(b.flow.category, b.flow.origName)
}

func (b *block) score(category, name string) {
	// ghetto analytics yall
	// :D
	score := 0
	scoring := make([]string, 0)

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
		scoring = append(scoring, "+4 exact match")
	} else {
		if b.queryRegEx.Match([]byte(name)) {
			scoring = append(scoring, "+1 fuzzy match")
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
		scoring = append(scoring, "+"+strconv.Itoa(b.scoring[category])+" category modifier")
	}

	// do we want a directory?
	if b.category == category {
		score++
		scoring = append(scoring, "+1 category match")
	}

	// boost if it ends with what we wanted
	if strings.HasSuffix(name, b.query) {
		score++
		scoring = append(scoring, "+1 suffix match")
	}

	// same directory? lets boost it
	if strings.HasPrefix(origName, b.rootDir) {
		score += 2
		scoring = append(scoring, "+2 same dir match")
	}

	// score equal? which one is shorter closer to where youy are? <-- really shitty folks. really shitty. lol
	if score == b.flow.score {
		if len(name) < len(b.flow.name) {
			score++
			scoring = append(scoring, "+1 len is shorter(tie breaker)")
		}
	}

	// we have a winner?
	if score > b.flow.score {
		b.lock.Lock()
		b.flow = flow{
			score:    score,
			category: category,
			origName: origName,
			name:     name,
			scoring:  scoring,
		}
		b.lock.Unlock()
		fmt.Println(score, category, origName, scoring)
	}
}
