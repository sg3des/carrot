package carrot

import "sync"

var index = &indexMap{items: make(map[int]Info)}

type indexMap struct {
	items map[int]Info
	sync.RWMutex
}

func (m *indexMap) Set(key int, v Info) {
	m.Lock()
	m.items[key] = v
	m.Unlock()
}

func (m *indexMap) Get(key int) (Info, bool) {
	m.RLock()
	item, ok := m.items[key]
	m.RUnlock()
	return item, ok
}

func (m *indexMap) Has(key int) bool {
	m.RLock()
	_, ok := m.items[key]
	m.RUnlock()
	return ok
}

func (m *indexMap) Del(key int) {
	m.Lock()
	delete(m.items, key)
	m.Unlock()
}

var cacheUsers = &cacheUsersMap{items: make(map[int]*Users)}

// var writeUsers = &cacheUsersMap{items: make(map[int]*Users)}

type cacheUsersMap struct {
	items map[int]*Users
	sync.RWMutex
}

func (m *cacheUsersMap) Set(key int, v *Users) {
	m.Lock()
	m.items[key] = v
	m.Unlock()
}

func (m *cacheUsersMap) Get(key int) (*Users, bool) {
	m.RLock()
	item, ok := m.items[key]
	m.RUnlock()
	return item, ok
}

func (m *cacheUsersMap) Has(key int) bool {
	m.RLock()
	_, ok := m.items[key]
	m.RUnlock()
	return ok
}

func (m *cacheUsersMap) Del(key int) {
	m.Lock()
	delete(m.items, key)
	m.Unlock()
}

func (m *cacheUsersMap) Truncate() {
	m.Lock()
	m.items = make(map[int]*Users)
	m.Unlock()
}
