package zgoerror

//import (
//	"fmt"
//)
//
//var HttpResponeError Error
//
//func init() {
//	//HttpResponeError =
//}
//
//type Error interface {
//	HttpStatus() int
//	ErrorCode() string
//	Message() string
//}
//
//type ServerError struct {
//	httpStatus int
//	errorCode  string
//	message    string
//}
//
//func (err *ServerError) Error() string {
//	return fmt.Sprintf("zgo.ServerError\nErrorCode: %s\nRecommend: %s\nRequestId: %s\nMessage: %s",
//		err.errorCode, err.message)
//}
//
//func NewServerError(httpStatus int, responseContent string) Error {
//	result := &ServerError{
//		httpStatus: httpStatus,
//		message:    responseContent,
//	}
//
//	return result
//}
//
//func (err *ServerError) HttpStatus() int {
//	return err.httpStatus
//}
//
//func (err *ServerError) ErrorCode() string {
//	return err.errorCode
//}
//
//func (err *ServerError) Message() string {
//	return err.message
//}
