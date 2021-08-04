package common

import (
	"sync"
)

// global map of objects
var (
	registry = &objTable{
		hash: make(map[ObjectID]interface{}),
	}
)

// ObjectID is the ID of some opaque object stored in Registry.
type ObjectID uint64

// Registry implements CRUD operations to map ID and objects.
type Registry interface {
	Create(interface{}) ObjectID
	Read(ObjectID) interface{}
	Update(ObjectID, interface{})
	Delete(ObjectID)
}

// NewRegistryMap creates new Registry which stores all objects in a
// map.
func NewRegistryMap() Registry {
	return &objTable{
		hash: make(map[ObjectID]interface{}),
	}
}

// map of objects references with lock
type objTable struct {
	sync.Mutex
	hash map[ObjectID]interface{}
	id   ObjectID
}

func (r *objTable) Create(obj interface{}) ObjectID {
	r.Lock()
	id := r.id
	r.hash[id] = obj
	r.id = id + 1
	r.Unlock()
	return id
}

func (r *objTable) Read(id ObjectID) interface{} {
	r.Lock()
	obj := r.hash[id]
	r.Unlock()
	return obj
}

func (r *objTable) Update(id ObjectID, obj interface{}) {
	r.Lock()
	r.hash[id] = obj
	r.Unlock()
}

func (r *objTable) Delete(id ObjectID) {
	r.Lock()
	delete(r.hash, id)
	r.Unlock()
}

type objArray struct {
	sync.Mutex
	array []interface{}
	cnt   uint64
}

func (r *objArray) Create(obj interface{}) ObjectID {
	r.Lock()
	id := ObjectID(len(r.array))
	r.array = append(r.array, obj)
	r.cnt++
	r.Unlock()
	return id
}

func (r *objArray) Read(id ObjectID) interface{} {
	r.Lock()
	obj := r.array[id]
	r.Unlock()
	return obj
}

func (r *objArray) Update(id ObjectID, obj interface{}) {
	r.Lock()
	r.array[id] = obj
	r.Unlock()
}

func (r *objArray) Delete(id ObjectID) {
	r.Lock()
	r.array[id] = nil
	if r.cnt--; r.cnt == 0 {
		r.array = make([]interface{}, 0, 10)
	}
	r.Unlock()
}

// NewRegistryArray creates new Registry as a linear array.
func NewRegistryArray() Registry {
	return &objArray{}
}
