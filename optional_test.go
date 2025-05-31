// Copyright (c) 2024 Matt Way
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to
// deal in the Software without restriction, including without limitation the
// rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
// sell copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS
// IN THE THE SOFTWARE.

package optional_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.mway.dev/optional"
)

func TestSome(t *testing.T) {
	requireOptionalHasValue(t, true, optional.Some(true))
	requireOptionalHasValue(t, false, optional.Some(false))
	requireOptionalHasValue(t, 123, optional.Some(123))
	requireOptionalHasValue(t, t.Name(), optional.Some(t.Name()))
	requireOptionalHasValue(t, 0, optional.Some(0))
	requireOptionalHasValue(t, "", optional.Some(""))
	requireOptionalHasValue(t, t, optional.Some(t))
	requireOptionalHasValue(t, nil, optional.Some((*testing.T)(nil)))
}

func TestNone(t *testing.T) {
	none := optional.None[bool]()
	require.False(t, none.IsSome())
	require.True(t, none.IsNone())

	x, ok := none.Get()
	require.False(t, ok)
	require.Zero(t, x)
	require.Equal(t, false, none.ValueOr(false))
	require.Equal(t, true, none.ValueOr(true))
	require.Panics(t, func() {
		none.Value()
	})
}

func TestOptional_Filter(t *testing.T) {
	t.Run("none", func(t *testing.T) {
		var (
			give = optional.None[int]()
			have = give.Filter(func(x int) bool { return x > 0 })
		)
		require.False(t, have.IsSome())
		require.True(t, have.IsNone())
	})

	t.Run("some predicate returns true", func(t *testing.T) {
		var (
			give = optional.Some(123)
			have = give.Filter(func(x int) bool { return x > 0 })
		)
		require.True(t, have.IsSome())
		require.False(t, have.IsNone())
		require.Equal(t, give.Value(), have.Value())
	})

	t.Run("some predicate returns false", func(t *testing.T) {
		var (
			give = optional.Some(123)
			have = give.Filter(func(x int) bool { return x%2 == 0 })
		)
		require.False(t, have.IsSome())
		require.True(t, have.IsNone())
	})
}

func TestOptional_Get(t *testing.T) {
	var opt optional.Optional[bool]

	value, ok := opt.Get()
	require.False(t, ok)
	require.Zero(t, value)

	opt = optional.None[bool]()
	value, ok = opt.Get()
	require.False(t, ok)
	require.Zero(t, value)

	opt = optional.Some(false)
	value, ok = opt.Get()
	require.True(t, ok)
	require.False(t, value)
}

func TestOptional_IsNone(t *testing.T) {
	t.Run("some", func(t *testing.T) {
		opt := optional.Some(true)
		require.False(t, opt.IsNone())
	})

	t.Run("none", func(t *testing.T) {
		var opt optional.Optional[bool]
		require.True(t, opt.IsNone())

		opt = optional.None[bool]()
		require.True(t, opt.IsNone())
	})
}

func TestOptional_IsNoneOr(t *testing.T) {
	t.Run("some", func(t *testing.T) {
		opt := optional.Some(true)
		require.False(t, opt.IsNoneOr(func(bool) bool {
			return false
		}))
		require.True(t, opt.IsNoneOr(func(x bool) bool {
			return x
		}))
	})

	t.Run("none", func(t *testing.T) {
		require.NotPanics(t, func() {
			opt := optional.None[bool]()
			require.True(t, opt.IsNoneOr(nil))
			require.True(t, opt.IsNoneOr(func(bool) bool {
				panic("never called")
			}))
		})
	})
}

func TestOptional_IsSome(t *testing.T) {
	t.Run("none", func(t *testing.T) {
		var opt optional.Optional[bool]
		require.False(t, opt.IsSome())

		opt = optional.None[bool]()
		require.False(t, opt.IsSome())
	})

	t.Run("some", func(t *testing.T) {
		opt := optional.Some(true)
		require.True(t, opt.IsSome())
	})
}

func TestOptional_IsSomeAnd(t *testing.T) {
	t.Run("none", func(t *testing.T) {
		opt := optional.None[int]()
		require.False(t, opt.IsSomeAnd(func(int) bool {
			return true
		}))
		require.NotPanics(t, func() {
			require.False(t, opt.IsSomeAnd(nil))
		})
	})

	t.Run("some", func(t *testing.T) {
		opt := optional.Some(123)
		require.True(t, opt.IsSomeAnd(func(int) bool {
			return true
		}))
		require.False(t, opt.IsSomeAnd(func(int) bool {
			return false
		}))
		require.Panics(t, func() {
			opt.IsSomeAnd(nil)
		})
	})
}

func TestOptional_Map(t *testing.T) {
	t.Run("some", func(t *testing.T) {
		var (
			give = optional.Some(123)
			have = give.Map(func(x int) int { return x * 2 })
		)
		require.True(t, have.IsSome())
		require.Equal(t, give.Value()*2, have.Value())
	})

	t.Run("none", func(t *testing.T) {
		var (
			give = optional.None[int]()
			have = give.Map(func(x int) int { return x * 2 })
		)
		require.True(t, have.IsNone())
	})
}

func TestOptional_MapOr(t *testing.T) {
	t.Run("some", func(t *testing.T) {
		var (
			give = optional.Some(123)
			have = give.MapOr(0, func(x int) int { return x * 2 })
		)
		require.Equal(t, give.Value()*2, have)
	})

	t.Run("none", func(t *testing.T) {
		var (
			give = optional.None[int]()
			have = give.MapOr(123, func(x int) int { return x * 2 })
		)
		require.Equal(t, 123, have)
	})
}

func TestOptional_MapOrElse(t *testing.T) {
	t.Run("some", func(t *testing.T) {
		var (
			give = optional.Some(123)
			have = give.MapOrElse(
				func() int { return 999 },
				func(x int) int { return x * 2 },
			)
		)
		require.Equal(t, 246, have)
	})

	t.Run("none", func(t *testing.T) {
		var (
			give = optional.None[int]()
			have = give.MapOrElse(
				func() int { return 123 },
				func(x int) int { return x * 2 },
			)
		)
		require.Equal(t, 123, have)
	})
}

func TestOptional_Or(t *testing.T) {
	t.Run("some", func(t *testing.T) {
		var (
			give = optional.Some(123)
			have = give.Or(optional.Some(456))
		)
		require.True(t, have.IsSome())
		require.Equal(t, give.Value(), have.Value())
	})

	t.Run("none", func(t *testing.T) {
		var (
			give = optional.None[int]()
			have = give.Or(optional.Some(123))
		)
		require.True(t, have.IsSome())
		require.Equal(t, 123, have.Value())
	})
}

func TestOptional_OrElse(t *testing.T) {
	t.Run("some", func(t *testing.T) {
		var (
			give = optional.Some(123)
			have = give.OrElse(func() optional.Optional[int] {
				return optional.Some(456)
			})
		)
		require.True(t, have.IsSome())
		require.Equal(t, give.Value(), have.Value())
	})

	t.Run("none", func(t *testing.T) {
		var (
			give = optional.None[int]()
			have = give.OrElse(func() optional.Optional[int] {
				return optional.Some(123)
			})
		)
		require.True(t, have.IsSome())
		require.Equal(t, 123, have.Value())
	})
}

func TestOptional_Swap(t *testing.T) {
	t.Run("some", func(t *testing.T) {
		var (
			give = optional.Some(123)
			have = give.Swap(999)
		)

		require.True(t, give.IsSome())
		require.True(t, have.IsSome())
		require.Equal(t, 123, have.Value())
		require.Equal(t, 999, give.Value())
	})

	t.Run("none", func(t *testing.T) {
		var (
			give = optional.None[int]()
			have = give.Swap(999)
		)

		require.True(t, give.IsSome())
		require.True(t, have.IsNone())
		require.Equal(t, 999, give.Value())
	})
}

func TestOptional_Value(t *testing.T) {
	var opt optional.Optional[bool]
	require.Panics(t, func() {
		opt.Value()
	})

	opt = optional.Some(false)
	require.NotPanics(t, func() {
		require.False(t, opt.Value())
	})
}

func TestOptional_ValueOr(t *testing.T) {
	var opt optional.Optional[int]
	require.Equal(t, 123, opt.ValueOr(123))

	opt = optional.None[int]()
	require.Equal(t, 234, opt.ValueOr(234))

	opt = optional.Some(345)
	require.Equal(t, 345, opt.ValueOr(-1))
}

func TestOptional_ValueOrFunc(t *testing.T) {
	var opt optional.Optional[int]
	require.Equal(t, 123, opt.ValueOrElse(func() int { return 123 }))

	opt = optional.None[int]()
	require.Equal(t, 234, opt.ValueOrElse(func() int { return 234 }))

	opt = optional.Some(345)
	require.Equal(t, 345, opt.ValueOrElse(func() int { return -1 }))
}

func TestMap(t *testing.T) {
	t.Run("some", func(t *testing.T) {
		var (
			give = optional.Some(123)
			have = optional.Map(give, func(v int) bool {
				return v%2 == 0
			})
		)

		require.True(t, have.IsSome())
		require.False(t, have.Value())
	})

	t.Run("none", func(t *testing.T) {
		var (
			give = optional.None[int]()
			have = optional.Map(give, func(v int) bool {
				return v%2 == 0
			})
		)

		require.True(t, have.IsNone())
	})
}

func TestMapOr(t *testing.T) {
	t.Run("some", func(t *testing.T) {
		var (
			give = optional.Some(123)
			have = optional.MapOr(give, true, func(v int) bool {
				return v%2 == 0
			})
		)

		require.False(t, have)
	})

	t.Run("none", func(t *testing.T) {
		var (
			give = optional.None[int]()
			have = optional.MapOr(give, true, func(v int) bool {
				return v%2 == 0
			})
		)

		require.True(t, have)
	})
}

func TestMapOrElse(t *testing.T) {
	t.Run("some", func(t *testing.T) {
		var (
			give = optional.Some(123)
			have = optional.MapOrElse(
				give,
				func() bool { return true },
				func(v int) bool { return v%2 == 0 },
			)
		)

		require.False(t, have)
	})

	t.Run("none", func(t *testing.T) {
		var (
			give = optional.None[int]()
			have = optional.MapOrElse(
				give,
				func() bool { return true },
				func(v int) bool { return v%2 == 0 },
			)
		)

		require.True(t, have)
	})
}

func requireOptionalHasValue[T any](
	t *testing.T,
	want T,
	maybe optional.Optional[T],
) {
	require.True(t, maybe.IsSome())
	require.False(t, maybe.IsNone())

	have, ok := maybe.Get()
	require.True(t, ok)
	require.Equal(t, want, have)
	require.NotPanics(t, func() {
		require.Equal(t, want, maybe.Value())
	})
}
