package storage

import (
	"context"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type Storage struct {
	storage *cloudinary.Cloudinary
}

func NewStorage(cloudName, apiKey, apiSecret string) (*Storage, error) {
	cld, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	return &Storage{
		storage: cld,
	}, err
}

func (s *Storage) UploadFile(file string, ctx context.Context) (*uploader.UploadResult, error) {
	resp, err := s.storage.Upload.Upload(ctx, file, uploader.UploadParams{})
	if err != nil {
		return nil, err
	}
	return resp, err
}
