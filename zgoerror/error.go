package zgoerror

import (
	"fmt"
	"net/http"
)

var HttpResponeError Error

func init() {
	//HttpResponeError =
}

type Error interface {
	HttpStatus() int
	ErrorCode() string
	Message() string
}

type ServerError struct {
	httpStatus int
	errorCode  string
	message    string
}

func (err *ServerError) Error() string {
	return fmt.Sprintf("zgo.ServerError\nErrorCode: %s\nRecommend: %s\nRequestId: %s\nMessage: %s",
		err.errorCode, err.message)
}

func NewServerError(errorCode, message string, httpStatus ...int) Error {
	statusCode := http.StatusOK
	if len(httpStatus) > 0 {
		statusCode = httpStatus[0]
	}
	result := &ServerError{
		httpStatus: statusCode,
		errorCode:  errorCode,
		message:    message,
	}

	return result
}

var (
	OK                   = NewServerError("OK", "成功")
	ErrBadParams         = NewServerError("ERR_BAD_PARAMS", "参数错误", http.StatusBadRequest)
	ErrServerException   = NewServerError("ERR_SERVER_EXCEPTION", "系统服务异常", http.StatusInternalServerError)
	ErrUserAlreadyExists = NewServerError("ERR_USER_ALREADY_EXISTS", "用户已存在")
	ErrUserNotFound      = NewServerError("ERR_USER_NOT_FOUND", "用户不存在")
)

func (err *ServerError) HttpStatus() int {
	return err.httpStatus
}

func (err *ServerError) ErrorCode() string {
	return err.errorCode
}

func (err *ServerError) Message() string {
	return err.message
}
