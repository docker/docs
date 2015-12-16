package utils

import (
	"github.com/stretchr/testify/assert"
	"sort"
	"strings"
	"testing"
)

func TestRoleListLen(t *testing.T) {
	rl := RoleList{"foo", "bar"}
	assert.Equal(t, 2, rl.Len())
}

func TestRoleListLess(t *testing.T) {
	rl := RoleList{"foo", "foo/bar", "bar/foo"}
	assert.True(t, rl.Less(0, 1))
	assert.False(t, rl.Less(1, 2))
	assert.False(t, rl.Less(2, 1))
}

func TestRoleListSwap(t *testing.T) {
	rl := RoleList{"foo", "bar"}
	rl.Swap(0, 1)
	assert.Equal(t, "bar", rl[0])
	assert.Equal(t, "foo", rl[1])
}

func TestRoleListSort(t *testing.T) {
	rl := RoleList{"foo/bar", "foo", "bar", "bar/foo/baz", "bar/foo"}
	sort.Sort(rl)
	for i, s := range rl {
		if i == 0 || i == 1 {
			segs := strings.Split(s, "/")
			assert.Len(t, segs, 1)
		} else if i == 2 || i == 3 {
			segs := strings.Split(s, "/")
			assert.Len(t, segs, 2)
		} else if i == 4 {
			segs := strings.Split(s, "/")
			assert.Len(t, segs, 3)
		} else {
			// there are elements present that shouldn't be there
			t.Fail()
		}
	}
}
