package models

type Result[T any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data T      `json:"data"`
}

func (r *Result[T]) Success() {
	r.Code = 1
}

func (r *Result[T]) SuccessWithObject(object T) {
	r.Code = 1
	r.Data = object
}

func (r *Result[T]) Error(msg string) {
	r.Code = 0
	r.Msg = msg
}
