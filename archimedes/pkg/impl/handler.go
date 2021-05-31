package impl

import (
	"fmt"
	"github.com/kpango/glg"
	"github.com/lexiko/plato/models"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"unicode"
)

func ParseLines(filePath, outDir string) {
	plan, _ := ioutil.ReadFile(filePath)
	wordList := strings.Split(string(plan),"\n")
	glg.Info(fmt.Sprintf("found %d words in %s", len(wordList), filePath))

	var biblos models.Biblos
	currentLetter := "α"

	for i, word := range wordList {
		var greek string
		var english string
		for j, char := range word {
			c := fmt.Sprintf("%c", char)
			if j == 0 {
				removedAccent := removeAccent(c)
				if currentLetter != removedAccent {
					jsonBiblos, err := biblos.Marshal()
					if err != nil {
						glg.Error(err)
					}

					outputFile := fmt.Sprintf("%s/%s.json", outDir, currentLetter)
					writeFile(jsonBiblos, outputFile)
					currentLetter = removedAccent
					biblos = models.Biblos{}
				}
			}
			matched, err := regexp.MatchString(`[A-Za-z]`, c)
			if err != nil {
				glg.Error(err)
			}
			if matched {
				greek = strings.TrimSpace(word[0:j-1])
				english = strings.TrimSpace(word[j-1:])
				glg.Debug(fmt.Sprintf("found the greek: %s and the english %s", greek, english))

				meros := models.Meros{
					Greek:      greek,
					English:    english,
				}

				biblos.Biblos = append(biblos.Biblos, meros)
				break
			}
		}
		if i == len(wordList) -1 {
			jsonBiblos, err := biblos.Marshal()
			if err != nil {
				glg.Error(err)
			}

			outputFile := fmt.Sprintf("%s/%s.json", outDir, currentLetter)
			writeFile(jsonBiblos, outputFile)
		}
	}

	glg.Info(fmt.Sprintf("all words have been parsed and saved to %s", outDir))
}

func removeAccent(s string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	output, _, e := transform.String(t, s)
	if e != nil {
		panic(e)
	}
	return output
}

func writeFile(jsonBiblos []byte, outputFile string) {
	openedFile, err := os.Create(outputFile)
	if err != nil {
		glg.Error(err)
	}
	defer openedFile.Close()

	outputFromWrite, err := openedFile.Write(jsonBiblos)
	if err != nil {
		glg.Error(err)
	}

	glg.Info(fmt.Sprintf("finished writing %d bytes", outputFromWrite))
	glg.Info(fmt.Sprintf("file written to %s", outputFile))
}
