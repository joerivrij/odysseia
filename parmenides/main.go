package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/kpango/glg"
)

func init() {
	errlog := glg.FileWriter("/tmp/error.log", 0666)
	defer errlog.Close()

	glg.Get().
		SetMode(glg.BOTH).
		AddLevelWriter(glg.ERR, errlog)
}

func main() {
	glg.Info("\n _______  _______  ______    __   __  _______  __    _  ___   ______   _______  _______ \n|       ||   _   ||    _ |  |  |_|  ||       ||  |  | ||   | |      | |       ||       |\n|    _  ||  |_|  ||   | ||  |       ||    ___||   |_| ||   | |  _    ||    ___||  _____|\n|   |_| ||       ||   |_||_ |       ||   |___ |       ||   | | | |   ||   |___ | |_____ \n|    ___||       ||    __  ||       ||    ___||  _    ||   | | |_|   ||    ___||_____  |\n|   |    |   _   ||   |  | || ||_|| ||   |___ | | |   ||   | |       ||   |___  _____| |\n|___|    |__| |__||___|  |_||_|   |_||_______||_|  |__||___| |______| |_______||_______|\n")
	glg.Info(strings.Repeat("~", 37))
	glg.Info("\"τό γάρ αυτο νοειν έστιν τε καί ειναι\"")
	glg.Info("\"for it is the same thinking and being\"")
	glg.Info(strings.Repeat("~", 37))


	elasticClient := createElasticClient("changeme", "elastic")
	healthy := checkHealthyStatusElasticSearch(elasticClient, 60)
	if !healthy {
		glg.Fatal("death has found me")
	}
	root := "sullego"
	rootDir, err := ioutil.ReadDir(root)
	if err != nil {
		glg.Fatal(err)
	}

	created := 0
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
				var logoi Logos
				err := json.Unmarshal(plan, &logoi)
				if err != nil {
					glg.Fatal(err)
				}

				deleteIndex(elasticClient, dir.Name())
				for _, word := range logoi {
					jsonifiedLogos, _ := word.Marshal()
					esRequest := esapi.IndexRequest{
						Body:        strings.NewReader(string(jsonifiedLogos)),
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
	os.Exit(0)
}

func deleteIndex(es *elasticsearch.Client, index string) {
	glg.Warnf("deleting index: %s", index)

	res, err := es.Indices.Delete([]string{index})
	if err != nil {
		glg.Errorf("Error getting response: %s", err)
	}

	glg.Infof("status: %s", strconv.Itoa(res.StatusCode))

	responseBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	responseBody := string(responseBytes)

	switch res.StatusCode {
	case 200:
		glg.Infof("delete success: %s", responseBody)
	case 404:
		glg.Warnf("could not find index: %s", responseBody)
	default:
		glg.Errorf("something else went wrong: %s", responseBody)
	}

	return
}

func checkHealthyStatusElasticSearch(es *elasticsearch.Client, ticks time.Duration) bool {
	healthy := false

	ticker := time.NewTicker(1 * time.Second)
	timeout := time.After(ticks * time.Second)

	for {
		select {
		case t := <-ticker.C:
			glg.Infof("tick: %s", t)
			res, err := es.Info()
			if err != nil {
				glg.Errorf("Error getting response: %s", err)
				continue
			}
			defer res.Body.Close()
			// Check response status
			if res.IsError() {
				glg.Errorf("Error: %s", res.String())
			}

			var r map[string]interface{}

			// Deserialize the response into a map.
			if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
				glg.Errorf("Error parsing the response body: %s", err)
			}

			glg.Infof("serverVersion: %s", r["version"].(map[string]interface{})["number"])
			glg.Infof("serverName: %s", r["name"])
			glg.Infof("clusterName: %s", r["cluster_name"])
			healthy = true
			ticker.Stop()

		case <- timeout:
			ticker.Stop()
		}
		break
	}

	return healthy
}

func createElasticClient(password, username string) *elasticsearch.Client {
	glg.Info("creating elasticClient")

	cfg := elasticsearch.Config{
		Username: username,
		Password: password,
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		glg.Fatalf("Error creating the client: %s", err)
	}

	// Print client and server version numbers.
	glg.Infof("elasticClient version: %s", elasticsearch.Version)
	glg.Info(strings.Repeat("~", 37))

	return es
}