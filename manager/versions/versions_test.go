package versions

import (
	"reflect"
	"sort"
	"testing"
)

func TestManagerVersionListSort(t *testing.T) {
	testCases := []struct {
		unsorted, sorted []string
	}{
		// Empty config
		{},
		{
			unsorted: []string{"0.0.1423182830"},
			sorted:   []string{"0.0.1423182830"},
		},
		{
			unsorted: []string{"0.1.0-alpha-000697_g23f938c"},
			sorted:   []string{"0.1.0-alpha-000697_g23f938c"},
		},
		{
			unsorted: []string{"0.1.0-alpha-000697_g23f938c", "0.0.1423182830"},
			sorted:   []string{"0.0.1423182830", "0.1.0-alpha-000697_g23f938c"},
		},
		{
			unsorted: []string{"0.1.0-alpha-000697_g23f938c", "0.2.0", "0.0.1423182830", "0.1.0"},
			sorted:   []string{"0.0.1423182830", "0.1.0-alpha-000697_g23f938c", "0.1.0", "0.2.0"},
		},
		{
			unsorted: []string{"0.1.0-alpha-000696_g91c3684", "0.1.0-beta-000697_g23f938c", "0.2.0-alpha-000697_g23f938c", "0.1.0", "0.2.0"},
			sorted:   []string{"0.1.0-alpha-000696_g91c3684", "0.1.0-beta-000697_g23f938c", "0.1.0", "0.2.0-alpha-000697_g23f938c", "0.2.0"},
		},
		// backcompat with pre-semantic versioning
		{
			unsorted: []string{"0.0.1423182830", "0.1.0-alpha-000697_g23f938c", "0.1.0", "0.2.0"},
			sorted:   []string{"0.0.1423182830", "0.1.0-alpha-000697_g23f938c", "0.1.0", "0.2.0"},
		},
		// backcompat with pre-numeric dev versioning
		{
			unsorted: []string{"0.1.0-alpha-000696_g91c3684", "0.1.0-beta-000697_g23f938c", "0.2.0-alpha-000697_g23f938c", "0.1.0", "0.2.0", "0.1.0-alpha_91c368401a6f", "0.1.0-beta_23f938c842c8"},
			sorted:   []string{"0.1.0-alpha_91c368401a6f", "0.1.0-alpha-000696_g91c3684", "0.1.0-beta_23f938c842c8", "0.1.0-beta-000697_g23f938c", "0.1.0", "0.2.0-alpha-000697_g23f938c", "0.2.0"},
		},
	}

	for _, testCase := range testCases {
		sort.Sort(ManagerVersionList(testCase.unsorted))
		assertManagerVersionListsEqual(t, testCase.unsorted, testCase.sorted)
	}
}

func assertManagerVersionListsEqual(t *testing.T, received, expected ManagerVersionList) {
	if !reflect.DeepEqual(received, expected) {
		t.Fatalf("Manager version lists not equal!\nReceived: %#v\nExpected: %#v", received, expected)
	}
}

func TestPatchFor(t *testing.T) {
	versionList := ManagerVersionList{
		"bogus", // make sure we don't break when finding an invalid semver in the list
		"1.0.0", "1.0.1", "1.0.3", "1.0.4", "1.0.12",
		"1.1.0", "1.1.1", "1.2.0", "1.3.0", "1.3.1",
		"2.0.0", "2.0.1", "2.1.0", "2.1.1", "2.1.2",
	}

	tests := []struct {
		value, expected string
	}{
		{value: "1.0.0", expected: "1.0.12"},
		{value: "1.0.1", expected: "1.0.12"},
		{value: "1.0.12", expected: "1.0.12"},
		{value: "1.1.0", expected: "1.1.1"},
		{value: "1.2.0", expected: "1.2.0"},
		{value: "2.0.0-alpha-000697_g23f938c", expected: "2.0.1"},
		{value: "1.2.0-rc2-000555_g2317eab", expected: "1.2.0-rc2-000555_g2317eab"},
		{value: "not.a.semver", expected: "not.a.semver"},
	}

	for _, testCase := range tests {
		if result := versionList.PatchFor(testCase.value); result != testCase.expected {
			t.Fatalf("Patch for version %s expected %s but found %s instead", testCase.value, testCase.expected, result)
		}
	}
}
