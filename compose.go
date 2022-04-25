package mixer

// Clone performs a deep copy.
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

// With clones the Mixer[T] and then calls setup on the fresh copy.
func (m *Mixer[T]) With(setup func(f *Mixer[T])) *Mixer[T] {
	n := m.Clone()
	if setup != nil {
		setup(n)
	}
	return n
}
