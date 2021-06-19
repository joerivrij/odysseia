package impl

import (
	"fmt"
	"github.com/kpango/glg"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
)

func BuildProject(odysseiaPath string) {
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

		// first action is to copy over the swagger files since they are needed for the build stage
		switch mode := fi.Mode(); {
		case mode.IsDir():
			charOne := dir.Name()[0]
			if string(charOne) == "." {
				continue
			}

			absolutePath, _ := filepath.Abs(dir.Name())
			lookForYamlFile(absolutePath, ploutarchosPath)
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
			glg.Infof("found swagger file %s copying to %s", swaggerSource, swaggerDestination)
			err = copyFileContents(swaggerSource, swaggerDestination)
			if err != nil {
				glg.Error(err)
			}
			glg.Info("****** 📋 Copied OpenApi Doc 📋 ******")
		}
	}
}
func lookForMakeFile(absolutePath string) {

}

func copyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}
