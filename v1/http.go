package core

import "encoding/json"

type APIError struct {
	Code    string      `json:"code,omitempty"`
	Message string      `json:"message,omitempty"`
	Details interface{} `json:"details,omitempty"`
}

type Metadata struct {
	SessionToken string `json:"session_token,omitempty"`
	RequestID    string `json:"request_id,omitempty"`
	RequestTime  string `json:"request_time,omitempty"`
}

type Response struct {
	Status   string      `json:"status,omitempty"`
	Error    *APIError   `json:"error,omitempty"`
	Data     interface{} `json:"data,omitempty"`
	Metadata Metadata    `json:"metadata,omitempty"`
}

func responseBytes(status string, data interface{}, apiErr *APIError, metadata Metadata) []byte {
	payload, err := json.Marshal(Response{
		Status:   status,
		Error:    apiErr,
		Data:     data,
		Metadata: metadata,
	})
	if err != nil {
		return []byte(`{"status":"error","error":{"code":"encoding_error","message":"failed to encode response"}}`)
	}
	return payload
}
