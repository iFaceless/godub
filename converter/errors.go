package converter

type EncodeError string

func (e EncodeError) Error() string {
	return string(e)
}

type InvalidCoverError string

func (e InvalidCoverError) Error() string {
	return string(e)
}

type InvalidID3TagVersionError string

func (e InvalidID3TagVersionError) Error() string {
	return string(e)
}
