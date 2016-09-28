package schema

import (
	"encoding/hex"
	"reflect"
	"testing"
)

func TestTagDigestMatches(t *testing.T) {
	tests := []struct {
		Tag     Tag
		Input   string
		IsMatch bool
		IsError bool
	}{
		{
			Tag:     Tag{Digest: "sha256:1f8d6e1edee77de035d79ca992df4e5cc8d358ec38f527077a84945a79907566"},
			Input:   "1f8d6e1edee77de035d79ca992df4e5cc8d358ec38f527077a84945a79907566",
			IsMatch: true,
			IsError: false,
		},
		{
			Tag:     Tag{Digest: "sha256:1f8d6e1edee77de035d79ca992df4e5cc8d358ec38f527077a84945a79907566"},
			Input:   "1",
			IsMatch: false,
			IsError: false,
		},
		// errors, no prefix for digest
		{
			Tag:     Tag{Digest: "1f8d6e1edee77de035d79ca992df4e5cc8d358ec38f527077a84945a79907566"},
			Input:   "1",
			IsMatch: false,
			IsError: true,
		},
		// errors, invalid digest
		{
			Tag:     Tag{Digest: "sha256:1"},
			Input:   "1",
			IsMatch: false,
			IsError: true,
		},
	}

	for _, test := range tests {
		input, _ := hex.DecodeString(test.Input)

		output, err := test.Tag.DigestMatches(input)
		if output != test.IsMatch {
			t.Fatalf("unexpected result '%t', expected %t with digest '%s' and comparator '%s'", output, test.IsMatch, test.Tag.Digest, test.Input)
		}
		if err == nil && test.IsError {
			t.Fail()
		}
		if err != nil && !test.IsError {
			t.Fail()
		}
	}
}

func TestTagsFilterValid(t *testing.T) {
	tests := []struct {
		Input    Tags
		Expected Tags
	}{
		// Test 0: Valid tag
		{
			Input: Tags{
				Tag{
					Digest:     "sha256:1f8d6e1edee77de035d79ca992df4e5cc8d358ec38f527077a84945a79907566",
					Repository: "repo/a",
					DigestPK:   "repo/a@sha256:1f8d6e1edee77de035d79ca992df4e5cc8d358ec38f527077a84945a79907566",
					Manifest: Manifest{
						PK:      "repo/a@sha256:1f8d6e1edee77de035d79ca992df4e5cc8d358ec38f527077a84945a79907566",
						Payload: []byte("1"),
					},
				},
			},
			Expected: Tags{
				Tag{
					Digest:     "sha256:1f8d6e1edee77de035d79ca992df4e5cc8d358ec38f527077a84945a79907566",
					Repository: "repo/a",
					DigestPK:   "repo/a@sha256:1f8d6e1edee77de035d79ca992df4e5cc8d358ec38f527077a84945a79907566",
					Manifest: Manifest{
						PK:      "repo/a@sha256:1f8d6e1edee77de035d79ca992df4e5cc8d358ec38f527077a84945a79907566",
						Payload: []byte("1"),
					},
				},
			},
		},
		// Test 1: deleted tag, undeleted manifest
		{
			Input: Tags{
				Tag{
					Digest:     "sha256:1f8d6e1edee77de035d79ca992df4e5cc8d358ec38f527077a84945a79907566",
					Repository: "repo/a",
					DigestPK:   "repo/a@sha256:1f8d6e1edee77de035d79ca992df4e5cc8d358ec38f527077a84945a79907566",
					IsDeleted:  true,
					Manifest: Manifest{
						PK:        "repo/a@sha256:1f8d6e1edee77de035d79ca992df4e5cc8d358ec38f527077a84945a79907566",
						Payload:   []byte("1"),
						IsDeleted: false,
					},
				},
			},
			Expected: Tags{},
		},
		// Test 2: valid tag, deleted manifest
		{
			Input: Tags{
				Tag{
					Digest:     "sha256:1f8d6e1edee77de035d79ca992df4e5cc8d358ec38f527077a84945a79907566",
					Repository: "repo/a",
					DigestPK:   "repo/a@sha256:1f8d6e1edee77de035d79ca992df4e5cc8d358ec38f527077a84945a79907566",
					IsDeleted:  false,
					Manifest: Manifest{
						PK:        "repo/a@sha256:1f8d6e1edee77de035d79ca992df4e5cc8d358ec38f527077a84945a79907566",
						Payload:   []byte("1"),
						IsDeleted: true,
					},
				},
			},
			Expected: Tags{},
		},
		// Test 3: valid tag, invalid manifest
		{
			Input: Tags{
				Tag{
					Digest:     "repo/a@sha256:1f8d6e1edee77de035d79ca992df4e5cc8d358ec38f527077a84945a79907566",
					DigestPK:   "repo/a@sha256:1f8d6e1edee77de035d79ca992df4e5cc8d358ec38f527077a84945a79907566",
					IsDeleted:  false,
					Repository: "repo/a",
					Manifest:   Manifest{},
				},
			},
			Expected: Tags{},
		},
		// Test 4: valid tag with a digest pointing to another repo
		{
			Input: Tags{
				Tag{
					Digest:     "sha256:1f8d6e1edee77de035d79ca992df4e5cc8d358ec38f527077a84945a79907566",
					Repository: "repo/b",
					IsDeleted:  false,
					Manifest: Manifest{
						PK:        "LOL/a@sha256:1f8d6e1edee77de035d79ca992df4e5cc8d358ec38f527077a84945a79907566",
						Payload:   []byte("1"),
						IsDeleted: false,
					},
				},
			},
			Expected: Tags{},
		},
	}

	for i, test := range tests {
		actual := test.Input.FilterValid()
		if !reflect.DeepEqual(test.Expected, actual) {
			t.Fatalf("unexepected result from FilterValid test %d (expected):\n%#v\nvs. (got):\n%#v", i, test.Expected, actual)
		}
	}
}
