package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/kpango/glg"
	"github.com/odysseia/aristoteles"
	"github.com/odysseia/aristoteles/configs"
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

	rootDir, err := ioutil.ReadDir(root)
	if err != nil {
		glg.Fatal(err)
	}
	elastic.DeleteIndex(&anaximanderConfig.ElasticClient, anaximanderConfig.Index)

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
				var declensions models.Declension
				err := json.Unmarshal(plan, &declensions)
				upload, _ := declensions.Marshal()
				esRequest := esapi.IndexRequest{
					Body:       strings.NewReader(string(upload)),
					Refresh:    "true",
					Index:      anaximanderConfig.Index,
					DocumentID: "",
				}

				// Perform the request with the client.
				res, err := esRequest.Do(context.Background(), &anaximanderConfig.ElasticClient)
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
						anaximanderConfig.Created++
					}
				}
			}
		}
	}
	glg.Infof("created: %s", strconv.Itoa(anaximanderConfig.Created))
	os.Exit(0)

}
