package errors

import (
	"errors"
	"fmt"
)

// Wrap 将底层错误包装成新的错误信息
func Wrap(err error, message string) error {
	return fmt.Errorf("%s: %w", message, err)
}

// Wrapf 将底层错误包装成格式化的错误信息
func Wrapf(err error, format string, args ...interface{}) error {
	return fmt.Errorf("%s: %w", fmt.Sprintf(format, args...), err)
}

// Is 判断错误是否为指定类型的错误
func Is(err, target error) bool {
	return errors.Is(err, target)
}

// As 将错误转换为指定类型的错误
func As(err error, target interface{}) bool {
	return errors.As(err, target)
}

// New 创建一个新的错误
func New(text string) error {
	return errors.New(text)
}
