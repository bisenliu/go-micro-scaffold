package response

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestUnifiedAPI_Handle(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		data           any
		err            error
		expectedStatus int
		expectedCode   int
	}{
		{
			name:           "success case",
			data:           map[string]string{"message": "success"},
			err:            nil,
			expectedStatus: http.StatusOK,
			expectedCode:   CodeSuccess,
		},
		{
			name:           "error case - not found",
			data:           nil,
			err:            NewNotFoundError("resource not found"),
			expectedStatus: http.StatusNotFound,
			expectedCode:   CodeNotFound,
		},
		{
			name:           "error case - validation",
			data:           nil,
			err:            NewValidationError("validation failed"),
			expectedStatus: http.StatusBadRequest,
			expectedCode:   CodeValidation,
		},
		{
			name:           "error case - unauthorized",
			data:           nil,
			err:            NewUnauthorizedError("unauthorized"),
			expectedStatus: http.StatusUnauthorized,
			expectedCode:   CodeUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.GET("/test", func(c *gin.Context) {
				Handle(c, tt.data, tt.err)
			})

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/test", nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response Response
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, response.Code)
		})
	}
}

func TestUnifiedAPI_HandleWith(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		data           any
		err            error
		options        []ErrorHandleOption
		expectedStatus int
		expectedCode   int
		expectedMsg    string
	}{
		{
			name:           "success with data",
			data:           map[string]string{"result": "ok"},
			err:            nil,
			options:        nil,
			expectedStatus: http.StatusOK,
			expectedCode:   CodeSuccess,
		},
		{
			name:           "error with custom code",
			data:           nil,
			err:            NewValidationError("validation failed"),
			options:        []ErrorHandleOption{WithCode(CodeBusinessError)},
			expectedStatus: http.StatusBadRequest,
			expectedCode:   CodeBusinessError,
		},
		{
			name:           "error with custom message",
			data:           nil,
			err:            NewNotFoundError("not found"),
			options:        []ErrorHandleOption{WithMessage("Custom not found message")},
			expectedStatus: http.StatusNotFound,
			expectedCode:   CodeNotFound,
			expectedMsg:    "Custom not found message",
		},
		{
			name:           "error with data",
			data:           nil,
			err:            NewValidationError("validation failed"),
			options:        []ErrorHandleOption{WithData(map[string]string{"field": "username"})},
			expectedStatus: http.StatusBadRequest,
			expectedCode:   CodeValidation,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.GET("/test", func(c *gin.Context) {
				HandleWith(c, tt.data, tt.err, tt.options...)
			})

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/test", nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response Response
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, response.Code)

			if tt.expectedMsg != "" {
				assert.Equal(t, tt.expectedMsg, response.Message)
			}
		})
	}
}

func TestUnifiedAPI_HandlePaging(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		data           any
		page           int
		pageSize       int
		total          int64
		err            error
		expectedStatus int
		expectedCode   int
	}{
		{
			name:           "successful paging",
			data:           []string{"item1", "item2", "item3"},
			page:           1,
			pageSize:       10,
			total:          3,
			err:            nil,
			expectedStatus: http.StatusOK,
			expectedCode:   CodeSuccess,
		},
		{
			name:           "paging with error",
			data:           nil,
			page:           1,
			pageSize:       10,
			total:          0,
			err:            NewInternalServerError("database error"),
			expectedStatus: http.StatusInternalServerError,
			expectedCode:   CodeInternalError,
		},
		{
			name:           "paging with zero page size",
			data:           []string{"item1"},
			page:           1,
			pageSize:       0,
			total:          1,
			err:            nil,
			expectedStatus: http.StatusOK,
			expectedCode:   CodeSuccess,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.GET("/test", func(c *gin.Context) {
				HandlePaging(c, tt.data, tt.page, tt.pageSize, tt.total, tt.err)
			})

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/test", nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response Response
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, response.Code)

			if tt.err == nil {
				// Check paging structure
				assert.NotNil(t, response.Data)
			}
		})
	}
}

func TestErrorHandleOptions(t *testing.T) {
	tests := []struct {
		name     string
		options  []ErrorHandleOption
		expected ErrorHandleConfig
	}{
		{
			name:    "WithCode option",
			options: []ErrorHandleOption{WithCode(1001)},
			expected: ErrorHandleConfig{
				Code: func() *int { i := 1001; return &i }(),
			},
		},
		{
			name:    "WithMessage option",
			options: []ErrorHandleOption{WithMessage("测试消息")},
			expected: ErrorHandleConfig{
				Message: "测试消息",
			},
		},
		{
			name:    "WithData option",
			options: []ErrorHandleOption{WithData("测试数据")},
			expected: ErrorHandleConfig{
				Data: "测试数据",
			},
		},
		{
			name: "multiple options",
			options: []ErrorHandleOption{
				WithCode(1002),
				WithMessage("多选项测试"),
				WithData(map[string]int{"count": 5}),
			},
			expected: ErrorHandleConfig{
				Code:    func() *int { i := 1002; return &i }(),
				Message: "多选项测试",
				Data:    map[string]int{"count": 5},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &ErrorHandleConfig{}
			for _, option := range tt.options {
				option(config)
			}

			if tt.expected.Code != nil {
				assert.NotNil(t, config.Code)
				assert.Equal(t, *tt.expected.Code, *config.Code)
			} else {
				assert.Nil(t, config.Code)
			}

			assert.Equal(t, tt.expected.Message, config.Message)

			if tt.expected.Data != nil {
				assert.Equal(t, tt.expected.Data, config.Data)
			} else {
				assert.Nil(t, config.Data)
			}
		})
	}
}

func TestDomainError_ErrorInterface(t *testing.T) {
	tests := []struct {
		name     string
		err      *DomainError
		expected string
	}{
		{
			name: "error without cause",
			err: &DomainError{
				Type:    ErrorTypeNotFound,
				Message: "resource not found",
				Cause:   nil,
			},
			expected: "resource not found",
		},
		{
			name: "error with cause",
			err: &DomainError{
				Type:    ErrorTypeInternalServer,
				Message: "database error",
				Cause:   NewValidationError("validation failed"),
			},
			expected: "validation failed: database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.err.Error())
		})
	}
}

func TestDomainError_Unwrap(t *testing.T) {
	cause := NewValidationError("validation failed")
	err := &DomainError{
		Type:    ErrorTypeInternalServer,
		Message: "database error",
		Cause:   cause,
	}

	unwrapped := err.Unwrap()
	assert.Equal(t, cause, unwrapped)

	// Test nil cause
	errNoCause := &DomainError{
		Type:    ErrorTypeNotFound,
		Message: "not found",
		Cause:   nil,
	}

	unwrappedNil := errNoCause.Unwrap()
	assert.Nil(t, unwrappedNil)
}

func TestResponseEngine_ErrorCreation(t *testing.T) {
	engine := NewResponseEngine()

	tests := []struct {
		name      string
		errorType ErrorType
		message   string
		expected  ErrorType
	}{
		{"NotFound", ErrorTypeNotFound, "not found", ErrorTypeNotFound},
		{"Validation", ErrorTypeValidationFailed, "validation failed", ErrorTypeValidationFailed},
		{"Unauthorized", ErrorTypeUnauthorized, "unauthorized", ErrorTypeUnauthorized},
		{"Forbidden", ErrorTypeForbidden, "forbidden", ErrorTypeForbidden},
		{"Internal", ErrorTypeInternalServer, "internal error", ErrorTypeInternalServer},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := engine.CreateError(tt.errorType, tt.message)
			assert.Equal(t, tt.expected, err.Type)
			assert.Equal(t, tt.message, err.Message)
			assert.Nil(t, err.Context) // Should be nil for lazy allocation
		})
	}
}

func TestResponseEngine_ErrorCreationWithContext(t *testing.T) {
	engine := NewResponseEngine()

	context := map[string]any{
		"user_id": 123,
		"action":  "create",
	}

	err := engine.CreateErrorWithContext(ErrorTypeValidationFailed, "validation failed", context)

	assert.Equal(t, ErrorTypeValidationFailed, err.Type)
	assert.Equal(t, "validation failed", err.Message)
	assert.NotNil(t, err.Context)
	assert.Equal(t, 123, err.Context["user_id"])
	assert.Equal(t, "create", err.Context["action"])
}

func TestResponseEngine_StandardErrorCreators(t *testing.T) {
	engine := NewResponseEngine()

	// Test a few standard error creators
	notFoundErr := engine.NewNotFoundError("resource not found")
	assert.Equal(t, ErrorTypeNotFound, notFoundErr.Type)
	assert.Equal(t, "resource not found", notFoundErr.Message)

	validationErr := engine.NewValidationError("validation failed")
	assert.Equal(t, ErrorTypeValidationFailed, validationErr.Type)
	assert.Equal(t, "validation failed", validationErr.Message)

	unauthorizedErr := engine.NewUnauthorizedError("unauthorized")
	assert.Equal(t, ErrorTypeUnauthorized, unauthorizedErr.Type)
	assert.Equal(t, "unauthorized", unauthorizedErr.Message)
}

func TestGlobalErrorCreators(t *testing.T) {
	// Test global error creation functions
	notFoundErr := NewNotFoundError("global not found")
	assert.Equal(t, ErrorTypeNotFound, notFoundErr.Type)
	assert.Equal(t, "global not found", notFoundErr.Message)

	validationErr := NewValidationError("global validation")
	assert.Equal(t, ErrorTypeValidationFailed, validationErr.Type)
	assert.Equal(t, "global validation", validationErr.Message)

	internalErr := NewInternalServerError("global internal")
	assert.Equal(t, ErrorTypeInternalServer, internalErr.Type)
	assert.Equal(t, "global internal", internalErr.Message)
}

func TestCodeRegistryIntegration(t *testing.T) {
	// Test that code registry works with error mapping
	engine := NewResponseEngine()

	// Create a domain error
	err := engine.NewNotFoundError("test not found")

	// Test error handling
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/test", func(c *gin.Context) {
		engine.Handle(c, nil, err)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response Response
	jsonErr := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, jsonErr)
	assert.Equal(t, CodeNotFound, response.Code)

	// Verify the message comes from the domain error (not overridden by code registry)
	// The domain error message should be preserved
	assert.Equal(t, "test not found", response.Message)
}
