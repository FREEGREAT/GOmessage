package miniorepo

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gomessage.com/media/internal/storage"
	"gomessage.com/media/pkg"
	client_minio "gomessage.com/media/pkg/minio"
)

const (
	empty_string = ""
)

var (
	cfg    = pkg.InitConfig()
	bucket = viper.GetString("minio.bucket")
	fileID string
)

type mediaRepository struct {
}

func SetFileID(id string) {
	fileID = id
}

// GetImgFile implements MediaRepository.
func (m *mediaRepository) GetImgFile(img_name string) (*minio.Object, error) {
	minioClient := client_minio.CreateConncet()

	object, err := minioClient.GetObject(context.Background(), bucket, img_name, minio.GetObjectOptions{})
	if err != nil {
		log.Println("Failed to get object:", err)
		return nil, err
	}

	return object, nil
}

// UploadImgFile implements MediaRepository.
func (m *mediaRepository) UploadImgFile(img_name string, file []byte, contentType string) (string, error) {
	minioClient := client_minio.CreateConncet()
	logrus.Infof("Uploading file %s to bucket %s", img_name, bucket)
	content, err := minioClient.PutObject(context.Background(), bucket, img_name, bytes.NewReader(file), int64(len(file)),
		minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.Println(err)
		return empty_string, err
	}

	logrus.Infof("Successfully uploaded %s to bucket %s", img_name, bucket)
	url, err := minioClient.PresignedGetObject(context.Background(), bucket, img_name, time.Hour*24*7, nil)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Uploaded", content.Key, "to", content.Bucket, content.ETag, content.VersionID, content.Size)
	fileID = img_name

	return url.Path, err

}

func NewMediaRepository() storage.MediaRepository {
	return &mediaRepository{}
}
