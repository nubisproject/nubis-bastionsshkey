package main

import (
	"sort"
)

func RemoveDuplicates(xs *[]string) {
	found := make(map[string]bool)
	j := 0
	for i, x := range *xs {
		if !found[x] {
			found[x] = true
			(*xs)[j] = (*xs)[i]
			j++
		}
	}
	*xs = (*xs)[:j]
}

func SortUsers(allEntries []string) []string {
	RemoveDuplicates(&allEntries)
	sort.Strings(allEntries)
	return allEntries

}
