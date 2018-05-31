package block

import (
	"fmt"
	"strconv"
	"strings"
)

func (b *Block) score(inventory Inventory) {
	score := 0
	scoring := make([]string, 0)

	// lets strip off the root directory if it exists
	name := strings.Replace(inventory.FileNameLowerCase, b.cwd, "", 1)

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
