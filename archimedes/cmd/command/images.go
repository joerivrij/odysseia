package command

import (
	"fmt"
	"github.com/kpango/glg"
	handler "github.com/odysseia/archimedes/pkg/impl"
	"github.com/spf13/cobra"
)

func CreateImages() *cobra.Command {
	var (
		filePath string
	)
	cmd := &cobra.Command{
		Use:   "images",
		Short: "create images for all apis",
		Long: `Allows you to create images for all apis
- Filepath
`,
		Run: func(cmd *cobra.Command, args []string) {
			glg.Green("creating")
			if filePath == "" {
				glg.Error(fmt.Sprintf("filepath is empty"))
				return
			}

			handler.CreateImages(filePath)
		},
	}
	cmd.PersistentFlags().StringVarP(&filePath, "filepath", "f", "", "where to find the source code")

	return cmd
}
