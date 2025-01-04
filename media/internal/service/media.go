package service

import (
	"context"

	media "github.com/FREEGREAT/protos/gen/go/media"
	"github.com/sirupsen/logrus"
	"gomessage.com/media/internal/storage"
)

type MediaService struct {
	repo storage.MediaRepository
	media.UnimplementedMediaServiceServer
}

func NewMediaService(repo storage.MediaRepository) *MediaService {
	return &MediaService{repo: repo}
}

func (s *MediaService) SavePhoto(ctx context.Context, req *media.SavePhotoRequest) (*media.SavePhotoResponse, error) {
	name := "test_img_grpc.jpg"

	url, err := s.repo.UploadImgFile(name, req.Photo, "image/jpg")
	if err != nil {
		logrus.Errorf("Error while uploading image. %w", err)
	}
	return &media.SavePhotoResponse{PhotoLink: url}, nil
}
