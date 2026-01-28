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

import (
	"fmt"
)

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

// FromPointer produces an [Optional] that holds the dereferenced pointer value
// if the pointer is not nil.
func FromPointer[T any](ptr *T) Optional[T] {
	if ptr == nil {
		return None[T]()
	}
	return Some(*ptr)
}

// FromValue produces an [Optional] that holds the given value if it is not the
// zero value for its type.
func FromValue[T comparable](value T) Optional[T] {
	var zero T
	if value == zero {
		return None[T]()
	}
	return Some(value)
}

// Filter returns a new [Optional] containing the held value if the given
// predicate returns true. Otherwise, None() is returned. Panics if the
// predicate is nil when a value is held.
func (o *Optional[T]) Filter(pred func(T) bool) Optional[T] {
	if !o.isset {
		return None[T]()
	}
	if pred == nil {
		panicOptionalNilValue[T]("Filter", "predicate func")
	}
	if pred(o.value) {
		return *o
	}
	return None[T]()
}

// Get returns a value of type T (either the held value, or the zero value of
// T), and a boolean indicating if the value was held.
func (o *Optional[T]) Get() (T, bool) {
	return o.value, o.isset
}

// IsSome indicates if a value is held.
func (o *Optional[T]) IsSome() bool {
	return o.isset
}

// IsSomeAnd indicates if a value is held and the given predicate returns true.
// Panics if the predicate is nil when a value is held.
func (o *Optional[T]) IsSomeAnd(pred func(T) bool) bool {
	if !o.isset {
		return false
	}
	if pred == nil {
		panicOptionalNilValue[T]("IsSomeAnd", "predicate func")
	}
	return pred(o.value)
}

// IsNone indicates if a value is not held.
func (o *Optional[T]) IsNone() bool {
	return !o.isset
}

// IsNoneOr indicates if a value is not held or the given predicate returns
// true. Panics if the predicate is nil when a value is held.
func (o *Optional[T]) IsNoneOr(pred func(T) bool) bool {
	if !o.isset {
		return true
	}
	if pred == nil {
		panicOptionalNilValue[T]("IsNoneOr", "predicate func")
	}
	return pred(o.value)
}

// Map applies the given function to the held value, if present, and returns a
// new [Optional] containing the result. If no value is present, None() is
// returned instead. Panics if the mapper is nil when a value is held.
func (o *Optional[T]) Map(fn func(T) T) Optional[T] {
	switch {
	case o.isset:
		if fn == nil {
			panicOptionalNilValue[T]("Map", "mapper func")
		}
		return Some(fn(o.value))
	default:
		return None[T]()
	}
}

// MapOr returns the result of passing the option's value to the given function
// if a value is held. If no value is held, fallback is returned. Panics if the
// mapper is nil when a value is held.
func (o *Optional[T]) MapOr(fallback T, fn func(T) T) T {
	switch {
	case o.isset:
		if fn == nil {
			return panicOptionalNilValue[T]("MapOr", "mapper func")
		}
		return fn(o.value)
	default:
		return fallback
	}
}

// MapOrElse returns the result of passing the option's value to the given
// function if a value is held. If no value is held, the value produced by the
// fallback function is returned. Panics if the mapper is nil when a value is
// held, or if the fallback is nil when no value is held.
func (o *Optional[T]) MapOrElse(fallback func() T, fn func(T) T) T {
	switch {
	case o.isset:
		if fn == nil {
			return panicOptionalNilValue[T]("MapOrElse", "mapper func")
		}
		return fn(o.value)
	default:
		if fallback == nil {
			return panicOptionalNilValue[T]("MapOrElse", "fallback func")
		}
		return fallback()
	}
}

// Or returns the current option if it holds a value, or otherwise returns the
// given [Optional].
func (o *Optional[T]) Or(other Optional[T]) Optional[T] {
	if o.isset {
		return *o
	}
	return other
}

// OrElse returns the current option if it holds a value, or otherwise returns
// the option produced by the given fallback function. Panics if the fallback
// is nil when no value is held.
func (o *Optional[T]) OrElse(fallback func() Optional[T]) Optional[T] {
	switch {
	case o.isset:
		return *o
	default:
		if fallback == nil {
			panicOptionalNilValue[T]("OrElse", "fallback func")
			return None[T]()
		}
		return fallback()
	}
}

// Swap stores the given value in the option. If the optionl previously held a
// value, it is returned in a new [Optional], otherwise None[T]() is
// returned.
func (o *Optional[T]) Swap(value T) Optional[T] {
	var (
		isset = o.isset
		prev  = o.value
	)

	o.value = value
	o.isset = true

	if isset {
		return Some(prev)
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
// value is held. Panics if the fallback is nil when no value is held.
func (o *Optional[T]) ValueOrElse(fallback func() T) T {
	switch {
	case o.isset:
		return o.value
	default:
		if fallback == nil {
			return panicOptionalNilValue[T]("ValueOrElse", "fallback func")
		}
		return fallback()
	}
}

// Map converts the given [Optional] using the given transform function. Panics
// if the mapper is nil when a value is held.
func Map[In any, Out any](o Optional[In], fn func(In) Out) Optional[Out] {
	switch {
	case o.isset:
		if fn == nil {
			panicNilValue[Out]("Map", "mapper func")
		}
		return Some(fn(o.value))
	default:
		return None[Out]()
	}
}

// MapOr converts the given [Optional] to [Out] using the given transform
// function if it holds a value. If no value is held, fallback is returned.
// Panics if the mapper is nil when a value is held.
func MapOr[In any, Out any](
	o Optional[In],
	fallback Out,
	fn func(In) Out,
) Out {
	switch {
	case o.isset:
		if fn == nil {
			return panicNilValue[Out]("MapOr", "mapper func")
		}
		return fn(o.value)
	default:
		return fallback
	}
}

// MapOrElse converts the given [Optional] to [Out] using the given transform
// function if it holds a value. If no value is held, the value produced by
// the given fallback function is returned. Panics if the mapper is nil when a
// value is held, or if the fallback is nil when no value is held.
func MapOrElse[In any, Out any](
	o Optional[In],
	fallback func() Out,
	fn func(In) Out,
) Out {
	switch {
	case o.isset:
		if fn == nil {
			return panicNilValue[Out]("MapOrElse", "mapper func")
		}
		return fn(o.value)
	default:
		if fallback == nil {
			return panicNilValue[Out]("MapOrElse", "fallback func")
		}
		return fallback()
	}
}

func panicOptionalNilValue[T any](method string, value string) T {
	var zero T
	return panicNilValue[T](
		fmt.Sprintf("Optional[%T].%s", zero, method),
		value,
	)
}

func panicNilValue[T any](method string, value string) T {
	panic(fmt.Errorf("%s: nil %s", method, value))
}
