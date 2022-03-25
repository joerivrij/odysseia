package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"github.com/kpango/glg"
	"github.com/odysseia/anaximander/app"
	"github.com/odysseia/aristoteles"
	"github.com/odysseia/aristoteles/configs"
	"github.com/odysseia/plato/models"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
)

func init() {
	errlog := glg.FileWriter("/tmp/error.log", 0666)
	defer errlog.Close()

	glg.Get().
		SetMode(glg.BOTH).
		AddLevelWriter(glg.ERR, errlog)
}

//go:embed arkho
var arkho embed.FS

func main() {
	//https://patorjk.com/software/taag/#p=display&f=Crawford2&t=ANAXIMANDER
	glg.Info("\n  ____  ____    ____  __ __  ____  ___ ___   ____  ____   ___      ___  ____  \n /    ||    \\  /    ||  |  ||    ||   |   | /    ||    \\ |   \\    /  _]|    \\ \n|  o  ||  _  ||  o  ||  |  | |  | | _   _ ||  o  ||  _  ||    \\  /  [_ |  D  )\n|     ||  |  ||     ||_   _| |  | |  \\_/  ||     ||  |  ||  D  ||    _]|    / \n|  _  ||  |  ||  _  ||     | |  | |   |   ||  _  ||  |  ||     ||   [_ |    \\ \n|  |  ||  |  ||  |  ||  |  | |  | |   |   ||  |  ||  |  ||     ||     ||  .  \\\n|__|__||__|__||__|__||__|__||____||___|___||__|__||__|__||_____||_____||__|\\_|\n                                                                              \n")
	glg.Info(strings.Repeat("~", 37))
	glg.Info("\"οὐ γὰρ ἐν τοῖς αὐτοῖς ἐκεῖνος ἰχθῦς καὶ ἀνθρώπους, ἀλλ' ἐν ἰχθύσιν ἐγγενέσθαι τὸ πρῶτον ἀνθρώπους ἀποφαίνεται καὶ τραφέντας, ὥσπερ οἱ γαλεοί, καὶ γενομένους ἱκανους ἑαυτοῖς βοηθεῖν ἐκβῆναι τηνικαῦτα καὶ γῆς λαβέσθαι.\"")
	glg.Info("\"He declares that at first human beings arose in the inside of fishes, and after having been reared like sharks, and become capable of protecting themselves, they were finally cast ashore and took to land\"")
	glg.Info(strings.Repeat("~", 37))

	baseConfig := configs.AnaximanderConfig{}
	unparsedConfig, err := aristoteles.NewConfig(baseConfig)
	if err != nil {
		glg.Error(err)
		glg.Fatal("death has found me")
	}
	anaximanderConfig, ok := unparsedConfig.(*configs.AnaximanderConfig)
	if !ok {
		glg.Fatal("could not parse config")
	}

	root := "arkho"

	rootDir, err := arkho.ReadDir(root)
	if err != nil {
		glg.Fatal(err)
	}

	handler := app.AnaximanderHandler{Config: anaximanderConfig}
	err = handler.DeleteIndexAtStartUp()
	if err != nil {
		glg.Fatal(err)
	}

	err = handler.CreateIndexAtStartup()
	if err != nil {
		glg.Fatal(err)
	}

	var wg sync.WaitGroup

	for _, dir := range rootDir {
		glg.Debug("working on the following directory: " + dir.Name())
		if dir.IsDir() {
			filePath := path.Join(root, dir.Name())
			files, err := arkho.ReadDir(filePath)
			if err != nil {
				glg.Fatal(err)
			}
			for _, f := range files {
				glg.Debug(fmt.Sprintf("found %s in %s", f.Name(), filePath))
				plan, _ := arkho.ReadFile(path.Join(filePath, f.Name()))
				var declension models.Declension
				err := json.Unmarshal(plan, &declension)
				if err != nil {
					glg.Fatal(err)
				}

				wg.Add(1)
				go func() {
					err := handler.AddToElastic(declension, &wg)
					if err != nil {
						glg.Error(err)
					}
				}()
			}
		}
	}
	wg.Wait()
	glg.Infof("created: %s", strconv.Itoa(handler.Config.Created))
	os.Exit(0)
}
