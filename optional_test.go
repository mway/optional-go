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

func requireOptionalHasValue[T any](
	t *testing.T,
	want T,
	maybe optional.Optional[T],
) {
	require.True(t, maybe.HasValue())

	have, ok := maybe.Get()
	require.True(t, ok)
	require.Equal(t, want, have)
	require.NotPanics(t, func() {
		require.Equal(t, want, maybe.Value())
	})
}

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
	require.False(t, none.HasValue())

	x, ok := none.Get()
	require.False(t, ok)
	require.Zero(t, x)
	require.Equal(t, false, none.ValueOr(false))
	require.Equal(t, true, none.ValueOr(true))
	require.Panics(t, func() {
		none.Value()
	})
}

func TestOptional_HasValue(t *testing.T) {
	var opt optional.Optional[bool]
	require.False(t, opt.HasValue())

	opt = optional.None[bool]()
	require.False(t, opt.HasValue())

	opt = optional.Some(false)
	require.True(t, opt.HasValue())
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
	require.Equal(t, 123, opt.ValueOrFunc(func() int { return 123 }))

	opt = optional.None[int]()
	require.Equal(t, 234, opt.ValueOrFunc(func() int { return 234 }))

	opt = optional.Some(345)
	require.Equal(t, 345, opt.ValueOrFunc(func() int { return -1 }))
}
