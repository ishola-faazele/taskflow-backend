package domain_errors

import (
	"errors"
	"fmt"
	"net/http"
)

// ErrorCode represents HTTP status codes for domain errors
type ErrorCode int

const (
	ErrCodeNotFound         ErrorCode = http.StatusNotFound          // 404
	ErrCodeValidation       ErrorCode = http.StatusBadRequest        // 400
	ErrCodeDatabase         ErrorCode = http.StatusInternalServerError // 500
	ErrCodeConflict         ErrorCode = http.StatusConflict          // 409
	ErrCodeUnauthorized     ErrorCode = http.StatusUnauthorized      // 401
	ErrCodeForbidden        ErrorCode = http.StatusForbidden         // 403
	ErrCodeInvalidOperation ErrorCode = http.StatusUnprocessableEntity // 422
	ErrCodeInternal         ErrorCode = http.StatusInternalServerError // 500
)

// DomainError is the base error interface for all domain errors
type DomainError interface {
	error
	Code() ErrorCode
	Message() string
	Details() map[string]interface{}
	Unwrap() error
}

// BaseError provides common functionality for all domain errors
type BaseError struct {
	code    ErrorCode
	message string
	details map[string]interface{}
	cause   error
}

func (e *BaseError) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("%s: %v", e.message, e.cause)
	}
	return e.message
}

func (e *BaseError) Code() ErrorCode {
	return e.code
}

func (e *BaseError) Message() string {
	return e.message
}

func (e *BaseError) Details() map[string]interface{} {
	if e.details == nil {
		return make(map[string]interface{})
	}
	return e.details
}

func (e *BaseError) Unwrap() error {
	return e.cause
}

// NotFoundError represents a resource not found error
type NotFoundError struct {
	*BaseError
	Resource string
	ID       string
}

func NewNotFoundError(resource, id string) *NotFoundError {
	return &NotFoundError{
		BaseError: &BaseError{
			code:    ErrCodeNotFound,
			message: fmt.Sprintf("%s with ID '%s' not found", resource, id),
			details: map[string]interface{}{
				"resource": resource,
				"id":       id,
			},
		},
		Resource: resource,
		ID:       id,
	}
}

// ValidationError represents a validation failure
type ValidationError struct {
	*BaseError
	Field  string
	Value  interface{}
	Reason string
}

func NewValidationError(field, reason string) *ValidationError {
	return &ValidationError{
		BaseError: &BaseError{
			code:    ErrCodeValidation,
			message: fmt.Sprintf("validation failed for field '%s': %s", field, reason),
			details: map[string]interface{}{
				"field":  field,
				"reason": reason,
			},
		},
		Field:  field,
		Reason: reason,
	}
}

func NewValidationErrorWithValue(field string, value any, reason string) *ValidationError {
	return &ValidationError{
		BaseError: &BaseError{
			code:    ErrCodeValidation,
			message: fmt.Sprintf("validation failed for field '%s': %s", field, reason),
			details: map[string]interface{}{
				"field":  field,
				"value":  value,
				"reason": reason,
			},
		},
		Field:  field,
		Value:  value,
		Reason: reason,
	}
}

// MultiValidationError represents multiple validation failures
type MultiValidationError struct {
	*BaseError
	Errors []*ValidationError
}

func NewMultiValidationError(errs []*ValidationError) *MultiValidationError {
	details := make(map[string]interface{})
	fieldErrors := make(map[string]string)

	for _, err := range errs {
		fieldErrors[err.Field] = err.Reason
	}
	details["fields"] = fieldErrors

	return &MultiValidationError{
		BaseError: &BaseError{
			code:    ErrCodeValidation,
			message: fmt.Sprintf("validation failed for %d field(s)", len(errs)),
			details: details,
		},
		Errors: errs,
	}
}

// DatabaseError represents a database operation failure
type DatabaseError struct {
	*BaseError
	Operation string
	Table     string
}

func NewDatabaseError(operation string, cause error) *DatabaseError {
	return &DatabaseError{
		BaseError: &BaseError{
			code:    ErrCodeDatabase,
			message: fmt.Sprintf("database operation '%s' failed", operation),
			details: map[string]interface{}{
				"operation": operation,
			},
			cause: cause,
		},
		Operation: operation,
	}
}

func NewDatabaseErrorWithTable(operation, table string, cause error) *DatabaseError {
	return &DatabaseError{
		BaseError: &BaseError{
			code:    ErrCodeDatabase,
			message: fmt.Sprintf("database operation '%s' on table '%s' failed", operation, table),
			details: map[string]interface{}{
				"operation": operation,
				"table":     table,
			},
			cause: cause,
		},
		Operation: operation,
		Table:     table,
	}
}

// ConflictError represents a resource conflict (e.g., duplicate entry)
type ConflictError struct {
	*BaseError
	Resource   string
	Constraint string
}

func NewConflictError(resource, constraint string) *ConflictError {
	return &ConflictError{
		BaseError: &BaseError{
			code:    ErrCodeConflict,
			message: fmt.Sprintf("%s already exists (constraint: %s)", resource, constraint),
			details: map[string]interface{}{
				"resource":   resource,
				"constraint": constraint,
			},
		},
		Resource:   resource,
		Constraint: constraint,
	}
}

// UnauthorizedError represents an authentication failure
type UnauthorizedError struct {
	*BaseError
	Reason string
}

func NewUnauthorizedError(reason string) *UnauthorizedError {
	return &UnauthorizedError{
		BaseError: &BaseError{
			code:    ErrCodeUnauthorized,
			message: fmt.Sprintf("unauthorized: %s", reason),
			details: map[string]interface{}{
				"reason": reason,
			},
		},
		Reason: reason,
	}
}

// ForbiddenError represents an authorization failure
type ForbiddenError struct {
	*BaseError
	Resource string
	Action   string
}

func NewForbiddenError(resource, action string) *ForbiddenError {
	return &ForbiddenError{
		BaseError: &BaseError{
			code:    ErrCodeForbidden,
			message: fmt.Sprintf("forbidden: insufficient permissions to %s %s", action, resource),
			details: map[string]interface{}{
				"resource": resource,
				"action":   action,
			},
		},
		Resource: resource,
		Action:   action,
	}
}

// InvalidOperationError represents an operation that cannot be performed
type InvalidOperationError struct {
	*BaseError
	Operation string
	Reason    string
}

func NewInvalidOperationError(operation, reason string) *InvalidOperationError {
	return &InvalidOperationError{
		BaseError: &BaseError{
			code:    ErrCodeInvalidOperation,
			message: fmt.Sprintf("invalid operation '%s': %s", operation, reason),
			details: map[string]interface{}{
				"operation": operation,
				"reason":    reason,
			},
		},
		Operation: operation,
		Reason:    reason,
	}
}

// InternalError represents an unexpected internal error
type InternalError struct {
	*BaseError
	Component string
}

func NewInternalError(component string, cause error) *InternalError {
	return &InternalError{
		BaseError: &BaseError{
			code:    ErrCodeInternal,
			message: fmt.Sprintf("internal error in %s", component),
			details: map[string]interface{}{
				"component": component,
			},
			cause: cause,
		},
		Component: component,
	}
}

// Error checking utilities

// IsNotFound checks if an error is a NotFoundError
func IsNotFound(err error) bool {
	var notFoundErr *NotFoundError
	return errors.As(err, &notFoundErr)
}

// IsValidation checks if an error is a ValidationError
func IsValidation(err error) bool {
	var validationErr *ValidationError
	return errors.As(err, &validationErr)
}

// IsDatabase checks if an error is a DatabaseError
func IsDatabase(err error) bool {
	var dbErr *DatabaseError
	return errors.As(err, &dbErr)
}

// IsConflict checks if an error is a ConflictError
func IsConflict(err error) bool {
	var conflictErr *ConflictError
	return errors.As(err, &conflictErr)
}

// IsUnauthorized checks if an error is an UnauthorizedError
func IsUnauthorized(err error) bool {
	var unauthorizedErr *UnauthorizedError
	return errors.As(err, &unauthorizedErr)
}

// IsForbidden checks if an error is a ForbiddenError
func IsForbidden(err error) bool {
	var forbiddenErr *ForbiddenError
	return errors.As(err, &forbiddenErr)
}

// GetDomainError extracts a DomainError from an error chain
func GetDomainError(err error) (DomainError, bool) {
	var domainErr DomainError
	if errors.As(err, &domainErr) {
		return domainErr, true
	}
	return nil, false
}

// GetErrorCode extracts the error code from an error
func GetErrorCode(err error) (ErrorCode, bool) {
	if domainErr, ok := GetDomainError(err); ok {
		return domainErr.Code(), true
	}
	return 0, false
}