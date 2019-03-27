package wav

type DecodeError string

func (e DecodeError) Error() string {
	return string(e)
}
