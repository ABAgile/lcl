package lcl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMinMax(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		tests := []struct {
			a, b    int
			wantMin int
			wantMax int
		}{
			{1, 5, 1, 5},
			{9, 3, 3, 9},
			{4, 4, 4, 4},
		}
		for _, tc := range tests {
			min, max := MinMax(tc.a, tc.b)
			assert.Equal(t, tc.wantMin, min)
			assert.Equal(t, tc.wantMax, max)
		}
	})

	t.Run("float64", func(t *testing.T) {
		tests := []struct {
			a, b    float64
			wantMin float64
			wantMax float64
		}{
			{3.14, 2.71, 2.71, 3.14},
			{-1.5, 0.5, -1.5, 0.5},
			{0.0, 0.0, 0.0, 0.0},
			{100.0, 0.001, 0.001, 100.0},
		}
		for _, tc := range tests {
			min, max := MinMax(tc.a, tc.b)
			assert.Equal(t, tc.wantMin, min)
			assert.Equal(t, tc.wantMax, max)
		}
	})
}

func TestClamp(t *testing.T) {
	tests := []struct {
		value, min, max int
		want            int
	}{
		{5, 1, 10, 5},
		{-3, 0, 10, 0},
		{20, 0, 10, 10},
		{0, 0, 10, 0},
		{10, 0, 10, 10},
	}
	for _, tc := range tests {
		assert.Equal(t, tc.want, Clamp(tc.value, tc.min, tc.max))
	}
}

func TestMean(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		tests := []struct {
			name   string
			input  []int
			want   int
			wantOk bool
		}{
			{"sum of range", []int{1, 2, 3, 4, 5}, 3, true},
			{"single element", []int{7}, 7, true},
			{"empty", []int{}, 0, false},
		}
		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				got, ok := Mean(tc.input)
				assert.Equal(t, tc.want, got)
				assert.Equal(t, tc.wantOk, ok)
			})
		}
	})

	t.Run("float64", func(t *testing.T) {
		tests := []struct {
			name   string
			input  []float64
			want   float64
			wantOk bool
		}{
			{"average of three", []float64{1.0, 2.0, 3.0}, 2.0, true},
			{"negative values", []float64{-3.0, -1.0, -2.0}, -2.0, true},
			{"single element", []float64{4.5}, 4.5, true},
			{"empty", []float64{}, 0.0, false},
		}
		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				got, ok := Mean(tc.input)
				assert.Equal(t, tc.want, got)
				assert.Equal(t, tc.wantOk, ok)
			})
		}
	})
}

func TestMode(t *testing.T) {
	tests := []struct {
		name   string
		input  []int
		want   []int
		wantOk bool
	}{
		{"single mode", []int{1, 2, 1, 3}, []int{1}, true},
		{"all tied", []int{1, 2}, []int{1, 2}, true},
		{"multiple modes", []int{1, 2, 1, 2, 3}, []int{1, 2}, true},
		{"single element", []int{42}, []int{42}, true},
		{"empty", []int{}, []int{}, false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := Mode(tc.input)
			assert.ElementsMatch(t, tc.want, got)
			assert.Equal(t, tc.wantOk, ok)
		})
	}
}
