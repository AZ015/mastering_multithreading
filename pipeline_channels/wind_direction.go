package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	windRegex     = regexp.MustCompile(`\d* METAR.*EGLL \d*Z [A-Z ]*(\d{5}KT|VRB\d{2}KT).*=`)
	tafValidation = regexp.MustCompile(`.*TAF.*`)
	comment       = regexp.MustCompile(`\w*#.*`)
	metarClose    = regexp.MustCompile(`.*=`)
	variableWind  = regexp.MustCompile(`VRB\d{2}KT`)
	validWind     = regexp.MustCompile(`\d{5}KT`)
	windDirOnly   = regexp.MustCompile(`\d{3}\d{2}KT`)
	windDist      [8]int
)

func parseToArray(textCh chan string, metarCh chan []string) {
	for text := range textCh {
		lines := strings.Split(text, "\n")
		metarSilce := make([]string, 0, len(lines))
		metarStr := ""

		for _, line := range lines {
			if tafValidation.MatchString(line) {
				break
			}
			if !comment.MatchString(line) {
				metarStr += strings.Trim(line, " ")
			}
			if metarClose.MatchString(line) {
				metarSilce = append(metarSilce, metarStr)
				metarStr = ""
			}
		}
		metarCh <- metarSilce
	}
	close(metarCh)
}

func extractWindDirection(metarsCh chan []string, windsCh chan []string) {
	for metars := range metarsCh {
		winds := make([]string, 0, len(metars))
		for _, metar := range metars {
			if windRegex.MatchString(metar) {
				winds = append(winds, windRegex.FindAllStringSubmatch(metar, -1)[0][1])
			}
		}
		windsCh <- winds
	}
	close(windsCh)
}

func mineWindDistribution(windsCh chan []string, distCh chan [8]int) {
	for winds := range windsCh {
		for _, wind := range winds {
			if variableWind.MatchString(wind) {
				for i := 0; i < 8; i++ {
					windDist[i]++
				}
			} else if variableWind.MatchString(wind) {
				windStr := windDirOnly.FindAllStringSubmatch(wind, -1)[0][1]
				if d, err := strconv.ParseFloat(windStr, 64); err == nil {
					dirIndex := int(math.Round(d/45.0)) % 8
					windDist[dirIndex]++
				}
			}
		}
	}
	distCh <- windDist
	close(distCh)
}

func main() {
	textCh := make(chan string)
	metarCh := make(chan []string)
	windsCh := make(chan []string)
	resultCh := make(chan [8]int)

	go parseToArray(textCh, metarCh)
	go extractWindDirection(metarCh, windsCh)
	go mineWindDistribution(windsCh, resultCh)

	absPath, _ := filepath.Abs("./metarfiles")
	files, _ := ioutil.ReadDir(absPath)
	start := time.Now()
	for _, file := range files {
		dat, err := ioutil.ReadFile(filepath.Join(absPath, file.Name()))
		if err != nil {
			panic(err)
		}
		text := string(dat)
		textCh <- text
	}
	close(textCh)
	results := <-resultCh
	elapsed := time.Since(start)

	fmt.Printf("%v\n", results)
	fmt.Printf("Processin took %s\n", elapsed)
}
