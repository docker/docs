package errors

import (
	"fmt"
	"testing"
)

func TestHTTPError(t *testing.T) {
	err := HTTPError{
		400,
		1234,
		fmt.Errorf("Test error"),
	}
	errStr := err.Error()

	if errStr != "1234: Test error" {
		t.Fatalf("Error did not create expected string")
	}
}
