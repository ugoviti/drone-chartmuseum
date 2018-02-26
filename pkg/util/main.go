package util

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"code.gitea.io/git"
	"k8s.io/helm/pkg/chartutil"
)

// Extensions : the [...]T syntax is sugar for [123]T. It creates a fixed size array, but lets the compiler figure out how many elements are in it.
var Extensions = [...]string{".yaml", ".yml"}

// DeleteEmpty : to clean empty element from slice. See: http://dabase.com/e/15006/
func DeleteEmpty(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}

// SortStringSlice : little technique to sort slice to use in unit test. See: https://godoc.org/github.com/google/go-cmp/cmp#example-Option--SortedSlice
func SortStringSlice(in []string) []string {
	out := append([]string(nil), in...) // Copy input to avoid mutating it
	sort.Strings(out)
	return out
}

// GetUnique : return only unique elements from a predefined slice
func GetUnique(input []string) []string {
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

// IsDir : check if the input is directory or not
func IsDir(filePath string) bool {
	info, err := os.Stat(filePath)
	if err != nil {
		log.Print(err)
		return false
	}
	return info.IsDir()
}

// GetParentFolders : Get files's parent folder
func GetParentFolders(files []string) []string {
	var resultSlice []string
	for _, file := range files {
		dir := strings.Split(file, "/")[0]
		if IsDir(dir) {
			resultSlice = append(resultSlice, dir)
		}

	}
	return GetUnique(resultSlice)
}

// FilterExtFiles : Try to find glitch
func FilterExtFiles(files []string) []string {
	var resultSlice []string
	for _, ext := range Extensions {
		for _, file := range files {
			if filepath.Ext(file) == ext {
				resultSlice = append(resultSlice, file)

			}
		}
	}
	return resultSlice
}

// GetDiffFiles : similar to git diff, get the file changes between 2 commits
func GetDiffFiles(repoPath, previousCommitID, commitID string) []string {
	fmt.Printf("Getting diff between %v and %v ...\n", previousCommitID, commitID)
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

// SaveChartToPackage : save helm chart folder to compressed package
func SaveChartToPackage(chartPath string, dstPath string) (string, error) {
	var message string
	var err error
	if _, err := os.Stat(dstPath); os.IsNotExist(err) {
		os.Mkdir(dstPath, os.ModePerm)
	}

	if ok, _ := chartutil.IsChartDir(chartPath); ok == true {
		c, _ := chartutil.LoadDir(chartPath)
		message, err = chartutil.Save(c, dstPath)
		if err != nil {
			log.Printf("%v : %v", chartPath, err)
		}
		fmt.Printf("packaging %v ...\n", message)
	}

	return message, err
}

// UploadToServer : Upload chart package to chartmuseum server
func UploadToServer(filePaths []string, serverEndpoint string) {
	filePaths = DeleteEmpty(filePaths)
	for _, filePath := range filePaths {
		fmt.Printf("Uploading %v ...\n", filePath)
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
		fmt.Printf("%v \n", string(message))
	}

}
