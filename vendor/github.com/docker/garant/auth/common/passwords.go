package common

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/docker/distribution/context"
	"golang.org/x/crypto/pbkdf2"
)

// I've adjusted this to take approximately 100ms to perform a password check
// using a single cpu core on a 2011 model MacBookPro11,2.
const pbkdf2WorkFactor = 54 * 1024

// CheckPasswordPbkdf2Sha256 checks the password using pbkdf2 and sha256.
func CheckPasswordPbkdf2Sha256(ctx context.Context, encodedHash, password string) bool {
	logger := context.GetLogger(ctx)
	parts := strings.SplitN(encodedHash, "$", 3)

	if len(parts) != 3 {
		logger.Debugf("encoded hashed password does not have '$' separated parts")
		return false
	}

	iterString, salt, hashedB64 := parts[0], parts[1], parts[2]

	iterations, err := strconv.Atoi(iterString)
	if err != nil {
		logger.Debugf("unable to convert number of iterations: %s", err)
		return false
	}

	givenHashedBytes := pbkdf2.Key([]byte(password), []byte(salt), iterations, sha256.Size, sha256.New)
	givenHashedB64 := base64.StdEncoding.EncodeToString(givenHashedBytes)

	return subtle.ConstantTimeCompare([]byte(hashedB64), []byte(givenHashedB64)) == 1
}

// HashPasswordPbkdf2Sha256 hashes the given password using PBKDF2 and SHA256.
// The salt is generated using 128 random bits crypto/rand.Reader and a 64K
// rounds work factor.
func HashPasswordPbkdf2Sha256(password string) (encodedHash string, err error) {
	buf := make([]byte, 16) // 128 bit salt buffer.
	if _, err = io.ReadFull(rand.Reader, buf); err != nil {
		return "", WithStackTrace(err)
	}

	salt := base64.StdEncoding.EncodeToString(buf)
	hashedBytes := pbkdf2.Key([]byte(password), []byte(salt), pbkdf2WorkFactor, sha256.Size, sha256.New)
	hashedB64 := base64.StdEncoding.EncodeToString(hashedBytes)

	return fmt.Sprintf("%s$%s$%s", strconv.Itoa(pbkdf2WorkFactor), salt, hashedB64), nil
}
