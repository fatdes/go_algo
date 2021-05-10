package lwwes

import "sync"

// ElementSet contains AddSet and RemoveSet
type ElementSet struct {
	AddSet    *Set
	RemoveSet *Set

	mutex *sync.Mutex
}

// NewElementSet creates with predefined add set and remove set
func NewElementSet(addSet *Set, removeSet *Set) *ElementSet {
	return &ElementSet{
		AddSet:    addSet,
		RemoveSet: removeSet,
		mutex:     &sync.Mutex{},
	}
}

// AddData to ElementSet
func (es *ElementSet) AddData(data interface{}, timestamp int64) {
	es.mutex.Lock()
	defer es.mutex.Unlock()
	es.AddSet.Add(&Element{Timestamp: timestamp, Data: data})
}

// RemoveData to ElementSet
func (es *ElementSet) RemoveData(data interface{}, timestamp int64) {
	es.mutex.Lock()
	defer es.mutex.Unlock()
	// only remove data if locally have the element
	if found := es.AddSet.Find(data); found != nil && found.Timestamp <= timestamp {
		es.RemoveSet.Add(&Element{Timestamp: timestamp, Data: data})
	}
}

// Lookup return element if data is a member of the element set
// bias to Add instead of Remove when equal timestamp
func (es *ElementSet) Lookup(data interface{}) *Element {
	es.mutex.Lock()
	defer es.mutex.Unlock()
	inAdd := es.AddSet.Find(data)
	inRemove := es.RemoveSet.Find(data)

	if inAdd == nil {
		return nil
	}

	if inRemove == nil {
		return inAdd
	}

	if inAdd.Timestamp >= inRemove.Timestamp {
		return inAdd
	}

	return nil
}

// LookupAll returns all members of the element set
func (es *ElementSet) LookupAll() []*Element {
	es.mutex.Lock()
	defer es.mutex.Unlock()

	result := []*Element{}
	for _, added := range es.AddSet.Elements {
		// bias add
		if removed := es.RemoveSet.Find(added.Data); removed == nil || removed.Timestamp <= added.Timestamp {
			result = append(result, added)
		}
	}
	return result
}

// Merge with other element set
func (es *ElementSet) Merge(other *ElementSet) {
	es.mutex.Lock()
	defer es.mutex.Unlock()

	es.AddSet = es.AddSet.Union(other.AddSet)
	es.RemoveSet = es.RemoveSet.Union(other.RemoveSet)
}
