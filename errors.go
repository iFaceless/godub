package godub

import "fmt"

type AudioSegmentError struct {
	inner string
}

func NewAudioSegmentError(format string, args ...interface{}) AudioSegmentError {
	return AudioSegmentError{inner: fmt.Sprintf(format, args...)}
}

func (e AudioSegmentError) Error() string {
	return e.inner
}
