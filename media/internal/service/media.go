package service

import (
	"context"

	media "github.com/FREEGREAT/protos/gen/go/media"
)

type MediaService struct {
		media.UnimplementedMediaServiceServer
}

// Зміна конструктора
func NewMediaService() *MediaService {
	return &MediaService{}
}

// Зміна методу SavePhoto
func (s *MediaService) SavePhoto(ctx context.Context, req *media.SavePhotoRequest) (*media.SavePhotoResponse, error) {
	// Логіка для збереження фото
	photoLink := "http://example.com/photo.jpg" // Згенероване посилання на фото

	return &media.SavePhotoResponse{PhotoLink: photoLink}, nil
}
