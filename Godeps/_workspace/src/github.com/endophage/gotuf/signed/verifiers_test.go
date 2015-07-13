package signed

import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"testing"
	"text/template"

	"github.com/endophage/gotuf/data"
	"github.com/stretchr/testify/assert"
)

type KeyTemplate struct {
	KeyType string
}

const baseRSAKey = `{"keytype":"{{.KeyType}}","keyval":{"public":"MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAyyvBtTg2xzYS+MTTIBqSpI4V78tt8Yzqi7Jki/Z6NqjiDvcnbgcTqNR2t6B2W5NjGdp/hSaT2jyHM+kdmEGaPxg/zIuHbL3NIp4e0qwovWiEgACPIaELdn8O/kt5swsSKl1KMvLCH1sM86qMibNMAZ/hXOwd90TcHXCgZ91wHEAmsdjDC3dB0TT+FBgOac8RM01Y196QrZoOaDMTWh0EQfw7YbXAElhFVDFxBzDdYWbcIHSIogXQmq0CP+zaL/1WgcZZIClt2M6WCaxxF1S34wNn45gCvVZiZQ/iKWHerSr/2dGQeGo+7ezMSutRzvJ+01fInD86RS/CEtBCFZ1VyQIDAQAB","private":"MIIEpAIBAAKCAQEAyyvBtTg2xzYS+MTTIBqSpI4V78tt8Yzqi7Jki/Z6NqjiDvcnbgcTqNR2t6B2W5NjGdp/hSaT2jyHM+kdmEGaPxg/zIuHbL3NIp4e0qwovWiEgACPIaELdn8O/kt5swsSKl1KMvLCH1sM86qMibNMAZ/hXOwd90TcHXCgZ91wHEAmsdjDC3dB0TT+FBgOac8RM01Y196QrZoOaDMTWh0EQfw7YbXAElhFVDFxBzDdYWbcIHSIogXQmq0CP+zaL/1WgcZZIClt2M6WCaxxF1S34wNn45gCvVZiZQ/iKWHerSr/2dGQeGo+7ezMSutRzvJ+01fInD86RS/CEtBCFZ1VyQIDAQABAoIBAHar8FFxrE1gAGTeUpOF8fG8LIQMRwO4U6eVY7V9GpWiv6gOJTHXYFxU/aL0Ty3eQRxwy9tyVRo8EJz5pRex+e6ws1M+jLOviYqW4VocxQ8dZYd+zBvQfWmRfah7XXJ/HPUx2I05zrmR7VbGX6Bu4g5w3KnyIO61gfyQNKF2bm2Q3yblfupx3URvX0bl180R/+QN2Aslr4zxULFE6b+qJqBydrztq+AAP3WmskRxGa6irFnKxkspJqUpQN1mFselj6iQrzAcwkRPoCw0RwCCMq1/OOYvQtgxTJcO4zDVlbw54PvnxPZtcCWw7fO8oZ2Fvo2SDo75CDOATOGaT4Y9iqECgYEAzWZSpFbN9ZHmvq1lJQg//jFAyjsXRNn/nSvyLQILXltz6EHatImnXo3v+SivG91tfzBI1GfDvGUGaJpvKHoomB+qmhd8KIQhO5MBdAKZMf9fZqZofOPTD9xRXECCwdi+XqHBmL+l1OWz+O9Bh+Qobs2as/hQVgHaoXhQpE0NkTcCgYEA/Tjf6JBGl1+WxQDoGZDJrXoejzG9OFW19RjMdmPrg3t4fnbDtqTpZtCzXxPTCSeMrvplKbqAqZglWyq227ksKw4p7O6YfyhdtvC58oJmivlLr6sFaTsER7mDcYce8sQpqm+XQ8IPbnOk0Z1l6g56euTwTnew49uy25M6U1xL0P8CgYEAxEXv2Kw+OVhHV5PX4BBHHj6we88FiDyMfwM8cvfOJ0datekf9X7ImZkmZEAVPJpWBMD+B0J0jzU2b4SLjfFVkzBHVOH2Ob0xCH2MWPAWtekin7OKizUlPbW5ZV8b0+Kq30DQ/4a7D3rEhK8UPqeuX1tHZox1MAqrgbq3zJj4yvcCgYEAktYPKPm4pYCdmgFrlZ+bA0iEPf7Wvbsd91F5BtHsOOM5PQQ7e0bnvWIaEXEad/2CG9lBHlBy2WVLjDEZthILpa/h6e11ao8KwNGY0iKBuebT17rxOVMqqTjPGt8CuD2994IcEgOPFTpkAdUmyvG4XlkxbB8F6St17NPUB5DGuhsCgYA//Lfytk0FflXEeRQ16LT1YXgV7pcR2jsha4+4O5pxSFw/kTsOfJaYHg8StmROoyFnyE3sg76dCgLn0LENRCe5BvDhJnp5bMpQldG3XwcAxH8FGFNY4LtV/2ZKnJhxcONkfmzQPOmTyedOzrKQ+bNURsqLukCypP7/by6afBY4dA=="}}`
const baseRSAx509Key = `{"keytype":"{{.KeyType}}","keyval":{"public":"LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUZLekNDQXhXZ0F3SUJBZ0lRSGZoeWdIbWFkenNMRW9vR0tUbzNuekFMQmdrcWhraUc5dzBCQVFzd09ERWEKTUJnR0ExVUVDaE1SWkc5amEyVnlMbU52YlM5dWIzUmhjbmt4R2pBWUJnTlZCQU1URVdSdlkydGxjaTVqYjIwdgpibTkwWVhKNU1CNFhEVEUxTURjeE16QTBNell4TTFvWERURTNNRGN4TWpBME16WXhNMW93T0RFYU1CZ0dBMVVFCkNoTVJaRzlqYTJWeUxtTnZiUzl1YjNSaGNua3hHakFZQmdOVkJBTVRFV1J2WTJ0bGNpNWpiMjB2Ym05MFlYSjUKTUlJQ0lqQU5CZ2txaGtpRzl3MEJBUUVGQUFPQ0FnOEFNSUlDQ2dLQ0FnRUFuVUZoelBSeUgyOG90SWRJSnlEdApXZDBMcURqQkZMUXNxZXRiTC90QS9hdmxVNE1UQk44eFBJQmJrazNjWDU2bTdOQVBwWDBaZkUzMGc3UXBkVElNCjJteUpNMUtLN2lnQkJzd3czMkpUOVhHRW15K0lWb1Nwc1lCdzJkMWF5dGdxWUI4UXZhZ01zamc4eEc2aWVhUGwKcG9tcUVYdEt1YzBoOTEyaTQ4YURpUzlIK3ExMmlvcmlkVDRmazFrcm1sZ1orMHMrSlZobUFlQ0FiMmZvTFc5YworUDErUnlEQ3FZN2NyaXhZcUJ3c3ZIZ00zbUw4SitmWlZVUWZLYTVmQlA1dFp5MGk3UE9QVFZpdVl3R20rSHlYCmhyQnRpalF0b0R3Y1U4VEVEdDAyelJSd0N3elZKMFhwdGhrVmRqZUNrSkFFcGpyOHVGQ1ZKYlJXOWgrWXRLQlYKMCtzMWl5elFqVWwydklEczRiSVc0RzVaeVp0OHNSaTAzRFhHTnNtNDhrRWlaVWswd0RuNGpzMW8vdUJEVUN6YwphdHdrN2t1aVhrcFFNMVdkRmF6TCtmYWJueWR3Z285bWI2c1FKQlRxMDdvNEI0M0JWYTBHZm5ZSFRsVUtWSHZ6CmNwb1pNWTMyb1AyN0t5TXlybkxETzducUlBQnA1UEFvMUpNU09GWWdKa3R1Sk5LT2h0Sm9qcUgyV21wajRvbzkKQmZMY2d6TFNQd2ZTbytXS0FaVmQzYU1FcnFCQ3RBcVN2aUdmdVRaT3FkK2JKZGY4aW1jZ3ZCeWdacVVRb0J2aAo4Q1hSWGxUNTdKSUFrVkY3aUxrVUZoUkhxY2lwVjZqVzFWeEFXVzJiZ0xrMEhzTnpRQkN2NjQ2YzkwU2d3cGZvCmxLTEJPNFE0QUdsaFFQUmxNQUNPMFRFQ0F3RUFBYU0xTURNd0RnWURWUjBQQVFIL0JBUURBZ0NnTUJNR0ExVWQKSlFRTU1Bb0dDQ3NHQVFVRkJ3TURNQXdHQTFVZEV3RUIvd1FDTUFBd0N3WUpLb1pJaHZjTkFRRUxBNElDQVFCbQo4QWU2RWw5WHlNWHlyRzN0Vkd3clZBZWFYUkNiTFllNDh2b3d1QTA2Ykx1VTh0L0dXcVBRMHhZVFBtRzdsdS9qCjJNalVIeXphZ2hpVUNOdWFvNDhDbGwyemJEajlHZkMvQWJKQUFybGRHc2lReWMwbDY1QUJJaHo5aml1dXlXQ0YKMnBsWFc4RCtldlQxSm5RanRiUXB4c2Q0Um1UOC9NRjVnK29mN0RJU0dGekFIQkNicFFjbTJWRytIZ3NSOEFGcgp6VTg4YU1uakJSNm9CN0IvU0tuaytHNDFrczZLWVJqcmNCS2tBMjlIYUVNUVk5eVNEN2pYUmdJb1pqY2FMR3hlCjAyYldnZTJ2d2hGRkZoYVhaZCtDSWFVWXhvcEVBM3ZCUzlTS1N3UFNQNEpuWDFCZU1KRS8zWElIUVFXdFZuREoKL05YbnFxUTJCNkF1azhMZGRsREpQSDRiNnpZMmdzNmVvVlFRU2FSdUEyd1Q2bkY4WHVIa2dEcUttQ2E4WHVMTgo5bFV0Y0dBeHc0WitUVXlSK2lyRVQwWk14TkNwU01zcUJieGtwU29DaFd2ekgyQTMrMklmSXhielNxWnZoaVF3Ck5zVlpSZTVWNVBSQlE4TVZ3L0FBUE96V0hzWjJCZEw4UXNFQ3Y5dDBlWWxEb3BwMlp5K3RSMkM1SDFQYTg4Y0kKbFFycEs4NGlhVnRYN1ZLek1nZ3hJK0ZsczZaRVR6WnlnT1dvZ0JKMUp5MnJsZ0Z6eFFRYks5S2dCWnl4RnkvZQp3VEVDdW1SSExPN0RucmR2ZU1LY1ZnVTlsaGViQ2ZaNlZiWERUSWFYcGZXYVZSYmpnS1ZwanJSdnZPZTZHVUsyClN3S005dG4wcGRIM09iczV3RzlSZ3pTUkxSUFByMU9TalhTSTI1UGlpUT09Ci0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0K","private":null}}`
const baseECDSAKey = `
{"keytype":"{{.KeyType}}","keyval":{"public":"MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEgl3rzMPMEKhS1k/AX16MM4PdidpjJr+z4pj0Td+30QnpbOIARgpyR1PiFztU8BZlqG3cUazvFclr2q/xHvfrqw==","private":"MHcCAQEEIDqtcdzU7H3AbIPSQaxHl9+xYECt7NpK7B1+6ep5cv9CoAoGCCqGSM49AwEHoUQDQgAEgl3rzMPMEKhS1k/AX16MM4PdidpjJr+z4pj0Td+30QnpbOIARgpyR1PiFztU8BZlqG3cUazvFclr2q/xHvfrqw=="}}`
const baseECDSAx509Key = `{"keytype":"ecdsa-x509","keyval":{"public":"LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUJwRENDQVVtZ0F3SUJBZ0lRQlBWc1NoUmRMdG45WEtRZ29JaDIvREFLQmdncWhrak9QUVFEQWpBNE1Sb3cKR0FZRFZRUUtFeEZrYjJOclpYSXVZMjl0TDI1dmRHRnllVEVhTUJnR0ExVUVBeE1SWkc5amEyVnlMbU52YlM5dQpiM1JoY25rd0hoY05NVFV3TnpFek1EVXdORFF4V2hjTk1UY3dOekV5TURVd05EUXhXakE0TVJvd0dBWURWUVFLCkV4RmtiMk5yWlhJdVkyOXRMMjV2ZEdGeWVURWFNQmdHQTFVRUF4TVJaRzlqYTJWeUxtTnZiUzl1YjNSaGNua3cKV1RBVEJnY3Foa2pPUFFJQkJnZ3Foa2pPUFFNQkJ3TkNBQVI3SjNSOGpWODV5Rnp0dGFTV3FMRDFHa042UHlhWAowUUdmOHh2Rzd6MUYwUG5DQUdSWk9QQ01aWWpZSGVkdzNXY0FmQWVVcDY5OVExSjNEYW9kbzNBcm96VXdNekFPCkJnTlZIUThCQWY4RUJBTUNBS0F3RXdZRFZSMGxCQXd3Q2dZSUt3WUJCUVVIQXdNd0RBWURWUjBUQVFIL0JBSXcKQURBS0JnZ3Foa2pPUFFRREFnTkpBREJHQWlFQWppVkJjaTBDRTBaazgwZ2ZqbytYdE9xM3NURGJkSWJRRTZBTQpoL29mN1RFQ0lRRGxlbXB5MDRhY0RKODNnVHBvaFNtcFJYdjdJbnRLc0lRTU1oLy9VZzliU2c9PQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg==","private":null}}`

func TestRSAVerifier(t *testing.T) {
	// Unmarshal our private RSA Key
	var testRSAKey data.PrivateKey
	var jsonKey bytes.Buffer

	// Execute our template
	templ, _ := template.New("KeyTemplate").Parse(baseRSAKey)
	templ.Execute(&jsonKey, KeyTemplate{KeyType: data.RSAKey})

	json.Unmarshal(jsonKey.Bytes(), &testRSAKey)

	// Sign some data using RSAPSS
	message := []byte("test data for signing")
	hash := crypto.SHA256
	hashed := sha256.Sum256(message)
	signedData, err := rsaSign(&testRSAKey, hash, hashed[:])
	assert.NoError(t, err)

	// Create and call Verify on the verifier
	rsaVerifier := RSAPSSVerifier{}
	err = rsaVerifier.Verify(&testRSAKey, signedData, message)
	assert.NoError(t, err, "expecting success but got error while verifying data using RSA PSS")
}

func TestRSAx509Verifier(t *testing.T) {
	// Unmarshal our private RSA Key
	var testRSAKey data.PublicKey
	var jsonKey bytes.Buffer

	// Execute our template
	templ, _ := template.New("KeyTemplate").Parse(baseRSAx509Key)
	templ.Execute(&jsonKey, KeyTemplate{KeyType: data.RSAx509Key})

	json.Unmarshal(jsonKey.Bytes(), &testRSAKey)

	// Valid signed message
	signedData, _ := hex.DecodeString("9a196a3458e0a9077772c1b3cf6f1605ac69576711d55a89ba390d68723a135aa851cf7a574074ae35fa5c22b0f8e28d31ab05ef66c96456707be2bfa3487edd4531996593bd3f0dd2a6d2034bf4adc1828f5502240a1c4a70506e2b218419d2498487725c22917455617c659087de2a6cb73023bd40dfa868c7e70f1e22e86a4c588f97294f0da1ba537c20a6f06692c6de34c305d3be0bfbeaabb712531d9b52e3118f252c87b27467587b457ae906f73183bec68ae2b56fda41757193b0b7f97fe27cf9efb6be101cad2edd014f5862df6b8fdcd939504f846624349bc480ef3b074b69d5096796c480bf8c6e41b95c2aefa54c6c34d22742c93e82e6dd42080a8d9841057130306f194b07b60c9cb54e5a16b1755f5a1180ab86c2bb244f17c9ccc9c326debacc35dc14a4d8226d75e7cd40b9843e7eacc138d59406d1a5e5f907c8bea588346441f2c464f74e18a0c063bd3ee27ec475929929dd248bcb2972812dc7ce3ab1513bc445f00a43fb98321cae75da6bf8f07ac4f26dd782db57338aa97350814eea55f160ba5c6c893d064edaf8f31d98d2fb544f0b54b5b4e30786dca9f8ef8ea4fa3d1a07335ae2a252079f1ffcadc8f9c53b8c8e32e0e4f9677ef781dba894a49442008d209d3a9b89a03f1ecb191fb1e56f4b894e2c073fe41d41d06a8261804e7321feb095d6da97b4c41aee4180718ee0d9bd964a4e")
	message := []byte("test data for signing")

	// Create and call Verify on the verifier
	rsaVerifier := RSAPSSVerifier{}
	err := rsaVerifier.Verify(&testRSAKey, signedData, message)
	assert.NoError(t, err, "expecting success but got error while verifying data using RSAPSS and an X509 encoded Key")
}

func TestRSAVerifierWithInvalidKeyType(t *testing.T) {
	var testRSAKey data.PrivateKey
	var jsonKey bytes.Buffer

	// Execute our template
	templ, _ := template.New("KeyTemplate").Parse(baseRSAKey)
	templ.Execute(&jsonKey, KeyTemplate{KeyType: "rsa-invalid"})

	json.Unmarshal(jsonKey.Bytes(), &testRSAKey)

	// Valid signed data with invalidRsaKeyJSON
	signedData, _ := hex.DecodeString("2741a57a5ef89f841b4e0a6afbcd7940bc982cd919fbd11dfc21b5ccfe13855b9c401e3df22da5480cef2fa585d0f6dfc6c35592ed92a2a18001362c3a17f74da3906684f9d81c5846bf6a09e2ede6c009ae164f504e6184e666adb14eadf5f6e12e07ff9af9ad49bf1ea9bcfa3bebb2e33be7d4c0fabfe39534f98f1e3c4bff44f637cff3dae8288aea54d86476a3f1320adc39008eae24b991c1de20744a7967d2e685ac0bcc0bc725947f01c9192ffd3e9300eba4b7faa826e84478493fdf97c705dd331dd46072050d6c5e317c2d63df21694dbaf909ebf46ce0ff04f3979fe13723ae1a823c65f27e56efa19e88f9e7b8ee56eac34353b944067deded3a")
	message := []byte("test data for signing")

	// Create and call Verify on the verifier
	rsaVerifier := RSAPSSVerifier{}
	err := rsaVerifier.Verify(&testRSAKey, signedData, message)
	assert.Error(t, err, "invalid key type for RSAPSS verifier: rsa-invalid")
}

func TestRSAVerifierWithInvalidKey(t *testing.T) {
	var testRSAKey data.PrivateKey
	var jsonKey bytes.Buffer

	// Execute our template
	templ, _ := template.New("KeyTemplate").Parse(baseECDSAKey)
	templ.Execute(&jsonKey, KeyTemplate{KeyType: "ecdsa"})

	json.Unmarshal(jsonKey.Bytes(), &testRSAKey)

	// Valid signed data with invalidRsaKeyJSON
	signedData, _ := hex.DecodeString("2741a57a5ef89f841b4e0a6afbcd7940bc982cd919fbd11dfc21b5ccfe13855b9c401e3df22da5480cef2fa585d0f6dfc6c35592ed92a2a18001362c3a17f74da3906684f9d81c5846bf6a09e2ede6c009ae164f504e6184e666adb14eadf5f6e12e07ff9af9ad49bf1ea9bcfa3bebb2e33be7d4c0fabfe39534f98f1e3c4bff44f637cff3dae8288aea54d86476a3f1320adc39008eae24b991c1de20744a7967d2e685ac0bcc0bc725947f01c9192ffd3e9300eba4b7faa826e84478493fdf97c705dd331dd46072050d6c5e317c2d63df21694dbaf909ebf46ce0ff04f3979fe13723ae1a823c65f27e56efa19e88f9e7b8ee56eac34353b944067deded3a")
	message := []byte("test data for signing")

	// Create and call Verify on the verifier
	rsaVerifier := RSAPSSVerifier{}
	err := rsaVerifier.Verify(&testRSAKey, signedData, message)
	assert.Error(t, err, "invalid key type for RSAPSS verifier: ecdsa")
}

func TestRSAVerifierWithInvalidSignature(t *testing.T) {
	var testRSAKey data.PrivateKey
	var jsonKey bytes.Buffer

	// Execute our template
	templ, _ := template.New("KeyTemplate").Parse(baseRSAKey)
	templ.Execute(&jsonKey, KeyTemplate{KeyType: data.RSAKey})

	json.Unmarshal(jsonKey.Bytes(), &testRSAKey)

	// Sign some data using RSAPSS
	message := []byte("test data for signing")
	hash := crypto.SHA256
	hashed := sha256.Sum256(message)
	signedData, err := rsaSign(&testRSAKey, hash, hashed[:])
	assert.NoError(t, err)

	// Modify the signature
	signedData[0] = []byte("a")[0]

	// Create and call Verify on the verifier
	rsaVerifier := RSAPSSVerifier{}
	err = rsaVerifier.Verify(&testRSAKey, signedData, message)
	assert.Error(t, err, "signature verification failed")
}

func TestECDSAVerifier(t *testing.T) {
	var testECDSAKey data.PrivateKey
	var jsonKey bytes.Buffer

	// Execute our template
	templ, _ := template.New("KeyTemplate").Parse(baseECDSAKey)
	templ.Execute(&jsonKey, KeyTemplate{KeyType: data.ECDSAKey})

	json.Unmarshal(jsonKey.Bytes(), &testECDSAKey)

	// Sign some data using ECDSA
	message := []byte("test data for signing")
	hashed := sha256.Sum256(message)
	signedData, err := ecdsaSign(&testECDSAKey, hashed[:])
	assert.NoError(t, err)

	// Create and call Verify on the verifier
	ecdsaVerifier := ECDSAVerifier{}
	err = ecdsaVerifier.Verify(&testECDSAKey, signedData, message)
	assert.NoError(t, err, "expecting success but got error while verifying data using ECDSA")
}

func TestECDSAx509Verifier(t *testing.T) {
	var testECDSAKey data.PrivateKey
	var jsonKey bytes.Buffer

	// Execute our template
	templ, _ := template.New("KeyTemplate").Parse(baseECDSAx509Key)
	templ.Execute(&jsonKey, KeyTemplate{KeyType: data.ECDSAx509Key})

	json.Unmarshal(jsonKey.Bytes(), &testECDSAKey)

	// Valid signature for message
	signedData, _ := hex.DecodeString("b82e0ed5c5dddd74c8d3602bfd900c423511697c3cfe54e1d56b9c1df599695c53aa0caafcdc40df3ef496d78ccf67750ba9413f1ccbd8b0ef137f0da1ee9889")
	message := []byte("test data for signing")

	// Create and call Verify on the verifier
	ecdsaVerifier := ECDSAVerifier{}
	err := ecdsaVerifier.Verify(&testECDSAKey, signedData, message)
	assert.NoError(t, err, "expecting success but got error while verifying data using ECDSA and an x509 encoded key")
}

func TestECDSAVerifierWithInvalidKeyType(t *testing.T) {
	var testECDSAKey data.PrivateKey
	var jsonKey bytes.Buffer

	// Execute our template
	templ, _ := template.New("KeyTemplate").Parse(baseECDSAKey)
	templ.Execute(&jsonKey, KeyTemplate{KeyType: "ecdsa-invalid"})

	json.Unmarshal(jsonKey.Bytes(), &testECDSAKey)

	// Valid signature using invalidECDSAx509Key
	signedData, _ := hex.DecodeString("7b1c45a4dd488a087db46ee459192d890d4f52352620cb84c2c10e0ce8a67fd6826936463a91ffdffab8e6f962da6fc3d3e5735412f7cd161a9fcf97ba1a7033")
	message := []byte("test data for signing")

	// Create and call Verify on the verifier
	ecdsaVerifier := ECDSAVerifier{}
	err := ecdsaVerifier.Verify(&testECDSAKey, signedData, message)
	assert.Error(t, err, "invalid key type for ECDSA verifier: ecdsa-invalid")
}

func TestECDSAVerifierWithInvalidKey(t *testing.T) {
	var testECDSAKey data.PrivateKey
	var jsonKey bytes.Buffer

	// Execute our template
	templ, _ := template.New("KeyTemplate").Parse(baseRSAKey)
	templ.Execute(&jsonKey, KeyTemplate{KeyType: "rsa"})

	json.Unmarshal(jsonKey.Bytes(), &testECDSAKey)

	// Valid signature using invalidECDSAx509Key
	signedData, _ := hex.DecodeString("7b1c45a4dd488a087db46ee459192d890d4f52352620cb84c2c10e0ce8a67fd6826936463a91ffdffab8e6f962da6fc3d3e5735412f7cd161a9fcf97ba1a7033")
	message := []byte("test data for signing")

	// Create and call Verify on the verifier
	ecdsaVerifier := ECDSAVerifier{}
	err := ecdsaVerifier.Verify(&testECDSAKey, signedData, message)
	assert.Error(t, err, "invalid key type for ECDSA verifier: rsa")
}

func TestECDSAVerifierWithInvalidSignature(t *testing.T) {
	var testECDSAKey data.PrivateKey
	var jsonKey bytes.Buffer

	// Execute our template
	templ, _ := template.New("KeyTemplate").Parse(baseECDSAKey)
	templ.Execute(&jsonKey, KeyTemplate{KeyType: data.ECDSAKey})

	json.Unmarshal(jsonKey.Bytes(), &testECDSAKey)

	// Sign some data using ECDSA
	message := []byte("test data for signing")
	hashed := sha256.Sum256(message)
	signedData, err := ecdsaSign(&testECDSAKey, hashed[:])
	assert.NoError(t, err)

	// Modify the signature
	signedData[0] = []byte("a")[0]

	// Create and call Verify on the verifier
	ecdsaVerifier := ECDSAVerifier{}
	err = ecdsaVerifier.Verify(&testECDSAKey, signedData, message)
	assert.Error(t, err, "signature verification failed")

}

func rsaSign(privKey *data.PrivateKey, hash crypto.Hash, hashed []byte) ([]byte, error) {
	if privKey.Cipher() != data.RSAKey {
		return nil, fmt.Errorf("private key type not supported: %s", privKey.Cipher())
	}

	// Create an rsa.PrivateKey out of the private key bytes
	rsaPrivKey, err := x509.ParsePKCS1PrivateKey(privKey.Private())
	if err != nil {
		return nil, err
	}

	// Use the RSA key to RSASSA-PSS sign the data
	sig, err := rsa.SignPSS(rand.Reader, rsaPrivKey, hash, hashed[:], &rsa.PSSOptions{SaltLength: rsa.PSSSaltLengthEqualsHash})
	if err != nil {
		return nil, err
	}

	return sig, nil
}

func ecdsaSign(privKey *data.PrivateKey, hashed []byte) ([]byte, error) {
	if privKey.Cipher() != data.ECDSAKey {
		return nil, fmt.Errorf("private key type not supported: %s", privKey.Cipher())
	}

	// Create an ecdsa.PrivateKey out of the private key bytes
	ecdsaPrivKey, err := x509.ParseECPrivateKey(privKey.Private())
	if err != nil {
		return nil, err
	}

	// Use the ECDSA key to sign the data
	r, s, err := ecdsa.Sign(rand.Reader, ecdsaPrivKey, hashed[:])
	if err != nil {
		return nil, err
	}

	rBytes, sBytes := r.Bytes(), s.Bytes()
	octetLength := (ecdsaPrivKey.Params().BitSize + 7) >> 3

	// MUST include leading zeros in the output
	rBuf := make([]byte, octetLength-len(rBytes), octetLength)
	sBuf := make([]byte, octetLength-len(sBytes), octetLength)

	rBuf = append(rBuf, rBytes...)
	sBuf = append(sBuf, sBytes...)

	return append(rBuf, sBuf...), nil
}
