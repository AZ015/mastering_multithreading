package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Point2D struct {
	x int
	y int
}

const numberOfThreads = 8

var (
	r  = regexp.MustCompile(`\((\d*),(\d*)\)`)
	wg = sync.WaitGroup{}
)

func findArea(inCh chan string) {
	for pointsStr := range inCh {
		var points []Point2D

		for _, p := range r.FindAllStringSubmatch(pointsStr, -1) {
			x, _ := strconv.Atoi(p[1])
			y, _ := strconv.Atoi(p[2])
			points = append(points, Point2D{x, y})
		}

		area := 0.0
		for i := 0; i < len(points); i++ {
			a, b := points[i], points[(i+1)%len(points)]
			area += float64(a.x*b.y) - float64(a.y*b.x)
		}
		fmt.Println(math.Abs(area) / 2.0)
	}
	wg.Done()
}

func main() {
	absPath, err := filepath.Abs("./threadpool")
	if err != nil {
		log.Fatalf("abs fileptah failed: %s", err)
	}
	data, err := ioutil.ReadFile(filepath.Join(absPath, "polygons.txt"))
	if err != nil {
		log.Fatalf("read file err: %s", err)
	}

	text := string(data)

	inCh := make(chan string, 100)
	for i := 0; i < numberOfThreads; i++ {
		go findArea(inCh)
	}
	wg.Add(numberOfThreads)

	start := time.Now()

	for _, line := range strings.Split(text, "\n") {
		inCh <- line
	}
	close(inCh)

	wg.Wait()

	elapsed := time.Since(start)
	fmt.Printf("\nProcessed took: %s\n", elapsed)

}
