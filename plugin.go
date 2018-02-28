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

	return
}

func (p *Plugin) exec() (err error) {
	p.ValidateConfig()
	chartsMap := make(map[string]struct{})

	if p.Config.PreviousCommitID == "" {
		chartsMap = p.FindCharts(p.ExtractAllDirs())
	} else {
		chartsMap = p.FindCharts(p.ExtractModifiedDirs())
	}

	uploadPackages, err := p.SaveChartToPackage(chartsMap)
	cmclient.UploadToServer(uploadPackages, p.Config.RepoURL)
	return nil
}

// FindCharts : closure function, to return unique map of charts
func (p *Plugin) FindCharts(filesMap map[string]struct{}) map[string]struct{} {
	chartsMap := make(map[string]struct{})
	if p.Config.ChartPath != "" {
		if util.Contains(filesMap, p.Config.ChartPath) {
			for k := range filesMap {
				delete(filesMap, k)
			}
			filesMap[p.Config.ChartPath] = struct{}{}
		}
	}
	for file := range filesMap {
		if util.CheckValidChart(p.Config.ChartsDir + file) {
			chartsMap[file] = struct{}{}
		}
	}

	return chartsMap
}

// ExtractAllDirs : function to extract all folders
func (p *Plugin) ExtractAllDirs() map[string]struct{} {

	fileInfos, err := ioutil.ReadDir(p.Config.ChartsDir)
	if err != nil {
		log.Fatal(err)
	}

	return util.ExtractName(fileInfos, p.Config.ChartsDir)
}

// ExtractModifiedDirs : function to extract diff folders
func (p *Plugin) ExtractModifiedDirs() map[string]struct{} {
	filesDir := make(map[string]struct{})
	filesMap, err := p.GetDiffFiles()
	if err != nil {
		log.Fatal(err)
	}

	for file := range filesMap {
		filesDir[strings.Split(file, "/")[0]] = struct{}{}
	}

	return filesDir
}

// GetDiffFiles : similar to git diff, get the file changes between 2 commits
func (p *Plugin) GetDiffFiles() (map[string]struct{}, error) {
	fmt.Printf("Getting diff between %v and %v ...\n", p.Config.PreviousCommitID, p.Config.CurrentCommitID)
	filesMap := make(map[string]struct{})
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

	for _, file := range files {
		filesMap[file] = struct{}{}
	}

	return filesMap, nil
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
