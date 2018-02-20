package main

import (
	"flag"
	"strings"

	"github.com/kcmerrill/block/pkg/block"
)

func main() {
	debug := flag.Bool("d", false, "Debug results!")
	flag.Parse()
	query := flag.Args()
	if len(query) == 0 {
		// use the recommended category
	} else {
		if len(query) == 1 {
			block.Search("open", "", query[0], *debug)
			return
		}
		// use their category
		block.Search("open", query[0], strings.Join(query[1:], " "), *debug)
		return
	}
}
