package main

import "testing"

func TestParseNoProxy(t *testing.T) {
	original := "*.local, 169.254/16"
	expected := "*.local,169.254/16"

	result := parseNoProxy(original)
	if result != expected {
		t.Fatalf("Expected %s \n Got: %s \n", expected, result)
	}
}

func TestParseEmptyNoProxy(t *testing.T) {
	original := ""
	expected := ""

	result := parseNoProxy(original)
	if result != expected {
		t.Fatalf("Expected %s \n Got: %s \n", expected, result)
	}
}

func TestParseHttpProxy(t *testing.T) {
	original := "http://proxy.mycorp.com:3128"
	expected := "http://proxy.mycorp.com:3128"

	result, err := parseHTTPProxy(original)
	if err != nil {
		t.Fatal(err)
	}
	if result != expected {
		t.Fatalf("Expected %s \n Got: %s \n", expected, result)
	}
}

func TestParseEmptyHttpProxy(t *testing.T) {
	original := ""
	expected := ""

	result, err := parseHTTPProxy(original)
	if err != nil {
		t.Fatal(err)
	}
	if result != expected {
		t.Fatalf("Expected %s \n Got: %s \n", expected, result)
	}
}

func TestParseHttpProxyNoScheme(t *testing.T) {
	original := "proxy.mycorp.com:3128"
	expected := "http://proxy.mycorp.com:3128"

	result, err := parseHTTPProxy(original)
	if err != nil {
		t.Fatal(err)
	}
	if result != expected {
		t.Fatalf("Expected %s \n Got: %s \n", expected, result)
	}
}

func TestParseHttpProxyIPAddress(t *testing.T) {
	original := "192.168.65.2:3128"
	expected := "http://192.168.65.2:3128"

	result, err := parseHTTPProxy(original)
	if err != nil {
		t.Fatal(err)
	}
	if result != expected {
		t.Fatalf("Expected %s \n Got: %s \n", expected, result)
	}
}

func TestParseHttpProxyLocalhost(t *testing.T) {
	original := "http://127.0.0.1:3128"
	expected := ""

	result, err := parseHTTPProxy(original)
	if err == nil {
		t.Fatal("Expected error to be thrown")
	}
	if result != expected {
		t.Fatalf("Expected %s \n Got: %s \n", expected, result)
	}
}

func TestParseHttpProxyLocalhost2(t *testing.T) {
	original := "http://localhost:3128"
	expected := ""

	result, err := parseHTTPProxy(original)
	if err == nil {
		t.Fatal("Expected error to be thrown")
	}
	if result != expected {
		t.Fatalf("Expected %s \n Got: %s \n", expected, result)
	}
}
