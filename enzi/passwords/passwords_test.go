package passwords

import (
	"crypto/rand"
	"io"
	"testing"
)

func benchPasswordHasher(b *testing.B, hasher func(string) (string, error)) {
	buf := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, buf); err != nil {
		b.Fatal(err)
	}

	password := string(buf)

	b.ResetTimer()
	for i := 0; i <= b.N; i++ {
		if _, err := hasher(password); err != nil {
			b.Fatalf("unable to hash password: %s", err)
		}
	}
}

func BenchmarkHashPasswordBcrypt(b *testing.B) {
	benchPasswordHasher(b, hashPasswordBcrypt)
}

func BenchmarkHashPasswordPbkdf2Sha256(b *testing.B) {
	benchPasswordHasher(b, hashPasswordPbkdf2Sha256)
}
