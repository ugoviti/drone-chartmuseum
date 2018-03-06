package cmclient

import (
	"context"
	"net/http"
	"os"

	"github.com/pkg/errors"
)

// ChartService handles communication with the Chart Manipulation
// related methods of the ChartMuseum API.
type ChartService service

// UploadChart uploads a Helm chart to the ChartMuseum server
func (s *ChartService) UploadChart(ctx context.Context, file *os.File) (*Response, error) {
	u := "api/charts"
	stat, err := file.Stat()
	if err != nil {
		return nil, errors.Wrap(err, "Unable to access file")
	}
	if stat.IsDir() {
		return nil, errors.New("Chart to upload can't be a directory")
	}
	//mediaType := mime.TypeByExtension(filepath.Ext(file.Name()))
	mediaType, _ := detectContentType(file)
	req, err := s.client.NewUploadRequest(u, file, stat.Size(), mediaType)
	if err != nil {
		return nil, errors.Wrap(err, "Failed creating uploading request")
	}
	resp, err := s.client.Do(ctx, req)
	if err != nil {
		return resp, errors.Wrap(err, "Failed to do request")
	}
	return resp, nil
}

// detectContentType returns a valid content-type and "application/octet-stream" if error or no match
func detectContentType(file *os.File) (string, error) {
	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)
	_, err := file.Read(buffer)
	if err != nil {
		return "application/octet-stream", err
	}

	// Reset the read pointer.
	file.Seek(0, 0)

	// Always returns a valid content-type and "application/octet-stream" if no others seemed to match.
	return http.DetectContentType(buffer), nil
}
