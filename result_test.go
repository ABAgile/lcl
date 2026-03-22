package lcl

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var errSample = errors.New("something went wrong")

func TestNewResult(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		tests := []struct {
			name    string
			value   int
			err     error
			wantErr bool
		}{
			{"with value", 42, nil, false},
			{"with error", 0, errSample, true},
		}
		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				r := NewResult(tc.value, tc.err)
				assert.Equal(t, tc.value, r.Value)
				if tc.wantErr {
					assert.ErrorIs(t, r.Error, tc.err)
				} else {
					assert.NoError(t, r.Error)
				}
			})
		}
	})

	t.Run("string", func(t *testing.T) {
		tests := []struct {
			name    string
			value   string
			err     error
			wantErr bool
		}{
			{"with value", "hello", nil, false},
			{"with error", "", errSample, true},
		}
		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				r := NewResult(tc.value, tc.err)
				assert.Equal(t, tc.value, r.Value)
				if tc.wantErr {
					assert.ErrorIs(t, r.Error, tc.err)
				} else {
					assert.NoError(t, r.Error)
				}
			})
		}
	})
}

func TestNewValueResult(t *testing.T) {
	r := NewValueResult("hello")
	assert.Equal(t, "hello", r.Value)
	assert.NoError(t, r.Error)
}

func TestNewErrorResult(t *testing.T) {
	tests := []struct {
		name    string
		err     error
		wantErr bool
	}{
		{"carries error", errSample, true},
		{"nil error", nil, false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := NewErrorResult(tc.err)
			if tc.wantErr {
				assert.ErrorIs(t, r.Error, tc.err)
			} else {
				assert.NoError(t, r.Error)
			}
		})
	}
}

func TestUnwrap(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		tests := []struct {
			name      string
			result    *Result[int]
			wantValue int
			wantErr   error
		}{
			{"returns value", NewResult(42, nil), 42, nil},
			{"returns error", NewResult(0, errSample), 0, errSample},
		}
		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				got, err := tc.result.Unwrap()
				assert.Equal(t, tc.wantValue, got)
				assert.ErrorIs(t, err, tc.wantErr)
			})
		}
	})

	t.Run("string", func(t *testing.T) {
		tests := []struct {
			name      string
			result    *Result[string]
			wantValue string
			wantErr   error
		}{
			{"returns value", NewResult("hello", nil), "hello", nil},
			{"returns error", NewResult("", errSample), "", errSample},
		}
		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				got, err := tc.result.Unwrap()
				assert.Equal(t, tc.wantValue, got)
				assert.ErrorIs(t, err, tc.wantErr)
			})
		}
	})
}

func TestBind(t *testing.T) {
	double := func(v int) *Result[int] { return NewValueResult(v * 2) }

	tests := []struct {
		name      string
		result    *Result[int]
		wantValue int
		wantErr   bool
	}{
		{"applies on success, value changes", NewValueResult(5), 10, false},
		{"skips on error, value unchanged", NewResult(5, errSample), 5, true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := tc.result.Bind(double)
			assert.Equal(t, tc.wantValue, r.Value)
			if tc.wantErr {
				assert.ErrorIs(t, r.Error, errSample)
			} else {
				assert.NoError(t, r.Error)
			}
		})
	}
}

func TestMustPass(t *testing.T) {
	tests := []struct {
		name      string
		result    *Result[any]
		wantPanic bool
		wantMsg   string
	}{
		{"nil error does not panic", NewErrorResult(nil), false, ""},
		{"error panics with message", NewErrorResult(errSample), true, fmt.Sprintf("context: %s", errSample)},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.wantPanic {
				assert.PanicsWithValue(t, tc.wantMsg, func() { tc.result.MustPass("context") })
			} else {
				assert.NotPanics(t, func() { tc.result.MustPass("context") })
			}
		})
	}
}

func TestMustGet(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		tests := []struct {
			name      string
			result    *Result[int]
			want      int
			wantPanic bool
		}{
			{"returns value on success", NewResult(7, nil), 7, false},
			{"panics on error", NewResult(0, errSample), 0, true},
		}
		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				if tc.wantPanic {
					assert.Panics(t, func() { tc.result.MustGet("context") })
				} else {
					assert.Equal(t, tc.want, tc.result.MustGet("context"))
				}
			})
		}
	})

	t.Run("string", func(t *testing.T) {
		tests := []struct {
			name      string
			result    *Result[string]
			want      string
			wantPanic bool
		}{
			{"returns value on success", NewResult("hello", nil), "hello", false},
			{"panics on error", NewResult("", errSample), "", true},
		}
		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				if tc.wantPanic {
					assert.Panics(t, func() { tc.result.MustGet("context") })
				} else {
					assert.Equal(t, tc.want, tc.result.MustGet("context"))
				}
			})
		}
	})
}

func TestMustPresent(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		tests := []struct {
			name      string
			result    *Result[int]
			wantPanic bool
		}{
			{"non-zero does not panic", NewValueResult(1), false},
			{"zero value panics", NewValueResult(0), true},
		}
		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				if tc.wantPanic {
					assert.Panics(t, func() { tc.result.MustPresent("msg") })
				} else {
					assert.NotPanics(t, func() { tc.result.MustPresent("msg") })
				}
			})
		}
	})

	t.Run("string", func(t *testing.T) {
		tests := []struct {
			name      string
			result    *Result[string]
			wantPanic bool
		}{
			{"non-empty string does not panic", NewValueResult("hello"), false},
			{"empty string panics", NewValueResult(""), true},
		}
		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				if tc.wantPanic {
					assert.Panics(t, func() { tc.result.MustPresent("msg") })
				} else {
					assert.NotPanics(t, func() { tc.result.MustPresent("msg") })
				}
			})
		}
	})
}
