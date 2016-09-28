package passwords

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
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/pbkdf2"
)

const (
	// I've adjusted this to take approximately 100ms to perform a password
	// check using a single cpu core on a 2011 model MacBookPro11,2.
	pbkdf2WorkFactor  = 54 * 1024
	pbkdf2Sha256Label = "pbkdf2_sha256"
	bcryptCost        = bcrypt.DefaultCost
	bcryptLabel       = "bcrypt"
)

// HashPassword generates a hash of the given password with the default
// password hashing algorithm.
func HashPassword(password string) (hashedPassword string, err error) {
	hashedPassword, err = hashPasswordBcrypt(password)
	if err != nil {
		return "", fmt.Errorf("bcrypt: %s", err)
	}

	// Prepend the algorithm label.
	return fmt.Sprintf("%s$%s", bcryptLabel, hashedPassword), nil
}

// CheckPassword checks that the given password matches the hashed password
// when hashed using the same method. Returns true iff the password matches.
func CheckPassword(ctx context.Context, hashedPassword, password string) bool {
	parts := strings.SplitN(hashedPassword, "$", 2)
	if len(parts) != 2 {
		context.GetLogger(ctx).Warn("invalid hashed password encoding")
		return false
	}

	algLabel, encodedHash := parts[0], parts[1]

	switch algLabel {
	case pbkdf2Sha256Label:
		return checkPasswordPbkdf2Sha256(ctx, encodedHash, password)
	case bcryptLabel:
		return checkPasswordBcrypt(ctx, encodedHash, password)
	default:
		context.GetLogger(ctx).Warnf("unknown algorithm: %s", algLabel)
		return false
	}
}

// CheckPasswordHash checks that the given password hash is encoded in a way
// that could possibly match a password.
func CheckPasswordHash(passwordHash string) error {
	parts := strings.SplitN(passwordHash, "$", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid hashed password encoding")
	}

	algLabel, encodedHash := parts[0], parts[1]
	switch algLabel {
	case pbkdf2Sha256Label:
		parts := strings.SplitN(encodedHash, "$", 3)
		if len(parts) != 3 {
			return fmt.Errorf("invalid %s encoding", pbkdf2Sha256Label)
		}
	case bcryptLabel:
		compareToNilRes := bcrypt.CompareHashAndPassword([]byte(encodedHash), nil)
		if compareToNilRes == nil {
			return fmt.Errorf("%s password hash is of an empty password", bcryptLabel)
		} else if compareToNilRes != bcrypt.ErrMismatchedHashAndPassword {
			return fmt.Errorf("invalid %s encoding: %s", bcryptLabel, compareToNilRes)
		}
	default:
		return fmt.Errorf("unknown algorithm: %s", algLabel)
	}
	return nil
}

// checkPasswordPbkdf2Sha256 checks the password using pbkdf2 and sha256.
func checkPasswordPbkdf2Sha256(ctx context.Context, encodedHash, password string) bool {
	parts := strings.SplitN(encodedHash, "$", 3)
	if len(parts) != 3 {
		context.GetLogger(ctx).Warnf("invalid %s encoding", pbkdf2Sha256Label)
		return false
	}

	iterString, salt, hashedB64 := parts[0], parts[1], parts[2]

	iterations, err := strconv.Atoi(iterString)
	if err != nil {
		context.GetLogger(ctx).Warnf("unable to parse %s iterations", pbkdf2Sha256Label)
		return false
	}

	givenHashedBytes := pbkdf2.Key([]byte(password), []byte(salt), iterations, sha256.Size, sha256.New)
	givenHashedB64 := base64.StdEncoding.EncodeToString(givenHashedBytes)

	return subtle.ConstantTimeCompare([]byte(hashedB64), []byte(givenHashedB64)) == 1
}

// hashPasswordPbkdf2Sha256 hashes the given password using PBKDF2 and SHA256.
// The salt is generated using 128 random bits crypto/rand.Reader and a 64K
// rounds work factor.
func hashPasswordPbkdf2Sha256(password string) (encodedHash string, err error) {
	buf := make([]byte, 16) // 128 bit salt buffer.
	if _, err = io.ReadFull(rand.Reader, buf); err != nil {
		return "", fmt.Errorf("unable to generate random salt: %s", err)
	}

	salt := base64.StdEncoding.EncodeToString(buf)
	hashedBytes := pbkdf2.Key([]byte(password), []byte(salt), pbkdf2WorkFactor, sha256.Size, sha256.New)
	hashedB64 := base64.StdEncoding.EncodeToString(hashedBytes)

	return fmt.Sprintf("%s$%s$%s", strconv.Itoa(pbkdf2WorkFactor), salt, hashedB64), nil
}

// checkPasswordBcrypt checks the password using bcrypt.
func checkPasswordBcrypt(_ context.Context, encodedHash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(encodedHash), []byte(password)) == nil
}

// hashPasswordBcrypt hashes the given password using Bcrypt with the default
// cost factor.
func hashPasswordBcrypt(password string) (encodedHash string, err error) {
	encodedHashBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return "", fmt.Errorf("unable to generate hash: %s", err)
	}

	return string(encodedHashBytes), nil
}
