package response

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestResponseEngine_NewAPI(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		setupFunc      func(*gin.Context) (any, error)
		expectedStatus int
		expectedCode   int
	}{
		{
			name: "Handle success case",
			setupFunc: func(c *gin.Context) (any, error) {
				return map[string]string{"message": "success"}, nil
			},
			expectedStatus: http.StatusOK,
			expectedCode:   CodeSuccess,
		},
		{
			name: "Handle error case",
			setupFunc: func(c *gin.Context) (any, error) {
				return nil, CreateError(ErrorTypeNotFound, "resource not found")
			},
			expectedStatus: http.StatusNotFound,
			expectedCode:   CodeNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建测试引擎
			engine := NewResponseEngine()

			// 设置Gin路由
			router := gin.New()
			router.GET("/test", func(c *gin.Context) {
				data, err := tt.setupFunc(c)
				engine.Handle(c, data, err)
			})

			// 创建测试请求
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/test", nil)
			router.ServeHTTP(w, req)

			// 验证响应
			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestResponseEngine_HandleWith(t *testing.T) {
	gin.SetMode(gin.TestMode)

	engine := NewResponseEngine()

	router := gin.New()
	router.GET("/test", func(c *gin.Context) {
		err := CreateError(ErrorTypeValidationFailed, "validation failed")
		engine.HandleWith(c, nil, err, WithData(map[string]string{"field": "name"}))
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestResponseEngine_HandlePaging(t *testing.T) {
	gin.SetMode(gin.TestMode)

	engine := NewResponseEngine()

	router := gin.New()
	router.GET("/test", func(c *gin.Context) {
		data := []string{"item1", "item2", "item3"}
		engine.HandlePaging(c, data, 1, 10, 3, nil)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestResponseEngine_CreateError(t *testing.T) {
	engine := NewResponseEngine()

	err := engine.CreateError(ErrorTypeNotFound, "test error")

	assert.NotNil(t, err)
	assert.Equal(t, ErrorTypeNotFound, err.Type)
	assert.Equal(t, "test error", err.Message)
}

func TestResponseEngine_CreateErrorWithContext(t *testing.T) {
	engine := NewResponseEngine()

	context := map[string]any{"user_id": 123}
	err := engine.CreateErrorWithContext(ErrorTypeValidationFailed, "validation error", context)

	assert.NotNil(t, err)
	assert.Equal(t, ErrorTypeValidationFailed, err.Type)
	assert.Equal(t, "validation error", err.Message)
	assert.True(t, err.HasContext())

	value, exists := err.GetContextValue("user_id")
	assert.True(t, exists)
	assert.Equal(t, 123, value)
}

func TestResponseEngine_ComponentAccess(t *testing.T) {
	engine := NewResponseEngine()

	// 测试组件访问
	assert.NotNil(t, engine.GetErrorFactory())
	assert.NotNil(t, engine.GetContextManager())
	assert.NotNil(t, engine.GetCodeRegistry())
	assert.NotNil(t, engine.GetPool())
	assert.NotNil(t, engine.GetErrorHandler())
}

func TestResponseEngine_CodeRegistration(t *testing.T) {
	engine := NewResponseEngine()

	// 注册新的业务码
	customCode := 9999
	customInfo := &CodeInfo{
		Code:       customCode,
		Message:    "Custom error",
		HTTPStatus: http.StatusTeapot,
		Category:   "custom",
	}

	engine.RegisterCode(customCode, customInfo)

	// 验证注册成功
	registry := engine.GetCodeRegistry()
	info, exists := registry.GetInfo(customCode)
	assert.True(t, exists)
	assert.Equal(t, "Custom error", info.Message)
	assert.Equal(t, http.StatusTeapot, info.HTTPStatus)
}

func TestGlobalNewAPI(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/test", func(c *gin.Context) {
		data := map[string]string{"message": "success"}
		Handle(c, data, nil)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGlobalHandleWith(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/test", func(c *gin.Context) {
		err := CreateError(ErrorTypeNotFound, "not found")
		HandleWith(c, nil, err, WithCode(CodeNotFound))
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGlobalHandlePaging(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/test", func(c *gin.Context) {
		data := []int{1, 2, 3}
		HandlePaging(c, data, 1, 10, 3, nil)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestEngineManagement(t *testing.T) {
	// 测试默认引擎管理
	originalEngine := GetDefaultEngine()
	assert.NotNil(t, originalEngine)

	// 创建新引擎并设置为默认
	newEngine := NewResponseEngine()
	SetDefaultEngine(newEngine)

	// 验证默认引擎已更改
	currentEngine := GetDefaultEngine()
	assert.Equal(t, newEngine, currentEngine)

	// 恢复原始引擎
	SetDefaultEngine(originalEngine)
}

func TestGlobalCodeRegistration(t *testing.T) {
	customCode := 8888
	customInfo := &CodeInfo{
		Code:       customCode,
		Message:    "Global custom error",
		HTTPStatus: http.StatusBadRequest,
		Category:   "global_custom",
	}

	RegisterGlobalCode(customCode, customInfo)

	// 验证注册成功
	engine := GetDefaultEngine()
	registry := engine.GetCodeRegistry()
	info, exists := registry.GetInfo(customCode)
	assert.True(t, exists)
	assert.Equal(t, "Global custom error", info.Message)
}
func TestHandleSuccessWithPaging_ZeroPageSize(t *testing.T) {
	gin.SetMode(gin.TestMode)

	engine := NewResponseEngine()

	tests := []struct {
		name     string
		pageSize int
		total    int64
		expected int // expected totalPages
	}{
		{
			name:     "pageSize is zero",
			pageSize: 0,
			total:    100,
			expected: 0,
		},
		{
			name:     "pageSize is negative",
			pageSize: -5,
			total:    100,
			expected: 0,
		},
		{
			name:     "pageSize is positive",
			pageSize: 10,
			total:    100,
			expected: 10,
		},
		{
			name:     "pageSize is positive with remainder",
			pageSize: 7,
			total:    100,
			expected: 15, // ceil(100/7) = 15
		},
		{
			name:     "total is zero",
			pageSize: 10,
			total:    0,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 这个测试应该不会panic
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// 这个调用不应该引发panic
			assert.NotPanics(t, func() {
				engine.HandleSuccessWithPaging(c, []string{"item1", "item2"}, 1, tt.pageSize, tt.total)
			})

			// 验证响应状态码
			assert.Equal(t, 200, w.Code)
		})
	}
}

func TestOKWithPaging_ZeroPageSize(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 测试全局函数也不会panic
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	assert.NotPanics(t, func() {
		OKWithPaging(c, []string{"item1", "item2"}, 1, 0, 100)
	})

	assert.Equal(t, 200, w.Code)
}
