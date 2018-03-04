package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"code.gitea.io/git"
	cm "github.com/honestbee/drone-chartmuseum/pkg/cmclient"
	"github.com/honestbee/drone-chartmuseum/pkg/util"
	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/ignore"
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
		Config     *Config
		Repository *git.Repository
		Commit     *git.Commit
		Client     *cm.Client

		CombinedChartPath string
	}

	// Chart holds path and parsed helmignore Rules
	Chart struct {
		Path  string
		Rules *ignore.Rules
	}
)

// ValidateConfig validates plugin configuration
func (p *Plugin) ValidateConfig() error {
	var err error
	// validate ChartMuseum baseURL
	if p.Client, err = cm.NewClient(p.Config.RepoURL, nil); err != nil {
		return err
	}

	// validate charts-dir is a directory
	if fi, err := os.Stat(p.Config.ChartsDir); err == nil {
		if !fi.IsDir() {
			return errors.New("charts-dir should be a directory")
		}
	} else {
		return err
	}

	if p.Config.CurrentCommitID != "" {
		// validate ChartsDir is a valid repository
		if p.Repository, err = git.OpenRepository(p.Config.ChartsDir); err != nil {
			return err
		}

		// validate CurrentCommitID is a valid commit in the repository
		if p.Commit, err = p.Repository.GetCommit(p.Config.CurrentCommitID); err != nil {
			return err
		}
	}

	if p.Config.ChartPath != "" {
		p.CombinedChartPath = filepath.Join(p.Config.ChartsDir, p.Config.ChartPath)
		// validate chart-path is a valid chart
		if valid, err := chartutil.IsChartDir(p.CombinedChartPath); !valid {
			return err
		}
	}

	return nil
}

func (p *Plugin) exec() error {
	ctx := context.Background()
	//ctx, cancel := context.WithTimeout(ctx, 60*time.Second)

	if err := p.ValidateConfig(); util.Pass(err) {
		var charts []string
		if charts, err = p.discoverCharts(); err != nil {
			return err
		}

		os.MkdirAll(p.Config.SaveDir, os.ModePerm)
		for _, chart := range charts {
			if c, err := p.packageChart(chart); util.Pass(err) {
				if f, err := os.Open(c); util.Pass(err) {
					message, err := p.Client.ChartService.UploadChart(ctx, f)
					fmt.Printf(message)
					return err
				}
			}
		}
	}

	return nil
}

// PackageChart saves a helm chart directory to a compressed package
func (p *Plugin) packageChart(chart string) (string, error) {
	c, err := chartutil.LoadDir(chart)
	if err != nil {
		return "", err
	}
	return chartutil.Save(c, p.Config.SaveDir)
}

// discoverCharts finds charts based on plugin configuration
func (p *Plugin) discoverCharts() (charts []string, err error) {
	if p.Config.ChartPath != "" {
		charts, err = []string{p.CombinedChartPath}, nil
	}

	if p.Config.CurrentCommitID != "" {
		modifiedCharts, err := p.findModifiedCharts()
		if err != nil {
			return nil, err
		}
		if p.Config.ChartPath != "" {
			if _, modified := modifiedCharts[p.CombinedChartPath]; !modified {
				fmt.Printf("%s wasn't modified.. nothing to do", p.Config.ChartPath)
				return nil, nil
			}
		} else {
			charts = util.Keys(modifiedCharts)
		}
	} else if p.Config.ChartPath == "" {
		charts, err = p.findAllCharts()
	}
	return charts, err
}

// findAllCharts recursively finds all charts within the configured charts-dir
func (p *Plugin) findAllCharts() (charts []string, err error) {
	fmt.Printf("Finding all charts...\n")
	walk := func(path string, stat os.FileInfo, err error) error {
		fmt.Printf("testing %s\n", path)
		if stat != nil && stat.IsDir() {
			if ok, _ := chartutil.IsChartDir(path); ok {
				fmt.Println("\tchart! jumping!\n")
				charts = append(charts, path)
				return filepath.SkipDir
			}
		}
		return nil
	}
	err = filepath.Walk(p.Config.ChartsDir, walk)
	return charts, err
}

// findModifiedCharts returns a map of all modified Charts filtered by .helmignore
func (p *Plugin) findModifiedCharts() (map[string]struct{}, error) {
	fmt.Printf("Getting diff between %v and %v ...\n", p.Config.PreviousCommitID, p.Config.CurrentCommitID)
	lookupCache := make(map[string]*Chart)
	modifiedCharts := make(map[string]struct{})
	files, err := p.Commit.GetFilesChangedSinceCommit(p.Config.PreviousCommitID)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		fi, err := os.Stat(file)
		if err != nil {
			fmt.Printf("\tIgnoring %s due to error: %v\n", file, err)
			continue // with next modified file
		}
		dirName := file
		if !fi.IsDir() {
			dirName = filepath.Dir(dirName)
		}
		c, err := getChart(dirName, p.Config.ChartsDir, lookupCache)
		if err != nil {
			fmt.Printf("\tIgnoring %s due to error: %v\n", file, err)
			continue // with next modified file
		}

		// flag chart modified if modified file not helmignored
		if !c.Rules.Ignore(file, fi) {
			modifiedCharts[c.Path] = struct{}{}
		}
	}
	return modifiedCharts, nil
}

// getChart recursively walks up the file tree to find the chart a directory belongs to
// Bug(vincent) this expects chartsDir to be valid prefix of filepath (both relative or absolute?)
func getChart(dirName string, chartsDir string, cache map[string]*Chart) (*Chart, error) {
	if cachedChart, _ := cache[dirName]; cachedChart != nil {
		return cachedChart, nil
	}

	c := &Chart{}
	if ok, _ := chartutil.IsChartDir(dirName); !ok {
		// are we at root dir?
		if strings.TrimPrefix(dirName, chartsDir) == "" {
			cache[dirName] = c
			return c, fmt.Errorf("Bailing! No chart found up to %s", chartsDir)
		}
		// search parent dir
		c, err := getChart(filepath.Dir(dirName), chartsDir, cache)

		cache[dirName] = c
		return c, err
	}
	c.Path = dirName
	parseHelmIgnoreRules(c)

	cache[dirName] = c
	return c, nil
}

// parseHelmIgnoreRules detects and loads helmignore Rules
func parseHelmIgnoreRules(c *Chart) error {
	c.Rules = ignore.Empty()
	ifile := filepath.Join(c.Path, ignore.HelmIgnore)
	if _, err := os.Stat(ifile); err == nil {
		r, err := ignore.ParseFile(ifile)
		if err != nil {
			return err
		}
		c.Rules = r
	}
	c.Rules.AddDefaults()
	return nil
}
