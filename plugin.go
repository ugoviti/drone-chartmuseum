package main

import (
	"fmt"
	"log"

	"code.gitea.io/git"
	"k8s.io/helm/pkg/chartutil"
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
	c, _ := chartutil.LoadDir(chartPath)
	chartutil.Save(c, dstPath)
}
