package zgo

var HttpResponeError Error

func init() {
	//HttpResponeError =
}

type Error interface {
	HttpStatus() int
	ErrorCode() string
	Message() string
}
