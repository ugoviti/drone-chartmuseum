package util

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"code.gitea.io/git"
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

// ExtractDirs : Get Dir path from os.info
func ExtractDirs(fileInfos []os.FileInfo) []string {
	var resultList []string
	for _, fileInfo := range fileInfos {
		resultList = append(resultList, fileInfo.Name())
	}
	return resultList
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
