package httpclient

import (
	"crypto/tls"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

const (
	DefaultTimeout = 60 * time.Second
)

type RequestMethod string

const (
	MethodGet    RequestMethod = "GET"
	MethodPost   RequestMethod = "POST"
	MethodPut    RequestMethod = "PUT"
	MethodDelete RequestMethod = "DELETE"
)

// SendRequest 用于封装请求的所有参数
type SendRequest struct {
	URL         string
	Method      RequestMethod
	QueryParams map[string]string
	JSONData    map[string]interface{}
	FormData    map[string]string
	Headers     map[string]string
	Proxies     map[string]string
	Timeout     time.Duration
}

// NewRequest 创建一个新的 RequestConfig 实例，并设置默认值
func NewRequest(url string, method RequestMethod) *SendRequest {
	return &SendRequest{
		URL:         url,
		Method:      method,
		QueryParams: make(map[string]string),
		JSONData:    make(map[string]interface{}),
		FormData:    make(map[string]string),
		Headers:     make(map[string]string),
		Proxies:     make(map[string]string),
		Timeout:     DefaultTimeout,
	}
}

// Send 发送HTTP请求
func (req *SendRequest) Send() (bool, []byte, error) {
	client := resty.New()
	client.SetHeaders(req.Headers)
	client.SetTimeout(req.Timeout)
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})

	// 设置代理
	if len(req.Proxies) > 0 {
		client.SetProxy(fmt.Sprintf("%s://%s", req.Proxies["scheme"], req.Proxies["host"]))
	}

	// 准备请求
	request := client.R()
	request.SetQueryParams(req.QueryParams)

	var resp *resty.Response
	var err error

	// 注意：这里需要根据实际的请求方法常量进行调整
	switch req.Method {
	case MethodGet:
		resp, err = request.Get(req.URL)
	case MethodPost:
		if len(req.FormData) > 0 {
			request.SetFormData(req.FormData)
		} else {
			request.SetBody(req.JSONData)
		}
		resp, err = request.Post(req.URL)
	case MethodPut:
		if len(req.FormData) > 0 {
			request.SetFormData(req.FormData)
		} else {
			request.SetBody(req.JSONData)
		}
		resp, err = request.Put(req.URL)
	case MethodDelete:
		resp, err = request.Delete(req.URL)
	default:
		return false, nil, fmt.Errorf("不支持的请求方法: %s", req.Method)
	}

	if err != nil {
		return false, nil, fmt.Errorf("请求失败: %w", err)
	}

	// 检查响应状态码
	if resp.IsError() {
		return false, nil, fmt.Errorf("请求失败，状态代码: %d", resp.StatusCode())
	}

	return true, resp.Body(), nil
}
