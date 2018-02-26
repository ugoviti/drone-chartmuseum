package main

import (
	"io/ioutil"
	"log"

	"github.com/honestbee/drone-chartmuseum/pkg/util"
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

func (p *Plugin) defaultExec(files []string) {
	var resultList []string
	for _, file := range files {
		chart, err := util.SaveChartToPackage(file, p.Config.SaveDir)
		if err == nil {
			resultList = append(resultList, chart)
		}
	}
	util.UploadToServer(resultList, p.Config.RepoURL)
}

func (p *Plugin) exec() error {
	var files []string
	if p.Config.ChartPath != "" {
		files = []string{p.Config.ChartPath}
	} else if p.Config.PreviousCommitID != "" && p.Config.CurrentCommitID != "" {
		files = util.GetParentFolders(util.FilterExtFiles(util.GetDiffFiles(p.Config.ChartPath, p.Config.PreviousCommitID, p.Config.CurrentCommitID)))
	} else {
		dirs, err := ioutil.ReadDir(p.Config.ChartDir)
		if err != nil {
			log.Fatal(err)
		}
		files = util.ExtractDirs(dirs)
	}

	p.defaultExec(files)
	return nil
}
