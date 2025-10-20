package response

import (
	"encoding/json"
	"fmt"
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
	err := NewNotFoundError("test not found")

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

// TestNewArchitectureIntegration tests the integration of the new architecture
func TestNewArchitectureIntegration(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("ErrorFactory and ResponseEngine integration", func(t *testing.T) {
		// Create custom ErrorFactory
		customFactory := NewErrorFactory()
		engine := NewResponseEngineWithFactory(customFactory)

		// Test that errors created by ErrorFactory work with ResponseEngine
		router := gin.New()
		router.GET("/test", func(c *gin.Context) {
			err := CreateError(ErrorTypeValidationFailed, "custom validation error")
			engine.Handle(c, nil, err)
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response Response
		jsonErr := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, jsonErr)
		assert.Equal(t, CodeValidation, response.Code)
		assert.Equal(t, "custom validation error", response.Message)
	})

	t.Run("Global API functions delegate to default engine", func(t *testing.T) {
		// Test that global functions work correctly
		router := gin.New()
		router.GET("/test", func(c *gin.Context) {
			err := CreateErrorWithContext(ErrorTypeUnauthorized, "unauthorized access", map[string]any{"user": "test"})
			Handle(c, nil, err)
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response Response
		jsonErr := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, jsonErr)
		assert.Equal(t, CodeUnauthorized, response.Code)
		assert.Equal(t, "unauthorized access", response.Message)
	})

	t.Run("Multiple engines work independently", func(t *testing.T) {
		engine1 := NewResponseEngine()
		engine2 := NewResponseEngine()

		// Register different codes in each engine
		engine1.RegisterCode(9001, &CodeInfo{Code: 9001, Message: "Engine 1 Error", HTTPStatus: http.StatusBadRequest})
		engine2.RegisterCode(9002, &CodeInfo{Code: 9002, Message: "Engine 2 Error", HTTPStatus: http.StatusConflict})

		// Test engine1
		router1 := gin.New()
		router1.GET("/test", func(c *gin.Context) {
			engine1.HandleErrorWithCode(c, 9001, "")
		})

		w1 := httptest.NewRecorder()
		req1, _ := http.NewRequest("GET", "/test", nil)
		router1.ServeHTTP(w1, req1)

		assert.Equal(t, http.StatusBadRequest, w1.Code)

		// Test engine2
		router2 := gin.New()
		router2.GET("/test", func(c *gin.Context) {
			engine2.HandleErrorWithCode(c, 9002, "")
		})

		w2 := httptest.NewRecorder()
		req2, _ := http.NewRequest("GET", "/test", nil)
		router2.ServeHTTP(w2, req2)

		assert.Equal(t, http.StatusConflict, w2.Code)

		// Verify engines are independent - engine1 shouldn't have code 9002
		_, exists := engine1.GetCodeInfo(9002)
		assert.False(t, exists)
	})

	t.Run("All error types are properly mapped", func(t *testing.T) {
		errorTypesToTest := []struct {
			errorType      ErrorType
			expectedCode   int
			expectedStatus int
		}{
			{ErrorTypeNotFound, CodeNotFound, http.StatusNotFound},
			{ErrorTypeValidationFailed, CodeValidation, http.StatusBadRequest},
			{ErrorTypeAlreadyExists, CodeAlreadyExists, http.StatusConflict},
			{ErrorTypeUnauthorized, CodeUnauthorized, http.StatusUnauthorized},
			{ErrorTypeForbidden, CodeForbidden, http.StatusForbidden},
			{ErrorTypeInternalServer, CodeInternalError, http.StatusInternalServerError},
			{ErrorTypeTimeout, CodeTimeout, http.StatusRequestTimeout},
		}

		for i, tt := range errorTypesToTest {
			t.Run(fmt.Sprintf("ErrorType_%d", i), func(t *testing.T) {
				router := gin.New()
				router.GET("/test", func(c *gin.Context) {
					err := CreateError(tt.errorType, "test error")
					Handle(c, nil, err)
				})

				w := httptest.NewRecorder()
				req, _ := http.NewRequest("GET", "/test", nil)
				router.ServeHTTP(w, req)

				assert.Equal(t, tt.expectedStatus, w.Code)

				var response Response
				jsonErr := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, jsonErr)
				assert.Equal(t, tt.expectedCode, response.Code)
			})
		}
	})
}
