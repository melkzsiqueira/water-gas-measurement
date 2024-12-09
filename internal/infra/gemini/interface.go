package gemini

import (
	"context"

	"github.com/melkzsiqueira/water-gas-measurement/internal/dto"
)

type GeminiInterface interface {
	ProcessImage(request dto.ProcessImageRequest, ctx context.Context) (dto.ProcessImageResponse, error)
}
