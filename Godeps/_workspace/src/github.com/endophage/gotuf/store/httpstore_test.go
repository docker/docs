package store

import (
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"strings"
	"testing"

	"github.com/tent/canonical-json-go"

	"github.com/endophage/gotuf/data"
	"github.com/endophage/gotuf/signed"
)

func TestGetMeta(t *testing.T) {
	store, err := NewHTTPStore(
		"http://mirror1.poly.edu/test-pypi/",
		"metadata",
		"txt",
		"targets",
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
	rootPem := "-----BEGIN PUBLIC KEY-----\nMIIBojANBgkqhkiG9w0BAQEFAAOCAY8AMIIBigKCAYEArvqUPYb6JJROPJQglPTj\n5uDrsxQKl34Mo+3pSlBVuD6puE4lDnG649a2YksJy+C8ZIPJgokn5w+C3alh+dMe\nzbdWHHxrY1h9CLpYz5cbMlE16303ubkt1rvwDqEezG0HDBzPaKj4oP9YJ9x7wbsq\ndvFcy+Qc3wWd7UWcieo6E0ihbJkYcY8chRXVLg1rL7EfZ+e3bq5+ojA2ECM5JqzZ\nzgDpqCv5hTCYYZp72MZcG7dfSPAHrcSGIrwg7whzz2UsEtCOpsJTuCl96FPN7kAu\n4w/WyM3+SPzzr4/RQXuY1SrLCFD8ebM2zHt/3ATLhPnGmyG5I0RGYoegFaZ2AViw\nlqZDOYnBtgDvKP0zakMtFMbkh2XuNBUBO7Sjs0YcZMjLkh9gYUHL1yWS3Aqus1Lw\nlI0gHS22oyGObVBWkZEgk/Foy08sECLGao+5VvhmGpfVuiz9OKFUmtPVjWzRE4ng\niekEu4drSxpH41inLGSvdByDWLpcTvWQI9nkgclh3AT/AgMBAAE=\n-----END PUBLIC KEY-----"
	testKey, _ := pem.Decode([]byte(rootPem))
	k := data.NewPublicKey("rsa", testKey.Bytes)

	sigBytes, err := hex.DecodeString(p.Signatures[0].Signature.String())
	if err != nil {
		t.Fatal(err)
	}
	var decoded map[string]interface{}
	if err := json.Unmarshal(p.Signed, &decoded); err != nil {
		t.Fatal(err)
	}
	msg, err := cjson.Marshal(decoded)
	if err != nil {
		t.Fatal(err)
	}
	method := strings.ToLower(p.Signatures[0].Method)
	err = signed.Verifiers[method].Verify(k, sigBytes, msg)
	if err != nil {
		t.Fatal(err)
	}

}

func TestPyCryptoRSAPSSCompat(t *testing.T) {
	pubPem := "-----BEGIN PUBLIC KEY-----\nMIIBojANBgkqhkiG9w0BAQEFAAOCAY8AMIIBigKCAYEAnKuXZeefa2LmgxaL5NsM\nzKOHNe+x/nL6ik+lDBCTV6OdcwAhHQS+PONGhrChIUVR6Vth3hUCrreLzPO73Oo5\nVSCuRJ53UronENl6lsa5mFKP8StYLvIDITNvkoT3j52BJIjyNUK9UKY9As2TNqDf\nBEPIRp28ev/NViwGOEkBu2UAbwCIdnDXm8JQErCZA0Ydm7PKGgjLbFsFGrVzqXHK\n6pdzJXlhr9yap3UpgQ/iO9JtoEYB2EXsnSrPc9JRjR30bNHHtnVql3fvinXrAEwq\n3xmN4p+R4VGzfdQN+8Kl/IPjqWB535twhFYEG/B7Ze8IwbygBjK3co/KnOPqMUrM\nBI8ztvPiogz+MvXb8WvarZ6TMTh8ifZI96r7zzqyzjR1hJulEy3IsMGvz8XS2J0X\n7sXoaqszEtXdq5ef5zKVxkiyIQZcbPgmpHLq4MgfdryuVVc/RPASoRIXG4lKaTJj\n1ANMFPxDQpHudCLxwCzjCb+sVa20HBRPTnzo8LSZkI6jAgMBAAE=\n-----END PUBLIC KEY-----"
	//privPem := "-----BEGIN RSA PRIVATE KEY-----\n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\n-----END RSA PRIVATE KEY-----"
	testStr := "The quick brown fox jumps over the lazy dog."
	sigHex := "4e05ee9e435653549ac4eddbc43e1a6868636e8ea6dbec2564435afcb0de47e0824cddbd88776ddb20728c53ecc90b5d543d5c37575fda8bd0317025fc07de62ee8084b1a75203b1a23d1ef4ac285da3d1fc63317d5b2cf1aafa3e522acedd366ccd5fe4a7f02a42922237426ca3dc154c57408638b9bfaf0d0213855d4e9ee621db204151bcb13d4dbb18f930ec601469c992c84b14e9e0b6f91ac9517bb3b749dd117e1cbac2e4acb0e549f44558a2005898a226d5b6c8b9291d7abae0d9e0a16858b89662a085f74a202deb867acab792bdbd2c36731217caea8b17bd210c29b890472f11e5afdd1dd7b69004db070e04201778f2c49f5758643881403d45a58d08f51b5c63910c6185892f0b590f191d760b669eff2464456f130239bba94acf54a0cb98f6939ff84ae26a37f9b890be259d9b5d636f6eb367b53e895227d7d79a3a88afd6d28c198ee80f6527437c5fbf63accb81709925c4e03d1c9eaee86f58e4bd1c669d6af042dbd412de0d13b98b1111e2fadbe34b45de52125e9a"
	testKey, _ := pem.Decode([]byte(pubPem))
	k := data.NewPublicKey("rsa", testKey.Bytes)

	sigBytes, err := hex.DecodeString(sigHex)
	if err != nil {
		t.Fatal(err)
	}
	v := signed.RSAPSSVerifier{}
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
	k := data.NewPublicKey("ed25519", pub)

	sigBytes, _ := hex.DecodeString(sigHex)

	err := signed.Verifiers["ed25519"].Verify(k, sigBytes, []byte(testStr))
	if err != nil {
		t.Fatal(err)
	}
}
