package main

import (
	"fmt"
	"time"

	"github.com/kcmerrill/block/pkg/block"
)

func main() {
	files := make(chan string)
	finished := make(chan bool)
	go block.Inventory("./", files, finished, []string{}, []string{})

	for {
		var done bool
		select {
		case <-time.After(1 * time.Second):
			done = true
		case done = <-finished:
		case file := <-files:
			fmt.Println(file)
		}
		if done {
			break
		}
	}
}
