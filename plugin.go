package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/honestbee/drone-chartmuseum/pkg/util"
	"github.com/urfave/cli"
)

type (

	// Config struct map with drone plugin parameters
	Config struct {
		RepoURL          string `json:"repo_url,omitempty"`
		ChartPath        string `json:"chart_path,omitempty"`
		ChartDir         string `json:"chart_dir,omitempty"`
		SaveDir          string `json:"save_dir,omitempty"`
		PreviousCommitID string `json:"previous_commit_id,omitempty"`
		CurrentCommitID  string `json:"current_commit_id,omitempty"`
	}

	// Plugin struct
	Plugin struct {
		Config Config
	}
)

func extractDirs(fileInfos []os.FileInfo) []string {
	var resultList []string
	for _, fileInfo := range fileInfos {
		resultList = append(resultList, fileInfo.Name())
	}
	return resultList
}

func executeAction(files []string, conf Config) {
	var resultList []string
	for _, file := range files {
		chart, err := util.SaveChartToPackage(file, conf.SaveDir)
		if err == nil {
			resultList = append(resultList, chart)
		}
	}
	util.UploadToServer(resultList, conf.RepoURL)
}

func allMode(c *cli.Context, conf Config) error {
	dirs, err := ioutil.ReadDir(conf.ChartDir)
	if err != nil {
		log.Fatal(err)
	}

	executeAction(extractDirs(dirs), conf)
	return nil
}

func diffMode(c *cli.Context, conf Config) error {
	files := util.GetDiffFiles(conf.ChartDir, conf.PreviousCommitID, conf.CurrentCommitID)
	files = util.GetParentFolders(util.FilterExtFiles(files))
	if len(files) == 0 {
		fmt.Print("No chart needs to be updated! Exit ... \n")
		os.Exit(0)
	}
	executeAction(files, conf)
	return nil
}

func singleMode(c *cli.Context, conf Config) error {
	executeAction([]string{conf.ChartPath}, conf)
	return nil
}
