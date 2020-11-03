package models

import (
	"errors"
	"fmt"

	"github.com/0987363/mgo"
)

var ErrIsExist = errors.New("Is exist.")

type ErrorBody struct {
	Code int    `json:"code"`
	Msg  string `json:"msg,omitempty"`
}

func NewErrorBody(code int, msg ...interface{}) *ErrorBody {
	return &ErrorBody{
		Code: code,
		Msg:  fmt.Sprint(msg...),
	}
}

func Error(v ...interface{}) error {
	return errors.New(fmt.Sprint(v...))
}

func Errorf(format string, v ...interface{}) error {
	return errors.New(fmt.Sprintf(format, v...))
}

func ErrorConvert(err error) error {
	switch err {
	case mgo.ErrNotFound:
		return nil
	default:
		return err
	}
}
