package main

import (
	"fmt"
	"testing"
)

var fixture = []string{"group1", "group2", "group3"}
var ldapGroup1 = fmt.Sprintf("%s%s%s%s%s", fixture[0], delimiter, fixture[1], delimiter, fixture[2])

func TestConfigLDAPGroupExplode(t *testing.T) {
	ret := ExplodeLDAPGroup(ldapGroup1, delimiter)
	if ret == nil {
		t.Errorf("ret should not be nil")
	}
	if len(ret) != len(fixture) {
		t.Errorf("ret length should be 3")
	}
	for index, value := range fixture {
		if ret[index] != value {
			t.Errorf("ret doesn't match fixture at index: %d", index)
		}
	}
}
