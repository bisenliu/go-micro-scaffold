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

func TestResponseEngine_ComponentAccess(t *testing.T) {
	engine := NewResponseEngine()

	// 测试组件访问
	assert.NotNil(t, engine.GetContextManager())

	// 测试业务码查询功能
	info, exists := engine.GetCodeInfo(CodeSuccess)
	assert.True(t, exists)
	assert.Equal(t, "操作成功", info.Message)

	// 测试消息获取
	message := engine.GetCodeMessage(CodeSuccess)
	assert.Equal(t, "操作成功", message)

	// 测试HTTP状态码获取
	status := engine.GetHTTPStatus(CodeSuccess)
	assert.Equal(t, http.StatusOK, status)
}

func TestResponseEngine_CodeRegistration(t *testing.T) {
	engine := NewResponseEngine()

	// 注册新的业务码
	customCode := 9999
	customInfo := &CodeInfo{
		Code:       customCode,
		Message:    "Custom error",
		HTTPStatus: http.StatusTeapot,
	}

	engine.RegisterCode(customCode, customInfo)

	// 验证注册成功
	info, exists := engine.GetCodeInfo(customCode)
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
	}

	RegisterGlobalCode(customCode, customInfo)

	// 验证注册成功
	engine := GetDefaultEngine()
	info, exists := engine.GetCodeInfo(customCode)
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

func TestHandlePaging_ZeroPageSize(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 测试全局函数也不会panic
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	assert.NotPanics(t, func() {
		HandlePaging(c, []string{"item1", "item2"}, 1, 0, 100, nil)
	})

	assert.Equal(t, 200, w.Code)
}

// TestResponseEngine_ErrorFactoryIntegration tests the integration between ResponseEngine and ErrorFactory
func TestResponseEngine_ErrorFactoryIntegration(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Test default constructor creates engine with ErrorFactory
	engine1 := NewResponseEngine()
	assert.NotNil(t, engine1, "Engine should be created")

	// Test custom ErrorFactory constructor
	customFactory := NewErrorFactory()
	engine2 := NewResponseEngineWithFactory(customFactory)
	assert.NotNil(t, engine2, "Engine should be created with custom factory")

	// Test that ResponseEngine properly handles errors created by ErrorFactory
	router := gin.New()
	router.GET("/test", func(c *gin.Context) {
		// Create error using global ErrorFactory function
		err := CreateError(ErrorTypeNotFound, "resource not found")
		engine2.Handle(c, nil, err)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	// Test multiple engines work independently
	engine3 := NewResponseEngineWithFactory(NewErrorFactory())
	assert.NotNil(t, engine3, "Third engine should be created")

	// Test that all engines can handle the same error types correctly
	router2 := gin.New()
	router2.GET("/test", func(c *gin.Context) {
		err := CreateError(ErrorTypeValidationFailed, "validation error")
		engine3.Handle(c, nil, err)
	})

	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/test", nil)
	router2.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusBadRequest, w2.Code)
}
