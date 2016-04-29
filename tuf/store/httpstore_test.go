package store

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/docker/go/canonical/json"
	"github.com/docker/notary/tuf/data"
	"github.com/docker/notary/tuf/signed"
	"github.com/docker/notary/tuf/validation"
	"github.com/stretchr/testify/require"
)

const testRoot = `{"signed":{"_type":"Root","consistent_snapshot":false,"expires":"2025-07-17T16:19:21.101698314-07:00","keys":{"1ca15c7f4b2b0c6efce202a545e7267152da28ab7c91590b3b60bdb4da723aad":{"keytype":"ecdsa","keyval":{"private":null,"public":"MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEb0720c99Cj6ZmuDlznEZ52NA6YpeY9Sj45z51XvPnG63Bi2RSBezMJlPzbSfP39mXKXqOJyT+z9BZhi3FVWczg=="}},"b1d6813b55442ecbfb1f4b40eb1fcdb4290e53434cfc9ba2da24c26c9143873b":{"keytype":"ecdsa-x509","keyval":{"private":null,"public":"LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUJVekNCKzZBREFnRUNBaEFCWDNKLzkzaW8zbHcrZUsvNFhvSHhNQW9HQ0NxR1NNNDlCQU1DTUJFeER6QU4KQmdOVkJBTVRCbVY0Y0dseVpUQWVGdzB4TlRBM01qQXlNekU1TVRkYUZ3MHlOVEEzTVRjeU16RTVNVGRhTUJFeApEekFOQmdOVkJBTVRCbVY0Y0dseVpUQlpNQk1HQnlxR1NNNDlBZ0VHQ0NxR1NNNDlBd0VIQTBJQUJFTDhOTFhQCitreUJZYzhYY0FTMXB2S2l5MXRQUDlCZHJ1dEdrWlR3Z0dEYTM1THMzSUFXaWlrUmlPbGRuWmxVVEE5cG5JekoKOFlRQThhTjQ1TDQvUlplak5UQXpNQTRHQTFVZER3RUIvd1FFQXdJQW9EQVRCZ05WSFNVRUREQUtCZ2dyQmdFRgpCUWNEQXpBTUJnTlZIUk1CQWY4RUFqQUFNQW9HQ0NxR1NNNDlCQU1DQTBjQU1FUUNJRVJ1ZUVURG5xMlRqRFBmClhGRStqUFJqMEtqdXdEOG9HSmtoVGpMUDAycjhBaUI5cUNyL2ZqSXpJZ1NQcTJVSXZqR0hlYmZOYXh1QlpZZUUKYW8xNjd6dHNYZz09Ci0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0K"}},"fbddae7f25a6c23ca735b017206a849d4c89304a4d8de4dcc4b3d6f3eb22ce3b":{"keytype":"ecdsa","keyval":{"private":null,"public":"MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE/xS5fBHK2HKmlGcvAr06vwPITvmxWP4P3CMDCgY25iSaIiM21OiXA1/Uvo3Pa3xh5G3cwCtDvi+4FpflW2iB/w=="}},"fd75751f010c3442e23b3e3e99a1442a112f2f21038603cb8609d8b17c9e912a":{"keytype":"ed25519","keyval":{"private":null,"public":"rc+glN01m+q8jmX8SolGsjTfk6NMhUQTWyj10hjmne0="}}},"roles":{"root":{"keyids":["b1d6813b55442ecbfb1f4b40eb1fcdb4290e53434cfc9ba2da24c26c9143873b"],"threshold":1},"snapshot":{"keyids":["1ca15c7f4b2b0c6efce202a545e7267152da28ab7c91590b3b60bdb4da723aad"],"threshold":1},"targets":{"keyids":["fbddae7f25a6c23ca735b017206a849d4c89304a4d8de4dcc4b3d6f3eb22ce3b"],"threshold":1},"timestamp":{"keyids":["fd75751f010c3442e23b3e3e99a1442a112f2f21038603cb8609d8b17c9e912a"],"threshold":1}},"version":2},"signatures":[{"keyid":"b1d6813b55442ecbfb1f4b40eb1fcdb4290e53434cfc9ba2da24c26c9143873b","method":"ecdsa","sig":"A2lNVwxHBnD9ViFtRre8r5oG6VvcvJnC6gdvvxv/Jyag40q/fNMjllCqyHrb+6z8XDZcrTTDsFU1R3/e+92d1A=="}]}`

const testRootKey = "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUJVekNCKzZBREFnRUNBaEFCWDNKLzkzaW8zbHcrZUsvNFhvSHhNQW9HQ0NxR1NNNDlCQU1DTUJFeER6QU4KQmdOVkJBTVRCbVY0Y0dseVpUQWVGdzB4TlRBM01qQXlNekU1TVRkYUZ3MHlOVEEzTVRjeU16RTVNVGRhTUJFeApEekFOQmdOVkJBTVRCbVY0Y0dseVpUQlpNQk1HQnlxR1NNNDlBZ0VHQ0NxR1NNNDlBd0VIQTBJQUJFTDhOTFhQCitreUJZYzhYY0FTMXB2S2l5MXRQUDlCZHJ1dEdrWlR3Z0dEYTM1THMzSUFXaWlrUmlPbGRuWmxVVEE5cG5JekoKOFlRQThhTjQ1TDQvUlplak5UQXpNQTRHQTFVZER3RUIvd1FFQXdJQW9EQVRCZ05WSFNVRUREQUtCZ2dyQmdFRgpCUWNEQXpBTUJnTlZIUk1CQWY4RUFqQUFNQW9HQ0NxR1NNNDlCQU1DQTBjQU1FUUNJRVJ1ZUVURG5xMlRqRFBmClhGRStqUFJqMEtqdXdEOG9HSmtoVGpMUDAycjhBaUI5cUNyL2ZqSXpJZ1NQcTJVSXZqR0hlYmZOYXh1QlpZZUUKYW8xNjd6dHNYZz09Ci0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0K"

type TestRoundTripper struct{}

func (rt *TestRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return http.DefaultClient.Do(req)
}

func TestHTTPStoreGetMeta(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(testRoot))
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()
	store, err := NewHTTPStore(
		server.URL,
		"metadata",
		"txt",
		"key",
		&http.Transport{},
	)
	if err != nil {
		t.Fatal(err)
	}
	j, err := store.GetMeta("root", 4801)
	if err != nil {
		t.Fatal(err)
	}
	p := &data.Signed{}
	err = json.Unmarshal(j, p)
	if err != nil {
		t.Fatal(err)
	}
	rootKey, err := base64.StdEncoding.DecodeString(testRootKey)
	require.NoError(t, err)
	k := data.NewPublicKey("ecdsa-x509", rootKey)

	sigBytes := p.Signatures[0].Signature
	if err != nil {
		t.Fatal(err)
	}
	var decoded map[string]interface{}
	if err := json.Unmarshal(*p.Signed, &decoded); err != nil {
		t.Fatal(err)
	}
	msg, err := json.MarshalCanonical(decoded)
	if err != nil {
		t.Fatal(err)
	}
	method := p.Signatures[0].Method
	err = signed.Verifiers[method].Verify(k, sigBytes, msg)
	if err != nil {
		t.Fatal(err)
	}

}

// Test that passing -1 to httpstore's GetMeta will return all content
func TestHTTPStoreGetAllMeta(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(testRoot))
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()
	store, err := NewHTTPStore(
		server.URL,
		"metadata",
		"txt",
		"key",
		&http.Transport{},
	)
	if err != nil {
		t.Fatal(err)
	}
	j, err := store.GetMeta("root", NoSizeLimit)
	if err != nil {
		t.Fatal(err)
	}
	p := &data.Signed{}
	err = json.Unmarshal(j, p)
	if err != nil {
		t.Fatal(err)
	}
	rootKey, err := base64.StdEncoding.DecodeString(testRootKey)
	require.NoError(t, err)
	k := data.NewPublicKey("ecdsa-x509", rootKey)

	sigBytes := p.Signatures[0].Signature
	if err != nil {
		t.Fatal(err)
	}
	var decoded map[string]interface{}
	if err := json.Unmarshal(*p.Signed, &decoded); err != nil {
		t.Fatal(err)
	}
	msg, err := json.MarshalCanonical(decoded)
	if err != nil {
		t.Fatal(err)
	}
	method := p.Signatures[0].Method
	err = signed.Verifiers[method].Verify(k, sigBytes, msg)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSetMultiMeta(t *testing.T) {
	metas := map[string][]byte{
		"root":    []byte("root data"),
		"targets": []byte("targets data"),
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		reader, err := r.MultipartReader()
		if err != nil {
			t.Fatal(err)
		}
		updates := make(map[string][]byte)
		for {
			part, err := reader.NextPart()
			if err == io.EOF {
				break
			}
			role := strings.TrimSuffix(part.FileName(), ".json")
			updates[role], err = ioutil.ReadAll(part)
			if err != nil {
				t.Fatal(err)
			}
		}
		rd, rok := updates["root"]
		require.True(t, rok)
		require.Equal(t, rd, metas["root"])

		td, tok := updates["targets"]
		require.True(t, tok)
		require.Equal(t, td, metas["targets"])

	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()
	store, err := NewHTTPStore(server.URL, "metadata", "json", "key", http.DefaultTransport)
	if err != nil {
		t.Fatal(err)
	}

	store.SetMultiMeta(metas)
}

func TestPyCryptoRSAPSSCompat(t *testing.T) {
	pubPem := "-----BEGIN PUBLIC KEY-----\nMIIBojANBgkqhkiG9w0BAQEFAAOCAY8AMIIBigKCAYEAnKuXZeefa2LmgxaL5NsM\nzKOHNe+x/nL6ik+lDBCTV6OdcwAhHQS+PONGhrChIUVR6Vth3hUCrreLzPO73Oo5\nVSCuRJ53UronENl6lsa5mFKP8StYLvIDITNvkoT3j52BJIjyNUK9UKY9As2TNqDf\nBEPIRp28ev/NViwGOEkBu2UAbwCIdnDXm8JQErCZA0Ydm7PKGgjLbFsFGrVzqXHK\n6pdzJXlhr9yap3UpgQ/iO9JtoEYB2EXsnSrPc9JRjR30bNHHtnVql3fvinXrAEwq\n3xmN4p+R4VGzfdQN+8Kl/IPjqWB535twhFYEG/B7Ze8IwbygBjK3co/KnOPqMUrM\nBI8ztvPiogz+MvXb8WvarZ6TMTh8ifZI96r7zzqyzjR1hJulEy3IsMGvz8XS2J0X\n7sXoaqszEtXdq5ef5zKVxkiyIQZcbPgmpHLq4MgfdryuVVc/RPASoRIXG4lKaTJj\n1ANMFPxDQpHudCLxwCzjCb+sVa20HBRPTnzo8LSZkI6jAgMBAAE=\n-----END PUBLIC KEY-----"
	//privPem := "-----BEGIN RSA PRIVATE KEY-----\n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\n-----END RSA PRIVATE KEY-----"
	testStr := "The quick brown fox jumps over the lazy dog."
	sigHex := "4e05ee9e435653549ac4eddbc43e1a6868636e8ea6dbec2564435afcb0de47e0824cddbd88776ddb20728c53ecc90b5d543d5c37575fda8bd0317025fc07de62ee8084b1a75203b1a23d1ef4ac285da3d1fc63317d5b2cf1aafa3e522acedd366ccd5fe4a7f02a42922237426ca3dc154c57408638b9bfaf0d0213855d4e9ee621db204151bcb13d4dbb18f930ec601469c992c84b14e9e0b6f91ac9517bb3b749dd117e1cbac2e4acb0e549f44558a2005898a226d5b6c8b9291d7abae0d9e0a16858b89662a085f74a202deb867acab792bdbd2c36731217caea8b17bd210c29b890472f11e5afdd1dd7b69004db070e04201778f2c49f5758643881403d45a58d08f51b5c63910c6185892f0b590f191d760b669eff2464456f130239bba94acf54a0cb98f6939ff84ae26a37f9b890be259d9b5d636f6eb367b53e895227d7d79a3a88afd6d28c198ee80f6527437c5fbf63accb81709925c4e03d1c9eaee86f58e4bd1c669d6af042dbd412de0d13b98b1111e2fadbe34b45de52125e9a"
	k := data.NewPublicKey(data.RSAKey, []byte(pubPem))

	sigBytes, err := hex.DecodeString(sigHex)
	if err != nil {
		t.Fatal(err)
	}
	v := signed.RSAPyCryptoVerifier{}
	err = v.Verify(k, sigBytes, []byte(testStr))
	if err != nil {
		t.Fatal(err)
	}
}

func TestPyNaCled25519Compat(t *testing.T) {
	pubHex := "846612b43cef909a0e4ea9c818379bca4723a2020619f95e7a0ccc6f0850b7dc"
	//privHex := "bf3cdb9b2a664b0460e6755cb689ffca15b6e294f79f9f1fcf90b52e5b063a76"
	testStr := "The quick brown fox jumps over the lazy dog."
	sigHex := "166e7013e48f26dccb4e68fe4cf558d1cd3af902f8395534336a7f8b4c56588694aa3ac671767246298a59d5ef4224f02c854f41bfcfe70241db4be1546d6a00"

	pub, _ := hex.DecodeString(pubHex)
	k := data.NewPublicKey(data.ED25519Key, pub)

	sigBytes, _ := hex.DecodeString(sigHex)

	err := signed.Verifiers[data.EDDSASignature].Verify(k, sigBytes, []byte(testStr))
	if err != nil {
		t.Fatal(err)
	}
}

func testErrorCode(t *testing.T, errorCode int, errType error) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(errorCode)
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	store, err := NewHTTPStore(
		server.URL,
		"metadata",
		"txt",
		"key",
		&http.Transport{},
	)
	require.NoError(t, err)

	_, err = store.GetMeta("root", 4801)
	require.Error(t, err)
	require.IsType(t, errType, err,
		fmt.Sprintf("%d should translate to %v", errorCode, errType))
}

func Test404Error(t *testing.T) {
	testErrorCode(t, http.StatusNotFound, ErrMetaNotFound{})
}

func Test50XErrors(t *testing.T) {
	fiveHundreds := []int{
		http.StatusInternalServerError,
		http.StatusNotImplemented,
		http.StatusBadGateway,
		http.StatusServiceUnavailable,
		http.StatusGatewayTimeout,
		http.StatusHTTPVersionNotSupported,
	}
	for _, code := range fiveHundreds {
		testErrorCode(t, code, ErrServerUnavailable{})
	}
}

func Test400Error(t *testing.T) {
	testErrorCode(t, http.StatusBadRequest, ErrInvalidOperation{})
}

// If it's a 400, translateStatusToError attempts to parse the body into
// an error.  If successful (and a recognized error) that error is returned.
func TestTranslateErrorsParse400Errors(t *testing.T) {
	origErr := validation.ErrBadRoot{Msg: "bad"}

	serialObj, err := validation.NewSerializableError(origErr)
	require.NoError(t, err)
	serialization, err := json.Marshal(serialObj)
	require.NoError(t, err)
	errorBody := bytes.NewBuffer([]byte(fmt.Sprintf(
		`{"errors": [{"otherstuff": "what", "detail": %s}]}`,
		string(serialization))))
	errorResp := http.Response{
		StatusCode: http.StatusBadRequest,
		Body:       ioutil.NopCloser(errorBody),
	}

	finalError := translateStatusToError(&errorResp, "")
	require.Equal(t, origErr, finalError)
}

// If it's a 400, translateStatusToError attempts to parse the body into
// an error.  If parsing fails, an InvalidOperation is returned instead.
func TestTranslateErrorsWhenCannotParse400(t *testing.T) {
	invalids := []string{
		`{"errors": [{"otherstuff": "what", "detail": {"Name": "Muffin"}}]}`,
		`{"errors": [{"otherstuff": "what", "detail": {}}]}`,
		`{"errors": [{"otherstuff": "what"}]}`,
		`{"errors": []}`,
		`{}`,
		"400",
	}
	for _, body := range invalids {
		errorResp := http.Response{
			StatusCode: http.StatusBadRequest,
			Body:       ioutil.NopCloser(bytes.NewBuffer([]byte(body))),
		}

		err := translateStatusToError(&errorResp, "")
		require.IsType(t, ErrInvalidOperation{}, err)
	}
}

func TestHTTPStoreRemoveAll(t *testing.T) {
	// Set up a simple handler and server for our store
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(testRoot))
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()
	store, err := NewHTTPStore(server.URL, "metadata", "json", "key", http.DefaultTransport)
	require.NoError(t, err)

	// currently unsupported since there is no use case
	// check for the error
	err = store.RemoveAll()
	require.Error(t, err)
}

func TestHTTPOffline(t *testing.T) {
	s, err := NewHTTPStore("https://localhost/", "", "", "", nil)
	require.NoError(t, err)
	require.IsType(t, &OfflineStore{}, s)
}

func TestErrServerUnavailable(t *testing.T) {
	for i := 200; i < 600; i++ {
		err := ErrServerUnavailable{code: i}
		if i == 401 {
			require.Contains(t, err.Error(), "not authorized")
		} else {
			require.Contains(t, err.Error(), "unable to reach trust server")
		}
	}
}
