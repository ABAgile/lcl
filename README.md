# lcl - the "water of life" for mental synchronization...

[![Release](https://img.shields.io/github/v/release/abagile/lcl?sort=semver&logo=Go&color=%23007D9C)](https://github.com/abagile/lcl/releases)
[![Test](https://github.com/abagile/lcl/actions/workflows/test.yml/badge.svg)](https://github.com/abagile/lcl/actions/workflows/test.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/abagile/lcl.svg)](https://pkg.go.dev/github.com/abagile/lcl)
[![Go Report Card](https://goreportcard.com/badge/github.com/abagile/lcl)](https://goreportcard.com/report/github.com/abagile/lcl)
[![codecov](https://codecov.io/gh/abagile/lcl/branch/main/graph/badge.svg)](https://codecov.io/gh/abagile/lcl)

A generic Go utility library providing functional-style helpers for slices, math, result handling, and tuple operations.

**Requires Go 1.24+**

```
go get github.com/abagile/lcl
```

---

## Design philosophy

### Iterator-first collection manipulation

All slice operations in `slices.go` are built on Go 1.23's `iter.Seq` / `iter.Seq2`
push-iterator model rather than operating on concrete slices internally. Functions
accept a `[]T` for ergonomics at the call site, but immediately convert to iterators
via `slices.Values` / `slices.All` and delegate the actual traversal to
[`github.com/BooleanCat/go-functional`](https://github.com/BooleanCat/go-functional).

`iter.Seq[V]` is defined as `func(yield func(V) bool)` — the iterator owns the
loop and pushes each value into the consumer's `yield` callback (an internal
iterator). This is distinct from a pull model, where the consumer calls `next()` to
request one value at a time; Go's `iter.Pull` can convert a push iterator into that
style when needed, but ranging directly over an `iter.Seq` stays in push mode and is
the idiomatic choice here.

This means:
- **Composition is lazy by default** — `Filter`, `Map`, etc. chain without allocating intermediate slices inside the pipeline.
- **`Error` variants short-circuit cleanly** — iteration stops the moment a predicate or mapper returns an error, with no extra bookkeeping in the caller.
- **Custom slice types are preserved** — all functions are constrained on `S ~[]T`, so a named type such as `type UserList []User` comes out the other side unchanged.

The iterator model is most valuable for operations that can terminate early — `Filter`, `Exclude`, `MapError`, `FoldError`, `ForEachWhile`, and similar. For operations that must consume the entire input before producing a result, the benefit is marginal. `Mean` and `Mode` fall into this category: they require a full pass regardless, so they accept a concrete `[]T` instead, which lets them chain directly with the output of any `slices.go` helper and makes the full-consumption contract explicit at the call site.

---

## slices.go — Functional slice operations

All functions preserve the element type of custom slice types (`S ~[]T`).

### Function naming conventions

Each operation comes in up to three variants that mirror Go's iterator types directly:

| Suffix | Iterator type | Meaning |
|---|---|---|
| *(none)* | `iter.Seq[V]` | Plain value sequence |
| `2` | `iter.Seq2[K, V]` | `(index, value)` pair — index is the slice position |
| `Error` | — | Same as the plain variant, but the predicate/mapper may return an `error`; iteration stops immediately on the first non-nil error |

So `Filter`, `Filter2`, and `FilterError` are the same operation expressed at three
levels of richness. Reach for the plain variant by default; add `2` when you need
the position, add `Error` when the callback can fail.

These helpers cover the common case of working with a concrete `[]T`. They are
intentionally thin — each one converts the slice to an iterator, delegates to
[`go-functional/it`](https://github.com/BooleanCat/go-functional), and collects the
result back into a slice. For pipelines that compose multiple steps, operate on
non-slice sources, or need finer control over iteration, use the `it` package
directly rather than chaining these helpers.

### Filtering
| Function | Description |
|---|---|
| `Filter(xs, pred)` | Keep elements where `pred` is true |
| `Filter2(xs, pred)` | Like `Filter` but predicate receives `(index, value)` |
| `FilterError(xs, pred)` | Like `Filter` but predicate may return an error; stops on first error |
| `Exclude(xs, pred)` | Keep elements where `pred` is false |
| `Exclude2(xs, pred)` | Like `Exclude` with index-aware predicate |
| `ExcludeError(xs, pred)` | Like `Exclude` but predicate may return an error |

### Transformation
| Function | Description |
|---|---|
| `Map(xs, mapper)` | Transform each element |
| `Map2(xs, mapper)` | Transform with `(index, value)` mapper |
| `MapError(xs, mapper)` | Transform; stops on first error |

### Folding
| Function | Description |
|---|---|
| `Fold(xs, accum, initial)` | Left fold |
| `Fold2(xs, accum, initial)` | Left fold with index |
| `FoldError(xs, accum, initial)` | Left fold; stops on first error |
| `FoldRight(xs, accum, initial)` | Right fold |
| `FoldRight2(xs, accum, initial)` | Right fold with index |
| `FoldRightError(xs, accum, initial)` | Right fold; stops on first error |

### Iteration
| Function | Description |
|---|---|
| `ForEach(xs, f)` | Call `f` for every element |
| `ForEach2(xs, f)` | Call `f(index, value)` for every element |
| `ForEachWhile(xs, pred)` | Iterate while `pred` returns true |

### Grouping
| Function | Description |
|---|---|
| `GroupBy(xs, mapper)` | `map[R]S` — groups share the same key; order within each group is preserved |
| `PartitionBy(xs, mapper)` | `[]S` — like `GroupBy` but returns an ordered slice of groups in first-seen key order, rather than a map |
| `FrequenciesBy(xs, mapper)` | `map[R]int` — count occurrences by key |

Pass `Id` as the mapper when the element itself is the key:

```go
// group/count by the value itself
GroupBy(words, Id)        // map[string][]string
FrequenciesBy(words, Id)  // map[string]int — equivalent to a word-count
PartitionBy(runs, Id)     // ordered groups, e.g. [1,1,2,2,1] → [[1,1,1],[2,2]]
```

### Padding
| Function | Description |
|---|---|
| `Pad(xs, n)` | Extend slice to length `n` with zero values; panics if `n < 0` |
| `PadWith(xs, n, padding)` | Like `Pad` but fills with `padding` |

---

## result.go — Panic-on-failure monadic helper

`Result` is designed for CLI tools and startup paths where certain failures have no
sensible recovery action — a missing required environment variable, a database
connection that must succeed, a NATS dial that is a hard dependency. Rather than
threading errors through every call frame, you wrap the fallible operation in a
`Result` and call a `Must*` method: if something went wrong the process panics with
a clear, formatted message; if everything is fine you get the value and move on.

```go
type Result[T comparable] struct {
    Value T
    Error error
}
```

| Constructor | Description |
|---|---|
| `NewResult(v, err)` | Generic constructor |
| `NewValueResult(v)` | Success result with nil error |
| `NewErrorResult(err)` | Error result with zero value |

### Methods

```go
func (r *Result[T]) Unwrap() (T, error)
func (r *Result[T]) Bind(f func(T) *Result[T]) *Result[T]
func (r *Result[T]) MustPass(msg string, v ...any)
func (r *Result[T]) MustGet(msg string, v ...any) T
func (r *Result[T]) MustPresent(msg string, v ...any)
```

- **`Bind`** — chains operations; short-circuits on error, making it easy to compose multiple fallible steps before the final `Must*` call.
- **`MustPass`** — panics with a formatted message if the result holds an error.
- **`MustGet`** — panics on error, otherwise returns the value.
- **`MustPresent`** — panics if the value is the zero value (uses `IsEmpty`); useful when a successful operation that returns an empty string or zero integer is itself considered a configuration error.

### Examples

**`NewResult` + `MustGet`** — wrap any `(T, error)` return and extract the value,
panicking with context if there was an error:

```go
dsn := NewResult(os.LookupEnv("DATABASE_URL")).MustGet("DATABASE_URL is required")
db  := NewResult(sql.Open("pgx", dsn)).MustGet("failed to open database: %s", dsn)
NewResult(db.PingContext(ctx)).MustPass("database ping failed")
```

**`NewValueResult` + `MustPresent`** — use when you already have a value (no error path) but want to assert it is non-zero. Typical for config values that are present but must not be empty:

```go
// os.Getenv returns "" on missing — no error, but empty is still wrong
NewValueResult(os.Getenv("REDIS_HOST")).MustPresent("REDIS_HOST must not be empty")

// or when reading from a parsed config struct
NewValueResult(cfg.NATSUrl).MustPresent("nats_url is required in config")
NewValueResult(cfg.AppSecret).MustPresent("app_secret is required in config")
```

**`NewErrorResult` + `MustPass`** — use when an operation produces only an error (no meaningful return value), and any non-nil error is fatal:

```go
NewErrorResult(os.MkdirAll(cfg.DataDir, 0755)).MustPass("failed to create data dir %q", cfg.DataDir)
NewErrorResult(nats.Publish("startup", payload)).MustPass("failed to publish startup event")
```

**`Bind`** — chain multiple fallible steps before the final assertion, short-circuiting on the first error:

```go
token := NewResult(os.LookupEnv("API_TOKEN")).
    Bind(func(t string) *Result[string] {
        if len(t) < 32 {
            return NewErrorResult[string](fmt.Errorf("token too short"))
        }
        return NewValueResult(strings.TrimSpace(t))
    }).
    MustGet("invalid API_TOKEN")
```

---

## math.go — Statistical helpers

### `MinMax`
```go
func MinMax[T cmp.Ordered](a, b T) (min, max T)
```
Returns `(min, max)` regardless of argument order.

### `Clamp`
```go
func Clamp[T cmp.Ordered](value, min, max T) T
```
Clamps `value` into `[min, max]`. Caller must ensure `min <= max`.

### `Mean`
```go
func Mean[T Number](xs []T) (T, bool)
```
Returns the arithmetic mean of the slice. Returns `(zero, false)` for an empty slice. Accepts the result of any `slices.go` helper directly.

### `Mode`
```go
func Mode[T comparable](xs []T) ([]T, bool)
```
Returns all values that appear with the highest frequency. Returns `([], false)` for an empty slice. Order of results is non-deterministic when multiple values tie. Accepts the result of any `slices.go` helper directly.

---

## types.go — Generic type utilities

### Zero / emptiness

```go
func Empty[T any]() T                          // returns zero value of T
func IsEmpty[T comparable](v T) bool           // v == zero
func IsNotEmpty[T comparable](v T) bool        // v != zero
```

### Coalesce

```go
func Coalesce[T comparable](values ...T) (T, bool)   // first non-zero value
func Coalesced[T comparable](values ...T) T           // first non-zero, or zero
```

### Pointer helpers

```go
func ToPtr[T comparable](v T) *T                      // nil if v is zero
func FromPtr[T any](ptr *T, fallback ...T) T          // dereference or fallback/zero
```

### Slice conversion

```go
func ToAnySlice[T any](in []T) []any
func FromAnySlice[T any](in []any) ([]T, bool)   // false if any element cannot be asserted to T
```

### Nested data access

`GetIn` and `SetIn` navigate nested `map[string]any`, `map[any]any`, and `[]any` structures
using a dot-separated path string. Slice elements are accessed by numeric index.

```go
func GetIn(data any, path string) (any, error)
func SetIn(data any, path string, value any) error
```

```go
data := map[string]any{
    "user": map[string]any{
        "scores": []any{10, 20, 30},
    },
}

v, _ := GetIn(data, "user.scores.1")  // 20
SetIn(data, "user.scores.1", 99)      // scores[1] = 99
```

`SetIn` automatically creates intermediate `map[string]any` nodes for missing keys.
It does **not** grow slices — the index must already be in bounds.

### Identity

```go
func Id[T any](v T) T
```

Returns its argument unchanged. Primarily useful as a first-class mapper when the
element itself should be the key — avoiding the need to write `func(x T) T { return x }`
at every call site:

```go
FrequenciesBy(tags, Id)    // count each tag
GroupBy(events, Id)        // bucket identical events together
PartitionBy(tokens, Id)    // ordered groups of equal tokens
```

---

## tuples.go — Tuple types and zip/unzip

### Types

```go
type Tuple2[A, B any]       struct{ A A; B B }
type Tuple3[A, B, C any]    struct{ A A; B B; C C }
type Tuple4[A, B, C, D any] struct{ A A; B B; C C; D D }
```

Each type has an `Unpack()` method returning the fields as multiple return values.

### Constructors

```go
T2(a, b)          Tuple2[A, B]
T3(a, b, c)       Tuple3[A, B, C]
T4(a, b, c, d)    Tuple4[A, B, C, D]
```

### Zip

Combines multiple slices into a slice of tuples. Output length equals the longest input; shorter inputs are padded with zero values.

```go
Zip2(a []A, b []B)             []Tuple2[A, B]
Zip3(a []A, b []B, c []C)      []Tuple3[A, B, C]
Zip4(a, b, c, d)               []Tuple4[A, B, C, D]
```

### Unzip

Splits a slice of tuples back into individual slices.

```go
Unzip2(tuples []Tuple2[A, B])          ([]A, []B)
Unzip3(tuples []Tuple3[A, B, C])       ([]A, []B, []C)
Unzip4(tuples []Tuple4[A, B, C, D])    ([]A, []B, []C, []D)
```

### CrossJoin

Produces the Cartesian product of the input slices. Returns an empty slice immediately if any input is empty.

```go
CrossJoin2(a, b)       []Tuple2[A, B]
CrossJoin3(a, b, c)    []Tuple3[A, B, C]
CrossJoin4(a, b, c, d) []Tuple4[A, B, C, D]
```

### Usage patterns

Tuples are best suited for **local, short-lived grouping within a single function or pipeline** — situations where defining a one-off named struct would be more ceremony than it is worth. If the same pairing appears in more than one place, or is part of a public API, a named struct is almost always clearer.

**1. Pairing a sort key with its original value**

```go
// attach a score without losing the original user
ranked := Map(users, func(u User) Tuple2[int, User] {
    return T2(scoreOf(u), u)
})
```

**2. Zipping parallel slices from separate data sources**

```go
// two separate API responses that logically belong together
rows := Zip2(invoiceIDs, amounts)
ForEach(rows, func(row Tuple2[string, float64]) {
    id, amount := row.Unpack()
    process(id, amount)
})
```

**3. CrossJoin for generating combinations (e.g. test matrix)**

```go
envs    := []string{"staging", "prod"}
regions := []string{"us-east", "eu-west", "ap-south"}

ForEach(CrossJoin2(envs, regions), func(tc Tuple2[string, string]) {
    env, region := tc.Unpack()
    deploy(env, region)
})
```

**4. Unzip to split a parsed CSV into typed columns**

```go
// "alice,30"  "bob,25"  "carol,28"
names, ages := Unzip2(Map(rows, func(row string) Tuple2[string, int] {
    parts := strings.SplitN(row, ",", 2)
    age, _ := strconv.Atoi(parts[1])
    return T2(parts[0], age)
}))
```

**5. Carrying the original alongside a transformed value in a pipeline**

```go
// keep original for later comparison after normalisation
pairs := Map(inputs, func(s string) Tuple2[string, string] {
    return T2(s, normalize(s))
})
changed := Filter(pairs, func(p Tuple2[string, string]) bool {
    return p.A != p.B
})
```

---

## db.go — PostgreSQL / pgx utilities

### Pool construction

```go
func NewPgxPool(connStr string, opts ...DatabaseConfigOption) (*pgxpool.Pool, error)
```

Creates a `pgxpool.Pool` from a connection string. Options are applied to the parsed config before the pool is created.

```go
type DatabaseConfigOption func(*pgxpool.Config)
```

Functional option type for `NewPgxPool`.

### Options

```go
func WithDecimalRegister() DatabaseConfigOption
```

Registers [`govalues/decimal`](https://github.com/govalues/decimal) with the pgx type map on every new connection, enabling transparent encode/decode of `decimal.Decimal` for PostgreSQL `numeric` columns.

```go
pool, err := lcl.NewPgxPool(connStr, lcl.WithDecimalRegister())
```

### Connection string sanitization

```go
func SantizeDbConn(connStr string) string
```

Strips `user=` and `password=` key-value pairs from a DSN and collapses surrounding whitespace. Safe for logging.

```go
lcl.SantizeDbConn("host=localhost user=admin password=secret dbname=app")
// → "host=localhost dbname=app"
```

### Placeholder conversion

```go
func ConvertPgPlaceholders(sql string, args ...any) (string, []any, error)
```

Converts PostgreSQL-style positional placeholders (`$1`, `$2`, …) to `?` placeholders, reordering `args` to match. Useful when targeting drivers that use `?` syntax (e.g. MySQL, SQLite) while authoring queries in PostgreSQL style.

```go
sql, args, err := lcl.ConvertPgPlaceholders(
    "SELECT * FROM orders WHERE user_id = $1 AND status = $2",
    userID, status,
)
// sql  → "SELECT * FROM orders WHERE user_id = ? AND status = ?"
// args → []any{userID, status}
```

Returns an error if a placeholder index is out of range or malformed. Positional args may be repeated (`$1` used twice maps the same argument into both positions).

### Struct copy with dereferencing

```go
func CopyDeref[T any, U any](src T, dst *U) (*U, error)
```

Copies fields from `src` to `dst` by name, with three conversion rules applied in order:

| Source type | Destination type | Behaviour |
|---|---|---|
| `[]byte` | `string` | Converts bytes to string |
| `*T` | `T` | Dereferences pointer; zeroes dst field if nil |
| `T` | `T` | Direct assignment if types are assignable |

Fields with no matching name in `src`, unexported fields, and fields whose types don't satisfy any rule are silently skipped. `src` may itself be a pointer — it is dereferenced before field matching begins.

```go
type Row struct {
    Name  []byte
    Score *int
    Tag   string
}
type Model struct {
    Name  string
    Score int
    Tag   string
}

score := 42
row := Row{Name: []byte("alice"), Score: &score, Tag: "vip"}
model, err := lcl.CopyDeref(row, &Model{})
// model → Model{Name: "alice", Score: 42, Tag: "vip"}
```

---

## Acknowledgements

The slice and iterator utilities in this library are built on top of
[**BooleanCat/go-functional**](https://github.com/BooleanCat/go-functional), which
provides the composable `iter.Seq`-based primitives (`it.Filter`, `it.Map`,
`it.Fold`, `it.ForEach`, and their `2` / `Error` variants) that power the
collection functions here. go-functional deserves full credit for the iterator
plumbing that makes lazy, allocation-free pipelines possible in idiomatic Go.
