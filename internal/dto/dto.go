package dto

type CreateMeasurementInput struct {
	Value int    `json:"value"`
	Image string `json:"image"`
	Type  string `json:"type"`
	User  string `json:"user"`
}
