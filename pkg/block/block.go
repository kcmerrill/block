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
	ignore        []string
	scoring       map[string]int
	lock          sync.Mutex
	rootDir       string
	searchDir     string
	homeDir       string
	blockDir      string
	category      string
	debug         bool
	checked       map[string]bool
}

// Search something
func Search(cmd, category, query string, debug bool) {
	b := block{
		query:         strings.ToLower(query),
		queryRegExStr: strings.Join(strings.Split(query, ""), ".*?"), // fuzzy matching
		category:      category,
		checked:       make(map[string]bool),
		ignore: []string{
			// could get large  .... but I think it's worth it(for now)
			"/tmp/",
			"/Library",
			"/dev",
			"/cores",
			"/var",
			"/Network",
			"/System",
			"/Volumes",
			"/etc",
			"/net",
			"/private",
			"/opt",
			"/usr",
			"/sbin",
			"/Applications",
			"/.git/",
			"vendor/",
			"Downloads/",
			"/node_modules/",
			"/gems/",
			"/go/pkg/dep/",
			"/cache/",
		},
		scoring: map[string]int{
			"cmd": 2,
		},
		debug:    debug,
		lock:     sync.Mutex{},
		homeDir:  os.Getenv("HOME"),
		blockDir: os.Getenv("HOME") + "/block/",
		flow: flow{
			category: "echo",
			name:     "[BLOCK] No results found ... Try broadening your search",
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
	fmt.Println("#DEFAULT-REPSPONSE", b.flow.category, b.flow.name)
	fmt.Println(b.flow.category, b.flow.name)

	var wg sync.WaitGroup

	// this is seriously dumb ... lets refactor yo
	wg.Add(1)
	go func() {
		// the crappiest plugin system ever, but we can figure this out later
		b.filesystem("/bin/bash -c", b.blockDir, 1)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		b.filesystem("open", b.rootDir, 0)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		b.filesystem("open", b.searchDir, 0)
		wg.Done()
	}()

	dirs, _ := ioutil.ReadDir("/")
	for _, d := range dirs {
		dirPath := "/" + d.Name()
		if isDirIgnored(d.Name(), dirPath, b.ignore) {
			continue
		}
		wg.Add(1)
		go func(dir string) {
			b.filesystem("open", dir, 0)
			wg.Done()
		}(dirPath)
	}

	wg.Add(1)
	go func() {
		b.dirs("open", "/Applications/")
		wg.Done()
	}()

	wg.Wait()
}

func (b *block) score(category, basepath, name string, boost int) {
	if b.debug {
		fmt.Println("#SCORED", category, name)
	}

	if name == b.rootDir {
		// no need to do anything with the current directory
		return
	}

	if b.category == "cd" && category != "cd" {
		// no need to do anything if the user doesn't want anything but directories
		return
	}

	// ghetto analytics yall
	// :D
	score := 0
	scoring := make([]string, 0)

	origName := name

	// lets strip off the root directory if it exists
	name = strings.ToLower(strings.Replace(name, basepath, "", 1))

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
		score += 2
		scoring = append(scoring, "+2 suffix match")
	}

	// same directory? lets boost it
	if strings.HasPrefix(origName, b.rootDir) {
		score += 2
		scoring = append(scoring, "+2 same dir match")
	}

	// add the boost(usually based on the directory)
	if boost > 0 {
		score += boost
		scoring = append(scoring, "+"+strconv.Itoa(boost)+" boost")
	}

	if b.debug {
		fmt.Println("#DEBUG-RANKED", score, category, origName)
	}

	// we have a winner?
	if score >= b.flow.score {
		if score == b.flow.score {
			if len(name) >= len(b.flow.name) {
				return
			}
			scoring = append(scoring, "+1 len is shorter(tie breaker)")
		}

		b.lock.Lock()
		b.flow = flow{
			score:    score,
			category: category,
			origName: origName,
			name:     name,
			scoring:  scoring,
		}

		if b.category != "" {
			b.flow.category = b.category
		}

		fmt.Println("#RANKED", b.flow.score, b.flow.category, b.flow.origName, scoring)
		// winner?
		fmt.Println(b.flow.category, b.flow.origName)
		b.lock.Unlock()
	}
}
