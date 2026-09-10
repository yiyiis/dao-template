package errors

import (
	stderr "errors"
	"fmt"
	"github.com/pkg/errors"
)

// StackErr 堆栈错误类型
type StackErr struct {
	err error
}

func (e *StackErr) Error() string {
	return e.err.Error()
}

func (e *StackErr) Format(s fmt.State, verb rune) {
	if format, ok := e.err.(interface{ Format(s fmt.State, verb rune) }); ok {
		format.Format(s, verb)
	} else {
		fmt.Fprintf(s, "%+v", e.err)
	}
}

// New 创建一个携带堆栈信息的 error
func New(msg string) *StackErr {
	return &StackErr{
		errors.New(msg),
	}
}

// Join 将多个错误合成一个错误
func Join(errs ...error) error {
	return stderr.Join(errs...)
}

// Is 判断 error 是否为目标错误
func Is(err error, target error) bool {
	return stderr.Is(err, target)
}

// As 从错误中提取目标类型的错误实例
func As(err error, target any) bool {
	return stderr.As(err, target)
}
