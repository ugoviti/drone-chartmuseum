package chartmuseumclient

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"code.gitea.io/git"
	"github.com/honestbee/drone-chartmuseum/pkg/util"
	"k8s.io/helm/pkg/chartutil"
)

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
	filePaths = util.DeleteEmpty(filePaths)
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
