package lcl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// --- Constructors and Unpack ---

func TestTupleConstructorsAndUnpack(t *testing.T) {
	t.Run("T2", func(t *testing.T) {
		a, b := T2(1, "hello").Unpack()
		assert.Equal(t, 1, a)
		assert.Equal(t, "hello", b)
	})
	t.Run("T3", func(t *testing.T) {
		a, b, c := T3(1, 2.0, "x").Unpack()
		assert.Equal(t, 1, a)
		assert.Equal(t, 2.0, b)
		assert.Equal(t, "x", c)
	})
	t.Run("T4", func(t *testing.T) {
		a, b, c, d := T4(1, 2, 3, 4).Unpack()
		assert.Equal(t, 1, a)
		assert.Equal(t, 2, b)
		assert.Equal(t, 3, c)
		assert.Equal(t, 4, d)
	})
}

// --- Zip2 ---

func TestZip2(t *testing.T) {
	tests := []struct {
		name string
		a    []int
		b    []string
		want []Tuple2[int, string]
	}{
		{"equal length", []int{1, 2, 3}, []string{"a", "b", "c"}, []Tuple2[int, string]{{1, "a"}, {2, "b"}, {3, "c"}}},
		{"a longer", []int{1, 2, 3}, []string{"a"}, []Tuple2[int, string]{{1, "a"}, {2, ""}, {3, ""}}},
		{"b longer", []int{1}, []string{"a", "b", "c"}, []Tuple2[int, string]{{1, "a"}, {0, "b"}, {0, "c"}}},
		{"both empty", []int{}, []string{}, []Tuple2[int, string]{}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, Zip2(tc.a, tc.b))
		})
	}
}

// --- Zip3 ---

func TestZip3(t *testing.T) {
	got := Zip3([]int{1, 2}, []string{"a", "b"}, []bool{true, false})
	assert.Equal(t, []Tuple3[int, string, bool]{{1, "a", true}, {2, "b", false}}, got)
}

// --- Zip4 (was buggy: iterated b instead of d for .D) ---

func TestZip4(t *testing.T) {
	tests := []struct {
		name       string
		a, b, c, d []int
		want       []Tuple4[int, int, int, int]
	}{
		{
			"equal length",
			[]int{1, 2}, []int{3, 4}, []int{5, 6}, []int{7, 8},
			[]Tuple4[int, int, int, int]{{1, 3, 5, 7}, {2, 4, 6, 8}},
		},
		{
			"d longer than b",
			[]int{1, 2, 3}, []int{10}, []int{20, 21, 22}, []int{30, 31, 32},
			[]Tuple4[int, int, int, int]{{1, 10, 20, 30}, {2, 0, 21, 31}, {3, 0, 22, 32}},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, Zip4(tc.a, tc.b, tc.c, tc.d))
		})
	}
}

// --- Unzip2 / Unzip3 / Unzip4 ---

func TestUnzip2(t *testing.T) {
	a, b := Unzip2([]Tuple2[int, string]{{1, "a"}, {2, "b"}, {3, "c"}})
	assert.Equal(t, []int{1, 2, 3}, a)
	assert.Equal(t, []string{"a", "b", "c"}, b)
}

func TestUnzip3(t *testing.T) {
	a, b, c := Unzip3([]Tuple3[int, string, bool]{{1, "a", true}, {2, "b", false}})
	assert.Equal(t, []int{1, 2}, a)
	assert.Equal(t, []string{"a", "b"}, b)
	assert.Equal(t, []bool{true, false}, c)
}

func TestUnzip4(t *testing.T) {
	a, b, c, d := Unzip4([]Tuple4[int, int, int, int]{{1, 2, 3, 4}, {5, 6, 7, 8}})
	assert.Equal(t, []int{1, 5}, a)
	assert.Equal(t, []int{2, 6}, b)
	assert.Equal(t, []int{3, 7}, c)
	assert.Equal(t, []int{4, 8}, d)
}

// --- CrossJoin2 / CrossJoin3 / CrossJoin4 ---

func TestCrossJoin2(t *testing.T) {
	tests := []struct {
		name string
		a    []int
		b    []string
		want []Tuple2[int, string]
	}{
		{"full product", []int{1, 2}, []string{"a", "b"}, []Tuple2[int, string]{{1, "a"}, {1, "b"}, {2, "a"}, {2, "b"}}},
		{"empty a", []int{}, []string{"a"}, []Tuple2[int, string]{}},
		{"empty b", []int{1}, []string{}, []Tuple2[int, string]{}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, CrossJoin2(tc.a, tc.b))
		})
	}
}

func TestCrossJoin3(t *testing.T) {
	t.Run("full product", func(t *testing.T) {
		got := CrossJoin3([]int{1, 2}, []int{3}, []int{4, 5})
		assert.Len(t, got, 4)
	})
	t.Run("empty c returns empty", func(t *testing.T) {
		got := CrossJoin3([]int{1, 2}, []int{3, 4}, []int{})
		assert.Empty(t, got)
	})
}

func TestCrossJoin4(t *testing.T) {
	got := CrossJoin4([]int{1}, []int{2}, []int{3}, []int{4, 5})
	assert.Len(t, got, 2)
}
