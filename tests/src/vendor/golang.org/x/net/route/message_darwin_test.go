// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package route

import "testing"

func TestFetchAndParseRIBOnDarwin(t *testing.T) {
	for _, af := range []int{sysAF_UNSPEC, sysAF_INET, sysAF_INET6} {
		for _, typ := range []RIBType{sysNET_RT_FLAGS, sysNET_RT_DUMP2, sysNET_RT_IFLIST2} {
			ms, err := fetchAndParseRIB(af, typ)
			if err != nil {
				t.Error(err)
				continue
			}
			ss, err := msgs(ms).validate()
			if err != nil {
				t.Errorf("%v %d %v", addrFamily(af), typ, err)
				continue
			}
			for _, s := range ss {
				t.Log(s)
			}
		}
	}
}
