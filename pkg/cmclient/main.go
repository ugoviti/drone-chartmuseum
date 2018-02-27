package cmclient

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/honestbee/drone-chartmuseum/pkg/util"
)

// UploadToServer : Upload chart package to chartmuseum server
func UploadToServer(filePaths []string, serverEndpoint string) (err error) {
	if serverEndpoint == "" {
		err = errors.New("No valid RepoURL was defined")
		log.Print(err)
		return err
	}

	filePaths = util.DeleteEmpty(filePaths)
	for _, filePath := range filePaths {
		fmt.Printf("Uploading %v ...\n", filePath)
		file, err := os.Open(filePath)
		if err != nil {
			log.Print(err)
		}
		defer file.Close()
		resp, err := http.Post(serverEndpoint+"/api/charts", "application/octet-stream", file)
		if err != nil {
			log.Print(err)
		}
		defer resp.Body.Close()
		message, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Print(err)
		}
		fmt.Printf("%v \n", string(message))
	}
	return nil
}
