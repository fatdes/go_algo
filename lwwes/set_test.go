package lwwes_test

import (
	"fatdes/go_algo/lwwes"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ElementSet_RemoveData(t *testing.T) {
	type tc struct {
		name        string
		addSet      []*lwwes.Element
		removeThese []*lwwes.Element
		expect      []*lwwes.Element
	}

	tcs := []*tc{
		{
			name:   "not in add set - not remove",
			addSet: []*lwwes.Element{},
			removeThese: []*lwwes.Element{
				{Timestamp: 1, Data: "findme"},
			},
			expect: []*lwwes.Element{},
		},
		{
			name: "in add set and earlier timestamp - not remove",
			addSet: []*lwwes.Element{
				{Timestamp: 2, Data: "findme"},
			},
			removeThese: []*lwwes.Element{
				{Timestamp: 1, Data: "findme"},
			},
			expect: []*lwwes.Element{},
		},
		{
			name: "in add set and later timestamp - allow remove",
			addSet: []*lwwes.Element{
				{Timestamp: 1, Data: "findme"},
			},
			removeThese: []*lwwes.Element{
				{Timestamp: 2, Data: "findme"},
			},
			expect: []*lwwes.Element{
				{Timestamp: 2, Data: "findme"},
			},
		},
	}

	for _, tt := range tcs {
		es := lwwes.NewElementSet(
			lwwes.NewSetWithElements(tt.addSet, stringEqual),
			lwwes.NewSetWithElements([]*lwwes.Element{}, stringEqual),
		)

		for _, e := range tt.removeThese {
			es.RemoveData(e.Data, e.Timestamp)
		}

		assert.Equal(t, tt.expect, es.RemoveSet.Elements, tt.name)
	}
}
func Test_ElementSet_Lookup(t *testing.T) {
	type tc struct {
		name      string
		addSet    []*lwwes.Element
		removeSet []*lwwes.Element
		isMember  interface{}
		expect    *lwwes.Element
	}

	tcs := []*tc{
		{
			name:      "all empty - not a member",
			addSet:    nil,
			removeSet: nil,
			isMember:  "findme",
			expect:    nil,
		},
		{
			name: "not in both set - not a member",
			addSet: []*lwwes.Element{
				{Timestamp: 1, Data: "notme"},
			},
			removeSet: []*lwwes.Element{
				{Timestamp: 2, Data: "notme2"},
			},
			isMember: "findme",
			expect:   nil,
		},
		{
			name: "in add set only - a member",
			addSet: []*lwwes.Element{
				{Timestamp: 1, Data: "findme"},
			},
			removeSet: []*lwwes.Element{
				{Timestamp: 1, Data: "notme"},
			},
			isMember: "findme",
			expect:   &lwwes.Element{Timestamp: 1, Data: "findme"},
		},
		{
			name: "in remove set only - not a member",
			addSet: []*lwwes.Element{
				{Timestamp: 1, Data: "notme"},
			},
			removeSet: []*lwwes.Element{
				{Timestamp: 2, Data: "findme"},
			},
			isMember: "findme",
			expect:   nil,
		},
		{
			name: "in both set add ts > remove ts - a member",
			addSet: []*lwwes.Element{
				{Timestamp: 2, Data: "findme"},
			},
			removeSet: []*lwwes.Element{
				{Timestamp: 1, Data: "findme"},
			},
			isMember: "findme",
			expect:   &lwwes.Element{Timestamp: 2, Data: "findme"},
		},
		{
			name: "in both set add ts == remove ts - bias a member",
			addSet: []*lwwes.Element{
				{Timestamp: 1, Data: "findme"},
			},
			removeSet: []*lwwes.Element{
				{Timestamp: 1, Data: "findme"},
			},
			isMember: "findme",
			expect:   &lwwes.Element{Timestamp: 1, Data: "findme"},
		},
		{
			name: "in both set add ts < remove ts - not a member",
			addSet: []*lwwes.Element{
				{Timestamp: 1, Data: "findme"},
			},
			removeSet: []*lwwes.Element{
				{Timestamp: 2, Data: "findme"},
			},
			isMember: "findme",
			expect:   nil,
		},
	}

	for _, tt := range tcs {
		es := lwwes.NewElementSet(
			lwwes.NewSetWithElements(tt.addSet, stringEqual),
			lwwes.NewSetWithElements(tt.removeSet, stringEqual),
		)

		actual := es.Lookup(tt.isMember)
		assert.Equal(t, tt.expect, actual, tt.name)
	}
}

func Test_Concurrency(t *testing.T) {
	type tc struct {
		addThese    []*lwwes.Element
		removeThese []*lwwes.Element
	}

	tcs := []*tc{
		{
			addThese: []*lwwes.Element{
				{Timestamp: 1, Data: "notme"},
			},
			removeThese: []*lwwes.Element{
				{Timestamp: 1, Data: "notme2"},
			},
		},
		{
			addThese: []*lwwes.Element{
				{Timestamp: 1, Data: "findme"},
			},
			removeThese: []*lwwes.Element{
				{Timestamp: 1, Data: "notme"},
			},
		},
		{
			addThese: []*lwwes.Element{
				{Timestamp: 5, Data: "findme"},
			},
			removeThese: []*lwwes.Element{},
		},
		{
			addThese: []*lwwes.Element{
				{Timestamp: 5, Data: "notme"},
			},
			removeThese: []*lwwes.Element{
				{Timestamp: 3, Data: "findme"},
			},
		},
		{
			addThese: []*lwwes.Element{
				{Timestamp: 3, Data: "findme"},
			},
			removeThese: []*lwwes.Element{
				{Timestamp: 1, Data: "findme"},
			},
		},
	}

	var wg sync.WaitGroup

	es := lwwes.NewElementSet(
		lwwes.NewSet(stringEqual),
		lwwes.NewSet(stringEqual),
	)
	concurrentTest := func(tt *tc) {
		defer wg.Done()

		for _, e := range tt.addThese {
			es.AddData(e.Data, e.Timestamp)
		}
		for _, e := range tt.removeThese {
			es.RemoveData(e.Data, e.Timestamp)
		}
	}

	for _, tt := range tcs {
		wg.Add(1)
		go concurrentTest(tt)
	}

	wg.Wait()

	expect := &lwwes.Element{
		Data:      "findme",
		Timestamp: 5,
	}
	actual := es.Lookup("findme")
	assert.Equal(t, expect, actual)
}

func Test_ElementSet_Merge(t *testing.T) {
	type set struct {
		addSet    []*lwwes.Element
		removeSet []*lwwes.Element
	}

	type tc struct {
		name   string
		local  set
		remote set
		verify func(*lwwes.ElementSet)
	}

	tcs := []*tc{
		{
			name:   "empty local merged with empty remote",
			local:  set{addSet: []*lwwes.Element{}, removeSet: []*lwwes.Element{}},
			remote: set{addSet: []*lwwes.Element{}, removeSet: []*lwwes.Element{}},
			verify: func(actual *lwwes.ElementSet) {
				assert.Equal(t, 0, len(actual.AddSet.Elements))
				assert.Equal(t, 0, len(actual.RemoveSet.Elements))
			},
		},
		{
			name:  "empty local merged with remote",
			local: set{addSet: []*lwwes.Element{}, removeSet: []*lwwes.Element{}},
			remote: set{
				addSet: []*lwwes.Element{
					{Timestamp: 1, Data: "data1"},
				},
				removeSet: []*lwwes.Element{
					{Timestamp: 2, Data: "data2"},
				},
			},
			verify: func(actual *lwwes.ElementSet) {
				assert.Equal(t, 1, len(actual.AddSet.Elements))
				assert.Equal(t, &lwwes.Element{Timestamp: 1, Data: "data1"}, actual.AddSet.Find("data1"))
				assert.Equal(t, 1, len(actual.RemoveSet.Elements))
				assert.Equal(t, &lwwes.Element{Timestamp: 2, Data: "data2"}, actual.RemoveSet.Find("data2"))
			},
		},
		{
			name: "local merged with empty remote",
			local: set{
				addSet: []*lwwes.Element{
					{Timestamp: 1, Data: "data1"},
				},
				removeSet: []*lwwes.Element{
					{Timestamp: 2, Data: "data2"},
				},
			},
			remote: set{addSet: []*lwwes.Element{}, removeSet: []*lwwes.Element{}},
			verify: func(actual *lwwes.ElementSet) {
				assert.Equal(t, 1, len(actual.AddSet.Elements))
				assert.Equal(t, &lwwes.Element{Timestamp: 1, Data: "data1"}, actual.AddSet.Find("data1"))
				assert.Equal(t, 1, len(actual.RemoveSet.Elements))
				assert.Equal(t, &lwwes.Element{Timestamp: 2, Data: "data2"}, actual.RemoveSet.Find("data2"))
			},
		},
		{
			name: "local merged with remote",
			local: set{
				addSet: []*lwwes.Element{
					{Timestamp: 1, Data: "data1"},
				},
				removeSet: []*lwwes.Element{
					{Timestamp: 2, Data: "data2"},
				},
			},
			remote: set{
				addSet: []*lwwes.Element{
					{Timestamp: 1, Data: "data2"},
				},
				removeSet: []*lwwes.Element{
					{Timestamp: 2, Data: "data1"},
				},
			},
			verify: func(actual *lwwes.ElementSet) {
				assert.Equal(t, 2, len(actual.AddSet.Elements))
				assert.Equal(t, &lwwes.Element{Timestamp: 1, Data: "data1"}, actual.AddSet.Find("data1"))
				assert.Equal(t, &lwwes.Element{Timestamp: 1, Data: "data2"}, actual.AddSet.Find("data2"))
				assert.Equal(t, 2, len(actual.RemoveSet.Elements))
				assert.Equal(t, &lwwes.Element{Timestamp: 2, Data: "data1"}, actual.RemoveSet.Find("data1"))
				assert.Equal(t, &lwwes.Element{Timestamp: 2, Data: "data2"}, actual.RemoveSet.Find("data2"))
			},
		},
		{
			name: "local merged with remote - commutative",
			local: set{
				addSet: []*lwwes.Element{
					{Timestamp: 1, Data: "data2"},
				},
				removeSet: []*lwwes.Element{
					{Timestamp: 2, Data: "data1"},
				},
			},
			remote: set{
				addSet: []*lwwes.Element{
					{Timestamp: 1, Data: "data1"},
				},
				removeSet: []*lwwes.Element{
					{Timestamp: 2, Data: "data2"},
				},
			},
			verify: func(actual *lwwes.ElementSet) {
				assert.Equal(t, 2, len(actual.AddSet.Elements))
				assert.Equal(t, &lwwes.Element{Timestamp: 1, Data: "data1"}, actual.AddSet.Find("data1"))
				assert.Equal(t, &lwwes.Element{Timestamp: 1, Data: "data2"}, actual.AddSet.Find("data2"))
				assert.Equal(t, 2, len(actual.RemoveSet.Elements))
				assert.Equal(t, &lwwes.Element{Timestamp: 2, Data: "data1"}, actual.RemoveSet.Find("data1"))
				assert.Equal(t, &lwwes.Element{Timestamp: 2, Data: "data2"}, actual.RemoveSet.Find("data2"))
			},
		},
	}

	for _, tt := range tcs {
		local := lwwes.NewElementSet(
			lwwes.NewSetWithElements(tt.local.addSet, stringEqual),
			lwwes.NewSetWithElements(tt.local.removeSet, stringEqual),
		)
		remote := lwwes.NewElementSet(
			lwwes.NewSetWithElements(tt.remote.addSet, stringEqual),
			lwwes.NewSetWithElements(tt.remote.removeSet, stringEqual),
		)

		local.Merge(remote)

		tt.verify(local)
	}
}
