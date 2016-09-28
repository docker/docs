package licensing

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/docker/dhe-license-server/tiers"
)

type HourlyVerifier struct {
	configFile                string
	awsPublicKey              *rsa.PublicKey
	awsMetadataURL            string
	awsConfirmProductInstance bool
	awsAPIURL                 string
	awsAPIAccount             string
	awsAPISecret              string
	awsProductCode            string
}

var client = http.Client{
	Timeout: time.Second,
	Transport: &http.Transport{
		Proxy: nil,
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 10 * time.Second,
	},
}

func certToPubKey(certStr string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(certStr))
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}
	return cert.PublicKey.(*rsa.PublicKey), nil
}

func httpGet(url string) (string, error) {
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.WithField("error", err).Error("Failed to receive aws identity document")
		return "", err
	}
	return string(body), nil
}

func (hc *HourlyVerifier) VerifySignature(message, signature []byte) error {
	hasher := sha256.New()
	hasher.Write(message)
	hash := hasher.Sum(nil)
	err := rsa.VerifyPKCS1v15(hc.awsPublicKey, crypto.SHA256, hash, signature)
	return err
}

func (hc *HourlyVerifier) IsHourly() bool {
	_, err := os.Stat(hc.configFile)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

func (hc *HourlyVerifier) HourlyTier() string {
	hourlyType, err := ioutil.ReadFile(hc.configFile)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(hourlyType))
}

func (hc *HourlyVerifier) HourlyIsValid() error {
	switch hc.HourlyTier() {
	case tiers.HourlyAWS:
		return hc.AwsCheck()
	case tiers.HourlyAzure:
		// placeholder
		return errors.New("Azure marketplace is not supported yet.")
	}
	return errors.New("Unknown hourly type.")
}

// AwsCheck returns nil if it has verified that we are running on aws, otherwise it returns an error
func (hc *HourlyVerifier) AwsCheck() error {
	identityDocumentURL := fmt.Sprintf("%s/latest/dynamic/instance-identity/document", hc.awsMetadataURL)
	identityDocument, err := httpGet(identityDocumentURL)
	if err != nil {
		log.WithField("error", err).Error("Failed to fetch aws identity document")
		return err
	}
	identityDocumentSignatureURL := fmt.Sprintf("%s/latest/dynamic/instance-identity/signature", hc.awsMetadataURL)
	identityDocumentSignature, err := httpGet(identityDocumentSignatureURL)
	if err != nil {
		log.WithField("error", err).Error("Failed to fetch aws identity document signature")
		return err
	}
	productCodeURL := fmt.Sprintf("%s/latest/meta-data/product-codes", hc.awsMetadataURL)
	productCode, err := httpGet(productCodeURL)
	if err != nil {
		return err
	}
	if productCode != hc.awsProductCode {
		return fmt.Errorf("Product code doesn't match.")
	}

	signature, err := base64.StdEncoding.DecodeString(identityDocumentSignature)
	if err != nil {
		log.WithField("error", err).Error("Failed to base64 decode identity document")
		return err
	}
	err = hc.VerifySignature([]byte(identityDocument), signature)
	if err != nil {
		return err
	}

	if hc.awsConfirmProductInstance {
		identityDocumentObj := struct {
			InstanceId string `json:"instanceId"`
			Region     string `json:"region"`
		}{}

		err = json.Unmarshal([]byte(identityDocument), &identityDocumentObj)
		if err != nil {
			return err
		}

		awsClient := ec2.New(session.New(), &aws.Config{
			Endpoint:    &hc.awsAPIURL,
			Region:      &identityDocumentObj.Region,
			Credentials: credentials.NewStaticCredentials(hc.awsAPIAccount, hc.awsAPISecret, ""),
		})

		input := &ec2.ConfirmProductInstanceInput{
			InstanceId:  &identityDocumentObj.InstanceId,
			ProductCode: &productCode,
		}
		resp, err := awsClient.ConfirmProductInstance(input)
		if err != nil {
			return err
		}

		if resp.Return == nil || !(*resp.Return) {
			return fmt.Errorf("instance id is not running the product")
		}
	}

	return nil
}
