package util

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Extensions : the [...]T syntax is sugar for [123]T. It creates a fixed size array, but lets the compiler figure out how many elements are in it.
var Extensions = [...]string{".yaml", ".yml"}

// Contains : check if element exists in string slice
func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// DeleteEmpty : to clean empty element from slice. See: http://dabase.com/e/15006/
func DeleteEmpty(s []string) (r []string) {
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}

// ExtractDirs : Get Dir path from os.info
func ExtractDirs(fileInfos []os.FileInfo) (resultList []string) {
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
func IsDir(filePath string) (bool, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return false, err
	}
	return info.IsDir(), nil
}

// GetParentFolders : Get files's parent folder
func GetParentFolders(files []string) (resultSlice []string) {
	for _, file := range files {
		dir := strings.Split(file, "/")[0]
		if ok, _ := IsDir(dir); ok {
			resultSlice = append(resultSlice, dir)
		}

	}
	return GetUnique(resultSlice)
}

// FilterExtFiles : Try to find glitch
func FilterExtFiles(files []string) (resultSlice []string) {
	for _, ext := range Extensions {
		for _, file := range files {
			if filepath.Ext(file) == ext {
				resultSlice = append(resultSlice, file)

			}
		}
	}
	return resultSlice
}
