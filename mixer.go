package mixer

import (
	"net/http"
)

type Handler[T any] func(T)
type ContextFn[T any] func(Context) T

type Mixer[T any] struct {
	fn          ContextFn[T]
	beforeHooks []Handler[T]
	afterHooks  []Handler[T]
}

func New[T any](fn ContextFn[T]) *Mixer[T] {
	return &Mixer[T]{
		fn: fn,
	}
}

func Classic() *Mixer[Context] {
	return New(func(c Context) Context {
		return c
	})
}

func (m *Mixer[T]) Handler(handler Handler[T]) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := NewResponseWriter(w)

		// Create Context
		base := &context{w: resp, r: r}
		context := m.fn(base)

		// Before Hooks
		for _, before := range m.beforeHooks {
			if resp.Written() {
				break
			}

			before(context)
		}

		// Call Handler
		if !resp.Written() {
			handler(context)
		}

		// After Hooks
		for _, after := range m.afterHooks {
			after(context)
		}
	})
}

func (m *Mixer[T]) Before(h Handler[T]) {
	m.beforeHooks = append(m.beforeHooks, h)
}

func (m *Mixer[T]) After(h Handler[T]) {
	m.afterHooks = append(m.afterHooks, h)
}
