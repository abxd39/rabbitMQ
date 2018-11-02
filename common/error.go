package common

const (
	ErrorUnauthorized = 100001
	ErrorHeaderParam  = 100002
)

var errorText = map[int]string{
	ErrorUnauthorized: "失败",
	ErrorHeaderParam:  "非法访问，请传递正确参数(timestamp,signature,username)",
}

func ErrorText(code int) string {
	return errorText[code]
}

type BadRequestError struct {
	code int
}

func (e BadRequestError) Error() string {
	return ErrorText(e.code)
}

func (e BadRequestError) Code() int {
	return e.code
}

func NewError(code int) error {
	return BadRequestError{code: code}
}
