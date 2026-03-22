package lcl

import (
	"slices"

	"github.com/BooleanCat/go-functional/v2/it"
)

func Filter[T any, S ~[]T](xs S, pred func(T) bool) S {
	return slices.Collect(it.Filter(slices.Values(xs), pred))
}

func Filter2[T any, S ~[]T](xs S, pred func(int, T) bool) S {
	res := make(S, 0)
	for _, x := range it.Filter2(slices.All(xs), pred) {
		res = append(res, x)
	}
	return res
}

func FilterError[T any, S ~[]T](xs S, pred func(T) (bool, error)) (S, error) {
	res := make(S, 0)
	for x, err := range it.FilterError(slices.Values(xs), pred) {
		if err != nil {
			return res, err
		}
		res = append(res, x)
	}
	return res, nil
}

func Exclude[T any, S ~[]T](xs S, pred func(T) bool) S {
	return slices.Collect(it.Exclude(slices.Values(xs), pred))
}

func Exclude2[T any, S ~[]T](xs S, pred func(int, T) bool) S {
	res := make(S, 0)
	for _, x := range it.Exclude2(slices.All(xs), pred) {
		res = append(res, x)
	}
	return res
}

func ExcludeError[T any, S ~[]T](xs S, pred func(T) (bool, error)) (S, error) {
	res := make(S, 0)
	for x, err := range it.ExcludeError(slices.Values(xs), pred) {
		if err != nil {
			return res, err
		}
		res = append(res, x)
	}
	return res, nil
}

func Map[T, R any](xs []T, mapper func(T) R) []R {
	return slices.Collect(it.Map(slices.Values(xs), mapper))
}

func Map2[T, R any](xs []T, mapper func(int, T) R) []R {
	res := make([]R, 0)
	f := func(i int, x T) (any, R) { return nil, mapper(i, x) }
	for _, r := range it.Map2(slices.All(xs), f) {
		res = append(res, r)
	}
	return res
}

func MapError[T, R any](xs []T, mapper func(T) (R, error)) ([]R, error) {
	res := make([]R, 0)
	for r, err := range it.MapError(slices.Values(xs), mapper) {
		if err != nil {
			return res, err
		}
		res = append(res, r)
	}
	return res, nil
}

func Fold[T, R any](xs []T, accum func(R, T) R, initial R) R {
	return it.Fold(slices.Values(xs), accum, initial)
}

func Fold2[T, R any](xs []T, accum func(R, int, T) R, initial R) R {
	return it.Fold2(slices.All(xs), accum, initial)
}

func FoldError[T, R any](xs []T, accum func(R, T) (R, error), initial R) (R, error) {
	res := initial
	for x := range slices.Values(xs) {
		r, err := accum(res, x)
		if err != nil {
			return res, err
		}
		res = r
	}
	return res, nil
}

func FoldRight[T, R any](xs []T, accum func(R, T) R, initial R) R {
	res := initial
	for _, x := range slices.Backward(xs) {
		res = accum(res, x)
	}
	return res
}

func FoldRight2[T, R any](xs []T, accum func(R, int, T) R, initial R) R {
	res := initial
	for i, x := range slices.Backward(xs) {
		res = accum(res, i, x)
	}
	return res
}

func FoldRightError[T, R any](xs []T, accum func(R, T) (R, error), initial R) (R, error) {
	res := initial
	for _, x := range slices.Backward(xs) {
		r, err := accum(res, x)
		if err != nil {
			return res, err
		}
		res = r
	}
	return res, nil
}

func ForEach[T any](xs []T, f func(T)) {
	it.ForEach(slices.Values(xs), f)
}

func ForEach2[T any](xs []T, f func(int, T)) {
	it.ForEach2(slices.All(xs), f)
}

func ForEachWhile[T any](xs []T, pred func(T) bool) {
	for x := range slices.Values(xs) {
		if !pred(x) {
			return
		}
	}
}

func GroupBy[T any, R comparable, S ~[]T](xs S, mapper func(T) R) map[R]S {
	res := map[R]S{}
	for _, x := range xs {
		key := mapper(x)
		res[key] = append(res[key], x)
	}
	return res
}

func PartitionBy[T any, R comparable, S ~[]T](xs S, mapper func(T) R) []S {
	res := []S{}
	indices := map[R]int{}
	for _, x := range xs {
		k := mapper(x)
		if i, ok := indices[k]; ok {
			res[i] = append(res[i], x)
		} else {
			indices[k] = len(res)
			res = append(res, S{x})
		}
	}
	return res
}

func FrequenciesBy[T any, R comparable, S ~[]T](xs S, mapper func(T) R) map[R]int {
	freq := make(map[R]int)
	for _, x := range xs {
		freq[mapper(x)]++
	}
	return freq
}

func Pad[T any, S ~[]T](xs S, n int) S {
	if n < 0 {
		panic("Pad: n must not be negative")
	}
	if len(xs) >= n {
		return xs
	}
	r := make(S, n)
	copy(r, xs)
	return r
}

func PadWith[T any, S ~[]T](xs S, n int, padding T) S {
	out := Pad(xs, n)
	for i := len(xs); i < n; i++ {
		out[i] = padding
	}
	return out
}
