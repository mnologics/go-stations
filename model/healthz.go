package model

// A HealthzResponse expresses health check message.
type HealthzResponse struct {
	// Message contains the message returned in the health check response.
	Message string `json:"message"`
}
