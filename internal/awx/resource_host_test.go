package awx

import (
	"slices"
	"testing"
)

func Test_intSetDiff(t *testing.T) {
	tests := []struct {
		name        string
		oldIDs      []int
		newIDs      []int
		wantAdded   []int
		wantRemoved []int
	}{
		{
			name:        "no change",
			oldIDs:      []int{1, 2, 3},
			newIDs:      []int{1, 2, 3},
			wantAdded:   nil,
			wantRemoved: nil,
		},
		{
			name:        "add groups",
			oldIDs:      []int{1},
			newIDs:      []int{1, 2, 3},
			wantAdded:   []int{2, 3},
			wantRemoved: nil,
		},
		{
			name:        "remove groups",
			oldIDs:      []int{1, 2, 3},
			newIDs:      []int{1},
			wantAdded:   nil,
			wantRemoved: []int{2, 3},
		},
		{
			name:        "add and remove groups",
			oldIDs:      []int{1, 2},
			newIDs:      []int{2, 3},
			wantAdded:   []int{3},
			wantRemoved: []int{1},
		},
		{
			name:        "empty to some",
			oldIDs:      []int{},
			newIDs:      []int{1, 2},
			wantAdded:   []int{1, 2},
			wantRemoved: nil,
		},
		{
			name:        "some to empty",
			oldIDs:      []int{1, 2},
			newIDs:      []int{},
			wantAdded:   nil,
			wantRemoved: []int{1, 2},
		},
		{
			name:        "both empty",
			oldIDs:      []int{},
			newIDs:      []int{},
			wantAdded:   nil,
			wantRemoved: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			added, removed := intSetDiff(tt.oldIDs, tt.newIDs)
			if !slices.Equal(added, tt.wantAdded) {
				t.Errorf("added = %v, want %v", added, tt.wantAdded)
			}
			if !slices.Equal(removed, tt.wantRemoved) {
				t.Errorf("removed = %v, want %v", removed, tt.wantRemoved)
			}
		})
	}
}

func Test_extractIntList(t *testing.T) {
	tests := []struct {
		name string
		raw  []interface{}
		want []int
	}{
		{"empty", []interface{}{}, []int{}},
		{"single", []interface{}{42}, []int{42}},
		{"multiple", []interface{}{1, 2, 3}, []int{1, 2, 3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractIntList(tt.raw)
			if !slices.Equal(got, tt.want) {
				t.Errorf("extractIntList() = %v, want %v", got, tt.want)
			}
		})
	}
}
