package carrot

import "sync"

type indexMap struct {
	Items map[int]Info
	sync.RWMutex
}

func (m *indexMap) Set(key int, v Info) {
	m.Lock()
	m.Items[key] = v
	m.Unlock()
}

func (m *indexMap) Get(key int) (Info, bool) {
	m.RLock()
	item, ok := m.Items[key]
	m.RUnlock()
	return item, ok
}

func (m *indexMap) Has(key int) bool {
	m.RLock()
	_, ok := m.Items[key]
	m.RUnlock()
	return ok
}

func (m *indexMap) Del(key int) {
	m.Lock()
	delete(m.Items, key)
	m.Unlock()
}

type usersMap struct {
	Items map[int]*Users
	sync.RWMutex
}

func (m *usersMap) Set(key int, item *Users) {
	m.Lock()
	m.Items[key] = item
	m.Unlock()
}

func (m *usersMap) Get(key int) (*Users, bool) {
	m.RLock()
	item, ok := m.Items[key]
	m.RUnlock()
	return item, ok
}

func (m *usersMap) Has(key int) bool {
	m.RLock()
	_, ok := m.Items[key]
	m.RUnlock()
	return ok
}

func (m *usersMap) Del(key int) {
	m.Lock()
	delete(m.Items, key)
	m.Unlock()
}

func (m *usersMap) Truncate() {
	m.Lock()
	m.Items = make(map[int]*Users)
	m.Unlock()
}
