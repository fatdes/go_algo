package lwwes

import "fmt"

// Element is associated to a timestamp and data
type Element struct {
	Timestamp int64
	Data      interface{}
}

func (e *Element) String() string {
	return fmt.Sprintf("Data:%v Timestamp:%d", e.Data, e.Timestamp)
}

// Equal return true if the data are equal
type Equal func(interface{}, interface{}) bool

// Set is a list of Elements
// it is not thread-safe, classes use this should guard the operations if necessary
type Set struct {
	Elements []*Element

	DataEqual Equal
}

// NewSet to create a set
func NewSet(dataEqual Equal) *Set {
	if dataEqual == nil {
		panic("dataEqual must not be nil")
	}

	return &Set{
		Elements:  []*Element{},
		DataEqual: dataEqual,
	}
}

// NewSetWithElements to create a set with predefined elements
func NewSetWithElements(elements []*Element, dataEqual Equal) *Set {
	return &Set{
		Elements:  elements,
		DataEqual: dataEqual,
	}
}

// find return the index and element if any, otherwise -1 and nil
func (set *Set) find(data interface{}) (int, *Element) {
	for j := len(set.Elements) - 1; j >= 0; j-- {
		if set.DataEqual(set.Elements[j].Data, data) {
			return j, set.Elements[j]
		}
	}

	return -1, nil
}

// Find the element with the matched data
func (set *Set) Find(data interface{}) *Element {
	_, found := set.find(data)
	return found
}

// Add the element if not present, otherwise update to latest element
func (set *Set) Add(element *Element) {
	if index, found := set.find(element.Data); found == nil {
		set.Elements = append(set.Elements, element)
	} else {
		if found.Timestamp >= element.Timestamp {
			// ignore earlier and same update
			return
		}

		set.Elements[index] = element
	}
}

// Union with other set, duplicate elements of same timestamp are not added
func (set *Set) Union(other *Set) *Set {
	union := NewSet(set.DataEqual)

	for _, e := range set.Elements {
		union.Add(e)
	}
	for _, e := range other.Elements {
		union.Add(e)
	}

	return union
}
