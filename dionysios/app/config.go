package app

import (
	"encoding/json"
	"github.com/odysseia/plato/elastic"
	"github.com/odysseia/plato/models"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func QueryRuleSet(es elastic.Client, index string) (*models.DeclensionConfig, error) {
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
			return nil, err
		}

		return declension, nil
	} else {
		query := es.Builder().MatchAll()
		response, err := es.Query().Match(index, query)

		if err != nil {
			return nil, err
		}
		var declensionConfig models.DeclensionConfig
		for _, jsonHit := range response.Hits.Hits {
			byteJson, err := json.Marshal(jsonHit.Source)
			if err != nil {
				return nil, err
			}
			declension, err := models.UnmarshalDeclension(byteJson)
			if err != nil {
				return nil, err
			}

			declensionConfig.Declensions = append(declensionConfig.Declensions, declension)
		}
		return &declensionConfig, nil
	}
}

func getJsonFilesFromAnaximander() (*models.DeclensionConfig, error) {
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

	anaximanderDirPath := filepath.Join(l, "anaximander", "arkho")
	anaximanderDir, err := ioutil.ReadDir(anaximanderDirPath)
	if err != nil {
		return nil, err
	}
	for _, subDir := range anaximanderDir {
		anaximanderNounsDir := filepath.Join(anaximanderDirPath, subDir.Name(), "*.json")
		declensionFiles, err := filepath.Glob(anaximanderNounsDir)
		if err != nil {
			return nil, err
		}

		for _, fpath := range declensionFiles {
			f, err := ioutil.ReadFile(fpath)
			if err != nil {
				return nil, err
			}

			declension, err := models.UnmarshalDeclension(f)
			if err != nil {
				return nil, err
			}

			declensionConfig.Declensions = append(declensionConfig.Declensions, declension)
		}
	}

	return &declensionConfig, nil
}
