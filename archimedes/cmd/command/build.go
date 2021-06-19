package command

import (
	"fmt"
	"github.com/kpango/glg"
	handler "github.com/odysseia/archimedes/pkg/impl"
	"github.com/spf13/cobra"
)

func BuildProject() *cobra.Command {
	var (
		filePath  string
	)
	cmd := &cobra.Command{
		Use:   "build",
		Short: "build all projects",
		Long: `Allows you build odysseia
- Filepath
`,
		Run: func(cmd *cobra.Command, args []string) {
			glg.Green("building")
			if filePath == "" {
				glg.Error(fmt.Sprintf("filepath is empty"))
				return
			}

			handler.BuildProject(filePath)
		},
	}
	cmd.PersistentFlags().StringVarP(&filePath, "filepath", "f", "", "where to find the source code")

	return cmd
}