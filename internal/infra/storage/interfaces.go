package storage

import (
	"context"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type MeasurementStorageInterface interface {
	UploadFile(file string, ctx context.Context) (*uploader.UploadResult, error)
}
