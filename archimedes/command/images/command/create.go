package command

import (
	"fmt"
	"github.com/kpango/glg"
	"github.com/odysseia/archimedes/util"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
)

func CreateImages() *cobra.Command {
	var (
		filePath string
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
				glg.Error(fmt.Sprintf("filepath is empty"))
				return
			}

			createImages(filePath)
		},
	}
	cmd.PersistentFlags().StringVarP(&filePath, "filepath", "f", "", "where to find the source code")

	return cmd
}

func createImages(odysseiaPath string) {
	ploutarchosPath := fmt.Sprintf("%s/%s/yaml", odysseiaPath, "ploutarchos")
	directories, err := ioutil.ReadDir(odysseiaPath)
	if err != nil {
		glg.Fatal(err)
	}

	for _, dir := range directories {
		fi, err := os.Stat(dir.Name())
		if err != nil {
			fmt.Println(err)
			return
		}

		// first action is to copy over the swagger files since they are needed for the image stage
		switch mode := fi.Mode(); {
		case mode.IsDir():
			charOne := dir.Name()[0]
			if string(charOne) == "." {
				continue
			}

			absolutePath, _ := filepath.Abs(dir.Name())
			lookForYamlFile(absolutePath, ploutarchosPath)
			command := "create-image"
			lookForMakeFile(absolutePath, command)
		case mode.IsRegular():
		}
	}
}

func lookForYamlFile(absolutePath, ploutarchosPath string) {
	files, err := ioutil.ReadDir(absolutePath)
	if err != nil {
		glg.Fatal(err)
	}

	for _, f := range files {
		re := regexp.MustCompile(`-swagger.yaml`)
		if re.Match([]byte(f.Name())) {
			swaggerSource := fmt.Sprintf("%s/%s", absolutePath, f.Name())
			swaggerDestination := fmt.Sprintf("%s/%s", ploutarchosPath, f.Name())
			glg.Info("****** 📗 Getting OpenApi Doc 📗 ******")
			glg.Debug("found swagger file %s copying to %s", swaggerSource, swaggerDestination)
			err = util.CopyFileContents(swaggerSource, swaggerDestination)
			if err != nil {
				glg.Error(err)
			}
			glg.Info("****** 📋 Copied OpenApi Doc 📋 ******")
		}
	}
}
func lookForMakeFile(absolutePath, command string) {
	files, err := ioutil.ReadDir(absolutePath)
	if err != nil {
		glg.Fatal(err)
	}

	for _, f := range files {
		if f.Name() == "Makefile" {
			glg.Debugf("Makefile found in %s", absolutePath)
			glg.Info("****** 🚢 Building Container Image 🚢 ******")
			makeCommand := fmt.Sprintf("make %s", command)

			err := util.ExecCommand(makeCommand, absolutePath)
			if err != nil {
				glg.Error(err)
			}

			glg.Info("****** 🔱 Image Done 🔱 ******")
		}
	}
}
