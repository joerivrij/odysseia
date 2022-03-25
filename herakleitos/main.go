package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"github.com/kpango/glg"
	"github.com/odysseia/aristoteles"
	"github.com/odysseia/aristoteles/configs"
	"github.com/odysseia/herakleitos/app"
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

//go:embed rhema
var rhema embed.FS

func main() {
	//https://patorjk.com/software/taag/#p=display&f=Crawford2&t=HERAKLEITOS
	glg.Info("\n __ __    ___  ____    ____  __  _  _        ___  ____  ______   ___   _____\n|  |  |  /  _]|    \\  /    ||  |/ ]| |      /  _]|    ||      | /   \\ / ___/\n|  |  | /  [_ |  D  )|  o  ||  ' / | |     /  [_  |  | |      ||     (   \\_ \n|  _  ||    _]|    / |     ||    \\ | |___ |    _] |  | |_|  |_||  O  |\\__  |\n|  |  ||   [_ |    \\ |  _  ||     ||     ||   [_  |  |   |  |  |     |/  \\ |\n|  |  ||     ||  .  \\|  |  ||  .  ||     ||     | |  |   |  |  |     |\\    |\n|__|__||_____||__|\\_||__|__||__|\\_||_____||_____||____|  |__|   \\___/  \\___|\n                                                                            \n")
	glg.Info(strings.Repeat("~", 37))
	glg.Info("\"πάντα ῥεῖ\"")
	glg.Info("\"everything flows\"")
	glg.Info(strings.Repeat("~", 37))

	baseConfig := configs.HerakleitosConfig{}
	unparsedConfig, err := aristoteles.NewConfig(baseConfig)
	if err != nil {
		glg.Error(err)
		glg.Fatal("death has found me")
	}
	herakleitosConfig, ok := unparsedConfig.(*configs.HerakleitosConfig)
	if !ok {
		glg.Fatal("could not parse config")
	}

	root := "rhema"
	rootDir, err := rhema.ReadDir(root)
	if err != nil {
		glg.Fatal(err)
	}

	handler := app.HerakleitosHandler{Config: herakleitosConfig}

	err = handler.DeleteIndexAtStartUp()
	if err != nil {
		glg.Fatal(err)
	}
	err = handler.CreateIndexAtStartup()
	if err != nil {
		glg.Fatal(err)
	}

	var wg sync.WaitGroup
	documents := 0

	for _, dir := range rootDir {
		glg.Debug("working on the following directory: " + dir.Name())
		if dir.IsDir() {
			filePath := path.Join(root, dir.Name())
			files, err := rhema.ReadDir(filePath)
			if err != nil {
				glg.Fatal(err)
			}
			for _, f := range files {
				glg.Debug(fmt.Sprintf("found %s in %s", f.Name(), filePath))
				plan, _ := rhema.ReadFile(path.Join(filePath, f.Name()))
				var rhemai models.Rhema
				err := json.Unmarshal(plan, &rhemai)
				if err != nil {
					glg.Fatal(err)
				}

				documents += len(rhemai.Rhemai)

				wg.Add(1)
				go func() {
					err := handler.Add(rhemai, &wg)
					if err != nil {
						glg.Error(err)
					}
				}()
			}
		}
	}
	wg.Wait()
	glg.Infof("created: %s", strconv.Itoa(handler.Config.Created))
	glg.Infof("texts found in rhema: %s", strconv.Itoa(documents))
	os.Exit(0)
}
