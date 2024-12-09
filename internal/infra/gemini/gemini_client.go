package gemini

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"github.com/melkzsiqueira/water-gas-measurement/internal/dto"
	"google.golang.org/api/option"
)

type Gemini struct {
	model  string
	Gemini *genai.Client
}

func NewGeminiClient(apiKey, model string) (*Gemini, error) {
	client, err := genai.NewClient(context.Background(), option.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}
	return &Gemini{
		model:  model,
		Gemini: client,
	}, err
}

func (g *Gemini) ProcessImage(request dto.ProcessImageRequest, ctx context.Context) (dto.ProcessImageResponse, error) {
	defer g.Gemini.Close()

	model := g.Gemini.GenerativeModel(g.model)
	model.ResponseMIMEType = "application/json"
	model.ResponseSchema = &genai.Schema{
		Type:  genai.TypeArray,
		Items: &genai.Schema{Type: genai.TypeString},
	}

	image, err := base64.StdEncoding.DecodeString(request.Data)
	if err != nil {
		return dto.ProcessImageResponse{}, err
	}

	resp, err := model.GenerateContent(
		ctx,
		genai.Text("You are a meter reading expert. Extract the entire numeric value of a gas or water meter reading from this image in base64. Explicitly return only the integer numeric value."),
		genai.ImageData(strings.TrimPrefix(request.Mime, "image/"), image),
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
