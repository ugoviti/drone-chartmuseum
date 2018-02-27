package cmclient

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/honestbee/drone-chartmuseum/pkg/util"
)

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
