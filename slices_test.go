package lcl

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// --- Filter / Filter2 / FilterError ---

func TestFilter(t *testing.T) {
	isEven := func(n int) bool { return n%2 == 0 }
	tests := []struct {
		name  string
		input []int
		want  []int
	}{
		{"keeps matching", []int{1, 2, 3, 4, 5}, []int{2, 4}},
		{"all match", []int{2, 4, 6}, []int{2, 4, 6}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, Filter(tc.input, isEven))
		})
	}
	assert.Empty(t, Filter([]int{1, 3, 5}, isEven))
	assert.Empty(t, Filter([]int{}, isEven))
}

func TestFilter2(t *testing.T) {
	evenIndex := func(i int, _ string) bool { return i%2 == 0 }
	got := Filter2([]string{"a", "b", "c", "d"}, evenIndex)
	assert.Equal(t, []string{"a", "c"}, got)
}

func TestFilterError(t *testing.T) {
	sentinel := errors.New("stop")

	t.Run("no error", func(t *testing.T) {
		pred := func(n int) (bool, error) { return n > 2, nil }
		got, err := FilterError([]int{1, 2, 3, 4}, pred)
		assert.NoError(t, err)
		assert.Equal(t, []int{3, 4}, got)
	})
	t.Run("stops on error", func(t *testing.T) {
		pred := func(n int) (bool, error) {
			if n == 3 {
				return false, sentinel
			}
			return n > 1, nil
		}
		_, err := FilterError([]int{1, 2, 3, 4}, pred)
		assert.ErrorIs(t, err, sentinel)
	})
}

// --- Exclude / Exclude2 / ExcludeError ---

func TestExclude(t *testing.T) {
	isEven := func(n int) bool { return n%2 == 0 }
	tests := []struct {
		name  string
		input []int
		want  []int
	}{
		{"removes matching", []int{1, 2, 3, 4, 5}, []int{1, 3, 5}},
		{"none removed", []int{1, 3, 5}, []int{1, 3, 5}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, Exclude(tc.input, isEven))
		})
	}
	assert.Empty(t, Exclude([]int{2, 4, 6}, isEven))
	assert.Empty(t, Exclude([]int{}, isEven))
}

func TestExclude2(t *testing.T) {
	evenIndex := func(i int, _ string) bool { return i%2 == 0 }
	got := Exclude2([]string{"a", "b", "c", "d"}, evenIndex)
	assert.Equal(t, []string{"b", "d"}, got)
}

func TestExcludeError(t *testing.T) {
	sentinel := errors.New("stop")

	t.Run("excludes matching items", func(t *testing.T) {
		pred := func(n int) (bool, error) { return n%2 == 0, nil }
		got, err := ExcludeError([]int{1, 2, 3, 4, 5}, pred)
		assert.NoError(t, err)
		assert.Equal(t, []int{1, 3, 5}, got)
	})
	t.Run("stops on error", func(t *testing.T) {
		pred := func(n int) (bool, error) {
			if n == 3 {
				return false, sentinel
			}
			return n%2 == 0, nil
		}
		_, err := ExcludeError([]int{1, 2, 3, 4}, pred)
		assert.ErrorIs(t, err, sentinel)
	})
}

// --- Map / Map2 / MapError ---

func TestMap(t *testing.T) {
	tests := []struct {
		name   string
		input  []int
		mapper func(int) int
		want   []int
	}{
		{"double", []int{1, 2, 3}, func(n int) int { return n * 2 }, []int{2, 4, 6}},
		{"identity", []int{1, 2, 3}, func(n int) int { return n }, []int{1, 2, 3}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, Map(tc.input, tc.mapper))
		})
	}
	assert.Empty(t, Map([]int{}, func(n int) int { return n }))
}

func TestMap2(t *testing.T) {
	got := Map2([]string{"a", "b", "c"}, func(i int, s string) string {
		return s + string(rune('0'+i))
	})
	assert.Equal(t, []string{"a0", "b1", "c2"}, got)
}

func TestMapError(t *testing.T) {
	sentinel := errors.New("bad")

	t.Run("no error", func(t *testing.T) {
		got, err := MapError([]int{1, 2, 3}, func(n int) (int, error) { return n * 3, nil })
		assert.NoError(t, err)
		assert.Equal(t, []int{3, 6, 9}, got)
	})
	t.Run("stops on error", func(t *testing.T) {
		_, err := MapError([]int{1, 2, 3}, func(n int) (int, error) {
			if n == 2 {
				return 0, sentinel
			}
			return n, nil
		})
		assert.ErrorIs(t, err, sentinel)
	})
}

// --- Fold / Fold2 / FoldError ---

func TestFold(t *testing.T) {
	sum := func(acc, n int) int { return acc + n }
	tests := []struct {
		name    string
		input   []int
		initial int
		want    int
	}{
		{"sum", []int{1, 2, 3, 4}, 0, 10},
		{"empty uses initial", []int{}, 42, 42},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, Fold(tc.input, sum, tc.initial))
		})
	}
}

func TestFold2(t *testing.T) {
	// sum of index * value: 0*10 + 1*20 + 2*30 = 80
	got := Fold2([]int{10, 20, 30}, func(acc, i, v int) int { return acc + i*v }, 0)
	assert.Equal(t, 80, got)
}

func TestFoldError(t *testing.T) {
	sentinel := errors.New("bad")
	sum := func(acc, n int) (int, error) { return acc + n, nil }

	t.Run("no error", func(t *testing.T) {
		got, err := FoldError([]int{1, 2, 3}, sum, 0)
		assert.NoError(t, err)
		assert.Equal(t, 6, got)
	})
	t.Run("stops on error", func(t *testing.T) {
		_, err := FoldError([]int{1, 2, 3}, func(acc, n int) (int, error) {
			if n == 2 {
				return acc, sentinel
			}
			return acc + n, nil
		}, 0)
		assert.ErrorIs(t, err, sentinel)
	})
}

// --- FoldRight / FoldRight2 / FoldRightError ---

func TestFoldRight(t *testing.T) {
	got := FoldRight([]string{"a", "b", "c"}, func(acc, s string) string { return acc + s }, "")
	assert.Equal(t, "cba", got)
}

func TestFoldRight2(t *testing.T) {
	got := FoldRight2([]string{"a", "b", "c"}, func(acc string, i int, s string) string {
		return acc + string(rune('0'+i)) + ":" + s + " "
	}, "")
	assert.Equal(t, "2:c 1:b 0:a ", got)
}

func TestFoldRightError(t *testing.T) {
	sum := func(acc, n int) (int, error) { return acc + n, nil }
	got, err := FoldRightError([]int{1, 2, 3}, sum, 0)
	assert.NoError(t, err)
	assert.Equal(t, 6, got)
}

// --- ForEach / ForEach2 / ForEachWhile ---

func TestForEach(t *testing.T) {
	sum := 0
	ForEach([]int{1, 2, 3}, func(n int) { sum += n })
	assert.Equal(t, 6, sum)
}

func TestForEach2(t *testing.T) {
	result := 0
	ForEach2([]int{10, 20, 30}, func(i, v int) { result += i * v })
	assert.Equal(t, 80, result) // 0*10 + 1*20 + 2*30
}

func TestForEachWhile(t *testing.T) {
	count := 0
	ForEachWhile([]int{1, 2, 3, 4, 5}, func(n int) bool {
		if n == 3 {
			return false
		}
		count++
		return true
	})
	assert.Equal(t, 2, count)
}

// --- GroupBy / PartitionBy / FrequenciesBy ---

func TestGroupBy(t *testing.T) {
	parity := func(n int) string {
		if n%2 == 0 {
			return "even"
		}
		return "odd"
	}
	got := GroupBy([]int{1, 2, 3, 4, 5, 6}, parity)
	assert.Equal(t, []int{2, 4, 6}, got["even"])
	assert.Equal(t, []int{1, 3, 5}, got["odd"])
}

func TestPartitionBy(t *testing.T) {
	// groups all elements by key in first-seen order; non-adjacent equal keys
	// are merged into the same group (unlike GroupBy which returns a map)
	got := PartitionBy([]int{1, 1, 2, 2, 1}, Id[int])
	assert.Equal(t, [][]int{{1, 1, 1}, {2, 2}}, got)
}

func TestFrequenciesBy(t *testing.T) {
	got := FrequenciesBy([]string{"a", "b", "a", "c", "b", "a"}, Id[string])
	assert.Equal(t, map[string]int{"a": 3, "b": 2, "c": 1}, got)
}

// --- Pad / PadWith ---

func TestPad(t *testing.T) {
	tests := []struct {
		name  string
		input []int
		n     int
		want  []int
	}{
		{"extends with zeros", []int{1, 2}, 5, []int{1, 2, 0, 0, 0}},
		{"noop when already long enough", []int{1, 2, 3}, 2, []int{1, 2, 3}},
		{"noop when equal length", []int{1, 2}, 2, []int{1, 2}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, Pad(tc.input, tc.n))
		})
	}
}

func TestPad_negativePanics(t *testing.T) {
	assert.Panics(t, func() { Pad([]int{1}, -1) })
}

func TestPadWith(t *testing.T) {
	got := PadWith([]int{1, 2}, 5, 99)
	assert.Equal(t, []int{1, 2, 99, 99, 99}, got)
}
