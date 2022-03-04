package exception

import (
	"encoding/json"
	"fmt"
)

var (
	ErrNotFound = NewAppError("not found", 404, "")
	ErrAlreadyExist = NewAppError("already exists", 400, "")
)

type AppError struct {
	Err              error  `json:"-"`
	Message          string `json:"message,omitempty"`
	DeveloperMessage string `json:"developer_message,omitempty"`
	Code             int `json:"code,omitempty"`
}

func NewAppError(message string, code int , developerMessage string) *AppError {
	return &AppError{
		Err:              fmt.Errorf(message),
		Code:             code,
		Message:          message,
		DeveloperMessage: developerMessage,
	}
}

func (e *AppError) Error() string {
	return e.Err.Error()
}

func (e *AppError) Unwrap() error { return e.Err }

func (e *AppError) Marshal() []byte {
	bytes, err := json.Marshal(e)
	if err != nil {
		return nil
	}
	return bytes
}

func UnauthorizedError(message string) *AppError {
	return NewAppError(message, 403, "")
}

func BadRequestError(message string) *AppError {
	return NewAppError(message, 400, "some thing wrong with user data")
}