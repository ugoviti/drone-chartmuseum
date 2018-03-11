package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"

	"code.gitea.io/git"
	log "github.com/sirupsen/logrus"
)

// testReposDir is a directory with bare repositories for git tests
const testReposDir = "testdata/charts-repos/"

// testChartsDir is a simple charts directory
const testChartsDir = "testdata/charts"

func TestMain(m *testing.M) {
	// read log level from environment prior to running tests
	if logLevelString, ok := os.LookupEnv("LOG_LEVEL"); ok {
		if level, err := log.ParseLevel(logLevelString); err == nil {
			log.SetLevel(level)
		}
	}
	os.Exit(m.Run())
}

func Test_exec(t *testing.T) {
	serverURL, mux, teardown := setupServerMock()
	defer teardown()

	// Only testing how many requests are received by ChartMuseum Mock
	// not testing the functionality of the cmclient
	var reqNum int
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		// will break if these tests are flagged t.Parallel()
		reqNum++
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "appliction/json")
		io.WriteString(w, `{"saved": true}`)
	})

	saveDir := filepath.Join(testChartsDir, "uploads")
	// plugin Exec() doesn't delete saveDir when done
	// this test cleans up after itself
	defer func() { os.RemoveAll(saveDir) }()

	tests := []struct {
		name   string
		p      *Plugin
		reqNum int
		err    error
	}{
		{
			name: "Test package and upload 1 chart",
			p: &Plugin{
				Config: &Config{
					RepoURL:   serverURL,
					ChartsDir: testChartsDir,
					ChartPath: "chart1",
					SaveDir:   saveDir,
				},
			},
			reqNum: 1,
			err:    nil,
		},
		{
			name: "Test package and upload all(4) test charts",
			p: &Plugin{
				Config: &Config{
					RepoURL:   serverURL,
					ChartsDir: testChartsDir,
					SaveDir:   saveDir,
				},
			},
			reqNum: 4,
			err:    nil,
		},
	}

	for _, test := range tests {
		reqNum = 0
		t.Run(test.name, func(t *testing.T) {
			if err := test.p.exec(); err != test.err {
				t.Error(err)
				return
			}
		})
		if reqNum != test.reqNum {
			t.Errorf("Want %d upload(s), got %d instead", test.reqNum, reqNum)
		}
	}
}

func Test_ValidateConfig(t *testing.T) {
	clonedPath, teardown := setupGitRepo("repo1_bare", t)
	defer teardown()

	// test invalid plugin configurations
	tests := []struct {
		name  string
		p     *Plugin
		valid bool
	}{
		{
			name: "Missing repo-url",
			p: &Plugin{
				Config: &Config{
					ChartsDir: testChartsDir,
				},
			},
			valid: false,
		},
		{
			name: "Invalid charts-dir",
			p: &Plugin{
				Config: &Config{
					ChartsDir: "foo",
				},
			},
			valid: false,
		},
		{
			name: "Invalid chart-path",
			p: &Plugin{
				Config: &Config{
					ChartsDir: testChartsDir,
					ChartPath: "foo",
				},
			},
			valid: false,
		},
		{
			name: "Invalid git commit",
			p: &Plugin{
				Config: &Config{
					ChartsDir:       clonedPath,
					CurrentCommitID: "foo",
				},
			},
			valid: false,
		},
	}

	var err error
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err = test.p.ValidateConfig()
			if test.valid && err != nil {
				t.Errorf("expected valid config, but got %v", err)
				return
			} else if !test.valid && err == nil {
				t.Error("expected invalid config, but config passed")
				return
			}
		})
	}

}

func Test_discoverCharts(t *testing.T) {
	// repo1 is a charts mono repo similar to kubernetes/charts
	clonedPath, teardown := setupGitRepo("repo1_bare", t)
	defer teardown()

	tests := []struct {
		name string
		p    *Plugin
		want []string
	}{
		{
			name: "Discover all charts",
			p: &Plugin{
				Config: &Config{
					RepoURL:   "http://charts.mycompany.com/",
					ChartsDir: testChartsDir,
				},
			},
			want: []string{
				filepath.Join(testChartsDir, "chart1"),
				filepath.Join(testChartsDir, "chart2"),
				filepath.Join(testChartsDir, "nested_charts/chart4"),
				filepath.Join(testChartsDir, "nested_charts/more_nesting/chart3"),
			},
		},
		{
			name: "Restrict to single valid chart_path",
			p: &Plugin{
				Config: &Config{
					RepoURL:   "http://charts.mycompany.com/",
					ChartsDir: testChartsDir,
					ChartPath: "chart1",
				},
			},
			want: []string{
				filepath.Join(testChartsDir, "chart1"),
			},
		},
		{
			name: "Discover all modified charts",
			p: &Plugin{
				Config: &Config{
					RepoURL:          "http://charts.mycompany.com/",
					ChartsDir:        clonedPath,
					CurrentCommitID:  "e3f490bf2cc020c8586d95699b8f82f873ec3968",
					PreviousCommitID: "4b825dc642cb6eb9a060e54bf8d69288fbee4904", //git's empty tree
				},
			},
			want: []string{
				filepath.Join(clonedPath, "chart1"),
				filepath.Join(clonedPath, "chart2"),
			},
		},
		{
			name: "Discover chart_path is modified",
			p: &Plugin{
				Config: &Config{
					RepoURL:          "http://charts.mycompany.com/",
					ChartsDir:        clonedPath,
					ChartPath:        "chart1",
					CurrentCommitID:  "e3f490bf2cc020c8586d95699b8f82f873ec3968",
					PreviousCommitID: "4b825dc642cb6eb9a060e54bf8d69288fbee4904", //git's empty tree
				},
			},
			want: []string{
				filepath.Join(clonedPath, "chart1"),
			},
		},
		{
			name: "Discover chart_path modifications are ignored",
			p: &Plugin{
				Config: &Config{
					RepoURL:          "http://charts.mycompany.com/",
					ChartsDir:        clonedPath,
					ChartPath:        "chart1",
					CurrentCommitID:  "ebd0cb15b6f07ff48663640e1c587fd696bad660",
					PreviousCommitID: "95e9e9bc9ecc78cf5a72de44ecced89e449e2a73",
				},
			},
			want: []string{},
		},
	}

	var err error
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err = test.p.ValidateConfig(); err != nil {
				t.Error(err)
				return
			}

			var got []string
			if got, err = test.p.discoverCharts(); err != nil {
				t.Error(err)
				return
			}
			if len(got) != len(test.want) {
				t.Errorf("Incorrect chart count got %v want %v", len(got), len(test.want))
			}
			sort.Strings(got)

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("Incorrect charts list got %#v want %#v", got, test.want)
			}
		})
	}
}

func Test_findModifiedCharts_repo2(t *testing.T) {
	// repo2 is a single-project repo with a single chart in a nested "deploy" directory
	clonedPath, teardown := setupGitRepo("repo2_bare", t)
	defer teardown()

	tests := []struct {
		p      *Plugin
		want   map[string]bool
		reason string
	}{
		{
			reason: "Initial Commit",
			p: &Plugin{
				Config: &Config{
					RepoURL:          "http://charts.mycompany.com/",
					ChartsDir:        clonedPath,
					CurrentCommitID:  "6ee3ced14a2388fa90f6dd27fcbb549442016866",
					PreviousCommitID: "4b825dc642cb6eb9a060e54bf8d69288fbee4904", //git's empty tree
				},
			},
			want: map[string]bool{
				filepath.Join(clonedPath, "deploy/chart"): true,
			},
		},
		{
			reason: "No change in Chart",
			p: &Plugin{
				Config: &Config{
					RepoURL:          "http://charts.mycompany.com/",
					ChartsDir:        clonedPath,
					CurrentCommitID:  "b1badf73d0eee1241a664cd2defc69c104f12b3b",
					PreviousCommitID: "6ee3ced14a2388fa90f6dd27fcbb549442016866",
				},
			},
			want: map[string]bool{},
		},
		{
			reason: "Bump Chart version",
			p: &Plugin{
				Config: &Config{
					RepoURL:          "http://charts.mycompany.com/",
					ChartsDir:        clonedPath,
					CurrentCommitID:  "7213fbe8dd8cb8b4bbd11fbdbc2a58cbb046a7b6",
					PreviousCommitID: "b1badf73d0eee1241a664cd2defc69c104f12b3b",
				},
			},
			want: map[string]bool{
				filepath.Join(clonedPath, "deploy/chart"): true,
			},
		},
	}

	for _, test := range tests {
		if err := test.p.ValidateConfig(); err != nil {
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

// setupGitRepo prepares a bare git repo through cloning and
// returns the location of the clone as well as a teardown function
func setupGitRepo(bareRepo string, t *testing.T) (clonedPath string, teardown func()) {
	clonedPath, err := cloneRepo(
		filepath.Join(testReposDir, bareRepo),
		testReposDir,
		strings.TrimSuffix(bareRepo, "_bare"))
	if err != nil {
		t.Error(err)
		return "", func() {}
	}

	fmt.Printf("cloned to %s\n", clonedPath)
	return clonedPath, func() {
		os.RemoveAll(clonedPath)
		fmt.Printf("cleaned up to %s\n", clonedPath)
	}
}

// setupServerMock sets up a mock HTTP server and a teardown function
func setupServerMock() (baseURL string, mux *http.ServeMux, teardown func()) {
	mux = http.NewServeMux()
	srv := httptest.NewServer(mux)
	return srv.URL, mux, srv.Close
}

// cloneRepo is a helper function to prepare git repositories for testing
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
