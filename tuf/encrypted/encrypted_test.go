package encrypted

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

var plaintext = []byte("reallyimportant")

func TestRoundtrip(t *testing.T) {
	passphrase := []byte("supersecret")

	enc, err := Encrypt(plaintext, passphrase)
	assert.NoError(t, err)

	// successful decrypt
	dec, err := Decrypt(enc, passphrase)
	assert.NoError(t, err)
	assert.Equal(t, dec, plaintext)

	// wrong passphrase
	passphrase[0] = 0
	dec, err = Decrypt(enc, passphrase)
	assert.Error(t, err)
	assert.Nil(t, dec)
}

func TestTamperedRoundtrip(t *testing.T) {
	passphrase := []byte("supersecret")

	enc, err := Encrypt(plaintext, passphrase)
	assert.NoError(t, err)

	data := &data{}
	err = json.Unmarshal(enc, data)
	assert.NoError(t, err)

	data.Ciphertext[0] = 0
	data.Ciphertext[1] = 0

	enc, _ = json.Marshal(data)

	dec, err := Decrypt(enc, passphrase)
	assert.Error(t, err)
	assert.Nil(t, dec)
}

func TestDecrypt(t *testing.T) {
	enc := []byte(`{"kdf":{"name":"scrypt","params":{"N":32768,"r":8,"p":1},"salt":"N9a7x5JFGbrtB2uBR81jPwp0eiLR4A7FV3mjVAQrg1g="},"cipher":{"name":"nacl/secretbox","nonce":"2h8HxMmgRfuYdpswZBQaU3xJ1nkA/5Ik"},"ciphertext":"SEW6sUh0jf2wfdjJGPNS9+bkk2uB+Cxamf32zR8XkQ=="}`)
	passphrase := []byte("supersecret")

	dec, err := Decrypt(enc, passphrase)
	assert.NoError(t, err)
	assert.Equal(t, dec, plaintext)
}
