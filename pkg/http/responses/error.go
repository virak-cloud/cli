package responses

// ErrorResponse represents a standard error response from the API.
type ErrorResponse struct {
	Message string      `json:"message"`
	Errors  interface{} `json:"errors,omitempty"`
	Code    int         `json:"code,omitempty"`
}
