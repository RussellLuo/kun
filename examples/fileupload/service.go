package fileupload

import (
	"context"
	"log"

	"github.com/RussellLuo/kok/pkg/httpcodec"
)

//go:generate kungen ./service.go Service

// Service is used for uploading files.
type Service interface {
	// Upload uploads a file.
	//kun:op POST /upload
	//kun:success statusCode=204
	Upload(ctx context.Context, file *httpcodec.FormFile) (err error)
}

type Uploader struct{}

func (u *Uploader) Upload(ctx context.Context, file *httpcodec.FormFile) (err error) {
	if file != nil {
		log.Printf("saved file: %s", file.Name)
		return file.Save(file.Name)
	}
	return nil
}
