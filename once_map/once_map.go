package once_map

import "sync"

type OnceMap[K comparable, V any] struct {
	data     sync.Map
	initLock sync.Map
}

type initItem[V any] struct {
	initOnce sync.Once
	makeFn   func() V
	value    V
}

func itemFromExisting[V any](v V) *initItem[V] {
	item := &initItem[V]{
		initOnce: sync.Once{},
		makeFn: func() V {
			return v
		},
		value: v,
	}
	item.initOnce.Do(func() {}) // mark as initialized
	return item
}

func (item *initItem[V]) unwrap() V {
	item.initOnce.Do(func() {
		item.value = item.makeFn()
	})
	return item.value
}

func (m *OnceMap[K, V]) Load(k K) (V, bool) {
	v, ok := m.data.Load(k)
	if ok {
		return v.(*initItem[V]).unwrap(), ok
	}
	return v.(V), ok
}

func (m *OnceMap[K, V]) Store(k K, v V) {
	m.data.Store(k, itemFromExisting(v))
}

func (m *OnceMap[K, V]) LoadOrStore(k K, v V) (actual V, loaded bool) {
	vv, loaded := m.data.LoadOrStore(k, itemFromExisting(v))
	return vv.(*initItem[V]).unwrap(), loaded
}

func (m *OnceMap[K, V]) LoadOrMake(k K, makeFn func() V) (actual V, loaded bool) {
	var nilV V
	v, loaded := m.data.LoadOrStore(k, &initItem[V]{
		initOnce: sync.Once{},
		makeFn:   makeFn,
		value:    nilV,
	})
	return v.(*initItem[V]).unwrap(), loaded
}
