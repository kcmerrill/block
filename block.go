package main

import (
	"time"

	"github.com/kcmerrill/block/pkg/block"
)

func main() {

	ignoreIfContains := []string{
		"/.git/",
		"/vendor/",
		"/node_modules/",
		"/gems/",
		"/go/pkg/",
		"/cache/",
		"/library/",
		"downloads/",
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

	query := "inventory.go"
	action := "cd"
	block.New(query, action, "/", 1*time.Second, ignoreIfContains, ignoreIfStartsWith)
}
