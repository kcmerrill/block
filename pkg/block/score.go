package block

func score(category, basepath, name string, boost int) int {
	/*if b.debug {
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

	*/
	return 0
}
