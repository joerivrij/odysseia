package command

import (
	"github.com/kpango/glg"
	"github.com/odysseia/archimedes/util"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func CreateImages() *cobra.Command {
	var (
		filePath        string
		tag             string
		destinationRepo string
	)
	cmd := &cobra.Command{
		Use:   "create",
		Short: "create images for all apis",
		Long: `Allows you to create images for all apis
- Filepath
`,
		Run: func(cmd *cobra.Command, args []string) {
			glg.Green("creating")
			if filePath == "" {
				_, callingFile, _, _ := runtime.Caller(0)
				callingDir := filepath.Dir(callingFile)
				dirParts := strings.Split(callingDir, string(os.PathSeparator))
				var odysseiaPath []string
				for i, part := range dirParts {
					if part == "odysseia" {
						odysseiaPath = dirParts[0 : i+1]
						break
					}
				}
				l := "/"
				for _, path := range odysseiaPath {
					l = filepath.Join(l, path)
				}

				filePath = l
			}

			if tag == "" {
				glg.Warn("no tag set for image, using the git short hash")
				gitTag, err := util.ExecCommandWithReturn(`git rev-parse --short HEAD`, filePath)
				if err != nil {
					glg.Fatal(err)
				}

				tag = gitTag
			}

			if destinationRepo == "" {
				glg.Warnf("destination repo empty, default to %s", defaultRepo)
				destinationRepo = defaultRepo
			}

			glg.Infof("filepath set to: %s", filePath)

			LoopAndCreateImages(filePath, tag, destinationRepo)
		},
	}
	cmd.PersistentFlags().StringVarP(&filePath, "filepath", "f", "", "where to find the source code")
	cmd.PersistentFlags().StringVarP(&tag, "tag", "t", "", "image tag")
	cmd.PersistentFlags().StringVarP(&destinationRepo, "dest", "d", "", "destination repo address")

	return cmd
}

func LoopAndCreateImages(filePath, tag, destRepo string) {
	directories, err := ioutil.ReadDir(filePath)
	if err != nil {
		glg.Fatal(err)
	}

	for _, dir := range directories {
		BuildImageSet(filePath, dir.Name(), tag, destRepo, true)
	}
}
