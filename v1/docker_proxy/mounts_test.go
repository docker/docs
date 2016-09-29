package proxy

import "testing"

func TestAdjustMount(t *testing.T) {
	var tests = []struct {
		spec     string
		expected string
	}{
		{"/c/Users/john/test:/test", "/c/Users/john/test:/test"},
		{`C:\Users\john:/test`, "/C/Users/john:/test"},
		{`D:\Users\john:/test`, "/D/Users/john:/test"},
		{`Z:\Users\john:/test`, "/Z/Users/john:/test"},
		{"//c/Users/john/test:/test", "/c/Users/john/test:/test"},
		{"///c/Users/john/test:/test", "/c/Users/john/test:/test"},
		{"/c/Users/john/test:/test:rw", "/c/Users/john/test:/test:rw"},
		{"/c/Users/john/test/../data:/data", "/c/Users/john/data:/data"},
		{"c:/data", "/c/data"},
		{"//c/:/data", "/c:/data"},
		{"/c:/data", "/c:/data"},
		{"/Volumes:/Volumes", "/Volumes:/Volumes"},
		{"/Volumes:/Volumes:r", "/Volumes:/Volumes:r"},
		{"", ""},
		{"/c/Users/john/data", "/c/Users/john/data"},
		{":/test", ":/test"},
		{"/d/Users/john/test:/test", "/d/Users/john/test:/test"},
		{"//d/Users/john/test:/test", "/d/Users/john/test:/test"},
		{"///d/Users/john/test:/test", "/d/Users/john/test:/test"},
	}

	for _, test := range tests {
		actual, _ := adjustMount(test.spec, "")

		if actual != test.expected {
			t.Errorf("Invalid mount point. Expected: %s. Got %s", test.expected, actual)
		}
	}
}
