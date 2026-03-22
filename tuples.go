package lcl

type Entry[K comparable, V any] struct {
	Key   K
	Value V
}

type Tuple2[A, B any] struct {
	A A
	B B
}

type Tuple3[A, B, C any] struct {
	A A
	B B
	C C
}

type Tuple4[A, B, C, D any] struct {
	A A
	B B
	C C
	D D
}

func (t Tuple2[A, B]) Unpack() (A, B) {
	return t.A, t.B
}

func (t Tuple3[A, B, C]) Unpack() (A, B, C) {
	return t.A, t.B, t.C
}

func (t Tuple4[A, B, C, D]) Unpack() (A, B, C, D) {
	return t.A, t.B, t.C, t.D
}

func T2[A, B any](a A, b B) Tuple2[A, B] {
	return Tuple2[A, B]{A: a, B: b}
}

func T3[A, B, C any](a A, b B, c C) Tuple3[A, B, C] {
	return Tuple3[A, B, C]{A: a, B: b, C: c}
}

func T4[A, B, C, D any](a A, b B, c C, d D) Tuple4[A, B, C, D] {
	return Tuple4[A, B, C, D]{A: a, B: b, C: c, D: d}
}

func Zip2[A, B any](a []A, b []B) []Tuple2[A, B] {
	result := make([]Tuple2[A, B], max(len(a), len(b)))
	for i := range a {
		result[i].A = a[i]
	}
	for i := range b {
		result[i].B = b[i]
	}
	return result
}

func Zip3[A, B, C any](a []A, b []B, c []C) []Tuple3[A, B, C] {
	result := make([]Tuple3[A, B, C], max(len(a), len(b), len(c)))
	for i := range a {
		result[i].A = a[i]
	}
	for i := range b {
		result[i].B = b[i]
	}
	for i := range c {
		result[i].C = c[i]
	}
	return result
}

func Zip4[A, B, C, D any](a []A, b []B, c []C, d []D) []Tuple4[A, B, C, D] {
	result := make([]Tuple4[A, B, C, D], max(len(a), len(b), len(c), len(d)))
	for i := range a {
		result[i].A = a[i]
	}
	for i := range b {
		result[i].B = b[i]
	}
	for i := range c {
		result[i].C = c[i]
	}
	for i := range d {
		result[i].D = d[i]
	}
	return result
}

func Unzip2[A, B any](tuples []Tuple2[A, B]) ([]A, []B) {
	size := len(tuples)
	r1 := make([]A, 0, size)
	r2 := make([]B, 0, size)
	for i := range tuples {
		r1 = append(r1, tuples[i].A)
		r2 = append(r2, tuples[i].B)
	}
	return r1, r2
}

func Unzip3[A, B, C any](tuples []Tuple3[A, B, C]) ([]A, []B, []C) {
	size := len(tuples)
	r1 := make([]A, 0, size)
	r2 := make([]B, 0, size)
	r3 := make([]C, 0, size)
	for i := range tuples {
		r1 = append(r1, tuples[i].A)
		r2 = append(r2, tuples[i].B)
		r3 = append(r3, tuples[i].C)
	}
	return r1, r2, r3
}

func Unzip4[A, B, C, D any](tuples []Tuple4[A, B, C, D]) ([]A, []B, []C, []D) {
	size := len(tuples)
	r1 := make([]A, 0, size)
	r2 := make([]B, 0, size)
	r3 := make([]C, 0, size)
	r4 := make([]D, 0, size)
	for i := range tuples {
		r1 = append(r1, tuples[i].A)
		r2 = append(r2, tuples[i].B)
		r3 = append(r3, tuples[i].C)
		r4 = append(r4, tuples[i].D)
	}
	return r1, r2, r3, r4
}

func CrossJoin2[A, B any](listA []A, listB []B) []Tuple2[A, B] {
	size := len(listA) * len(listB)
	if size == 0 {
		return []Tuple2[A, B]{}
	}
	result := make([]Tuple2[A, B], 0, size)
	for _, a := range listA {
		for _, b := range listB {
			result = append(result, T2(a, b))
		}
	}
	return result
}

func CrossJoin3[A, B, C any](listA []A, listB []B, listC []C) []Tuple3[A, B, C] {
	size := len(listA) * len(listB)
	if size == 0 {
		return []Tuple3[A, B, C]{}
	}
	result := make([]Tuple3[A, B, C], 0, size)
	for _, a := range listA {
		for _, b := range listB {
			for _, c := range listC {
				result = append(result, T3(a, b, c))
			}
		}
	}
	return result
}

func CrossJoin4[A, B, C, D any](listA []A, listB []B, listC []C, listD []D) []Tuple4[A, B, C, D] {
	size := len(listA) * len(listB)
	if size == 0 {
		return []Tuple4[A, B, C, D]{}
	}
	result := make([]Tuple4[A, B, C, D], 0, size)
	for _, a := range listA {
		for _, b := range listB {
			for _, c := range listC {
				for _, d := range listD {
					result = append(result, T4(a, b, c, d))
				}
			}
		}
	}
	return result
}
