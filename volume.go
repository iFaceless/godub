package godub

import (
	"fmt"
	"math"
)

// Volume is the quantity of three-dimensional space enclosed by a closed surface.
// Volume unit is dBFS.
type Volume float64

func NewVolumeFromRatio(ratio float64, denominator float64, useAmplitude bool) Volume {
	if denominator != 0 {
		ratio = ratio / denominator
	}

	if ratio == 0 {
		return Volume(math.Inf(1))
	}

	if useAmplitude {
		return Volume(20 * math.Log10(ratio))
	} else {
		return Volume(10 * math.Log10(ratio))
	}
}

func (volume Volume) String() string {
	return fmt.Sprintf("%.3fdBFS", float64(volume))
}

// ToFloat64 converts db to a float, which represents
// the equivalent ratio in power
func (volume Volume) ToRatio(useAmplitude bool) float64 {
	v := float64(volume)
	if useAmplitude {
		return math.Pow(10, v/20)
	} else {
		return math.Pow(10, v/10)
	}
}
