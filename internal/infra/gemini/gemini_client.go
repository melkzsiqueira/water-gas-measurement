package gemini

import (
	"context"
	"encoding/base64"
	"encoding/json"

	"github.com/google/generative-ai-go/genai"
	"github.com/melkzsiqueira/water-gas-measurement/internal/dto"
	"google.golang.org/api/option"
)

type GeminiClient struct {
	APIKey string
	Model  string
	Gemini *genai.Client
}

func NewGeminiClient(apiKey, model string) *GeminiClient {
	return &GeminiClient{
		APIKey: apiKey,
		Model:  model,
		Gemini: &genai.Client{},
	}
}

func (g *GeminiClient) ProcessImage(ctx context.Context, request dto.ProcessImageRequest) (dto.ProcessImageResponse, error) {
	client, err := genai.NewClient(ctx, option.WithAPIKey(g.APIKey))
	if err != nil {
		return dto.ProcessImageResponse{}, err
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-pro-latest")
	model.ResponseMIMEType = "application/json"
	model.ResponseSchema = &genai.Schema{
		Type:  genai.TypeArray,
		Items: &genai.Schema{Type: genai.TypeString},
	}

	image, err := base64.StdEncoding.DecodeString(request.Image)
	if err != nil {
		return dto.ProcessImageResponse{}, err
	}

	resp, err := model.GenerateContent(
		ctx,
		genai.Text("You are a meter reading expert. Extract the entire numeric value of a gas or water meter reading from this image in base64. Explicitly return only the integer numeric value."),
		genai.ImageData("png", image),
	)
	if err != nil {
		return dto.ProcessImageResponse{}, err
	}

	var recipes []string
	for _, part := range resp.Candidates[0].Content.Parts {
		if txt, ok := part.(genai.Text); ok {
			if err := json.Unmarshal([]byte(txt), &recipes); err != nil {
				return dto.ProcessImageResponse{}, err
			}
		}
	}
	return dto.ProcessImageResponse{Value: recipes[0]}, err
}
