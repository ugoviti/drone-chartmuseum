package cmclient

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/pkg/errors"
)

const testuploadDir = "testdata/upload-test/"
const testuploadChart = "chart1-0.1.0.tgz"

func TestUploadChart(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	f, size, buffer, err := openTestChart()
	if err != nil {
		t.Error(err)
		return
	}

	// this is mocked ChartMuseum api contract
	mux.HandleFunc("/api/charts", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testHeader(t, r, "Content-Type", "application/x-gzip")
		testHeader(t, r, "Content-Length", size)
		testFormValues(t, r, values{})
		testBody(t, r, buffer)
		fmt.Fprint(w, `{"Saved": true}`)
	})

	response, err := client.ChartService.UploadChart(context.Background(), f)
	if err != nil {
		t.Error(err)
	} else if !response.Saved {
		t.Error(fmt.Errorf("Unexpected response from upload: %v", response))
	}
}

func openTestChart() (*os.File, string, []byte, error) {
	f, err := os.Open(filepath.Join(testuploadDir, testuploadChart))
	if err != nil {
		return nil, "", nil, errors.Wrap(err, "Unable to open testfile for test upload")
	}
	fi, err := f.Stat()
	if err != nil {
		return nil, "", nil, errors.Wrap(err, "Unable to get testfile size to verify test upload")
	}

	buffer := make([]byte, fi.Size())
	f.Read(buffer)
	// Reset the read pointer.
	f.Seek(0, 0)

	size := strconv.FormatInt(fi.Size(), 10)
	return f, size, buffer, nil
}
