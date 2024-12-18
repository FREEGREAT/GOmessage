package service

import (
	"context"

	proto_user_service "github.com/FREEGREAT/protos/gen/go/media"
)


type MediaService struct {

    // Тут можуть бути поля для зберігання стану

}


func NewMediaService() *MediaService {

    return &MediaService{}

}


func (s *MediaService) SavePhoto(ctx context.Context, req *media.SavePhotoRequest) (*media.SavePhotoResponse, error) {

    // Логіка для збереження фото в MinIO або інше сховище

    photoLink := "http://example.com/photo.jpg" // Згенероване посилання на фото
	

    return &media.SavePhotoResponse{PhotoLink: photoLink}, nil

}