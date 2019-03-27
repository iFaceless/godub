package utils

func MaxUint16(items ...uint16) uint16 {
	if len(items) == 0 {
		return 0
	}

	if len(items) == 1 {
		return items[0]
	}

	maxValue := items[0]
	for _, v := range items {
		if maxValue < v {
			maxValue = v
		}
	}

	return maxValue
}

func MaxUint32(items ...uint32) uint32 {
	if len(items) == 0 {
		return 0
	}

	if len(items) == 1 {
		return items[0]
	}

	maxValue := items[0]
	for _, v := range items {
		if maxValue < v {
			maxValue = v
		}
	}

	return maxValue
}
