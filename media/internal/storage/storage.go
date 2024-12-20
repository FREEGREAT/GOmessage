package storage

import (
	"github.com/minio/minio-go/v7"
)

type MediaRepository interface {
	UploadImgFile(img_name string, file []byte, contentType string) (string, error)
	GetImgFile(img_name string) (*minio.Object, error)
}
