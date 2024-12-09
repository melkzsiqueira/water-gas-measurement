package dto

type CreateMeasurementImageInput struct {
	Mime string `json:"mime"`
	Data string `json:"data"`
}

type CreateMeasurementInput struct {
	Value     int                         `json:"value"`
	Image     CreateMeasurementImageInput `json:"image"`
	Type      string                      `json:"type"`
	Confirmed bool                        `json:"confirmed"`
	User      string                      `json:"user"`
}

type CreateUserInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type GetTokenInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type GetTokenOutput struct {
	AccessToken string `json:"access_token"`
}

type ProcessImageRequest struct {
	Mime string `json:"mime"`
	Data string `json:"data"`
}

type ProcessImageResponse struct {
	Value string `json:"value"`
}
