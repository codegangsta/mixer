package mixer

import (
	"net/http"
)

// Handler is a generic representation of a http Handler in Mixer.
type Handler[T any] func(T)

// ContextFN is a generic function that is passed a mixer.Context and expects a
// T in return. Used for transforming mixer.Contexts into a custom type.
type ContextFn[T any] func(Context) T

// A Mixer repressents a container that can convert mix.Handlers into http.HandlerFuncs
type Mixer[T any] struct {
	fn          ContextFn[T]
	beforeHooks []Handler[T]
	afterHooks  []Handler[T]
}

// New Creates a new instance of a mixer, bound to context type T
func New[T any](fn ContextFn[T]) *Mixer[T] {
	return &Mixer[T]{
		fn: fn,
	}
}

// Classsic returns an instance of a mixer.Mixer that is bound to type mixer.Context
func Classic() *Mixer[Context] {
	return New(func(c Context) Context {
		return c
	})
}

// Handler wraps a mixer.Handler in a http.HandlerFunc, encapsulating all the
// functionality such as before/after funcs and context injection in a simple,
// portable package. You can use this function to interface with existing
// packages that utilize the http.Handler interface
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

// Adds a mixer.Handler to be executed before any primary handlers are
// executed. If a before func writes to the http.ResponseWriter on the context,
// execution is skipped for the rest of the handlers, and resumed with the
// first After handler.
func (m *Mixer[T]) Before(h Handler[T]) {
	m.beforeHooks = append(m.beforeHooks, h)
}

// Adds a mixer.Handler to be executed after any primary handlers are
// executed.
func (m *Mixer[T]) After(h Handler[T]) {
	m.afterHooks = append(m.afterHooks, h)
}
