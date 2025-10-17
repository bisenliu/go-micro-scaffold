package response

import (
	"errors"
	"net/http"
	"testing"
)

func TestUnifiedErrorHandler_HandleError(t *testing.T) {
	handler := NewUnifiedErrorHandler()

	tests := []struct {
		name     string
		err      error
		options  []ErrorHandleOption
		expected *ErrorResult
	}{
		{
			name:    "handle nil error without options",
			err:     nil,
			options: nil,
			expected: &ErrorResult{
				Code:       CodeSuccess,
				Message:    "操作成功",
				HTTPStatus: http.StatusOK,
			},
		},
		{
			name:    "handle error with code option",
			err:     nil,
			options: []ErrorHandleOption{WithCode(CodeValidation), WithMessage("自定义消息")},
			expected: &ErrorResult{
				Code:       CodeValidation,
				Message:    "自定义消息",
				HTTPStatus: http.StatusBadRequest,
			},
		},
		{
			name:    "handle domain error",
			err:     NewValidationError("验证失败"),
			options: nil,
			expected: &ErrorResult{
				Code:       CodeValidation,
				Message:    "验证失败",
				HTTPStatus: http.StatusBadRequest,
			},
		},
		{
			name:    "handle domain error with data",
			err:     NewNotFoundError("资源不存在"),
			options: []ErrorHandleOption{WithData(map[string]string{"resource": "user"})},
			expected: &ErrorResult{
				Code:       CodeNotFound,
				Message:    "资源不存在",
				HTTPStatus: http.StatusNotFound,
				Data:       map[string]string{"resource": "user"},
			},
		},
		{
			name:    "handle regular error",
			err:     errors.New("普通错误"),
			options: nil,
			expected: &ErrorResult{
				Code:       CodeInternalError,
				Message:    "普通错误",
				HTTPStatus: http.StatusInternalServerError,
			},
		},
		{
			name:    "handle error with custom message override",
			err:     NewValidationError("原始消息"),
			options: []ErrorHandleOption{WithMessage("覆盖消息")},
			expected: &ErrorResult{
				Code:       CodeValidation,
				Message:    "覆盖消息",
				HTTPStatus: http.StatusBadRequest,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := handler.HandleError(tt.err, tt.options...)

			if result.Code != tt.expected.Code {
				t.Errorf("expected code %d, got %d", tt.expected.Code, result.Code)
			}
			if result.Message != tt.expected.Message {
				t.Errorf("expected message %q, got %q", tt.expected.Message, result.Message)
			}
			if result.HTTPStatus != tt.expected.HTTPStatus {
				t.Errorf("expected HTTP status %d, got %d", tt.expected.HTTPStatus, result.HTTPStatus)
			}

			// Check data if expected
			if tt.expected.Data != nil {
				if result.Data == nil {
					t.Error("expected data but got nil")
				}
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
				if config.Code == nil || *config.Code != *tt.expected.Code {
					t.Errorf("expected code %v, got %v", tt.expected.Code, config.Code)
				}
			}
			if config.Message != tt.expected.Message {
				t.Errorf("expected message %q, got %q", tt.expected.Message, config.Message)
			}
			// For maps, we need to check differently since they can't be compared directly
			if tt.expected.Data != nil {
				if config.Data == nil {
					t.Error("expected data but got nil")
				}
				// For this test, we know it's a map[string]int, so we can type assert and compare
				if expectedMap, ok := tt.expected.Data.(map[string]int); ok {
					if actualMap, ok := config.Data.(map[string]int); ok {
						if len(expectedMap) != len(actualMap) {
							t.Errorf("expected data length %d, got %d", len(expectedMap), len(actualMap))
						}
						for k, v := range expectedMap {
							if actualMap[k] != v {
								t.Errorf("expected data[%s] = %d, got %d", k, v, actualMap[k])
							}
						}
					} else {
						t.Error("expected data to be map[string]int")
					}
				}
			} else if config.Data != nil {
				t.Errorf("expected nil data, got %v", config.Data)
			}
		})
	}
}

func TestBackwardCompatibility(t *testing.T) {
	// Test that unified error handler works with old API patterns
	oldHandler := NewUnifiedErrorHandler()

	// Test HandleError method
	result := oldHandler.HandleError(NewValidationError("测试"))
	if result.Code != CodeValidation {
		t.Errorf("expected code %d, got %d", CodeValidation, result.Code)
	}

	// Test HandleError with code option
	result = oldHandler.HandleError(nil, WithCode(CodeNotFound), WithMessage("自定义消息"))
	if result.Code != CodeNotFound {
		t.Errorf("expected code %d, got %d", CodeNotFound, result.Code)
	}
	if result.Message != "自定义消息" {
		t.Errorf("expected message %q, got %q", "自定义消息", result.Message)
	}

	// Test HandleError with data option
	testData := map[string]string{"key": "value"}
	result = oldHandler.HandleError(NewNotFoundError("测试"), WithData(testData))
	if result.Data == nil {
		t.Error("expected data to be preserved")
	}
}

func TestGlobalConvenienceFunctions(t *testing.T) {
	// Test HandleError function
	result := HandleError(NewValidationError("全局测试"))
	if result.Code != CodeValidation {
		t.Errorf("expected code %d, got %d", CodeValidation, result.Code)
	}

	// Test HandleErrorWithCode function
	result = HandleErrorWithCode(CodeNotFound, "全局代码测试")
	if result.Code != CodeNotFound {
		t.Errorf("expected code %d, got %d", CodeNotFound, result.Code)
	}

	// Test HandleErrorWithData function
	testData := "全局数据"
	result = HandleErrorWithData(NewInternalServerError("测试"), testData)
	if result.Data != testData {
		t.Error("expected data to be preserved")
	}
}
