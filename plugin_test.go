package main

import (
	"path/filepath"
	"testing"
)

const testReposDir = "fixtures/chart-repos/"

func Test_findModifiedCharts(t *testing.T) {
	bareRepo1Path := filepath.Join(testReposDir, "repo1_bare")
	p = &Plugin{
		&Config {
			RepoURL: "http://charts.mycompany.com/",
			ChartsDir: bareRepo1Path,
			CurrentCommitID: "HEAD",
			PreviousCommitID: "HEAD~1",
		}
	}
	m, err := p.findModifiedCharts()

	fmt.Printf("%v",m)
	
}
