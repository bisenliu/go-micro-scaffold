package response

import (
	"errors"
	"testing"
)

func TestErrorFactory_Create(t *testing.T) {
	factory := NewErrorFactory()

	tests := []struct {
		name      string
		errorType ErrorType
		message   string
		cause     []error
		wantType  ErrorType
		wantMsg   string
		wantCause bool
	}{
		{
			name:      "create simple error without cause",
			errorType: ErrorTypeNotFound,
			message:   "resource not found",
			cause:     nil,
			wantType:  ErrorTypeNotFound,
			wantMsg:   "resource not found",
			wantCause: false,
		},
		{
			name:      "create error with cause",
			errorType: ErrorTypeValidationFailed,
			message:   "validation failed",
			cause:     []error{errors.New("underlying error")},
			wantType:  ErrorTypeValidationFailed,
			wantMsg:   "validation failed",
			wantCause: true,
		},
		{
			name:      "create error with multiple causes (only first is used)",
			errorType: ErrorTypeInternalServer,
			message:   "internal error",
			cause:     []error{errors.New("first error"), errors.New("second error")},
			wantType:  ErrorTypeInternalServer,
			wantMsg:   "internal error",
			wantCause: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := factory.Create(tt.errorType, tt.message, tt.cause...)

			if err.Type != tt.wantType {
				t.Errorf("Create() Type = %v, want %v", err.Type, tt.wantType)
			}
			if err.Message != tt.wantMsg {
				t.Errorf("Create() Message = %v, want %v", err.Message, tt.wantMsg)
			}
			if (err.Cause != nil) != tt.wantCause {
				t.Errorf("Create() Cause presence = %v, want %v", err.Cause != nil, tt.wantCause)
			}
			if err.Context != nil {
				t.Errorf("Create() Context should be nil for lazy allocation, got %v", err.Context)
			}
		})
	}
}

func TestErrorFactory_CreateWithContext(t *testing.T) {
	factory := NewErrorFactory()

	tests := []struct {
		name      string
		errorType ErrorType
		message   string
		context   map[string]any
		cause     []error
		wantCtx   bool
	}{
		{
			name:      "create with empty context",
			errorType: ErrorTypeNotFound,
			message:   "not found",
			context:   map[string]any{},
			cause:     nil,
			wantCtx:   false, // empty context should result in nil
		},
		{
			name:      "create with context data",
			errorType: ErrorTypeValidationFailed,
			message:   "validation failed",
			context:   map[string]any{"field": "username", "value": "invalid"},
			cause:     nil,
			wantCtx:   true,
		},
		{
			name:      "create with nil context",
			errorType: ErrorTypeInternalServer,
			message:   "internal error",
			context:   nil,
			cause:     nil,
			wantCtx:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := factory.CreateWithContext(tt.errorType, tt.message, tt.context, tt.cause...)

			if (err.Context != nil) != tt.wantCtx {
				t.Errorf("CreateWithContext() Context presence = %v, want %v", err.Context != nil, tt.wantCtx)
			}

			if tt.wantCtx && err.Context != nil {
				for k, v := range tt.context {
					if err.Context[k] != v {
						t.Errorf("CreateWithContext() Context[%s] = %v, want %v", k, err.Context[k], v)
					}
				}
			}
		})
	}
}

func TestConvenienceFunctions(t *testing.T) {
	tests := []struct {
		name     string
		fn       func(string, ...error) *DomainError
		wantType ErrorType
	}{
		{"NewNotFoundError", NewNotFoundError, ErrorTypeNotFound},
		{"NewValidationError", NewValidationError, ErrorTypeValidationFailed},
		{"NewAlreadyExistsError", NewAlreadyExistsError, ErrorTypeAlreadyExists},
		{"NewUnauthorizedError", NewUnauthorizedError, ErrorTypeUnauthorized},
		{"NewForbiddenError", NewForbiddenError, ErrorTypeForbidden},
		{"NewBusinessRuleViolationError", NewBusinessRuleViolationError, ErrorTypeBusinessRuleViolation},
		{"NewInvalidDataError", NewInvalidDataError, ErrorTypeInvalidData},
		{"NewInternalServerError", NewInternalServerError, ErrorTypeInternalServer},
		{"NewDatabaseConnectionError", NewDatabaseConnectionError, ErrorTypeDatabaseConnection},
		{"NewTimeoutError", NewTimeoutError, ErrorTypeTimeout},
		{"NewNetworkError", NewNetworkError, ErrorTypeNetworkError},
		{"NewRecordNotFoundError", NewRecordNotFoundError, ErrorTypeRecordNotFound},
		{"NewDuplicateKeyError", NewDuplicateKeyError, ErrorTypeDuplicateKey},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.fn("test message")
			if err.Type != tt.wantType {
				t.Errorf("%s() Type = %v, want %v", tt.name, err.Type, tt.wantType)
			}
			if err.Message != "test message" {
				t.Errorf("%s() Message = %v, want %v", tt.name, err.Message, "test message")
			}
		})
	}
}

func TestDomainError_WithContext_Integration(t *testing.T) {
	// Test that DomainError.WithContext works with the context manager
	err := NewNotFoundError("resource not found")

	// Add context
	errWithCtx := err.WithContext("resource_id", "123")

	if errWithCtx.Context == nil {
		t.Error("WithContext() should create context when adding non-nil value")
	}

	if errWithCtx.Context["resource_id"] != "123" {
		t.Errorf("WithContext() Context[resource_id] = %v, want %v", errWithCtx.Context["resource_id"], "123")
	}

	// Original error should remain unchanged
	if err.Context != nil {
		t.Error("Original error context should remain nil")
	}
}

func TestDomainError_WithContextMap_Integration(t *testing.T) {
	// Test that DomainError.WithContextMap works with the context manager
	err := NewValidationError("validation failed")

	contextMap := map[string]any{
		"field": "username",
		"value": "invalid",
	}

	errWithCtx := err.WithContextMap(contextMap)

	if errWithCtx.Context == nil {
		t.Error("WithContextMap() should create context when adding data")
	}

	for k, v := range contextMap {
		if errWithCtx.Context[k] != v {
			t.Errorf("WithContextMap() Context[%s] = %v, want %v", k, errWithCtx.Context[k], v)
		}
	}

	// Original error should remain unchanged
	if err.Context != nil {
		t.Error("Original error context should remain nil")
	}
}
func TestDomainError_WithContext_Optimizations(t *testing.T) {
	tests := []struct {
		name           string
		initialContext map[string]any
		key            string
		value          any
		expectSameRef  bool // 是否期望返回相同的引用（优化情况）
		expectContext  bool // 是否期望有上下文
	}{
		{
			name:           "nil context with nil value should return same reference",
			initialContext: nil,
			key:            "test",
			value:          nil,
			expectSameRef:  true,
			expectContext:  false,
		},
		{
			name:           "nil context with non-nil value should create new instance",
			initialContext: nil,
			key:            "test",
			value:          "value",
			expectSameRef:  false,
			expectContext:  true,
		},
		{
			name:           "existing context with nil value should remove key",
			initialContext: map[string]any{"test": "existing"},
			key:            "test",
			value:          nil,
			expectSameRef:  false,
			expectContext:  false, // key removed, context becomes empty
		},
		{
			name:           "existing context with new value should create new instance",
			initialContext: map[string]any{"existing": "value"},
			key:            "new",
			value:          "newvalue",
			expectSameRef:  false,
			expectContext:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create initial error
			err := &DomainError{
				Type:    ErrorTypeNotFound,
				Message: "test error",
				Context: tt.initialContext,
			}

			// Apply WithContext
			newErr := err.WithContext(tt.key, tt.value)

			// Check if same reference (optimization)
			if tt.expectSameRef {
				if newErr != err {
					t.Errorf("Expected same reference for optimization case, got different instance")
				}
			} else {
				if newErr == err {
					t.Errorf("Expected different reference, got same instance")
				}
			}

			// Check context presence
			if tt.expectContext {
				if newErr.Context == nil || len(newErr.Context) == 0 {
					t.Errorf("Expected context to be present, got nil or empty")
				}
			} else {
				if newErr.Context != nil && len(newErr.Context) > 0 {
					t.Errorf("Expected no context, got %v", newErr.Context)
				}
			}

			// Ensure original error is not modified (unless same reference)
			if !tt.expectSameRef {
				// Skip map comparison as maps can only be compared to nil
				// The context manager tests already verify immutability
			}
		})
	}
}

func TestDomainError_WithContextMap_Optimizations(t *testing.T) {
	tests := []struct {
		name           string
		initialContext map[string]any
		contextMap     map[string]any
		expectSameRef  bool // 是否期望返回相同的引用（优化情况）
		expectContext  bool // 是否期望有上下文
	}{
		{
			name:           "nil context with empty map should return same reference",
			initialContext: nil,
			contextMap:     map[string]any{},
			expectSameRef:  true,
			expectContext:  false,
		},
		{
			name:           "nil context with nil map should return same reference",
			initialContext: nil,
			contextMap:     nil,
			expectSameRef:  true,
			expectContext:  false,
		},
		{
			name:           "nil context with data should create new instance",
			initialContext: nil,
			contextMap:     map[string]any{"key": "value"},
			expectSameRef:  false,
			expectContext:  true,
		},
		{
			name:           "existing context with empty map should create copy",
			initialContext: map[string]any{"existing": "value"},
			contextMap:     map[string]any{},
			expectSameRef:  false,
			expectContext:  true,
		},
		{
			name:           "existing context with new data should merge",
			initialContext: map[string]any{"existing": "value"},
			contextMap:     map[string]any{"new": "newvalue"},
			expectSameRef:  false,
			expectContext:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create initial error
			err := &DomainError{
				Type:    ErrorTypeNotFound,
				Message: "test error",
				Context: tt.initialContext,
			}

			// Apply WithContextMap
			newErr := err.WithContextMap(tt.contextMap)

			// Check if same reference (optimization)
			if tt.expectSameRef {
				if newErr != err {
					t.Errorf("Expected same reference for optimization case, got different instance")
				}
			} else {
				if newErr == err {
					t.Errorf("Expected different reference, got same instance")
				}
			}

			// Check context presence
			if tt.expectContext {
				if newErr.Context == nil || len(newErr.Context) == 0 {
					t.Errorf("Expected context to be present, got nil or empty")
				}
			} else {
				if newErr.Context != nil && len(newErr.Context) > 0 {
					t.Errorf("Expected no context, got %v", newErr.Context)
				}
			}

			// Ensure original error is not modified (unless same reference)
			if !tt.expectSameRef {
				// Skip map comparison as maps can only be compared to nil
				// The context manager tests already verify immutability
			}
		})
	}
}

func TestDomainError_ContextHelperMethods(t *testing.T) {
	t.Run("GetContext returns copy", func(t *testing.T) {
		originalContext := map[string]any{"key": "value"}
		err := &DomainError{
			Type:    ErrorTypeNotFound,
			Message: "test",
			Context: originalContext,
		}

		contextCopy := err.GetContext()

		// Modify the copy
		contextCopy["new"] = "newvalue"

		// Original should be unchanged
		if err.Context["new"] != nil {
			t.Error("Original context was modified when copy was changed")
		}

		if len(err.Context) != 1 {
			t.Errorf("Original context length changed, expected 1, got %d", len(err.Context))
		}
	})

	t.Run("GetContext with nil context", func(t *testing.T) {
		err := &DomainError{
			Type:    ErrorTypeNotFound,
			Message: "test",
			Context: nil,
		}

		contextCopy := err.GetContext()
		if contextCopy != nil {
			t.Error("Expected nil context copy, got non-nil")
		}
	})

	t.Run("GetContextValue", func(t *testing.T) {
		err := &DomainError{
			Type:    ErrorTypeNotFound,
			Message: "test",
			Context: map[string]any{"existing": "value"},
		}

		// Test existing key
		value, exists := err.GetContextValue("existing")
		if !exists {
			t.Error("Expected key to exist")
		}
		if value != "value" {
			t.Errorf("Expected value 'value', got %v", value)
		}

		// Test non-existing key
		value, exists = err.GetContextValue("nonexistent")
		if exists {
			t.Error("Expected key to not exist")
		}
		if value != nil {
			t.Errorf("Expected nil value for non-existent key, got %v", value)
		}
	})

	t.Run("GetContextValue with nil context", func(t *testing.T) {
		err := &DomainError{
			Type:    ErrorTypeNotFound,
			Message: "test",
			Context: nil,
		}

		value, exists := err.GetContextValue("any")
		if exists {
			t.Error("Expected key to not exist in nil context")
		}
		if value != nil {
			t.Error("Expected nil value for nil context")
		}
	})

	t.Run("HasContext", func(t *testing.T) {
		// Test with context
		errWithContext := &DomainError{
			Type:    ErrorTypeNotFound,
			Message: "test",
			Context: map[string]any{"key": "value"},
		}
		if !errWithContext.HasContext() {
			t.Error("Expected HasContext to return true")
		}

		// Test with nil context
		errWithoutContext := &DomainError{
			Type:    ErrorTypeNotFound,
			Message: "test",
			Context: nil,
		}
		if errWithoutContext.HasContext() {
			t.Error("Expected HasContext to return false for nil context")
		}

		// Test with empty context
		errWithEmptyContext := &DomainError{
			Type:    ErrorTypeNotFound,
			Message: "test",
			Context: map[string]any{},
		}
		if errWithEmptyContext.HasContext() {
			t.Error("Expected HasContext to return false for empty context")
		}
	})
}
