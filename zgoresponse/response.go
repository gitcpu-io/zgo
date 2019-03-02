package zgoresponse

import "fmt"

type Response interface {
	GetHttpStatus() int
	GetErrorCode() int
	GetMessage() string
	GetData() interface{}
	GetCallBack() string
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////

type ServerJsonPResponse struct {
	HttpStatus int
	ErrorCode  int
	Message    string
	Data       interface{}
	CallBack   string
}

func (this *ServerJsonPResponse) DataMessage() string {
	return fmt.Sprintf("zgo.HttpStatus\nErrorCode: %s\nData: %s\nMessage: %s",
		this.HttpStatus, this.ErrorCode, this.Data, this.Message)
}

func NewServerJsonPResponse(httpStatus int, data interface{}, callBack string) Response {
	result := &ServerJsonPResponse{
		HttpStatus: httpStatus,
		ErrorCode: 0,
		Message: "success",
		Data:       data,
		CallBack: callBack,
	}

	return result
}

func (this *ServerJsonPResponse) GetHttpStatus() int {
	return this.HttpStatus
}

func (this *ServerJsonPResponse) GetErrorCode() int {
	return this.ErrorCode
}

func (this *ServerJsonPResponse) GetMessage() string {
	return this.Message
}

func (this *ServerJsonPResponse) GetData() interface{} {
	return this.Data
}

func (this *ServerJsonPResponse) GetCallBack() string {
	return this.CallBack
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////



//////////////////////////////////////////////////////////////////////////////////////////////////////////////

type ServerResponse struct {
	HttpStatus int
	ErrorCode  int
	Message    string
	Data       interface{}
}

func (this *ServerResponse) DataMessage() string {
	return fmt.Sprintf("zgo.HttpStatus\nErrorCode: %s\nData: %s\nMessage: %s",
		this.HttpStatus, this.ErrorCode, this.Data, this.Message)
}

func NewServerResponse(httpStatus int, data interface{}) Response {
	result := &ServerResponse{
		HttpStatus: httpStatus,
		ErrorCode: 0,
		Message: "success",
		Data:       data,
	}

	return result
}

func (this *ServerResponse) GetHttpStatus() int {
	return this.HttpStatus
}

func (this *ServerResponse) GetErrorCode() int {
	return this.ErrorCode
}

func (this *ServerResponse) GetMessage() string {
	return this.Message
}

func (this *ServerResponse) GetData() interface{} {
	return this.Data
}

func (this *ServerResponse) GetCallBack() string {
	return ""
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////



type ServerError struct {
	HttpStatus int
	ErrorCode  int
	Message    string
}

func (err *ServerError) Error() string {
	return fmt.Sprintf("zgo.ServerError\nErrorCode: %s\nRecommend: %s\nRequestId: %s\nMessage: %s",
		err.ErrorCode, err.Message)
}

func NewServerError(httpStatus int, responseContent string) Response {
	result := &ServerError{
		HttpStatus: httpStatus,
		Message:    responseContent,
	}

	return result
}

func (err *ServerError) GetHttpStatus() int {
	return err.HttpStatus
}

func (err *ServerError) GetErrorCode() int {
	return err.ErrorCode
}

func (err *ServerError) GetMessage() string {
	return err.Message
}

func (err *ServerError) GetData() interface{} {
	return nil
}

func (this *ServerError) GetCallBack() string {
	return ""
}