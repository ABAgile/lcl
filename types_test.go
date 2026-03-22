package lcl

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// --- Empty / IsEmpty / IsNotEmpty ---

func TestEmpty(t *testing.T) {
	assert.Equal(t, 0, Empty[int]())
	assert.Equal(t, "", Empty[string]())
	assert.Equal(t, false, Empty[bool]())
}

func TestIsEmpty(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		tests := []struct {
			value int
			want  bool
		}{
			{0, true},
			{1, false},
			{-1, false},
		}
		for _, tc := range tests {
			assert.Equal(t, tc.want, IsEmpty(tc.value))
		}
	})

	t.Run("string", func(t *testing.T) {
		tests := []struct {
			value string
			want  bool
		}{
			{"", true},
			{"hello", false},
		}
		for _, tc := range tests {
			assert.Equal(t, tc.want, IsEmpty(tc.value))
		}
	})
}

func TestIsNotEmpty(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		tests := []struct {
			value int
			want  bool
		}{
			{42, true},
			{0, false},
		}
		for _, tc := range tests {
			assert.Equal(t, tc.want, IsNotEmpty(tc.value))
		}
	})

	t.Run("string", func(t *testing.T) {
		tests := []struct {
			value string
			want  bool
		}{
			{"hello", true},
			{"", false},
		}
		for _, tc := range tests {
			assert.Equal(t, tc.want, IsNotEmpty(tc.value))
		}
	})
}

// --- Coalesce / Coalesced ---

func TestCoalesce(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		tests := []struct {
			name   string
			values []int
			want   int
			wantOk bool
		}{
			{"first non-zero", []int{0, 0, 3, 4}, 3, true},
			{"first is non-zero", []int{5, 0, 3}, 5, true},
			{"all zero", []int{0, 0, 0}, 0, false},
			{"no args", []int{}, 0, false},
		}
		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				got, ok := Coalesce(tc.values...)
				assert.Equal(t, tc.want, got)
				assert.Equal(t, tc.wantOk, ok)
			})
		}
	})

	t.Run("string", func(t *testing.T) {
		tests := []struct {
			name   string
			values []string
			want   string
			wantOk bool
		}{
			{"first non-empty", []string{"", "", "hello", "world"}, "hello", true},
			{"first is non-empty", []string{"hi", "", "hello"}, "hi", true},
			{"all empty", []string{"", ""}, "", false},
			{"no args", []string{}, "", false},
		}
		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				got, ok := Coalesce(tc.values...)
				assert.Equal(t, tc.want, got)
				assert.Equal(t, tc.wantOk, ok)
			})
		}
	})
}

func TestCoalesced(t *testing.T) {
	tests := []struct {
		name   string
		values []string
		want   string
	}{
		{"returns first non-empty", []string{"", "hello", "world"}, "hello"},
		{"all empty returns zero", []string{"", ""}, ""},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, Coalesced(tc.values...))
		})
	}
}

// --- ToPtr / FromPtr ---

func TestToPtr(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		tests := []struct {
			value   int
			wantNil bool
		}{
			{42, false},
			{0, true},
		}
		for _, tc := range tests {
			p := ToPtr(tc.value)
			if tc.wantNil {
				assert.Nil(t, p)
			} else {
				assert.NotNil(t, p)
				assert.Equal(t, tc.value, *p)
			}
		}
	})

	t.Run("string", func(t *testing.T) {
		tests := []struct {
			value   string
			wantNil bool
		}{
			{"hello", false},
			{"", true},
		}
		for _, tc := range tests {
			p := ToPtr(tc.value)
			if tc.wantNil {
				assert.Nil(t, p)
			} else {
				assert.NotNil(t, p)
				assert.Equal(t, tc.value, *p)
			}
		}
	})
}

func TestFromPtr(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		v := 99
		tests := []struct {
			name     string
			ptr      *int
			fallback []int
			want     int
		}{
			{"non-nil pointer", &v, nil, 99},
			{"nil no fallback", nil, nil, 0},
			{"nil with fallback", nil, []int{42}, 42},
		}
		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				assert.Equal(t, tc.want, FromPtr(tc.ptr, tc.fallback...))
			})
		}
	})

	t.Run("string", func(t *testing.T) {
		v := "hello"
		tests := []struct {
			name     string
			ptr      *string
			fallback []string
			want     string
		}{
			{"non-nil pointer", &v, nil, "hello"},
			{"nil no fallback", nil, nil, ""},
			{"nil with fallback", nil, []string{"default"}, "default"},
		}
		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				assert.Equal(t, tc.want, FromPtr(tc.ptr, tc.fallback...))
			})
		}
	})
}

// --- ToAnySlice / FromAnySlice ---

func TestToAnySlice(t *testing.T) {
	tests := []struct {
		name string
		input []int
		want  []any
	}{
		{"int slice", []int{1, 2, 3}, []any{1, 2, 3}},
		{"empty slice", []int{}, []any{}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, ToAnySlice(tc.input))
		})
	}
}

func TestFromAnySlice(t *testing.T) {
	tests := []struct {
		name    string
		input   []any
		want    []int
		wantOk  bool
	}{
		{"all same type", []any{1, 2, 3}, []int{1, 2, 3}, true},
		{"type mismatch", []any{1, "two", 3}, nil, false},
		{"empty slice", []any{}, []int{}, true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := FromAnySlice[int](tc.input)
			assert.Equal(t, tc.wantOk, ok)
			if tc.wantOk {
				assert.Equal(t, tc.want, got)
			}
		})
	}
}

// --- GetIn ---

func TestGetIn(t *testing.T) {
	data := map[string]any{
		"user": map[string]any{
			"name":   "alice",
			"scores": []any{10, 20, 30},
		},
	}

	tests := []struct {
		name    string
		path    string
		want    any
		wantErr bool
	}{
		{"nested map key", "user.name", "alice", false},
		{"slice by index", "user.scores.1", 20, false},
		{"missing key", "user.age", nil, true},
		{"index out of bounds", "user.scores.5", nil, true},
		{"invalid index", "user.scores.abc", nil, true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := GetIn(data, tc.path)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.want, got)
			}
		})
	}
}

func TestGetIn_mapAny(t *testing.T) {
	data := map[any]any{"key": "value"}
	got, err := GetIn(data, "key")
	assert.NoError(t, err)
	assert.Equal(t, "value", got)
}

// --- SetIn ---

func TestSetIn(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tests := []struct {
			name     string
			data     func() map[string]any
			path     string
			value    any
			validate func(t *testing.T, data map[string]any)
		}{
			{
				"updates existing map key",
				func() map[string]any { return map[string]any{"a": map[string]any{"b": 0}} },
				"a.b", 99,
				func(t *testing.T, data map[string]any) {
					assert.Equal(t, 99, data["a"].(map[string]any)["b"])
				},
			},
			{
				"creates missing intermediate key",
				func() map[string]any { return map[string]any{} },
				"x.y", "hello",
				func(t *testing.T, data map[string]any) {
					assert.Equal(t, "hello", data["x"].(map[string]any)["y"])
				},
			},
			{
				"updates slice element",
				func() map[string]any { return map[string]any{"list": []any{1, 2, 3}} },
				"list.1", 99,
				func(t *testing.T, data map[string]any) {
					assert.Equal(t, 99, data["list"].([]any)[1])
				},
			},
		}
		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				data := tc.data()
				assert.NoError(t, SetIn(data, tc.path, tc.value))
				tc.validate(t, data)
			})
		}
	})

	t.Run("error", func(t *testing.T) {
		tests := []struct {
			name string
			data map[string]any
			path string
		}{
			{"index out of bounds", map[string]any{"list": []any{1, 2}}, "list.5"},
			{"non-container leaf", map[string]any{"a": 42}, "a.b"},
		}
		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				assert.Error(t, SetIn(tc.data, tc.path, 99))
			})
		}
	})
}

// --- Id ---

func TestId(t *testing.T) {
	tests := []struct {
		name  string
		value any
	}{
		{"int", 42},
		{"string", "hello"},
		{"error", errors.New("x")},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.value, Id(tc.value))
		})
	}
}
