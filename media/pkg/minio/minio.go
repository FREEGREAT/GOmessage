package client_minio

import (
	"context"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gomessage.com/media/pkg"
)

func CreateConncet() *minio.Client {
	if err_cfg := pkg.InitConfig(); err_cfg != nil {
		logrus.Fatalf("error init config: %s", err_cfg.Error())
	}
	useSSL := false
	storage := viper.GetString("minio.bucket")
	endpoint := strings.TrimSpace(viper.GetString("minio.endpoint"))
	logrus.Infof("Connecting to MinIO at %s", endpoint)

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(viper.GetString("minio.accessKeyID"), viper.GetString("minio.secretKeyID"), ""),
		Secure: useSSL,
	})
	if err != nil {
		logrus.Fatalln(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	exists, err := minioClient.BucketExists(ctx, storage)
	if err != nil {
		logrus.Fatalf("Error while checking bucket: %v.   Bucketname: %s", err, storage)
	}
	if !exists {
		logrus.Printf("Bucket %s does not exist. Creating...", storage)
		err = minioClient.MakeBucket(ctx, storage, minio.MakeBucketOptions{})
		if err != nil {
			logrus.Fatalf("Error while creating bucket %s: %v", storage, err)
		}
		logrus.Printf("Bucket %s successfully created", storage)
	}

	return minioClient
}

func CreateBucket() {
	if err_cfg := pkg.InitConfig(); err_cfg != nil {
		logrus.Fatalf("error init config: %s", err_cfg.Error())
	}
	storage := viper.GetString("minio.bucket")

	minioClient := CreateConncet()
	bucketName := storage
	err := minioClient.MakeBucket(context.Background(), storage, minio.MakeBucketOptions{})
	if err != nil {
		logrus.Println(err)
		return
	}
	bucket_are_exist, err := minioClient.BucketExists(context.Background(), storage)
	if bucket_are_exist {
		logrus.Infof("Bucket: %s created", bucketName)
	} else {
		logrus.Error("Bucket does not exist.")
	}

}
