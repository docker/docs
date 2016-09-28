package licensing

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/docker/dhe-license-server/tiers"
	"gopkg.in/check.v1"
)

// test with our own cert, signature and message
const testAwsCert = `-----BEGIN CERTIFICATE-----
MIICWDCCAcGgAwIBAgIJAL6MQ0+nXsl3MA0GCSqGSIb3DQEBCwUAMEUxCzAJBgNV
BAYTAkFVMRMwEQYDVQQIDApTb21lLVN0YXRlMSEwHwYDVQQKDBhJbnRlcm5ldCBX
aWRnaXRzIFB0eSBMdGQwHhcNMTUwODE3MTYwMzM5WhcNMTUwOTE2MTYwMzM5WjBF
MQswCQYDVQQGEwJBVTETMBEGA1UECAwKU29tZS1TdGF0ZTEhMB8GA1UECgwYSW50
ZXJuZXQgV2lkZ2l0cyBQdHkgTHRkMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKB
gQCxuyEgQ6KT9QwcVMqHm5OL+ZgWfO7FHnO/sLoGFfE9eZ8a07AREHtpvfWWhOBo
eauH/cUfw2QYvWJ8xQ1mWvwTkRaU60ehGTIzGFOElMfsDddtolC+fQItaeAnpWxe
6upHlbT9YKY6mFTtbibj43IO0RayKEieS+sGzrU1SbSwIwIDAQABo1AwTjAdBgNV
HQ4EFgQUt6HJvZ/PlIPCDznNSW7VDXTbdZIwHwYDVR0jBBgwFoAUt6HJvZ/PlIPC
DznNSW7VDXTbdZIwDAYDVR0TBAUwAwEB/zANBgkqhkiG9w0BAQsFAAOBgQADlUPb
IUEvzinjbzwIg8ZifJEepR/01a9Dps6MxAG5QZshWNoLthoXA79NniCOtzMpZV65
vm7YukVZM8lQB+7tnF+DF1pjVx15R5717eGg1C6SmAvLP/dIu5paq/+N+pmwSYVR
eHIDaXsSQ0K0cofIM7oGdjfouGeuiudqHodL0g==
-----END CERTIFICATE-----`

const testSignature = "acf5c5203c6b81cdc863d9ae3bd557019c85a6c914f9a8f801251d20af958bd72140219fae18e6bc6656d4683b1a94df53d137cc49c5f28a924e5adc292f87ce5373dd7969f81b886655373eea0a0d5bd110f09c2190f864a5ca650fe79ed1111435f0352027bbca16b67bcec8ecc6e939f5e17e0648763cdd46e35717808624"

const testMessage = "hi\n"

const realAwsCert = `-----BEGIN CERTIFICATE-----
MIICSzCCAbSgAwIBAgIETJC2dzANBgkqhkiG9w0BAQUFADBqMQswCQYDVQQGEwJV
UzETMBEGA1UECBMKV2FzaGluZ3RvbjEQMA4GA1UEBxMHU2VhdHRsZTEYMBYGA1UE
ChMPQW1hem9uLmNvbSBJbmMuMRowGAYDVQQDExFlYzIuYW1hem9uYXdzLmNvbTAe
Fw0xMDA5MTUxMjA1MTFaFw0xMDEyMTQxMjA1MTFaMGoxCzAJBgNVBAYTAlVTMRMw
EQYDVQQIEwpXYXNoaW5ndG9uMRAwDgYDVQQHEwdTZWF0dGxlMRgwFgYDVQQKEw9B
bWF6b24uY29tIEluYy4xGjAYBgNVBAMTEWVjMi5hbWF6b25hd3MuY29tMIGfMA0G
CSqGSIb3DQEBAQUAA4GNADCBiQKBgQCHvRjf/0kStpJ248khtIaN8qkDN3tkw4Vj
vA9nvPl2anJO+eIBUqPfQG09kZlwpWpmyO8bGB2RWqWxCwuB/dcnIob6w420k9WY
5C0IIGtDRNauN3kuvGXkw3HEnF0EjYr0pcyWUvByWY4KswZV42X7Y7XSS13hOIcL
6NLA+H94/QIDAQABMA0GCSqGSIb3DQEBBQUAA4GBADrBEztYdJwz3bsEwqTsgMRC
IN6EHK95v5/x1DDmzlHV7SH+ok9By/zDlCXRyUctKgxhVXnR0ZIF0tssl7TpAHGs
VFBW4j8M1aZGC0AtTtnhUvc3u6Pu2cf1re4o+I6MmoqJRKcNYh1ejzrYdfx9ebY7
M/apafynee2Beib6aqo6
-----END CERTIFICATE-----`

// test with a real AWS cert, signature and message
const realDocument = `{
  "devpayProductCodes" : null,
  "availabilityZone" : "us-east-1c",
  "privateIp" : "172.31.19.96",
  "version" : "2010-08-31",
  "instanceId" : "i-b8bcb719",
  "billingProducts" : null,
  "instanceType" : "m3.2xlarge",
  "accountId" : "763518987884",
  "pendingTime" : "2015-08-21T18:32:14Z",
  "imageId" : "ami-d1a703ba",
  "kernelId" : null,
  "ramdiskId" : null,
  "architecture" : "x86_64",
  "region" : "us-east-1"
}`

const realSignature = `EFfJ2+7MubY34wQ/uqybAl8WCb4EaZJ4m47ZgNvDl2WmfN0rsP1xgJAEF7s5W2E8f3jiS3fPqHvx
eE1oho+h+eKNy+0KEMwQbKVrCdLLsg1m9iCDkJ4gPPc/w5gcz/6I+lih9E0hYpcl4HIfMdtn2ES0
gqn/A6/Miu+yHBMtPYA=`

const realProductCode = `eq134879o7551p2avbk7uf4lo`

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { check.TestingT(t) }

type HourlySuite struct {
}

var _ = check.Suite(&HourlySuite{})

func (s *HourlySuite) TestSignatureVerification(c *check.C) {
	publicKey, err := certToPubKey(testAwsCert)
	c.Assert(err, check.IsNil)
	c.Assert(publicKey, check.Not(check.IsNil))

	fakeHourlyVerifier := HourlyVerifier{awsPublicKey: publicKey}
	signature, err := hex.DecodeString(testSignature)
	c.Assert(err, check.IsNil)
	success := fakeHourlyVerifier.VerifySignature([]byte(testMessage), signature)
	c.Assert(success, check.IsNil)
}

func (s *HourlySuite) TestIsHourly(c *check.C) {
	fakeHourlyVerifier := HourlyVerifier{}

	fakeHourlyVerifier.configFile = "non-existent-file-name"
	hourly := fakeHourlyVerifier.IsHourly()
	c.Assert(hourly, check.Equals, false)

	file, err := ioutil.TempFile("", "dhe-license-hourly-test")
	c.Assert(err, check.IsNil)
	fakeHourlyVerifier.configFile = file.Name()
	hourly = fakeHourlyVerifier.IsHourly()
	c.Assert(hourly, check.Equals, true)
}

func (s *HourlySuite) TestHourlyTier(c *check.C) {
	fakeHourlyVerifier := HourlyVerifier{}
	file, err := ioutil.TempFile("", "dhe-license-hourly-test")
	c.Assert(err, check.IsNil)
	fakeHourlyVerifier.configFile = file.Name()
	typ := fakeHourlyVerifier.HourlyTier()
	c.Assert(typ, check.Equals, "")

	file.Write([]byte("HourlyAWS"))
	typ = fakeHourlyVerifier.HourlyTier()
	c.Assert(typ, check.Equals, tiers.HourlyAWS)

	file.Write([]byte(" \n"))
	typ = fakeHourlyVerifier.HourlyTier()
	c.Assert(typ, check.Equals, tiers.HourlyAWS)
}

func (s *HourlySuite) TestAwsCheck(c *check.C) {
	metadataServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/latest/dynamic/instance-identity/document" {
			fmt.Fprintf(w, realDocument)
		} else if req.URL.Path == "/latest/dynamic/instance-identity/signature" {
			fmt.Fprintf(w, realSignature)
		} else if req.URL.Path == "/latest/meta-data/product-codes" {
			fmt.Fprintf(w, realProductCode)
		}
	}))
	defer metadataServer.Close()

	awsServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, `<ConfirmProductInstanceResponse xmlns="http://ec2.amazonaws.com/doc/2015-04-15/">
  <requestId>59dbff89-35bd-4eac-99ed-be587EXAMPLE</requestId>
  <return>true</return>
  <ownerId>111122223333</ownerId>
</ConfirmProductInstanceResponse>`)
	}))
	defer awsServer.Close()

	publicKey, _ := certToPubKey(realAwsCert)
	fakeHourlyVerifier := HourlyVerifier{
		awsPublicKey:              publicKey,
		awsMetadataURL:            metadataServer.URL,
		awsConfirmProductInstance: true,
		awsAPIURL:                 awsServer.URL,
		awsAPIAccount:             "not-empty",
		awsAPISecret:              "not-empty",
		awsProductCode:            realProductCode,
	}
	c.Assert(fakeHourlyVerifier.AwsCheck(), check.Equals, nil)
}
func (s *HourlySuite) TestAwsCheckNoConfirm(c *check.C) {
	metadataServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/latest/dynamic/instance-identity/document" {
			fmt.Fprintf(w, realDocument)
		} else if req.URL.Path == "/latest/dynamic/instance-identity/signature" {
			fmt.Fprintf(w, realSignature)
		} else if req.URL.Path == "/latest/meta-data/product-codes" {
			fmt.Fprintf(w, realProductCode)
		}
	}))
	defer metadataServer.Close()

	publicKey, _ := certToPubKey(realAwsCert)
	fakeHourlyVerifier := HourlyVerifier{
		awsPublicKey:              publicKey,
		awsMetadataURL:            metadataServer.URL,
		awsConfirmProductInstance: false,
		awsProductCode:            realProductCode,
	}
	c.Assert(fakeHourlyVerifier.AwsCheck(), check.Equals, nil)
}
