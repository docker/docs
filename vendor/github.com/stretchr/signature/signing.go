package signature

import (
	"errors"
	"fmt"
	stewstrings "github.com/stretchr/stew/strings"
	"github.com/stretchr/tracer"
	"net/url"
	"regexp"
	"strings"
)

// FailedSignature is the string that will be used if signing fails.
const FailedSignature string = ":-("

// ErrNoSignatureFound is the error that is thrown when no signature could be found.
var ErrNoSignatureFound = errors.New("No signature was found.")

// SignatureRegex is the regex used to remove the signature from the URL string
var SignatureRegex = regexp.MustCompile("(.*)[&?](~|%7E)?sign=([0-9a-zA-Z]+)(.*)")

// trace writes some trace (if there is a Tracer set).
func trace(t *tracer.Tracer, format string, args ...interface{}) {

	if t.Should(tracer.LevelDebug) {

		// add the 'signature' prefix to trace
		if len(format) > 0 {
			format = stewstrings.MergeStrings("signature: ", format)
		}

		// trace this
		t.Trace(tracer.LevelDebug, format, args...)
	}

}

// GetSignature gets the signature of a request based on the given parameters.
func GetSignature(method, requestUrl, body, privateKey string) (string, error) {
	return GetSignatureWithTrace(method, requestUrl, body, privateKey, Tracer)
}

// GetSignatureWithTrace gets the signature of a request based on the given parameters.
func GetSignatureWithTrace(method, requestUrl, body, privateKey string, tracer *tracer.Tracer) (string, error) {

	trace(tracer, "GetSignature: method=%s", method)
	trace(tracer, "GetSignature: requestUrl=%s", requestUrl)
	trace(tracer, "GetSignature: body=%s", body)
	trace(tracer, "GetSignature: privateKey=%s", privateKey)

	// parse the URL
	u, parseErr := url.ParseRequestURI(requestUrl)

	if parseErr != nil {
		trace(tracer, "GetSignature: FAILED to parse the URL: %s", parseErr)
		return FailedSignature, parseErr
	}

	trace(tracer, "GetSignature: Parsed the URL as: %s", u.String())

	// get the query values
	values := u.Query()

	// add the private key parameter
	values.Set(PrivateKeyKey, privateKey)

	trace(tracer, "GetSignature: Set the private key (%s): %s", PrivateKeyKey, privateKey)

	if len(body) > 0 {
		bodyHash := Hash(body)
		trace(tracer, "GetSignature: Set the body hash (%s): %s", BodyHashKey, bodyHash)
		values.Set(BodyHashKey, bodyHash)
	} else {
		trace(tracer, "GetSignature: Skipping body hash as there's no body (%s).", BodyHashKey)
	}

	// get the ordered params
	orderedParams := OrderParams(values)

	trace(tracer, "GetSignature: Ordered parameters: %s", orderedParams)

	base := strings.Split(u.String(), "?")[0]
	combined := stewstrings.MergeStrings(strings.ToUpper(method), "&", base, "?", orderedParams)

	trace(tracer, "GetSignature: Base    : %s", base)
	trace(tracer, "GetSignature: Combined: %s", combined)

	theHash := Hash(combined)

	trace(tracer, "GetSignature: Output: %s", theHash)

	return theHash, nil

}

// GetSignedURL gets the URL with the sign parameter added based on the given parameters.
func GetSignedURL(method, requestUrl, body, privateKey string) (string, error) {
	return GetSignedURLWithTrace(method, requestUrl, body, privateKey, Tracer)
}

// GetSignedURL gets the URL with the sign parameter added based on the given parameters.
func GetSignedURLWithTrace(method, requestUrl, body, privateKey string, tracer *tracer.Tracer) (string, error) {

	trace(tracer, "GetSignedURL: method=%s", method)
	trace(tracer, "GetSignedURL: requestUrl=%s", requestUrl)
	trace(tracer, "GetSignedURL: body=%s", body)
	trace(tracer, "GetSignedURL: privateKey=%s", privateKey)

	hash, hashErr := GetSignatureWithTrace(method, requestUrl, body, privateKey, tracer)

	if hashErr != nil {
		trace(tracer, "GetSignedURL: FAILED to get the signature: %s", hashErr)
		return FailedSignature, hashErr
	}

	var signedUrl string
	if strings.Contains(requestUrl, "?") {
		signedUrl = stewstrings.MergeStrings(requestUrl, "&", url.QueryEscape(SignatureKey), "=", url.QueryEscape(hash))
	} else {
		signedUrl = stewstrings.MergeStrings(requestUrl, "?", url.QueryEscape(SignatureKey), "=", url.QueryEscape(hash))
	}

	trace(tracer, "GetSignedURL: Output: %s", signedUrl)

	return signedUrl, nil

}

// ValidateSignature validates the signature in a URL to ensure it is correct based on
// the specified parameters.
func ValidateSignature(method, requestUrl, body, privateKey string) (bool, error) {
	return ValidateSignatureWithTrace(method, requestUrl, body, privateKey, Tracer)
}

// ValidateSignature validates the signature in a URL to ensure it is correct based on
// the specified parameters.
func ValidateSignatureWithTrace(method, requestUrl, body, privateKey string, tracer *tracer.Tracer) (bool, error) {

	trace(tracer, "ValidateSignature: method=%s", method)
	trace(tracer, "ValidateSignature: requestUrl=%s", requestUrl)
	trace(tracer, "ValidateSignature: body=%s", body)
	trace(tracer, "ValidateSignature: privateKey=%s", privateKey)

	if !strings.Contains(requestUrl, "?") {
		trace(tracer, "ValidateSignature: FAILED because there was no signature found.")
		return false, ErrNoSignatureFound
	}

	segments := strings.Split(requestUrl, SignatureKey)

	if len(segments) < 2 {
		trace(tracer, "ValidateSignature: Failed to get signature: %s", ErrNoSignatureFound)
		return false, ErrNoSignatureFound
	} else {
		trace(tracer, "Segments: %s", segments)
	}

	modifiedURL := strings.TrimRight(segments[0], "&")
	signature := strings.TrimLeft(segments[1], "=")

	trace(tracer, "ValidateSignature: Modified URL (without signature): %s", modifiedURL)

	expectedSignature, signErr := GetSignatureWithTrace(method, modifiedURL, body, privateKey, tracer)

	if signErr != nil {
		trace(tracer, "ValidateSignature: FAILED to GetSignature: %s", signErr)
		return false, signErr
	}

	if signature != expectedSignature {
		err := errors.New(fmt.Sprintf("Signature \"%s\" is incorrect when \"%s\" is expected.", signature, expectedSignature))
		trace(tracer, "ValidateSignature: Signatures do not match: %s", err)
		return false, err
	}

	trace(tracer, "ValidateSignature: Happy because the signatures match: %s", signature)

	return true, nil

}
