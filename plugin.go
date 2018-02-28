package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
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
		if p.Config.ChartPath != "" {
			chartsMap = p.GetValidCharts(map[string]struct{}{
				p.Config.ChartPath: struct{}{},
			})
		} else {
			chartsMap = p.FindAllCharts()
		}
	} else {
		if p.Config.ChartPath != "" {
			if util.Contains(p.FindModifiedCharts(), p.Config.ChartPath) {
				chartsMap = p.GetValidCharts(map[string]struct{}{
					p.Config.ChartPath: struct{}{},
				})
			}
		} else {
			chartsMap = p.FindModifiedCharts()

		}
	}

	uploadPackages, err := p.SaveChartToPackage(chartsMap)
	cmclient.UploadToServer(uploadPackages, p.Config.RepoURL)
	return nil
}

// GetValidCharts : Get valid helm charts map
func (p *Plugin) GetValidCharts(filesMap map[string]struct{}) map[string]struct{} {
	chartsMap := make(map[string]struct{})
	for file := range filesMap {
		if ok, _ := chartutil.IsChartDir(filepath.Join(p.Config.ChartsDir, file)); ok == true {
			chartsMap[file] = struct{}{}
		}
	}

	return chartsMap
}

// FindAllCharts : function to extract all folders
func (p *Plugin) FindAllCharts() map[string]struct{} {

	fileInfos, err := ioutil.ReadDir(p.Config.ChartsDir)
	if err != nil {
		log.Fatal(err)
	}

	filesMap := util.ExtractName(fileInfos, p.Config.ChartsDir)
	return p.GetValidCharts(filesMap)
}

// FindModifiedCharts : function to extract diff folders
func (p *Plugin) FindModifiedCharts() map[string]struct{} {
	filesDir := make(map[string]struct{})
	filesDiffMap, err := p.GetDiffFiles()
	if err != nil {
		log.Fatal(err)
	}

	for file := range filesDiffMap {
		if strings.Contains(file, "terraform") == false {
			filesDir[strings.Split(file, "/")[0]] = struct{}{}
		}

	}

	return p.GetValidCharts(filesDir)
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

	if len(files) == 0 {
		return nil, nil
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
		c, _ := chartutil.LoadDir(filepath.Join(p.Config.ChartsDir, chart))
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
