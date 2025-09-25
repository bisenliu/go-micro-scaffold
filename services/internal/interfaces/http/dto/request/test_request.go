package request

// TestRequest 测试请求DTO
type TestRequest struct {
	Name  string `json:"name" binding:"required" label:"姓名"`
	Email string `json:"email" binding:"required,email" label:"邮箱"`
	Age   int    `json:"age" binding:"required,min=1,max=150" label:"年龄"`
}
