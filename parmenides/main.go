package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/kpango/glg"
	"github.com/odysseia/plato/elastic"
	"github.com/odysseia/plato/models"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
)

func init() {
	errlog := glg.FileWriter("/tmp/error.log", 0666)
	defer errlog.Close()

	glg.Get().
		SetMode(glg.BOTH).
		AddLevelWriter(glg.ERR, errlog)
}

func main() {
	//https://patorjk.com/software/taag/#p=display&f=Crawford2&t=PARMENIDES
	glg.Info("\n ____   ____  ____   ___ ___    ___  ____   ____  ___      ___  _____\n|    \\ /    ||    \\ |   |   |  /  _]|    \\ |    ||   \\    /  _]/ ___/\n|  o  )  o  ||  D  )| _   _ | /  [_ |  _  | |  | |    \\  /  [_(   \\_ \n|   _/|     ||    / |  \\_/  ||    _]|  |  | |  | |  D  ||    _]\\__  |\n|  |  |  _  ||    \\ |   |   ||   [_ |  |  | |  | |     ||   [_ /  \\ |\n|  |  |  |  ||  .  \\|   |   ||     ||  |  | |  | |     ||     |\\    |\n|__|  |__|__||__|\\_||___|___||_____||__|__||____||_____||_____| \\___|\n                                                                     \n")
	glg.Info(strings.Repeat("~", 37))
	glg.Info("\"τό γάρ αυτο νοειν έστιν τε καί ειναι\"")
	glg.Info("\"for it is the same thinking and being\"")
	glg.Info(strings.Repeat("~", 37))

	elasticClient, err := elastic.CreateElasticClientFromEnvVariables()
	if err != nil {
		glg.Fatal("failed to create client")
	}
	healthy := elastic.CheckHealthyStatusElasticSearch(elasticClient, 180)
	if !healthy {
		glg.Fatal("death has found me")
	}

	root := "sullego"
	rootDir, err := ioutil.ReadDir(root)
	if err != nil {
		glg.Fatal(err)
	}

	created := 0
	documents := 0
	for _, dir := range rootDir {
		glg.Debug("working on the following directory: " + dir.Name())
		if dir.IsDir() {
			filePath := path.Join(root, dir.Name())
			files, err := ioutil.ReadDir(filePath)
			if err != nil {
				glg.Fatal(err)
			}
			for _, f := range files {
				glg.Debug(fmt.Sprintf("found %s in %s", f.Name(), filePath))
				plan, _ := ioutil.ReadFile(path.Join(filePath, f.Name()))
				var logoi models.Logos
				err := json.Unmarshal(plan, &logoi)
				if err != nil {
					glg.Fatal(err)
				}

				documents += len(logoi.Logos)

				elastic.DeleteIndex(elasticClient, dir.Name())
				for _, word := range logoi.Logos {
					jsonifiedLogos, _ := word.Marshal()
					esRequest := esapi.IndexRequest{
						Body:       strings.NewReader(string(jsonifiedLogos)),
						Refresh:    "true",
						Index:      dir.Name(),
						DocumentID: "",
					}

					// Perform the request with the client.
					res, err := esRequest.Do(context.Background(), elasticClient)
					if err != nil {
						glg.Fatalf("Error getting response: %s", err)
					}
					defer res.Body.Close()

					if res.IsError() {
						glg.Debugf("[%s]", res.Status())
					} else {
						// Deserialize the response into a map.
						var r map[string]interface{}
						if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
							glg.Errorf("Error parsing the response body: %s", err)
						} else {
							// Print the response status and indexed document version.
							created++
						}
					}
				}
			}
		}
	}
	glg.Infof("created: %s", strconv.Itoa(created))
	glg.Infof("words found in sullego: %s", strconv.Itoa(documents))
	os.Exit(0)
}
