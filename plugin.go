package main

import (
	"fmt"
	"log"

	"k8s.io/helm/pkg/chartutil"

	"code.gitea.io/git"
)

func getDiff(repoPath, previousCommitID, commitID string) {
	repository, err := git.OpenRepository(repoPath)
	if err != nil {
		log.Fatal(err)
	}

	commit, err := repository.GetCommit(commitID)
	if err != nil {
		log.Fatal(err)
	}

	files, err := commit.GetFilesChangedSinceCommit(previousCommitID)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		fmt.Println(file)
	}
}

func saveChart(chartPath string, dstPath string) {
	if ok, err := chartutil.IsChartDir(chartPath); ok != true {
		log.Fatal(err)
	}
	c, _ := chartutil.LoadDir(chartPath)
	chartutil.Save(c, dstPath)
}
