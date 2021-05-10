package lwwes_test

import (
	"fatdes/go_algo/lwwes"
	"testing"

	"github.com/stretchr/testify/assert"
)

var stringEqual = func(d1 interface{}, d2 interface{}) bool {
	return d1.(string) == d2.(string)
}

func Test_Set_Add(t *testing.T) {
	type tc struct {
		name     string
		elements []*lwwes.Element
		expect   []*lwwes.Element
	}

	tcs := []*tc{
		{
			name: "different elements - all added",
			elements: []*lwwes.Element{
				{Timestamp: 1, Data: "data1"},
				{Timestamp: 2, Data: "data2"},
				{Timestamp: 3, Data: "data3"},
			},
			expect: []*lwwes.Element{
				{Timestamp: 1, Data: "data1"},
				{Timestamp: 2, Data: "data2"},
				{Timestamp: 3, Data: "data3"},
			},
		},
		{
			name: "duplicated elements - idempotent",
			elements: []*lwwes.Element{
				{Timestamp: 1, Data: "data1"},
				{Timestamp: 1, Data: "data1"},
				{Timestamp: 1, Data: "data1"},
			},
			expect: []*lwwes.Element{
				{Timestamp: 1, Data: "data1"},
			},
		},
		{
			name: "duplicated elements - updated timestamp",
			elements: []*lwwes.Element{
				{Timestamp: 1, Data: "data1"},
				{Timestamp: 2, Data: "data2"},
				{Timestamp: 3, Data: "data3"},
				{Timestamp: 4, Data: "data1"},
			},
			expect: []*lwwes.Element{
				{Timestamp: 4, Data: "data1"},
				{Timestamp: 2, Data: "data2"},
				{Timestamp: 3, Data: "data3"},
			},
		},
	}

	for _, tt := range tcs {
		set := lwwes.NewSet(stringEqual)

		for _, e := range tt.elements {
			set.Add(e)
		}

		assert.Equal(t, len(tt.expect), len(set.Elements))
		assert.Equal(t, tt.expect, set.Elements, tt.name)
	}
}

func Test_Set_Find(t *testing.T) {
	type tc struct {
		name      string
		elements  []*lwwes.Element
		dataEqual lwwes.Equal
		find      interface{}
		expect    *lwwes.Element
	}

	tcs := []*tc{
		{
			name:      "no element - not found",
			elements:  nil,
			dataEqual: stringEqual,
			find:      "findme",
			expect:    nil,
		},
		{
			name: "element - not found",
			elements: []*lwwes.Element{
				{Timestamp: 1, Data: "notme"},
				{Timestamp: 2, Data: "notme2"},
				{Timestamp: 3, Data: "notme3"},
			},
			dataEqual: stringEqual,
			find:      "findme",
			expect:    nil,
		},
		{
			name: "multiple elements - found",
			elements: []*lwwes.Element{
				{Timestamp: 1, Data: "notme"},
				{Timestamp: 2, Data: "findme"},
				{Timestamp: 3, Data: "notme2"},
			},
			dataEqual: stringEqual,
			find:      "findme",
			expect:    &lwwes.Element{Timestamp: 2, Data: "findme"},
		},
		{
			name: "multiple int elements - found",
			elements: []*lwwes.Element{
				{Timestamp: 1, Data: 999},
				{Timestamp: 2, Data: 888},
				{Timestamp: 3, Data: 777},
			},
			dataEqual: func(d1 interface{}, d2 interface{}) bool {
				return d1.(int) == d2.(int)
			},
			find:   888,
			expect: &lwwes.Element{Timestamp: 2, Data: 888},
		},
	}

	for _, tt := range tcs {
		set := lwwes.NewSetWithElements(tt.elements, tt.dataEqual)

		actual := set.Find(tt.find)
		assert.Equal(t, tt.expect, actual, tt.name)
	}
}

func Test_Set_Union(t *testing.T) {
	type tc struct {
		name   string
		set1   []*lwwes.Element
		set2   []*lwwes.Element
		expect []*lwwes.Element
	}

	tcs := []*tc{
		{
			name:   "empty union empty - empty",
			set1:   []*lwwes.Element{},
			set2:   []*lwwes.Element{},
			expect: []*lwwes.Element{},
		},
		{
			name: "empty union something - something",
			set1: []*lwwes.Element{},
			set2: []*lwwes.Element{
				{Timestamp: 1, Data: "data1"},
			},
			expect: []*lwwes.Element{
				{Timestamp: 1, Data: "data1"},
			},
		},
		{
			name: "something union empty - something",
			set1: []*lwwes.Element{
				{Timestamp: 1, Data: "data1"},
			},
			set2: []*lwwes.Element{},
			expect: []*lwwes.Element{
				{Timestamp: 1, Data: "data1"},
			},
		},
		{
			name: "something union something - something + something",
			set1: []*lwwes.Element{
				{Timestamp: 1, Data: "data1"},
			},
			set2: []*lwwes.Element{
				{Timestamp: 2, Data: "data2"},
			},
			expect: []*lwwes.Element{
				{Timestamp: 1, Data: "data1"},
				{Timestamp: 2, Data: "data2"},
			},
		},
		{
			name: "union later timestamp - updated timestamp",
			set1: []*lwwes.Element{
				{Timestamp: 1, Data: "data1"},
			},
			set2: []*lwwes.Element{
				{Timestamp: 4, Data: "data1"},
			},
			expect: []*lwwes.Element{
				{Timestamp: 4, Data: "data1"},
			},
		},
		{
			name: "union later timestamp - updated timestamp - commutative",
			set1: []*lwwes.Element{
				{Timestamp: 4, Data: "data1"},
			},
			set2: []*lwwes.Element{
				{Timestamp: 1, Data: "data1"},
			},
			expect: []*lwwes.Element{
				{Timestamp: 4, Data: "data1"},
			},
		},
	}

	for _, tt := range tcs {
		set1 := lwwes.NewSetWithElements(tt.set1, stringEqual)
		set2 := lwwes.NewSetWithElements(tt.set2, stringEqual)

		actual := set1.Union(set2)
		assert.Equal(t, tt.expect, actual.Elements, tt.name)
	}
}
