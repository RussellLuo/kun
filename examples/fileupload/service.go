package fileupload

import (
	"context"
	"log"

	"github.com/RussellLuo/kok/pkg/codec/httpcodec"
)

//go:generate kokgen ./service.go Service

type Service interface {
	// @kok(op): POST /upload
	// @kok(success): statusCode:204
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
