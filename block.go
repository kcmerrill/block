package main

import (
	"flag"
	"os"
	"strings"
	"time"

	"github.com/kcmerrill/block/pkg/block"
)

func main() {
	debug := flag.Bool("debug", false, "Debug")
	flag.Parse()

	var query string
	var action string

	args := flag.Args()

	switch len(args) {
	case 0:
		os.Exit(1)
	case 1:
		action = ""
		query = args[0]
	case 2:
		action = args[0]
		query = args[1]
	}

	ignoreIfContains := []string{
		"/.git/",
		"/vendor/",
		"/node_modules/",
		"/gems/",
		"/go/pkg/", // hmmm ... just me?
		"/cache/",
		"/library/",
		"downloads/", // could be iffy ...
		"/applications/",
		"/album artwork/",
		".app/",
		"/.", // controversial. Don't @ me
	}

	ignoreIfStartsWith := []string{
		"/network",
		"/system",
		"/volumes",
		"/bin",
		"/cores",
		"/dev",
		"/keybase",
		"/net",
		"/opt",
		"/private",
		"/usr",
		"/var",
		"/sbin",
	}

	b := &block.Block{
		IgnoreIfContains:   ignoreIfContains,
		IgnoreIfStartsWith: ignoreIfStartsWith,
		Timeout:            1 * time.Second,
		Query:              strings.ToLower(query),
		Action:             action,
		Debug:              *debug,
	}

	block.New(b)
}
