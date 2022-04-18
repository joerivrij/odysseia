package command

import (
	"fmt"
	"github.com/kpango/glg"
	"github.com/odysseia/archimedes/util"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

type ValuesImages struct {
	Images struct {
		Tag        string `yaml:"tag"`
		HarborRepo string `yaml:"imageRepo"`
		PullSecret string `yaml:"pullSecret"`
		Sidecar    struct {
			Repo string `yaml:"repo"`
		} `yaml:"sidecar"`
		Init struct {
			Repo string `yaml:"repo"`
		} `yaml:"init"`
		Api struct {
			Repo string `yaml:"repo"`
		} `yaml:"odysseiaapi"`
		Seeder struct {
			Repo string `yaml:"repo"`
		} `yaml:"seeder"`
		Job struct {
			Repo string `yaml:"repo"`
		} `yaml:"job"`
		JobInit struct {
			Repo string `yaml:"repo"`
		} `yaml:"jobinit"`
	} `yaml:"images"`
}

type ValuesImagesTests struct {
	Images struct {
		Tag        string `yaml:"tag"`
		HarborRepo string `yaml:"imageRepo"`
		PullSecret string `yaml:"pullSecret"`
		System     struct {
			Repo string `yaml:"repo"`
		} `yaml:"system"`
		Load struct {
			Repo string `yaml:"repo"`
		} `yaml:"load"`
	} `yaml:"images"`
}

const (
	distDirectory string = "dist"
	binDirectory  string = "bin"
	defaultRepo   string = "core.harbor.domain:30003/odysseia"
	sidecarType   string = "sidecar"
	initType      string = "init"
	seederType    string = "seeder"
	apiType       string = "odysseiaapi"
	docsDest      string = "ploutarchos"
	yamlDest      string = "yaml"
	defaultTests  string = "tests"
)

var differentFlow = [...]string{"docs", "pheidias"}

func CreateImageSet() *cobra.Command {
	var (
		filePath        string
		name            string
		tag             string
		destinationRepo string
		fullBuild       bool
	)
	cmd := &cobra.Command{
		Use:   "set",
		Short: "create images for a set",
		Long: `Allows you to create images for all apis
- Filepath
`,
		Run: func(cmd *cobra.Command, args []string) {
			glg.Green("creating")

			if name == "" {
				glg.Fatal("a name is needed to build a set")
			}

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

			BuildImageSet(filePath, name, tag, destinationRepo, fullBuild)
		},
	}
	cmd.PersistentFlags().StringVarP(&filePath, "filepath", "f", "", "where to find the source code")
	cmd.PersistentFlags().StringVarP(&name, "name", "n", "", "chart name for the set")
	cmd.PersistentFlags().StringVarP(&tag, "tag", "t", "", "image tag")
	cmd.PersistentFlags().StringVarP(&destinationRepo, "dest", "d", "", "destination repo address")
	cmd.PersistentFlags().BoolVarP(&fullBuild, "full-build", "b", true, "whether to build all images or just the core (job and api)")

	return cmd
}

func BuildImageSet(filePath, name, tag, destRepo string, fullBuild bool) {
	themistoklesPath := filepath.Join(filePath, "themistokles", "odysseia", "charts")

	files, err := ioutil.ReadDir(themistoklesPath)
	if err != nil {
		glg.Error(err)
	}

	glg.Infof("****** 👷‍️ Working on Project %s 👷‍ ******", name)

	var valuesPath string
	for _, f := range files {
		if f.Name() == name {
			valuesPath = filepath.Join(themistoklesPath, f.Name(), "values.yaml")
			break
		}
	}

	glg.Info("****** ⚓ Getting Config From Themistokles ⚓ ******")

	readValues, err := ioutil.ReadFile(valuesPath)
	if err != nil {
		glg.Error(err)
	}

	var values ValuesImages
	err = yaml.Unmarshal(readValues, &values)
	if err != nil {
		glg.Error(err)
	}

	if name == defaultTests {
		var values ValuesImagesTests
		err = yaml.Unmarshal(readValues, &values)
		if err != nil {
			glg.Error(err)
		}

		systemName := values.Images.System.Repo
		projectPath := filepath.Join(filePath, systemName)

		err = buildImageWithLocalFile(projectPath, systemName, tag)

		err = pushImage(systemName, tag, destRepo, filePath)
		if err != nil {
			glg.Error(err)
			return
		}

		loadName := values.Images.System.Repo
		projectPathLoad := filepath.Join(filePath, loadName)

		err = buildImageWithLocalFile(projectPathLoad, loadName, tag)

		err = pushImage(loadName, tag, destRepo, filePath)
		if err != nil {
			glg.Error(err)
			return
		}

		return
	}

	for _, api := range differentFlow {
		if api == name {
			projectName := values.Images.Api.Repo
			projectPath := filepath.Join(filePath, projectName)

			err = buildImageWithLocalFile(projectPath, projectName, tag)

			err = pushImage(projectName, tag, destRepo, filePath)
			if err != nil {
				glg.Error(err)
				return
			}
			return
		}
	}

	if fullBuild {
		if values.Images.Sidecar.Repo != "" {
			sideCar := values.Images.Sidecar.Repo
			sideCarPath := filepath.Join(filePath, sideCar)

			err := runImageBuildFlow(filePath, sideCarPath, sideCar, tag, destRepo, sidecarType)
			if err != nil {
				glg.Error(err)
				return
			}
		}

		if values.Images.Init.Repo != "" {
			init := values.Images.Init.Repo
			initPath := filepath.Join(filePath, init)

			err := runImageBuildFlow(filePath, initPath, init, tag, destRepo, initType)
			if err != nil {
				glg.Error(err)
				return
			}
		}
	}

	if values.Images.Seeder.Repo != "" {
		seeder := values.Images.Seeder.Repo
		seederPath := filepath.Join(filePath, seeder)

		err := runImageBuildFlow(filePath, seederPath, seeder, tag, destRepo, seederType)
		if err != nil {
			glg.Error(err)
			return
		}
	}

	if values.Images.Api.Repo != "" {
		api := values.Images.Api.Repo
		apiPath := filepath.Join(filePath, api)

		err := runImageBuildFlow(filePath, apiPath, api, tag, destRepo, apiType)
		if err != nil {
			glg.Error(err)
			return
		}
	}

	if values.Images.Job.Repo != "" {
		job := values.Images.Job.Repo
		jobPath := filepath.Join(filePath, job)

		err := runImageBuildFlow(filePath, jobPath, job, tag, destRepo, apiType)
		if err != nil {
			glg.Error(err)
			return
		}
	}

	if values.Images.JobInit.Repo != "" {
		job := values.Images.JobInit.Repo
		jobPath := filepath.Join(filePath, job)

		err := runImageBuildFlow(filePath, jobPath, job, tag, destRepo, initType)
		if err != nil {
			glg.Error(err)
			return
		}
	}

	return
}

func runImageBuildFlow(filePath, projectPath, project, tag, destRepo, buildType string) error {
	err := checkForSwaggerFiles(filePath, projectPath)
	if err != nil {
		return err
	}

	err = runUnitTests(projectPath)
	if err != nil {
		return err
	}

	binPath, err := buildLocal(projectPath, project)
	if err != nil {
		return err
	}

	err = buildImage(binPath, filePath, project, tag, buildType)
	if err != nil {
		return err
	}

	err = pushImage(project, tag, destRepo, filePath)
	if err != nil {
		return err
	}

	return nil
}

func checkForSwaggerFiles(rootPath, projectPath string) error {
	files, err := ioutil.ReadDir(projectPath)
	if err != nil {
		glg.Fatal(err)
	}

	for _, f := range files {
		re := regexp.MustCompile(`-swagger.yaml`)
		if re.Match([]byte(f.Name())) {
			swaggerSource := filepath.Join(projectPath, f.Name())
			swaggerDestination := filepath.Join(rootPath, docsDest, yamlDest, f.Name())
			glg.Info("****** 🗄️ Getting OpenApi Doc 🗄️ ******")
			err = util.CopyFileContents(swaggerSource, swaggerDestination)
			if err != nil {
				glg.Error(err)
			}
			glg.Info("****** 📋 Copied OpenApi Doc 📋 ******")
		}
	}

	return nil
}

func buildImage(binPath, rootPath, projectName, tag, projectType string) error {
	projectPath := filepath.Join(rootPath, projectName)

	dockerFile := fmt.Sprintf("%s.Dockerfile", projectType)
	dockerSrc := filepath.Join(rootPath, dockerFile)
	dockerDest := filepath.Join(projectPath, dockerFile)
	err := util.CopyFileContents(dockerSrc, dockerDest)
	if err != nil {
		return err
	}

	binDest := filepath.Join(projectPath, projectName)
	err = util.CopyFileContents(binPath, binDest)
	if err != nil {
		return err
	}

	err = os.Chmod(binDest, os.ModePerm)
	if err != nil {
		return err
	}

	glg.Info("****** 🔨 Building Container Image 🔨 ******")
	imageName := fmt.Sprintf("%s:%s", projectName, tag)
	buildCommand := fmt.Sprintf("docker build --build-arg project_name=%s -f %s -t %s . --no-cache", projectName, dockerFile, imageName)
	err = util.ExecCommand(buildCommand, projectPath)
	if err != nil {
		return err
	}

	glg.Info("****** 🔱 Image Done 🔱 ******")

	err = os.Remove(binDest)
	if err != nil {
		return err
	}

	err = os.Remove(dockerDest)
	if err != nil {
		return err
	}

	return nil
}

func buildImageWithLocalFile(projectPath, projectName, tag string) error {
	glg.Info("****** 🔨 Building Container Image 🔨 ******")
	imageName := fmt.Sprintf("%s:%s", projectName, tag)

	buildCommand := fmt.Sprintf("docker build -t %s .", imageName)
	err := util.ExecCommand(buildCommand, projectPath)
	if err != nil {
		return err
	}

	glg.Info("****** 🔱 Image Done 🔱 ******")

	return err
}

func runUnitTests(path string) error {
	cmd := "go test ./... -cover"

	glg.Info("****** 🔎 Running Unittests 🔎 ******")
	err := util.ExecCommand(cmd, path)
	if err != nil {
		glg.Error("****** ❌ Unittests Failed ❌ ******")
		return err
	}

	glg.Info("****** ✅ Unittests Passed ✅ ******")

	return nil
}

func buildLocal(path, projectName string) (string, error) {
	fmtCommand := "go fmt ./..."
	err := util.ExecCommand(fmtCommand, path)
	if err != nil {
		return "", err
	}

	buildFor := runtime.GOOS
	binPath := filepath.Join(path, distDirectory, binDirectory, buildFor)
	err = os.MkdirAll(binPath, os.ModePerm)
	if err != nil {
		return "", err
	}

	projectBinPath := filepath.Join(binPath, projectName)
	buildCommand := fmt.Sprintf("GO111MODULE=on GOOS=linux CGO_ENABLED=0 go build main.go;mv main %s", projectBinPath)

	glg.Info("****** 🏗️ Building Golang Bin 🏗️ ******")
	err = util.ExecCommand(buildCommand, path)
	if err != nil {
		return "", err
	}

	glg.Info("****** 🏛️ Building Complete 🏛️ ******")

	return projectBinPath, nil
}

func pushImage(projectName, tag, destRepo, rootPath string) error {
	newTag := fmt.Sprintf("%s/%s:%s", destRepo, projectName, tag)

	glg.Info("****** 🖊️ Tagging Container Image 🖊️ ******")
	imageName := fmt.Sprintf("%s:%s", projectName, tag)
	tagCommand := fmt.Sprintf("docker tag %s %s", imageName, newTag)
	err := util.ExecCommand(tagCommand, rootPath)
	if err != nil {
		return err
	}

	glg.Infof("****** 📗 Tagged Image %s 📗 ******", newTag)

	glg.Info("****** 🚢️ Pushing Container Image 🚢 ******")
	pushCommand := fmt.Sprintf("docker push %s", newTag)
	err = util.ExecCommand(pushCommand, rootPath)
	if err != nil {
		return err
	}

	glg.Info("****** 🚀 Pushed Container Image 🚀 ******")
	glg.Info("image can be pulled as:")
	glg.Info(newTag)

	return nil
}
