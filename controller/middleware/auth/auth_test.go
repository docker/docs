package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/docker/orca/controller/mock_test"
)

var testHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("testing"))
})

func getManager() *mock_test.MockManager {
	return &mock_test.MockManager{}
}

func TestNoAuthToken(t *testing.T) {
	require := require.New(t)
	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)

	mgr := getManager()
	a := NewAuthRequired(mgr, []string{})
	_, err := a.Initializer(res, req)
	require.NotNil(err)
	require.Equal(res.Code, http.StatusUnauthorized)
}

func TestWhiteListAny(t *testing.T) {
	require := require.New(t)
	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)

	mgr := getManager()
	a := NewAuthRequired(mgr, []string{})
	_, err := a.Initializer(res, req)
	require.NotNil(t, err)
	require.Equal(res.Code, http.StatusUnauthorized)

	res = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/", nil)

	a = NewAuthRequired(nil, []string{"0.0.0.0/0"})
	_, err = a.Initializer(res, req)
	require.Nil(err)
	require.Equal(res.Code, http.StatusOK)
}

func TestWhiteListInvalid(t *testing.T) {
	require := require.New(t)
	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)

	mgr := getManager()
	a := NewAuthRequired(mgr, []string{"1.2.3.4/32"})
	_, err := a.Initializer(res, req)
	require.NotNil(err)
	require.Equal(res.Code, http.StatusUnauthorized)
}
