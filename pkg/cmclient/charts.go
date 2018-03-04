package cmclient

import (
	"context"
	"errors"
	"log"
	"mime"
	"os"
	"path/filepath"
)

// ChartService handles communication with the Chart Manipulation
// related methods of the ChartMuseum API.
type ChartService service

// UploadChart uploads a Helm chart to the ChartMuseum server
func (s *ChartService) UploadChart(ctx context.Context, file *os.File) (string, error) {
	log.Printf("Uploading Chart %v ...\n", file)
	u := "api/charts"
	stat, err := file.Stat()
	if err != nil {
		return "", err
	}
	if stat.IsDir() {
		return "", errors.New("Chart to upload can't be a directory")
	}
	mediaType := mime.TypeByExtension(filepath.Ext(file.Name()))
	req, err := s.client.NewUploadRequest(u, file, stat.Size(), mediaType)
	if err != nil {
		return "", err
	}
	resp, err := s.client.Do(ctx, req)
	if err != nil {
		return resp, err
	}
	return resp, nil
}
