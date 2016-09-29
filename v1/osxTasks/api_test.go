package osxTasks

import (
	"encoding/json"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestAPIRestartVMRequest(t *testing.T) {
	r := `{"action":"restartvm"}`
	var res APIRestartVMRequest
	err := json.Unmarshal([]byte(r), &res)
	if err != nil {
		t.Fatal(err)
	}
	if res.Action != "restartvm" {
		t.Fatal("Incorrect decoding of request")
	}
}

func TestAPISetVMSettingsRequest(t *testing.T) {
	r := `{"action":"setvmsettings", "args": {"memory":2,"cpus":2,"daemonjson":"{\"foo\":\"bar\"}"}}`
	var res APISetVMSettingsRequest
	err := json.Unmarshal([]byte(r), &res)
	if err != nil {
		t.Fatal(err)
	}
	if res.Action != "setvmsettings" {
		t.Fatal("Incorrect decoding of request")
	}
	if res.Args.Memory != 2 {
		t.Fatalf("Expected Memory %d, Got %d.", 2, res.Args.Memory)
	}
	if res.Args.Cpus != 2 {
		t.Fatalf("Expected Cpus %d, Got %d.", 2, res.Args.Cpus)
	}
	if res.Args.DaemonJSON != `{"foo":"bar"}` {
		t.Fatalf("Expected DaemonJSONy %s, Got %s.", `{"foo":"bar"}`, res.Args.DaemonJSON)
	}
}

func TestAPIGetVMSettingsRequest(t *testing.T) {
	r := `{"action":"getvmsettings"}`
	var res APIGetVMSettingsRequest
	err := json.Unmarshal([]byte(r), &res)
	if err != nil {
		t.Fatal(err)
	}
	if res.Action != "getvmsettings" {
		t.Fatal("Incorrect decoding of request")
	}
}

func TestAPIGetSharedDirectoriesRequest(t *testing.T) {
	r := `{"action":"getshareddirectories"}`
	var res GetSharedDirectoriesRequest
	err := json.Unmarshal([]byte(r), &res)
	if err != nil {
		t.Fatal(err)
	}
	if res.Action != "getshareddirectories" {
		t.Fatal("Incorrect decoding of request")
	}
}

func TestAPISetSharedDirectoriesRequest(t *testing.T) {
	r := `{"action":"setshareddirectories", "args": { "directories" : [ "foo", "bar", "baz" ] }}`
	var res SetSharedDirectoriesRequest
	err := json.Unmarshal([]byte(r), &res)
	if err != nil {
		t.Fatal(err)
	}
	if res.Action != "setshareddirectories" {
		t.Fatal("Incorrect decoding of request")
	}
	if !reflect.DeepEqual(res.Args.Directories, []string{"foo", "bar", "baz"}) {
		t.Fatal("Directory list is not equal!")
	}
}

func TestAPIVMStateEventRequest(t *testing.T) {
	r := `{"action":"vmstateevent", "args" : { "vmstate": "foo" }}`
	var res VMStateEventRequest
	err := json.Unmarshal([]byte(r), &res)
	if err != nil {
		t.Fatal(err)
	}
	if res.Action != "vmstateevent" {
		t.Fatal("Incorrect decoding of request")
	}
	if res.Args.VMState != "foo" {
		t.Fatalf("Expected VMState %s, Got %s", "foo", res.Args.VMState)
	}
}

func TestErrorResponse(t *testing.T) {
	r := NewErrorResponse("test error")

	s, err := json.Marshal(r)
	if err != nil {
		t.Fatal(err)
	}

	expected := `{"status":"error","message":"test error"}`
	if string(s) != expected {
		t.Fatalf("Expected: %s, Got: %s", expected, string(s))
	}
}

func TestNoErrorResponse(t *testing.T) {
	r := NoError
	s, err := json.Marshal(r)
	if err != nil {
		t.Fatal(err)
	}

	expected := `{"status":"ok"}`
	if string(s) != expected {
		t.Fatalf("Expected: %s, Got: %s", expected, string(s))
	}
}

func TestGetVMSettingsResponse(t *testing.T) {
	r := GetVMSettingsResponse{
		Status: statusOK,
		VMSettings: &VMSettings{
			Memory:               2048,
			Cpus:                 2,
			DaemonJSON:           `{"foo":"bar"}`,
			SystemProxyHTTP:      "http://example.com:3128",
			SystemProxyHTTPS:     "http://example.com:3129",
			SystemProxyExclude:   "169.* , *.local",
			OverrideProxyHTTP:    "http://mycorp.com:3128",
			OverrideProxyHTTPS:   "http://mycorp.com:3129",
			OverrideProxyExclude: "*.mycorp.com",
		},
	}
	s, err := json.Marshal(r)
	if err != nil {
		t.Fatal(err)
	}

	expected := `{"status":"ok","memory":2048,"cpus":2,"daemonjson":"{\"foo\":\"bar\"}","systemProxyHttp":"http://example.com:3128","systemProxyHttps":"http://example.com:3129","systemProxyExclude":"169.* , *.local","overrideProxyHttp":"http://mycorp.com:3128","overrideProxyHttps":"http://mycorp.com:3129","overrideProxyExclude":"*.mycorp.com"}`

	if string(s) != expected {
		t.Fatalf("Expected: \n %s \n Got: \n %s \n", expected, string(s))
	}
}

func TestGetSharedDirectoriesResponse(t *testing.T) {
	directories := []string{"/Users/docker", "/Applications", "/tmp"}
	r := SharedDirectoriesResponse{
		Status:      statusOK,
		Directories: directories,
	}

	s, err := json.Marshal(r)
	if err != nil {
		t.Fatal(err)
	}

	expected := `{"status":"ok","directories":["/Users/docker","/Applications","/tmp"]}`
	if string(s) != expected {
		t.Fatalf("Expected: \n %s \n Got: \n %s \n", expected, string(s))
	}

}

func TestDecodeTwice(t *testing.T) {
	r := `{"action":"vmstateevent", "args" : { "vmstate": "foo" }}`

	var api APIRequest
	err := json.Unmarshal([]byte(r), &api)
	if err != nil {
		t.Fatal(err)
	}

	var res VMStateEventRequest
	err = json.Unmarshal([]byte(r), &res)
	if err != nil {
		t.Fatal(err)
	}
	if res.Action != "vmstateevent" {
		t.Fatal("Incorrect decoding of request")
	}
	if res.Args.VMState != "foo" {
		t.Fatalf("Expected VMState %s, Got %s", "foo", res.Args.VMState)
	}

}

func TestWriteInterface(t *testing.T) {
	r := NewErrorResponse("test error")

	_, err := json.Marshal(r)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	r.Write(w)

	expected := `{"status":"error","message":"test error"}`
	if expected != w.Body.String() {
		t.Fatalf("Expected: %s, \n Got: %s\n", expected, w.Body.String())
	}
}
