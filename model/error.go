package model

import "time"

type (
	// A Error expresses ...
	ErrNotFound struct {
		Message   string    `json:"message"`
		CreatedAt time.Time `json:"created_at"`
	}

	// An Error interface expresses ...
	Error interface {
		Error() string
	}
)

func NewErrNotFound(message string) *ErrNotFound {
	return &ErrNotFound{
		Message:   message,
		CreatedAt: time.Now(),
	}
}

func (e *ErrNotFound) Error() string {
	return e.Message
}
