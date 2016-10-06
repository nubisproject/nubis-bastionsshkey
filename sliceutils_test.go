package main

import (
	"testing"
)

var listWithDupes = []string{"foo", "foo", "bar", "baz", "baz", "bar"}

func TestRemoveDuplicates(t *testing.T) {
	RemoveDuplicates(&listWithDupes)
	properLength := 3
	dupeFreeLength := len(listWithDupes)
	if dupeFreeLength != properLength {
		t.Errorf("dupeFreeList size: %d, should be %d", dupeFreeLength, properLength)
	}
}

func TestSortUsers(t *testing.T) {
	properSortedUsers := []string{"bar", "baz", "foo"}
	sortedUsers := SortUsers(listWithDupes)
	properLength := 3
	sortedUsersLength := len(sortedUsers)
	if sortedUsersLength != properLength {
		t.Errorf("dupeFreeList size: %d, should be %d", sortedUsersLength, properLength)
	}
	for i, entry := range sortedUsers {
		if properSortedUsers[i] != entry {
			t.Error("slice did not sort correctly")
		}
	}
}
