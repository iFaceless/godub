package godub

type AudioSegmentOption func(*AudioSegment)

func SampleWidth(v uint16) AudioSegmentOption {
	return func(s *AudioSegment) {
		s.sampleWidth = v
	}
}

func FrameRate(v uint32) AudioSegmentOption {
	return func(s *AudioSegment) {
		s.frameRate = v
	}
}

func FrameWidth(v uint32) AudioSegmentOption {
	return func(s *AudioSegment) {
		s.frameWidth = v
	}
}

func Channels(v uint16) AudioSegmentOption {
	return func(s *AudioSegment) {
		s.channels = v
	}
}
