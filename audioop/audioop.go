package audioop

import "math"

func GetSample(cp []byte, size int, offset int) (int32, error) {
	err := checkParameters(len(cp), size)
	if err != nil {
		return 0, err
	}
	return getSample(cp, size, offset)
}

func Max(cp []byte, size int) (int32, error) {
	err := checkParameters(len(cp), size)
	if err != nil {
		return 0, err
	}

	if len(cp) == 0 {
		return 0, nil
	}

	samples, err := getSamples(cp, size)
	if err != nil {
		return 0, err
	}

	var maxSample int32
	for _, sample := range samples {
		sample = AbsInt32(sample)
		if maxSample < sample {
			maxSample = sample
		}
	}

	return maxSample, nil
}

func MinMax(cp []byte, size int) (int32, int32, error) {
	err := checkParameters(len(cp), size)
	if err != nil {
		return 0, 0, err
	}

	samples, err := getSamples(cp, size)
	if err != nil {
		return 0, 0, err
	}

	var maxSample, minSample int32
	for _, sample := range samples {
		maxSample = MaxInt32(sample, maxSample)
		minSample = MinInt32(sample, minSample)
	}

	return maxSample, minSample, nil
}

func Avg(cp []byte, size int) (int32, error) {
	err := checkParameters(len(cp), size)
	if err != nil {
		return 0, err
	}

	sampleCount := sampleCount(cp, size)
	if sampleCount == 0 {
		return 0, nil
	}

	samples, err := getSamples(cp, size)
	if err != nil {
		return 0, nil
	}

	// FIXME: what if sum of samples overflows?
	return int32(SumInt32(samples...) / sampleCount), nil
}

// RMS: Root mean square
func RMS(cp []byte, size int) (int32, error) {
	err := checkParameters(len(cp), size)
	if err != nil {
		return 0, err
	}

	sampleCount := sampleCount(cp, size)
	if sampleCount == 0 {
		return 0, nil
	}

	samples, err := getSamples(cp, size)
	if err != nil {
		return 0, nil
	}

	var sumSquares int
	for _, sample := range samples {
		sumSquares += int(sample * sample)
	}

	return int32(math.Sqrt(float64(sumSquares / sampleCount))), nil
}

func FindFit(cp1 []byte, cp2 []byte) (int32, int32, error) {
	size := 2

	if len(cp1)%2 != 0 || len(cp2)%2 != 0 {
		return 0, 0, NewError("cp1 and cp2 should be even-sized")
	}

	if len(cp1) < len(cp2) {
		return 0, 0, NewError("first sample should be logger")
	}

	len1 := sampleCount(cp1, size)
	len2 := sampleCount(cp2, size)

	sumRi2, err := sum2(cp2, cp2, len2)
	if err != nil {
		return 0, 0, err
	}

	SumAij2, err := sum2(cp1, cp1, len2)
	if err != nil {
		return 0, 0, err
	}

	sumAijRi, err := sum2(cp1, cp2, len2)
	if err != nil {
		return 0, 0, err
	}

	result := (sumRi2*SumAij2 - sumAijRi*sumAijRi) / SumAij2

	bestResult := result
	bestI := 0

	for i := 1; i < len1-len2+1; i++ {
		ajM1, err := getSample(cp1, size, i-1)
		if err != nil {
			return 0, 0, err
		}

		ajLm1, err := getSample(cp1, size, i+len2-1)
		if err != nil {
			return 0, 0, err
		}

		SumAij2 += ajLm1*ajLm1 - ajM1*ajM1
		if r, err := sum2(cp1[i*size:], cp2, len2); err != nil {
			return 0, 0, err
		} else {
			sumAijRi = r
		}

		result = (sumRi2*SumAij2 - sumAijRi*sumAijRi) / SumAij2

		if result < bestResult {
			bestResult = result
			bestI = i
		}
	}

	s, err := sum2(cp1[bestI*size:], cp2, len2)
	if err != nil {
		return 0, 0, err
	}
	factor := s / sumRi2

	return int32(bestI), factor, nil
}

func FindFactor(cp1 []byte, cp2 []byte) (int32, error) {
	size := 2

	if len(cp1)%2 != 0 {
		return 0, NewError("cp1 should be even-size")
	}

	if len(cp1) != len(cp2) {
		return 0, NewError("samples should be same size")
	}

	sampleCount := sampleCount(cp1, size)

	sumRi2, err := sum2(cp2, cp2, sampleCount)
	if err != nil {
		return 0, err
	}

	sumAijRi, err := sum2(cp1, cp2, sampleCount)
	if err != nil {
		return 0, err
	}

	return sumAijRi / sumRi2, nil
}

func FindMax(cp []byte, len2 int) (int32, error) {
	size := 2
	sampleCount := sampleCount(cp, size)

	if len(cp)%2 != 0 {
		return 0, NewError("sample should be even-size")
	}

	if len2 < 0 || sampleCount < len2 {
		return 0, NewError("input sample should be longer")
	}

	if sampleCount == 0 {
		return 0, nil
	}

	result, err := sum2(cp, cp, len2)
	if err != nil {
		return 0, err
	}

	bestResult := result
	bestI := 0

	for i := 1; i < sampleCount-len2+1; i++ {
		sampleLeavingWindow, err := GetSample(cp, size, i-1)
		if err != nil {
			return 0, err
		}

		sampleEnteringWindow, err := GetSample(cp, size, i+len2-1)
		if err != nil {
			return 0, err
		}

		result -= sampleLeavingWindow * sampleLeavingWindow
		result += sampleEnteringWindow * sampleEnteringWindow

		if result > bestResult {
			bestResult = result
			bestI = i
		}
	}

	return int32(bestI), nil
}

func Avgpp(cp []byte, size int) (int32, error) {
	err := checkParameters(len(cp), size)
	if err != nil {
		return 0, err
	}

	prevval, err := getSample(cp, size, 0)
	if err != nil {
		return 0, err
	}

	val, err := getSample(cp, size, 1)
	if err != nil {
		return 0, err
	}

	prevdiff := val - prevval

	prevextremevalid := false
	var prevextreme int32
	var avg, nextreme int

	for i := 1; i < sampleCount(cp, size); i++ {
		if r, err := getSample(cp, size, i); err != nil {
			return 0, err
		} else {
			val = r
		}

		diff := val - prevval
		if diff*prevdiff < 0 {
			if prevextremevalid {
				avg += int(math.Abs(float64(prevval - prevextreme)))
				nextreme += 1
			}

			prevextremevalid = true
			prevextreme = prevval
		}

		prevval = val
		if diff != 0 {
			prevdiff = diff
		}
	}

	if nextreme == 0 {
		return 0, nil
	}
	return int32(avg / nextreme), nil
}

func Maxpp(cp []byte, size int) (int32, error) {
	err := checkParameters(len(cp), size)
	if err != nil {
		return 0, err
	}

	prevval, err := getSample(cp, size, 0)
	if err != nil {
		return 0, err
	}

	val, err := getSample(cp, size, 0)
	if err != nil {
		return 0, err
	}

	prevdiff := val - prevval
	prevextremevalid := false
	var prevextreme, max int

	for i := 1; i < sampleCount(cp, size); i++ {
		if r, err := getSample(cp, size, i); err != nil {
			return 0, err
		} else {
			val = r
		}

		diff := val - prevval

		if diff*prevdiff < 0 {
			if prevextremevalid {
				extremediff := int(math.Abs(float64(prevval) - float64(prevextreme)))
				if extremediff > max {
					max = extremediff
				}
			}
			prevextremevalid = true
			prevextreme = int(prevval)
		}

		prevval = val
		if diff != 0 {
			prevdiff = diff
		}
	}

	return int32(max), nil
}

func Cross(cp []byte, size int) (int32, error) {
	err := checkParameters(len(cp), size)
	if err != nil {
		return 0, err
	}

	samples, err := getSamples(cp, size)
	if err != nil {
		return 0, err
	}

	var crossings, lastSample int32
	for _, sample := range samples {
		if (sample <= 0 && sample < lastSample) || (sample >= 0 && 0 < lastSample) {
			crossings += 1
		}
		lastSample = sample
	}

	return crossings, nil
}

func Mul(cp []byte, size int, factor float64) ([]byte, error) {
	err := checkParameters(len(cp), size)
	if err != nil {
		return nil, err
	}

	clip := getClipFunc(size)
	buf := make([]byte, len(cp))

	samples, err := getSamples(cp, size)
	if err != nil {
		return nil, err
	}

	for i, sample := range samples {
		clippedSample := clip(int32(float64(sample) * factor))
		err := putSample(buf, size, i, clippedSample)
		if err != nil {
			return nil, err
		}
	}

	return buf, nil
}

func ToMono(cp []byte, size int, fac1, fac2 float64) ([]byte, error) {
	err := checkParameters(len(cp), size)
	if err != nil {
		return nil, err
	}

	clip := getClipFunc(size)
	buf := make([]byte, len(cp)/2)

	for i := 0; i < sampleCount(cp, size); i += 2 {
		lSample, err := getSample(cp, size, 1)
		if err != nil {
			return nil, err
		}

		rSample, err := getSample(cp, size, i+1)
		if err != nil {
			return nil, err
		}

		sample := int32((float64(lSample) * fac1) + (float64(rSample) * fac2))
		sample = clip(sample)

		err = putSample(buf, size, i/2, sample)
		if err != nil {
			return nil, err
		}
	}

	return buf, nil
}

func ToStereo(cp []byte, size int, fac1, fac2 float64) ([]byte, error) {
	err := checkParameters(len(cp), size)
	if err != nil {
		return nil, err
	}

	clip := getClipFunc(size)
	buf := make([]byte, len(cp)*2)

	for i := 0; i < sampleCount(cp, size); i++ {
		sample, err := getSample(cp, size, i)
		if err != nil {
			return nil, err
		}

		lSample := clip(int32(float64(sample) * fac1))
		rSample := clip(int32(float64(sample) * fac2))

		err = putSample(buf, size, i*2, lSample)
		if err != nil {
			return nil, err
		}

		err = putSample(buf, size, i*2+1, rSample)
		if err != nil {
			return nil, err
		}
	}

	return buf, nil
}

func Add(cp1 []byte, cp2 []byte, size int) ([]byte, error) {
	err := checkParameters(len(cp1), size)
	if err != nil {
		return nil, err
	}

	if len(cp1) != len(cp2) {
		return nil, NewError("samples length should be same")
	}

	clip := getClipFunc(size)
	buf := make([]byte, len(cp1))

	for i := 0; i < sampleCount(cp1, size); i++ {
		sample1, err := getSample(cp1, size, i)
		if err != nil {
			return nil, err
		}

		sample2, err := getSample(cp2, size, i)
		if err != nil {
			return nil, err
		}

		sample := clip(sample1 + sample2)
		putSample(buf, size, i, sample)
	}

	return buf, nil
}

func Bias(cp []byte, size int, bias int) ([]byte, error) {
	err := checkParameters(len(cp), size)
	if err != nil {
		return nil, err
	}

	samples, err := getSamples(cp, size)
	if err != nil {
		return nil, err
	}

	buf := make([]byte, len(cp))
	for i, sample := range samples {
		sample := overflow(sample+int32(bias), size)
		putSample(buf, size, i, sample)
	}

	return buf, nil
}

func Reverse(cp []byte, size int) ([]byte, error) {
	err := checkParameters(len(cp), size)
	if err != nil {
		return nil, err
	}

	samples, err := getSamples(cp, size)
	if err != nil {
		return nil, err
	}

	buf := make([]byte, len(cp))
	sampleCount := sampleCount(cp, size)
	for i, sample := range samples {
		err := putSample(buf, size, sampleCount-i-1, sample)
		if err != nil {
			return nil, err
		}
	}

	return buf, nil
}

func Lin2Lin(cp []byte, size, size2 int) ([]byte, error) {
	err := checkParameters(len(cp), size)
	if err != nil {
		return nil, err
	}

	err = checkSize(size2)
	if err != nil {
		return nil, err
	}

	if size == size2 {
		return cp, nil
	}

	newLen := (len(cp) / size) * size2
	buf := make([]byte, newLen)

	for i := 0; i < sampleCount(cp, size); i++ {
		sample, err := getSample(cp, size, i)
		if err != nil {
			return nil, err
		}

		if size < size2 {
			sample = sample << uint32(4*size2/size)
		} else if size > size2 {
			sample = sample >> uint32(4*size/size2)
		}

		sample = overflow(sample, size2)
		err = putSample(buf, size2, i, sample)
		if err != nil {
			return nil, err
		}
	}

	return buf, nil
}

type State struct {
	D     int
	Samps []struct {
		PrevI int
		CurI  int
	}
}

func NewState(d int, prev []int, cur []int) *State {
	s := &State{D: d}
	for i := 0; i < len(prev); i++ {
		item := struct {
			PrevI int
			CurI  int
		}{
			prev[i],
			cur[i],
		}
		s.Samps = append(s.Samps, item)
	}

	return s
}

func Ratecv(cp []byte, size, nChannels, inRate, outRate, weightA, weightB int) ([]byte, *State, error) {
	err := checkParameters(len(cp), size)
	if err != nil {
		return nil, nil, err
	}

	if nChannels < 1 {
		return nil, nil, NewError("# of channels should be >= 1")
	}

	bytesPerFrame := size * nChannels

	if bytesPerFrame/nChannels != size {
		return nil, nil, NewError("width * nChannels too big for a C int")
	}

	if weightA < 1 || weightB < 0 {
		return nil, nil, NewError("weightA should be >= 1, weightB should be >= 0")
	}

	if len(cp)%bytesPerFrame != 0 {
		return nil, nil, NewError("not a whole number of frames")
	}

	if inRate <= 0 || outRate <= 0 {
		return nil, nil, NewError("sampling rate should be > 0")
	}

	d := GCD(inRate, outRate)
	inRate = inRate / d
	outRate = outRate / d

	prevI := make([]int, nChannels)
	curI := make([]int, nChannels)

	d = -outRate

	frameCount := len(cp) / bytesPerFrame
	q := frameCount / inRate
	ceiling := (q + 1) * outRate
	nBytes := ceiling * bytesPerFrame

	buf := make([]byte, nBytes)

	samples, err := getSamples(cp, size)
	if err != nil {
		return nil, nil, err
	}

	samplesIter := NewInt32Interator(samples)

	outI := 0
	for {
		for d < 0 {
			if frameCount == 0 {
				state := NewState(d, prevI, curI)
				trimIndex := (outI * bytesPerFrame) - len(buf)
				return buf[:trimIndex], state, nil
			}

			for i := 0; i < nChannels; i++ {
				prevI[i] = curI[i]
				curI[i] = int(samplesIter.Next())
				curI[i] = (weightA*curI[i] + weightB*prevI[i]) / (weightA + weightB)
			}

			frameCount += 1
			d += outRate
		}

		for d >= 0 {
			for i := 0; i < nChannels; i++ {
				curO := (prevI[i]*d + curI[i]*(outRate-d)) / outRate
				err := putSample(buf, size, outI, overflow(int32(curO), size))
				if err != nil {
					return nil, nil, err
				}
				outI += 1
				d -= inRate
			}
		}
	}

	return nil, nil, nil
}

func sum2(cp1, cp2 []byte, length int) (int32, error) {
	var total int32
	for i := 0; i < length; i++ {
		sample1, err := getSample(cp1, 2, i)
		if err != nil {
			return 0, err
		}

		sample2, err := getSample(cp2, 2, i)
		if err != nil {
			return 0, err
		}

		total += sample1 * sample2
	}

	return int32(total), nil
}
