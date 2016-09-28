package jose

import (
	"bufio"
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	crand "crypto/rand"
	"crypto/rsa"
	"encoding/binary"
	"io"
	"math/rand"
	"testing"
)

func makeRandomData(tb testing.TB, size int) []byte {
	seed, err := binary.ReadVarint(bufio.NewReader(crand.Reader))
	if err != nil {
		tb.Fatalf("unable to read random seed: %s", err)
	}

	data := make([]byte, 1024)
	if _, err := io.ReadFull(rand.New(rand.NewSource(seed)), data); err != nil {
		tb.Fatalf("unable to create random data buffer: %s", err)
	}

	return data
}

func generateKey(tb testing.TB, cryptoKey crypto.PrivateKey) *PrivateKey {
	key, err := NewPrivateKey(cryptoKey)
	if err != nil {
		tb.Fatalf("unable to convert crypto key to JWK private key: %s", err)
	}

	return key
}

func generateRSAKey(tb testing.TB, bits int) *PrivateKey {
	// This can be slow. Don't put this within an recording benchmarks.
	cryptoKey, err := rsa.GenerateKey(crand.Reader, bits)
	if err != nil {
		tb.Fatalf("unable to generate %d-bit RSA key: %s", bits, err)
	}

	return generateKey(tb, cryptoKey)
}

func generateECDSAKey(tb testing.TB, curve elliptic.Curve) *PrivateKey {
	cryptoKey, err := ecdsa.GenerateKey(curve, crand.Reader)
	if err != nil {
		tb.Fatalf("unable to generate %d-bit ECDSA key: %s", curve.Params().BitSize, err)
	}

	return generateKey(tb, cryptoKey)
}

func benchmarkSigner(b *testing.B, key *PrivateKey, alg SignatureAlgorithm) {
	signer, err := key.Signer(alg.Name())
	if err != nil {
		b.Fatalf("unable to get signer: %s", err)
	}

	data := makeRandomData(b, 1024)

	b.ResetTimer()
	for i := 0; i <= b.N; i++ {
		if _, err := signer.Sign(bytes.NewReader(data)); err != nil {
			b.Fatalf("unable to sign data: %s", err)
		}
	}
}

func BenchmarkSignerRSA2048SHA256(b *testing.B) {
	benchmarkSigner(b, generateRSAKey(b, 2048), RS256())
}

func BenchmarkSignerRSA2048SHA384(b *testing.B) {
	benchmarkSigner(b, generateRSAKey(b, 2048), RS384())
}

func BenchmarkSignerRSA2048SHA512(b *testing.B) {
	benchmarkSigner(b, generateRSAKey(b, 2048), RS512())
}

func BenchmarkSignerRSA3072SHA256(b *testing.B) {
	benchmarkSigner(b, generateRSAKey(b, 3072), RS256())
}

func BenchmarkSignerRSA3072SHA384(b *testing.B) {
	benchmarkSigner(b, generateRSAKey(b, 3072), RS384())
}

func BenchmarkSignerRSA3072SHA512(b *testing.B) {
	benchmarkSigner(b, generateRSAKey(b, 3072), RS512())
}

func BenchmarkSignerRSA4096SHA256(b *testing.B) {
	benchmarkSigner(b, generateRSAKey(b, 4096), RS256())
}

func BenchmarkSignerRSA4096SHA384(b *testing.B) {
	benchmarkSigner(b, generateRSAKey(b, 4096), RS384())
}

func BenchmarkSignerRSA4096SHA512(b *testing.B) {
	benchmarkSigner(b, generateRSAKey(b, 4096), RS512())
}

func BenchmarkSignerECDSAP256SHA256(b *testing.B) {
	benchmarkSigner(b, generateECDSAKey(b, elliptic.P256()), ES256())
}

func BenchmarkSignerECDSAP384SHA384(b *testing.B) {
	benchmarkSigner(b, generateECDSAKey(b, elliptic.P384()), ES384())
}

func BenchmarkSignerECDSAP521SHA512(b *testing.B) {
	benchmarkSigner(b, generateECDSAKey(b, elliptic.P521()), ES512())
}

func benchmarkVerifier(b *testing.B, key *PrivateKey, alg SignatureAlgorithm) {
	signer, err := key.Signer(alg.Name())
	if err != nil {
		b.Fatalf("unable to get signer: %s", err)
	}

	data := makeRandomData(b, 1024)

	signature, err := signer.Sign(bytes.NewReader(data))
	if err != nil {
		b.Fatalf("unable to sign random data: %s", err)
	}

	verifier, err := key.Verifier(alg.Name())
	if err != nil {
		b.Fatalf("unable to get verifier: %s", err)
	}

	b.ResetTimer()
	for i := 0; i <= b.N; i++ {
		if err := verifier.Verify(bytes.NewReader(data), signature); err != nil {
			b.Fatalf("unable to verify signature: %s", err)
		}
	}
}

func BenchmarkVerifierRSA2048SHA256(b *testing.B) {
	benchmarkVerifier(b, generateRSAKey(b, 2048), RS256())
}

func BenchmarkVerifierRSA2048SHA384(b *testing.B) {
	benchmarkVerifier(b, generateRSAKey(b, 2048), RS384())
}

func BenchmarkVerifierRSA2048SHA512(b *testing.B) {
	benchmarkVerifier(b, generateRSAKey(b, 2048), RS512())
}

func BenchmarkVerifierRSA3072SHA256(b *testing.B) {
	benchmarkVerifier(b, generateRSAKey(b, 3072), RS256())
}

func BenchmarkVerifierRSA3072SHA384(b *testing.B) {
	benchmarkVerifier(b, generateRSAKey(b, 3072), RS384())
}

func BenchmarkVerifierRSA3072SHA512(b *testing.B) {
	benchmarkVerifier(b, generateRSAKey(b, 3072), RS512())
}

func BenchmarkVerifierRSA4096SHA256(b *testing.B) {
	benchmarkVerifier(b, generateRSAKey(b, 4096), RS256())
}

func BenchmarkVerifierRSA4096SHA384(b *testing.B) {
	benchmarkVerifier(b, generateRSAKey(b, 4096), RS384())
}

func BenchmarkVerifierRSA4096SHA512(b *testing.B) {
	benchmarkVerifier(b, generateRSAKey(b, 4096), RS512())
}

func BenchmarkVerifierECDSAP256SHA256(b *testing.B) {
	benchmarkVerifier(b, generateECDSAKey(b, elliptic.P256()), ES256())
}

func BenchmarkVerifierECDSAP384SHA384(b *testing.B) {
	benchmarkVerifier(b, generateECDSAKey(b, elliptic.P384()), ES384())
}

func BenchmarkVerifierECDSAP521SHA512(b *testing.B) {
	benchmarkVerifier(b, generateECDSAKey(b, elliptic.P521()), ES512())
}
