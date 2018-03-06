package main

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"code.gitea.io/git"
)

const testReposDir = "testdata/charts-repos/"
const testChartsDir = "testdata/charts"

// cloneRepo is a helper function to prepare repositories for testing
func cloneRepo(url, dir, name string) (string, error) {
	repoDir := filepath.Join(dir, name)
	if _, err := os.Stat(repoDir); err == nil {
		return repoDir, nil
	}
	return repoDir, git.Clone(url, repoDir, git.CloneRepoOptions{
		Mirror:  false,
		Bare:    false,
		Quiet:   true,
		Timeout: 5 * time.Minute,
	})
}

func Test_repo1_modifiedCharts(t *testing.T) {
	// repo1 is a mono repo like Honestbee
	bareRepo1Path := filepath.Join(testReposDir, "repo1_bare")
	clonedPath, err := cloneRepo(bareRepo1Path, testReposDir, "repo1")
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("cloned to %s\n", clonedPath)

	tests := []struct {
		p      *Plugin
		want   map[string]bool
		reason string
	}{
		{
			reason: "Modifications are all under ignored paths",
			p: &Plugin{
				Config: &Config{
					RepoURL:          "http://charts.mycompany.com/",
					ChartsDir:        clonedPath,
					CurrentCommitID:  "ebd0cb15b6f07ff48663640e1c587fd696bad660",
					PreviousCommitID: "95e9e9bc9ecc78cf5a72de44ecced89e449e2a73",
				},
			},
			want: map[string]bool{},
		},
		{
			reason: "Initial commit of both chart1 and chart2",
			p: &Plugin{
				Config: &Config{
					RepoURL:          "http://charts.mycompany.com/",
					ChartsDir:        clonedPath,
					CurrentCommitID:  "e3f490bf2cc020c8586d95699b8f82f873ec3968",
					PreviousCommitID: "4b825dc642cb6eb9a060e54bf8d69288fbee4904", //git's empty tree
				},
			},
			want: map[string]bool{
				filepath.Join(clonedPath, "chart1"): true,
				filepath.Join(clonedPath, "chart2"): true,
			},
		},
	}

	for _, test := range tests {
		err = test.p.ValidateConfig()

		if err != nil {
			t.Error(err)
			return
		}

		got, err := test.p.findModifiedCharts()
		if err != nil {
			t.Error(err)
			return
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("Incorrect modified charts map got %#v want %#v - Reason: %s", got, test.want, test.reason)
		}
	}

}
