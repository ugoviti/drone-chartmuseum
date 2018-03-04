package cmclient

import (
	"context"
	"errors"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
)

// ChartService handles communication with the Chart Manipulation
// related methods of the ChartMuseum API.
type ChartService service

// UploadChart uploads a Helm chart to the ChartMuseum server
func (s *ChartService) UploadChart(ctx context.Context, file *os.File) (*http.Response, error) {
	log.Printf("Uploading Chart %v ...\n", file)
	u := "api/charts"
	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}
	if stat.IsDir() {
		return nil, errors.New("Chart to upload can't be a directory")
	}
	mediaType := mime.TypeByExtension(filepath.Ext(file.Name()))
	req, err := s.client.NewUploadRequest(u, file, stat.Size(), mediaType)
	if err != nil {
		return nil, err
	}
	resp, err := s.client.Do(ctx, req)
	if err != nil {
		return resp, err
	}
	return resp, nil
}
