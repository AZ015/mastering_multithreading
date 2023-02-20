package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"sync"
)

var (
	matches []string
	wg      = sync.WaitGroup{}
	m       = sync.Mutex{}
)

func fileSearch(root string, filename string) {
	fmt.Println("Searching in", root)

	files, _ := ioutil.ReadDir(root)
	for _, file := range files {
		if strings.Contains(file.Name(), filename) {
			m.Lock()
			matches = append(matches, filepath.Join(root, file.Name()))
			m.Unlock()
		}

		if file.IsDir() {
			wg.Add(1)
			go fileSearch(filepath.Join(root, file.Name()), filename)
		}
	}
	wg.Done()
}

func main() {
	wg.Add(1)
	go fileSearch("../../../", "README.md")
	wg.Wait()

	fmt.Println("\n------------------------ RESULTS ---------------------------\n")

	for _, file := range matches {
		fmt.Println("Matched", file)
	}
}
