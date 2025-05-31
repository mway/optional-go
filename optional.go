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

// Package optional provides optional types and helpers.
package optional

// An Optional is a wrapper type that may or may not hold a value of type T.
type Optional[T any] struct {
	value T
	isset bool
}

// Some produces an [Optional] that holds the given value.
func Some[T any](value T) Optional[T] {
	return Optional[T]{
		value: value,
		isset: true,
	}
}

// None produces an [Optional] that holds no value.
func None[T any]() Optional[T] {
	return Optional[T]{}
}

// Get returns a value of type T (either the held value, or the zero value of
// T), and a boolean indicating if the value was held.
func (o *Optional[T]) Get() (T, bool) {
	return o.value, o.isset
}

// HasValue indicates whether a value is held.
func (o *Optional[T]) HasValue() bool {
	return o.isset
}

// Map applies the given function to the held value, if present, and returns a
// new [Optional] containing the result. If no value is present, None() is
// returned instead.
func (o Optional[T]) Map(fn func(T) T) Optional[T] {
	if o.isset {
		return Some(fn(o.value))
	}
	return None[T]()
}

// Filter returns a new [Optional] containing the held value if the given
// predicate returns true. Otherwise, None() is returned.
func (o Optional[T]) Filter(pred func(T) bool) Optional[T] {
	if o.isset && pred(o.value) {
		return o
	}
	return None[T]()
}

// Value returns the held value of type T, or panics if no value is held.
func (o *Optional[T]) Value() T {
	if !o.isset {
		panic("Optional[%T].Value() called with no held value")
	}
	return o.value
}

// ValueOr returns a value of type T, either the held value or fallback if no
// value is held.
func (o *Optional[T]) ValueOr(fallback T) T {
	if o.isset {
		return o.value
	}
	return fallback
}

// ValueOrElse returns a value of type T, either the held value or the result
// of fallback if no value is held. The given function is only evaluated if no
// value is held.
func (o *Optional[T]) ValueOrElse(fallback func() T) T {
	if o.isset {
		return o.value
	}
	return fallback()
}
