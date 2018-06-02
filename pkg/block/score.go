package block

import (
	"fmt"
	"strings"
)

func (b *Block) score(inventory *Inventory) {
	inventory.Scoring = make([]string, 0)

	// lets strip off the root directory if it exists
	inventory.FileNameWithoutBaseDir = strings.Replace(inventory.FileNameLowerCase, b.cwd, "", 1)

	if strings.Contains(inventory.FileNameWithoutBaseDir, b.Query) {
		// exact matches should get a boost
		inventory.Score += 4
		inventory.Scoring = append(inventory.Scoring, "+4 exact match")
	} else {
		if b.queryRegEx.Match([]byte(inventory.FileNameWithoutBaseDir)) {
			inventory.Scoring = append(inventory.Scoring, "+1 fuzzy match")
			inventory.Score++
		}
	}

	if inventory.Score == 0 {
		// no need to go on ... drop it on the floor
		if b.Debug {
			fmt.Println(inventory.FileName, "did not match exactly or fuzzy.")
		}
		return
	}

	// are we trying to switch directories? If so .. lets boost it
	if b.Action == "cd" && inventory.Type == "directory" {
		inventory.Score++
		inventory.Scoring = append(inventory.Scoring, "+1 action match")
	}

	if b.Action != "cd" && inventory.Type != "directory" {
		inventory.Score++
		inventory.Scoring = append(inventory.Scoring, "+1 action match")
	}

	// boost if it ends with what we wanted
	if strings.HasSuffix(inventory.FileName, b.Query) {
		inventory.Score += 2
		inventory.Scoring = append(inventory.Scoring, "+2 suffix match")
	}

	// same directory? lets boost it
	if strings.HasPrefix(inventory.FileName, b.cwd) {
		inventory.Score += 2
		inventory.Scoring = append(inventory.Scoring, "+2 same dir match")
	}

	// we have a winner?
	if inventory.Score >= b.maxInventory.Score {
		if inventory.Score == b.maxInventory.Score {
			// TODO: shortness is messed up here I believe
			if len(inventory.FileNameWithoutBaseDir) >= len(b.maxInventory.FileNameWithoutBaseDir) {
				return
			}
			inventory.Scoring = append(inventory.Scoring, "+1 len is shorter(tie breaker)")
		}

		b.lock.Lock()
		b.maxInventory = inventory
		b.lock.Unlock()
	}

	if b.Debug {
		fmt.Println("#", inventory.FileName, strings.Join(inventory.Scoring, "\n"))
	}
}
