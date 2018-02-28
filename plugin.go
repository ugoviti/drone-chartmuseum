package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"code.gitea.io/git"
	"github.com/honestbee/drone-chartmuseum/pkg/cmclient"
	"github.com/honestbee/drone-chartmuseum/pkg/util"
	"k8s.io/helm/pkg/chartutil"
)

type (

	// Config struct map with drone plugin parameters
	Config struct {
		RepoURL          string `json:"repo_url,omitempty"`
		ChartPath        string `json:"chart_path,omitempty"`
		ChartsDir        string `json:"charts_dir,omitempty"`
		SaveDir          string `json:"save_dir,omitempty"`
		PreviousCommitID string `json:"previous_commit_id,omitempty"`
		CurrentCommitID  string `json:"current_commit_id,omitempty"`
		IsSingle         bool   `json:"is_single,omitempty"`
	}

	// Plugin struct
	Plugin struct {
		Config Config
	}
)

// ValidateConfig :
func (p *Plugin) ValidateConfig() (err error) {
	if p.Config.RepoURL == "" {
		err = errors.New("RepoURL is not valid")
	}

	if p.Config.ChartPath == "" {
		p.Config.IsSingle = false
	} else {
		p.Config.IsSingle = true
	}
	return
}

func (p *Plugin) exec() (err error) {
	p.ValidateConfig()
	chartsMap := make(map[string]struct{})

	if p.Config.PreviousCommitID == "" {
		execute := p.FindCharts(p.ExtractAllCharts())
		chartsMap = execute()
	} else {
		execute := p.FindCharts(p.ExtractModifiedCharts())
		chartsMap = execute()
	}
	uploadPackages, err := p.SaveChartToPackage(chartsMap)
	cmclient.UploadToServer(uploadPackages, p.Config.RepoURL)
	return nil
}

// FindCharts : closure function, to return unique map of charts
func (p *Plugin) FindCharts(filesList []string) func() map[string]struct{} {

	foo := func() map[string]struct{} {
		chartsMap := make(map[string]struct{})
		if p.Config.IsSingle {
			if util.Contains(filesList, p.Config.ChartPath) {
				filesList = []string{p.Config.ChartPath}
			}
		}
		for _, dir := range filesList {
			if util.CheckValidChart(p.Config.ChartsDir + dir) {
				chartsMap[dir] = struct{}{}
			}
		}
		return chartsMap
	}
	return foo
}

// ExtractAllCharts : function to extract all folders
func (p *Plugin) ExtractAllCharts() []string {

	fileInfos, err := ioutil.ReadDir(p.Config.ChartsDir)
	if err != nil {
		log.Fatal(err)
	}

	return util.ExtractName(fileInfos)
}

// ExtractModifiedCharts : function to extract diff folders
func (p *Plugin) ExtractModifiedCharts() (filesList []string) {
	files, err := p.GetDiffFiles()
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		filesList = append(filesList, strings.Split(file, "/")[0])
	}
	return filesList
}

// GetDiffFiles : similar to git diff, get the file changes between 2 commits
func (p *Plugin) GetDiffFiles() ([]string, error) {
	fmt.Printf("Getting diff between %v and %v ...\n", p.Config.PreviousCommitID, p.Config.CurrentCommitID)
	repository, err := git.OpenRepository(p.Config.ChartsDir)
	if err != nil {
		return nil, err
	}

	commit, err := repository.GetCommit(p.Config.CurrentCommitID)
	if err != nil {
		return nil, err
	}

	files, err := commit.GetFilesChangedSinceCommit(p.Config.PreviousCommitID)
	if err != nil {
		return nil, err
	}

	return files, nil
}

// SaveChartToPackage : save helm chart folder to compressed package
func (p *Plugin) SaveChartToPackage(chartsMap map[string]struct{}) (messages []string, err error) {
	if _, err := os.Stat(p.Config.SaveDir); os.IsNotExist(err) {
		os.Mkdir(p.Config.SaveDir, os.ModePerm)
	}

	for chart := range chartsMap {
		c, _ := chartutil.LoadDir(p.Config.ChartsDir + chart)
		message, err := chartutil.Save(c, p.Config.SaveDir)

		if err != nil {
			log.Printf("%v : %v", chart, err)
		} else {
			messages = append(messages, message)
		}
		fmt.Printf("packaging %v ...\n", message)
	}

	return messages, err
}
