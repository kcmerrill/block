package block

import (
	"regexp"
	"time"
)

// Block ...
type Block struct {
	query              string
	action             string
	inventory          chan *Inventory
	completed          chan bool
	ignoreIfContains   []string
	ignoreIfStartsWith []string
	cwd                string
	timeout            time.Duration
	queryRegEx         *regexp.Regexp
	queryRegExStr      string
}
