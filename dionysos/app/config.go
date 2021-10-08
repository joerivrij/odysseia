package app

import (
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/kpango/glg"
	"github.com/odysseia/plato/elastic"
	"github.com/odysseia/plato/models"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type DionysosConfig struct {
	ElasticClient elasticsearch.Client
	DictionaryIndex string
	Index string
	DeclensionConfig models.DeclensionConfig
}

const dictionaryIndexDefault = "alexandros"
const elasticIndexDefault = "dionysos"

func Get(es *elasticsearch.Client, declensionConfig *models.DeclensionConfig) (*DionysosConfig) {
	dictIndex := os.Getenv("DICTIONARY_INDEX")
	if dictIndex == "" {
		glg.Debugf("setting DICTIONARY_INDEX to default: %s", dictionaryIndexDefault)
		dictIndex = dictionaryIndexDefault
	}

	elasticIndex := os.Getenv("ELASTIC_INDEX")
	if elasticIndex == "" {
		glg.Debugf("setting ELASTIC_INDEX to default: %s", elasticIndexDefault)
		elasticIndex = elasticIndexDefault
	}

	config := &DionysosConfig{
		DictionaryIndex: dictIndex,
		Index: elasticIndex,
		ElasticClient: *es,
		DeclensionConfig: *declensionConfig,
	}

	return config
}

func QueryRuleSet(es *elasticsearch.Client, index string) *models.DeclensionConfig {
	var local bool

	localString := os.Getenv("ENV")
	if localString == "" || localString == "development" {
		local = true
	} else {
		local = false
	}

	if local {
		declension, err := getJsonFilesFromAnaximander()
		if err != nil {
			panic(fmt.Sprintf("Cannot start Dionysos encountered error when creating config: %s", err))
		}

		return declension
	} else {
		response, err := elastic.QueryWithMatchAll(*es, index)

		if err != nil {
			panic(fmt.Sprintf("Cannot start Dionysos encountered error when creating config: %s", err))
			return nil
		}
		var declensionConfig models.DeclensionConfig
		for _, jsonHit := range response.Hits.Hits {
			byteJson, err := json.Marshal(jsonHit.Source)
			if err != nil {
				panic(fmt.Sprintf("Cannot start Dionysos encountered error when creating config: %s", err))
				return nil
			}
			declension, err := models.UnmarshalDeclension(byteJson)
			if err != nil {
				panic(fmt.Sprintf("Cannot start Dionysos encountered error when creating config: %s", err))
				return nil
			}
			switch declension.Name {
			case "firstDeclension":
				declensionConfig.FirstDeclension = declension
			case "secondDeclension":
				declensionConfig.SecondDeclension = declension
			default:
				continue
			}
		}
		return &declensionConfig
	}
	return nil
}

func getJsonFilesFromAnaximander() (*models.DeclensionConfig, error){
	var declensionConfig models.DeclensionConfig
	_, callingFile, _, _ := runtime.Caller(0)
	callingDir := filepath.Dir(callingFile)
	dirParts := strings.Split(callingDir, string(os.PathSeparator))
	var odysseiaPath []string
	for i, part := range dirParts {
		if part == "odysseia" {
			odysseiaPath = dirParts[0 : i+1]
		}
	}
	l := "/"
	for _, path := range odysseiaPath {
		l = filepath.Join(l, path)
	}
	anaximanderDir := filepath.Join(l, "anaximander", "arkho", "nouns", "*.json")
	declensionFiles, err := filepath.Glob(anaximanderDir)
	if err != nil {
		panic(fmt.Sprintf("Cannot glob fixture files: %s", err))
	}

	for _, fpath := range declensionFiles {
		f, err := ioutil.ReadFile(fpath)
		if err != nil {
			panic(fmt.Sprintf("Cannot read fixture file: %s", err))
		}

		declension, err := models.UnmarshalDeclension(f)
		 if err != nil {
		 	return nil, err
		 }

		switch declension.Name {
		case "firstDeclension":
			declensionConfig.FirstDeclension = declension
		case "secondDeclension":
			declensionConfig.SecondDeclension = declension
		default:
			continue
		}
	}

	return &declensionConfig, nil
}