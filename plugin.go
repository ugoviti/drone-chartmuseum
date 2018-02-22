package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"k8s.io/helm/pkg/chartutil"

	"code.gitea.io/git"
)

// Just for clarification: the [...]T syntax is sugar for [123]T. It creates a fixed size array, but lets the compiler figure out how many elements are in it.
var extensions = [...]string{".yaml", ".yml"}

type config struct {
	RepoURL          string
	ChartPath        string
	ChartDir         string
	SaveDir          string
	PreviousCommitID string
	CurrentCommitID  string
}

func getUnique(input []string) []string {
	u := make([]string, 0, len(input))
	m := make(map[string]bool)

	for _, val := range input {
		if _, ok := m[val]; !ok {
			m[val] = true
			u = append(u, val)
		}
	}

	return u
}

func getUniqueParentFolders(files []string) []string {
	var resultSlice []string
	for _, file := range files {
		dir := strings.Split(file, "/")[0]
		if info, _ := os.Stat(dir); info.IsDir() {
			resultSlice = append(resultSlice, dir)
		}

	}
	return getUnique(resultSlice)
}

func filterExtFiles(files []string) []string {
	var resultSlice []string
	for _, ext := range extensions {
		for _, file := range files {
			if filepath.Ext(file) == ext {
				resultSlice = append(resultSlice, file)

			}
		}
	}
	return resultSlice
}

func getDiffFiles(repoPath, previousCommitID, commitID string) []string {
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

	return files
}

func saveChartToPackage(chartPath string, dstPath string) string {
	if _, err := os.Stat(dstPath); os.IsNotExist(err) {
		os.Mkdir(dstPath, os.ModePerm)
	}

	if ok, _ := chartutil.IsChartDir(chartPath); ok != true {
		log.Println("chart is not valid!")
	}

	c, _ := chartutil.LoadDir(chartPath)
	message, err := chartutil.Save(c, dstPath)
	if err != nil {
		log.Fatal(err)
	}
	return message
}

func uploadToServer(filePaths []string, serverEndpoint string) {
	for _, filePath := range filePaths {
		fmt.Printf("Uploading %v ...", filePath)
		file, err := os.Open(filePath)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		resp, err := http.Post(serverEndpoint+"/api/charts", "application/octet-stream", file)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		message, _ := ioutil.ReadAll(resp.Body)
		fmt.Printf(string(message))
	}

}
