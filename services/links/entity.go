package links

import (
	"errors"
	"net/url"
	"time"
)

var ErrInternalServerError = errors.New("internal server error")

// ProcessBatchRequest ...
type ProcessBatchRequest struct {
	URLs []*url.URL
}

// ProcessBatchResponse ...
type ProcessBatchResponse struct {
	Results []Result
}

// GetBatchRequest ...
type GetBatchRequest struct {
	BatchID string `json:"batch_id"`
}

// GetBatchResponse ...
type GetBatchResponse struct {
	Results []Result
}

// Result model
type Result struct {
	ID               string    `json:"id"`
	BatchID          string    `json:"batch_id"`
	PageURL          string    `json:"page_url"`
	InternalLinksNum uint      `json:"internal_links_num"`
	ExternalLinksNum uint      `json:"external_links_num"`
	Success          bool      `json:"success"`
	Error            error     `json:"error"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// Response - generic http response structure
type Response struct {
	Errors []string    `json:"errors,omitempty"`
	Data   interface{} `json:"data,omitempty"`
}
