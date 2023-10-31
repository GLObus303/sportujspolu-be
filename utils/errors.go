package utils

import "github.com/globus303/sportujspolu/models"

func GetDefaultedError(err error, fallbackMessage string) models.ErrorResponse {
	errorMessage := err.Error()
	if errorMessage != "" {
		return GetError(errorMessage)
	}

	return GetError(fallbackMessage)
}

func GetError(errMessage string) models.ErrorResponse {
	return models.ErrorResponse{Error: errMessage}
}
