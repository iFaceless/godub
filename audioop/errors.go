package audioop

import "fmt"

type Error struct {
	inner string
}

func NewError(format string, args ...interface{}) Error {
	return Error{inner: fmt.Sprintf(format, args...)}
}

func (e Error) Error() string {
	return e.inner
}
