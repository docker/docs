package encrypted

import (
	"encoding/json"
	"testing"

	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type EncryptedSuite struct{}

var _ = Suite(&EncryptedSuite{})

var plaintext = []byte("reallyimportant")

func (EncryptedSuite) TestRoundtrip(c *C) {
	passphrase := []byte("supersecret")

	enc, err := Encrypt(plaintext, passphrase)
	c.Assert(err, IsNil)

	// successful decrypt
	dec, err := Decrypt(enc, passphrase)
	c.Assert(err, IsNil)
	c.Assert(dec, DeepEquals, plaintext)

	// wrong passphrase
	passphrase[0] = 0
	dec, err = Decrypt(enc, passphrase)
	c.Assert(err, NotNil)
	c.Assert(dec, IsNil)
}

func (EncryptedSuite) TestTamperedRoundtrip(c *C) {
	passphrase := []byte("supersecret")

	enc, err := Encrypt(plaintext, passphrase)
	c.Assert(err, IsNil)

	data := &data{}
	err = json.Unmarshal(enc, data)
	c.Assert(err, IsNil)

	data.Ciphertext[0] = 0
	data.Ciphertext[1] = 0

	enc, _ = json.Marshal(data)

	dec, err := Decrypt(enc, passphrase)
	c.Assert(err, NotNil)
	c.Assert(dec, IsNil)
}

func (EncryptedSuite) TestDecrypt(c *C) {
	enc := []byte(`{"kdf":{"name":"scrypt","params":{"N":32768,"r":8,"p":1},"salt":"N9a7x5JFGbrtB2uBR81jPwp0eiLR4A7FV3mjVAQrg1g="},"cipher":{"name":"nacl/secretbox","nonce":"2h8HxMmgRfuYdpswZBQaU3xJ1nkA/5Ik"},"ciphertext":"SEW6sUh0jf2wfdjJGPNS9+bkk2uB+Cxamf32zR8XkQ=="}`)
	passphrase := []byte("supersecret")

	dec, err := Decrypt(enc, passphrase)
	c.Assert(err, IsNil)
	c.Assert(dec, DeepEquals, plaintext)
}
