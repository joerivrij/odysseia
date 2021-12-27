package command

import (
	"fmt"
	"github.com/kpango/glg"
	"github.com/odysseia/archimedes/util"
	"github.com/odysseia/plato/models"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

func ListToWords() *cobra.Command {
	var (
		filePath string
		outDir   string
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
				outDir = fmt.Sprintf("%s/go/src/github.com/odysseia/demokritos/odysseia/perseus", homeDir)
			}

			parseLines(filePath, outDir)
		},
	}
	cmd.PersistentFlags().StringVarP(&filePath, "filepath", "f", "", "where to find the txt file")
	cmd.PersistentFlags().StringVarP(&outDir, "outdir", "o", "", "demokritos dir")

	return cmd
}

func parseLines(filePath, outDir string) {
	plan, _ := ioutil.ReadFile(filePath)
	wordList := strings.Split(string(plan), "\n")
	glg.Info(fmt.Sprintf("found %d words in %s", len(wordList), filePath))

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
