package mixer

func (m *Mixer[T]) Clone() *Mixer[T] {
	beforeHooks := make([]Handler[T], len(m.beforeHooks))
	afterHooks := make([]Handler[T], len(m.afterHooks))
	copy(beforeHooks, m.beforeHooks)
	copy(afterHooks, m.afterHooks)
	return &Mixer[T]{
		fn:          m.fn,
		beforeHooks: beforeHooks,
		afterHooks:  afterHooks,
	}
}

func (m *Mixer[T]) WithBefore(h Handler[T]) *Mixer[T] {
	n := m.Clone()
	n.Before(h)
	return n
}

func (m *Mixer[T]) WithAfter(h Handler[T]) *Mixer[T] {
	n := m.Clone()
	n.After(h)
	return n
}
