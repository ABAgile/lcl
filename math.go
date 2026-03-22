package lcl

import "cmp"

// Number is a type constraint that permits all integer and floating-point types.
type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr | ~float32 | ~float64
}

func MinMax[T cmp.Ordered](mIn, mAx T) (T, T) {
	if mIn <= mAx {
		return mIn, mAx
	}
	return mAx, mIn
}

func Clamp[T cmp.Ordered](value, mIn, mAx T) T {
	if value < mIn {
		return mIn
	} else if value > mAx {
		return mAx
	}
	return value
}

func Mean[T Number](xs []T) (T, bool) {
	if len(xs) == 0 {
		return T(0), false
	}
	var sum T
	for _, item := range xs {
		sum += item
	}
	return sum / T(len(xs)), true
}

func Mode[T comparable](xs []T) ([]T, bool) {
	if len(xs) == 0 {
		return []T{}, false
	}
	freq := make(map[T]int)
	maxFreq := 0
	for _, item := range xs {
		freq[item]++
		if freq[item] > maxFreq {
			maxFreq = freq[item]
		}
	}
	mode := make([]T, 0)
	for item, count := range freq {
		if count == maxFreq {
			mode = append(mode, item)
		}
	}
	return mode, true
}
