package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ishola-faazele/taskflow/internal/shared/domain_errors"
	"github.com/ishola-faazele/taskflow/internal/shared/logger"
)

// APIResponder handles all HTTP responses with consistent formatting
type APIResponder struct {
	logger *logger.StdLogger
}

// NewAPIResponder creates a new API responder with a logger
func NewAPIResponder() *APIResponder {
	return &APIResponder{
		logger: logger.NewStdLogger(),
	}
}

// ErrorResponsePayload represents a detailed error response
type ErrorResponsePayload struct {
	Success   bool                   `json:"success"`
	Error     string                 `json:"error"`
	Message   string                 `json:"message"`
	Details   map[string]interface{} `json:"details,omitempty"`
	Timestamp string                 `json:"timestamp"`
	Path      string                 `json:"path,omitempty"`
}

// SuccessResponsePayload represents a successful API response
type SuccessResponsePayload struct {
	Success   bool        `json:"success"`
	Data      interface{} `json:"data"`
	Message   string      `json:"message,omitempty"`
	Timestamp string      `json:"timestamp"`
}

// PaginatedResponsePayload represents a paginated API response
type PaginatedResponsePayload struct {
	Success    bool        `json:"success"`
	Data       interface{} `json:"data"`
	Message    string      `json:"message,omitempty"`
	Pagination Pagination  `json:"pagination"`
	Timestamp  string      `json:"timestamp"`
}

// Pagination contains pagination metadata
type Pagination struct {
	Page       int   `json:"page"`
	PerPage    int   `json:"per_page"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// Error sends a detailed error response
func (a *APIResponder) Error(w http.ResponseWriter, r *http.Request, statusCode int, message string, err error) {
	// Check if it's a domain error first
	if domainErr, ok := domain_errors.GetDomainError(err); ok {
		// Override status code with the one from domain error
		statusCode = int(domainErr.Code())
		
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)

		errorResponse := ErrorResponsePayload{
			Success:   false,
			Error:     http.StatusText(statusCode),
			Message:   domainErr.Message(),
			Details:   domainErr.Details(),
			Timestamp: time.Now().UTC().Format(time.RFC3339),
			Path:      r.URL.Path,
		}

		if encodeErr := json.NewEncoder(w).Encode(errorResponse); encodeErr != nil {
			a.logger.Error(fmt.Sprintf("Failed to encode error response: %v (original error: %v)", encodeErr, err))
		}

		// Log domain errors
		a.logger.Error(fmt.Sprintf("%s %s - StatusCode: %d,  Message: %s, Details: %v",
			r.Method, r.URL.Path, statusCode, domainErr.Message(), domainErr.Details()))
		return
	}

	// Handle non-domain errors
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	errorResponse := ErrorResponsePayload{
		Success:   false,
		Error:     http.StatusText(statusCode),
		Message:   message,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Path:      r.URL.Path,
	}

	// Add error details if available
	if err != nil {
		errorResponse.Details = map[string]interface{}{
			"error": err.Error(),
		}
	}

	if encodeErr := json.NewEncoder(w).Encode(errorResponse); encodeErr != nil {
		a.logger.Error(fmt.Sprintf("Failed to encode error response: %v (original error: %v)", encodeErr, err))
	}
	
	// Log non-domain errors
	a.logger.Error(fmt.Sprintf("%s %s - Status: %d, Message: %s, Error: %v",
		r.Method, r.URL.Path, statusCode, message, err))
}


// Success sends a success response with a custom message
func (a *APIResponder) Success(w http.ResponseWriter, r *http.Request, statusCode int, message string, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := SuccessResponsePayload{
		Success:   true,
		Data:      payload,
		Message:   message,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		a.logger.Error(fmt.Sprintf("Failed to encode JSON response: %v", err))
	}
}

// Paginated sends a paginated response
func (a *APIResponder) Paginated(w http.ResponseWriter, r *http.Request, message string, payload interface{}, page, perPage int, total int64) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	totalPages := int(total) / perPage
	if int(total)%perPage != 0 {
		totalPages++
	}

	response := PaginatedResponsePayload{
		Success: true,
		Data:    payload,
		Message: message,
		Pagination: Pagination{
			Page:       page,
			PerPage:    perPage,
			Total:      total,
			TotalPages: totalPages,
		},
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		a.logger.Error(fmt.Sprintf("Failed to encode paginated response: %v", err))
	}
}

// NoContent sends a 204 No Content response
func (a *APIResponder) NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

// Created sends a 201 Created response with location header
func (a *APIResponder) Created(w http.ResponseWriter, r *http.Request, location string, payload interface{}) {
	if location != "" {
		w.Header().Set("Location", location)
	}
	a.Success(w, r, http.StatusCreated, "Resource created successfully", payload)
}