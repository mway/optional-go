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

// Filter returns a new [Optional] containing the held value if the given
// predicate returns true. Otherwise, None() is returned.
func (o *Optional[T]) Filter(pred func(T) bool) Optional[T] {
	if o.isset && pred(o.value) {
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
func (o *Optional[T]) IsSomeAnd(pred func(T) bool) bool {
	return o.isset && pred(o.value)
}

// IsNone indicates if a value is not held.
func (o *Optional[T]) IsNone() bool {
	return !o.isset
}

// IsNoneOr indicates if a value is not held or the given predicate returns
// true.
func (o *Optional[T]) IsNoneOr(pred func(T) bool) bool {
	return !o.isset || pred(o.value)
}

// Map applies the given function to the held value, if present, and returns a
// new [Optional] containing the result. If no value is present, None() is
// returned instead.
func (o *Optional[T]) Map(fn func(T) T) Optional[T] {
	if o.isset {
		return Some(fn(o.value))
	}
	return None[T]()
}

// MapOr returns the result of passing the option's value to the given function
// if a value is held. If no value is held, fallback is returned.
func (o *Optional[T]) MapOr(fallback T, fn func(T) T) T {
	if o.isset {
		return fn(o.value)
	}
	return fallback
}

// MapOrElse returns the result of passing the option's value to the given
// function if a value is held. If no value is held, the value produced by the
// fallback function is returned.
func (o *Optional[T]) MapOrElse(fallback func() T, fn func(T) T) T {
	if o.isset {
		return fn(o.value)
	}
	return fallback()
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
// the option produced by the given function.
func (o *Optional[T]) OrElse(other func() Optional[T]) Optional[T] {
	if o.isset {
		return *o
	}
	return other()
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
// value is held.
func (o *Optional[T]) ValueOrElse(fallback func() T) T {
	if o.isset {
		return o.value
	}
	return fallback()
}

// Map converts the given [Optional] using the given transform function.
func Map[In any, Out any](
	o Optional[In],
	transform func(In) Out,
) Optional[Out] {
	if !o.isset {
		return None[Out]()
	}
	return Some(transform(o.value))
}

// MapOr converts the given [Optional] to [Out] using the given transform
// function if it holds a value. If no value is held, fallback is returned.
func MapOr[In any, Out any](
	o Optional[In],
	fallback Out,
	transform func(In) Out,
) Out {
	if !o.isset {
		return fallback
	}
	return transform(o.value)
}

// MapOrElse converts the given [Optional] to [Out] using the given transform
// function if it holds a value. If no value is held, the value produced by
// the given fallback function is returned.
func MapOrElse[In any, Out any](
	o Optional[In],
	fallback func() Out,
	transform func(In) Out,
) Out {
	if !o.isset {
		return fallback()
	}
	return transform(o.value)
}
