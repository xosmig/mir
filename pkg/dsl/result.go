package dsl

import (
	"errors"
	"fmt"
)

type Result struct {
	err     error
	matched bool
	repeat  bool
}

func (r Result) Err() error {
	return r.err
}

func (r Result) Matched() bool {
	return r.matched
}

func (r Result) Repeat() bool {
	return r.repeat
}

func FromError(err error) Result {
	return Result{
		err:     err,
		matched: true,
		repeat:  false,
	}
}

func Error(s string) Result {
	return FromError(errors.New(s))
}

func Errorf(format string, args ...any) Result {
	return FromError(fmt.Errorf(format, args...))
}

func NoMatch() Result {
	return Result{
		err:     nil,
		matched: false,
		repeat:  true,
	}
}

func OK() Result {
	return Result{
		err:     nil,
		matched: true,
		repeat:  true,
	}
}

func DontRepeat() Result {
	return Result{
		err:     nil,
		matched: true,
		repeat:  false,
	}
}
