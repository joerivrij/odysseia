package main

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/kpango/glg"
	"github.com/odysseia/aristoteles"
	"github.com/odysseia/aristoteles/configs"
	"github.com/odysseia/plato/elastic"
	"github.com/odysseia/plato/models"
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

	documents := 0
	var authors models.Authors
	for _, dir := range rootDir {
		glg.Debug("working on the following directory: " + dir.Name())
		if dir.IsDir() {
			authors.Authors = append(authors.Authors, models.Author{
				Author: dir.Name(),
			})
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

				elastic.DeleteIndex(&herakleitosConfig.ElasticClient, herakleitosConfig.Index)
				for _, logos := range rhemai.Rhemai {
					jsonifiedLogos, _ := logos.Marshal()
					esRequest := esapi.IndexRequest{
						Body:       strings.NewReader(string(jsonifiedLogos)),
						Refresh:    "true",
						Index:      herakleitosConfig.Index,
						DocumentID: "",
					}

					// Perform the request with the client.
					res, err := esRequest.Do(context.Background(), &herakleitosConfig.ElasticClient)
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
							herakleitosConfig.Created++
						}
					}
				}
			}
		}
	}

	authorIndex := "authors"
	elastic.DeleteIndex(&herakleitosConfig.ElasticClient, authorIndex)

	for _, author := range authors.Authors {
		jsonifiedAuthor, _ := author.Marshal()
		esRequest := esapi.IndexRequest{
			Body:       strings.NewReader(string(jsonifiedAuthor)),
			Refresh:    "true",
			Index:      authorIndex,
			DocumentID: "",
		}

		// Perform the request with the client.
		res, err := esRequest.Do(context.Background(), &herakleitosConfig.ElasticClient)
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
				herakleitosConfig.Created++
			}
		}
	}

	glg.Infof("created: %s", strconv.Itoa(herakleitosConfig.Created))
	glg.Infof("words found in rhema: %s", strconv.Itoa(documents))
	glg.Infof("authors added: %s", strconv.Itoa(len(authors.Authors)))

	os.Exit(0)
}
