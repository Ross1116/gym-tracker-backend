package models

type ErrorResponse struct {
	Error string `json:"error" example:"Error message"`
}

type SuccessResponse struct {
	Success string `json:"success" example:"Success message"`
}
