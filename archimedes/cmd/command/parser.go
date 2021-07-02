package command

import (
	"fmt"
	"github.com/kpango/glg"
	handler "github.com/odysseia/archimedes/pkg/impl"
	"github.com/spf13/cobra"
	"os"
)

func ParseListToWords() *cobra.Command {
	var (
		filePath string
		outDir   string
	)
	cmd := &cobra.Command{
		Use:   "parse",
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

			handler.ParseLines(filePath, outDir)
		},
	}
	cmd.PersistentFlags().StringVarP(&filePath, "filepath", "f", "", "where to find the txt file")
	cmd.PersistentFlags().StringVarP(&outDir, "outdir", "o", "", "demokritos dir")

	return cmd
}
