package misc

import "fmt"

type CodedError struct {
	Code int
	Msg  string
}

func (e *CodedError) Error() string {
	return fmt.Sprintf("(%d) %s", e.Code, e.Msg)
}

func NewCodedError(code int, msg string) error {
	return &CodedError{
		Code: code,
		Msg:  msg,
	}
}

type SilentCodedError struct {
	Code int
}

func (e *SilentCodedError) Error() string {
	return fmt.Sprintf("Code(%d)", e.Code)
}

func NewSilentCodedError(code int) error {
	return &SilentCodedError{
		Code: code,
	}
}
