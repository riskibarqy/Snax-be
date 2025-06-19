package domain

type APIResponse struct {
	Data  any       `json:"data,omitempty"`
	Error *APIError `json:"error,omitempty"`
}

type APIError struct {
	Message string `json:"message"`
	Detail  string `json:"detail,omitempty"`
	Code    int    `json:"code"`
}
