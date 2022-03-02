package command

import (
	"fmt"
	"github.com/kpango/glg"
	"github.com/odysseia/archimedes/util"
	"github.com/odysseia/plato/helpers"
	"github.com/odysseia/plato/models"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

const (
	ByFile   string = "by-file"
	ByLetter string = "by-letter"
)

func ListToWords() *cobra.Command {
	var (
		filePath string
		outDir   string
		mode     string
	)
	cmd := &cobra.Command{
		Use:   "list",
		Short: "parse a list of words",
		Long: `Allows you to parse a list of words to be used by demokritos
- Filepath
- OutDir
`,
		Run: func(cmd *cobra.Command, args []string) {
			glg.Green("parsing")
			if filePath == "" {
				glg.Error(fmt.Sprintf("filepath is empty"))
				return
			}

			if outDir == "" {
				glg.Debug(fmt.Sprintf("no outdir set assuming one"))
				homeDir, _ := os.UserHomeDir()
				outDir = fmt.Sprintf("%s/go/src/github.com/odysseia/eratosthenes/bibliotheke", homeDir)
			}

			if mode == "" {
				mode = ByFile
			}

			parse(filePath, outDir, mode)

		},
	}
	cmd.PersistentFlags().StringVarP(&filePath, "filepath", "f", "", "where to find the txt file")
	cmd.PersistentFlags().StringVarP(&outDir, "outdir", "o", "", "demokritos dir")
	cmd.PersistentFlags().StringVarP(&mode, "mode", "m", "", "mode for parsing valid options: by-letter, by-file")

	return cmd
}

func parse(filePath, outDir, mode string) {
	plan, _ := ioutil.ReadFile(filePath)
	wordList := strings.Split(string(plan), "\n")
	glg.Info(fmt.Sprintf("found %d words in %s", len(wordList), filePath))

	if mode == ByFile {
		pathParts := strings.Split(filePath, "/")
		name := strings.Split(pathParts[len(pathParts)-1], ".")[0]
		glg.Info(name)

		parseLinesByFile(outDir, name, wordList)
	} else if mode == ByLetter {
		parseLinesByLetter(outDir, wordList)
	} else {
		glg.Fatal("No mode provided")
	}
}

func parseLinesByLetter(outDir string, wordList []string) {
	var biblos models.Biblos
	currentLetter := "α"

	for i, word := range wordList {
		var greek string
		var english string
		for j, char := range word {
			c := fmt.Sprintf("%c", char)
			if j == 0 {
				removedAccent := util.RemoveAccent(c)
				if currentLetter != removedAccent {
					jsonBiblos, err := biblos.Marshal()
					if err != nil {
						glg.Error(err)
					}

					outputFile := fmt.Sprintf("%s/%s.json", outDir, currentLetter)
					util.WriteFile(jsonBiblos, outputFile)
					currentLetter = removedAccent
					biblos = models.Biblos{}
				}
			}
			matched, err := regexp.MatchString(`[A-Za-z]`, c)
			if err != nil {
				glg.Error(err)
			}
			if matched {
				greek = strings.TrimSpace(word[0 : j-1])
				english = strings.TrimSpace(word[j-1:])
				glg.Debug(fmt.Sprintf("found the greek: %s and the english %s", greek, english))

				meros := models.Meros{
					Greek:   greek,
					English: english,
				}

				biblos.Biblos = append(biblos.Biblos, meros)
				break
			}
		}
		if i == len(wordList)-1 {
			jsonBiblos, err := biblos.Marshal()
			if err != nil {
				glg.Error(err)
			}

			outputFile := fmt.Sprintf("%s/%s.json", outDir, currentLetter)
			util.WriteFile(jsonBiblos, outputFile)
		}
	}

	glg.Info(fmt.Sprintf("all words have been parsed and saved to %s", outDir))
}

func parseLinesByFile(outDir, name string, wordList []string) {
	var logos models.Logos

	for _, word := range wordList {
		var greek string
		var translation string
		for j, char := range word {
			c := fmt.Sprintf("%c", char)
			matched, err := regexp.MatchString(`[A-Za-z]`, c)
			if err != nil {
				glg.Error(err)
			}
			if matched {
				greek = strings.TrimSpace(word[0 : j-1])
				translation = strings.TrimSpace(word[j-1:])
				glg.Debug(fmt.Sprintf("found the greek: %s and the translation %s", greek, translation))

				meros := models.Word{
					Greek:       greek,
					Translation: translation,
				}

				logos.Logos = append(logos.Logos, meros)
				break
			}
		}
	}

	numberOfWords := len(logos.Logos)
	var wordsPerChapter int
	switch {
	case numberOfWords < 500:
		wordsPerChapter = 10
	case numberOfWords > 501:
		wordsPerChapter = 20
	}

	chaptersLength := numberOfWords / wordsPerChapter
	lastChapter := chaptersLength + 1
	var randonNumbers []int

	for i := 1; i <= chaptersLength; i++ {
		for j := 1; j <= wordsPerChapter; j++ {
			randomNumber := helpers.GenerateRandomNumber(numberOfWords)
			numberIsUnique := uniqueNumber(randonNumbers, randomNumber)

			for !numberIsUnique {
				randomNumber = helpers.GenerateRandomNumber(numberOfWords)
				numberIsUnique = uniqueNumber(randonNumbers, randomNumber)
			}

			logos.Logos[randomNumber].Chapter = int64(i)
			randonNumbers = append(randonNumbers, randomNumber)
		}
	}

	for i, word := range logos.Logos {
		if word.Chapter == int64(0) {
			logos.Logos[i].Chapter = int64(lastChapter)
		}
	}

	jsonLogos, err := logos.Marshal()
	if err != nil {
		glg.Error(err)
	}

	outputFile := fmt.Sprintf("%s/%s.json", outDir, name)
	util.WriteFile(jsonLogos, outputFile)

}

func uniqueNumber(numberList []int, number int) bool {
	numberIsUnique := true
	for _, n := range numberList {
		if n == number {
			numberIsUnique = false

		}
	}

	return numberIsUnique
}
