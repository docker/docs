package client

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/Sirupsen/logrus"
	tuf "github.com/endophage/gotuf"
	"github.com/endophage/gotuf/data"
	"github.com/endophage/gotuf/keys"
	"github.com/endophage/gotuf/store"
)

func TestClientUpdate(t *testing.T) {
	pypiRoot := `{
 "signatures": [
  {
   "keyid": "92db731bf30e31c13e775360453f0adc8bfd3107f5000b99431c4bdbcebb31ed", 
   "method": "PyCrypto-PKCS#1 PSS", 
   "sig": "31d1fa4712eaf591b367740f69ba7577fb5565863639d7e89145abe159f047d5f848e7696fbaf16f5c4214e1f2295d14e3078f5c0a6e2cbc015c13f8557836a039208970b436bb13b921f86f3e4d4518ce2f731bd7a55083d45634f206dc92e886daeb15e65a513f6451575811e9e44b5573d6999d4b69f86a27f01d0d9a868a535a1f0f6534bd9e555d4df95f019ea6859c83fca30e95e0d1c2ce1dcb2b19c91facd98a7cae9f4c81b5ff4c12980f333e38eac99f4561a8f9e5342382443f165cf0af840d5c61b62698b27413d7f5e1bdba714b98759bcfc7c3d65c567459ce093c66c88ebae836665bcdf2efc1e6921bd792406c0529dab2678a922e5fa6ef39e92f89d3072a22bd755a49b2e2e2c801ce73006fd8a7f56595fae01af3e2a49a3bce4ce9fcaec43e1aae5b0c0fc807bc4caca9fae7c34ff026a417262b5435cfd6cbf17b70aae041eae30a6bd6a857db88566c89f2c02171674b9195f22a84ec5839a4e46e2e63fbb5d4821ae66237ea261846c7d5a341d5f4ab2e12412d4e"
  }
 ], 
 "signed": {
  "_type": "Root", 
  "expires": "2014-08-31 00:49:33 UTC", 
  "keys": {
   "07f12f6c470e60d49fe6a60cd893dfba870db387083d50fe4fc43c6171a0be59": {
    "keytype": "rsa", 
    "keyval": {
     "private": "", 
     "public": "-----BEGIN PUBLIC KEY-----\nMIIBojANBgkqhkiG9w0BAQEFAAOCAY8AMIIBigKCAYEAtb4vuMJF1kOc1deExdqV\nZ3o9W6zulfRDX6u8Exd3NEL0/9bCYluqpshFzgmDYpCsuS79Ve5S8WgXAnEXyWOf\nHL/FnAVqFSnmgmYr52waMYV57BsNMY8pXxOnJm1opUbt0PdnF2D2yfLLY8IgZ/0m\niJ+gWojwEBMRlTOnHFx+l/UVvZAhpVsua6C8C2ThFLrXlmmtkg5BAIm1kWPsZ0nB\n3Cczcwh8zE0l2ytVi458AuRGKwi1TMIGXl064ekjsWexH1zCuoEcv/CKob0BkL4c\nLFUKe6WD6XlEp/xnMh6lNG8LT+UKzSb1hPkTtt23RntFB4Qx7UUT8LoO79bv6coE\njeEqHltmeohHpVmTLWsbTdaX7W3clWPUErGh5kAO0SJu1EM94p5nGWlZ+kwASxtO\nz1qR8AqQ02HBBQU+cY24CNnOwKDq/Kgsg1Aw7bqglvtUwBUQkuuCuisjHI0bSMFr\ntGIhlswxItfQ709/OMrN44Vw/H/Z50UzGFtRlu1h07/vAgMBAAE=\n-----END PUBLIC KEY-----"
    }
   }, 
   "2d9f41a1b79429e9d950a687fe00da0bb4fd751da98aeade1524b8d28968eb89": {
    "keytype": "rsa", 
    "keyval": {
     "private": "", 
     "public": "-----BEGIN PUBLIC KEY-----\nMIIBojANBgkqhkiG9w0BAQEFAAOCAY8AMIIBigKCAYEAmjt8CIbohZxZQ1+I8CQn\nvugLsGXq1G2OiMZvbN281q1/E1futll3+EnfxPLP0SBlyYxixma/ozoMj94lPyOX\nBbiF/U1WJ0Wp5D1kpT0Jzt9Ar0bkRxWoPhubeJ7D4k8Br2m7aG+wchfozdbMmUwK\n/MiAZ1fmpKQAr1ek3/hJiN/dURw+mQEdgXwgA4raDy4Ty3AkG7SDCG1cYoYYMJa3\nKg82AWISQQEHUO1MwRVBon2B5d2UriUEzsYYi+2whDOekchjgyd2xdcRvdCBbdGv\nJBCtzVZpd52lCwqAMJUyDGre6Mb6NmKC2nuk+TYEujqRQK97LnjmPwI42b4cBHqv\nDtg7Z8K3rEQIyD+VvrcWlu+1cE8drjh3y+r7oTtjRPr1M8xaCWn/dh8huSJmaFlk\nmWKDX9KI3/5pxFjgpry20eBRYkZHJWwByc9GVvwhRsIF61QdKvA6uGqkFHy9YQBW\nnzWTPo/4UTUwcpfjneoyMOVKx2K05ZePa5UNAgJVDgFfAgMBAAE=\n-----END PUBLIC KEY-----"
    }
   }, 
   "92db731bf30e31c13e775360453f0adc8bfd3107f5000b99431c4bdbcebb31ed": {
    "keytype": "rsa", 
    "keyval": {
     "private": "", 
     "public": "-----BEGIN PUBLIC KEY-----\nMIIBojANBgkqhkiG9w0BAQEFAAOCAY8AMIIBigKCAYEArvqUPYb6JJROPJQglPTj\n5uDrsxQKl34Mo+3pSlBVuD6puE4lDnG649a2YksJy+C8ZIPJgokn5w+C3alh+dMe\nzbdWHHxrY1h9CLpYz5cbMlE16303ubkt1rvwDqEezG0HDBzPaKj4oP9YJ9x7wbsq\ndvFcy+Qc3wWd7UWcieo6E0ihbJkYcY8chRXVLg1rL7EfZ+e3bq5+ojA2ECM5JqzZ\nzgDpqCv5hTCYYZp72MZcG7dfSPAHrcSGIrwg7whzz2UsEtCOpsJTuCl96FPN7kAu\n4w/WyM3+SPzzr4/RQXuY1SrLCFD8ebM2zHt/3ATLhPnGmyG5I0RGYoegFaZ2AViw\nlqZDOYnBtgDvKP0zakMtFMbkh2XuNBUBO7Sjs0YcZMjLkh9gYUHL1yWS3Aqus1Lw\nlI0gHS22oyGObVBWkZEgk/Foy08sECLGao+5VvhmGpfVuiz9OKFUmtPVjWzRE4ng\niekEu4drSxpH41inLGSvdByDWLpcTvWQI9nkgclh3AT/AgMBAAE=\n-----END PUBLIC KEY-----"
    }
   }, 
   "e1ccde549849bb6aa548412e79c407c93f303f3a3ca0ab1e5923ff7d5c4de769": {
    "keytype": "rsa", 
    "keyval": {
     "private": "", 
     "public": "-----BEGIN PUBLIC KEY-----\nMIIBojANBgkqhkiG9w0BAQEFAAOCAY8AMIIBigKCAYEAqnNAzjIb73u2Nk+r2AdU\nnK+xvKDcZgIzjSzRVjtRJgu3MVffzPcGIsjv2RS8/LLl8nSf48V+tWZf/PnTzkMn\n1iJhbdOQTt6bABiizm5dLP9Jm/AIMTTUpbu+fpFbV8vNH+qM5/Z0WNOptQnOEfNs\n84MNh919lJjHc5VPTw86h68Mkn7W5RChZqnwEv75M1XfWdAnUGeLfZh6BKxFMnaB\nciYxieUvc1bnCtqsdaDE2Ab86WXM7cmNCqyLAh4JkTV+RcMzqEnZCAm68TkO6gM/\n5g4A7fbnn9Jc4O3fJvu80fYOg63nS6hAFldmN74oYu2h5PV75xyTIEERoXqJbSaj\n1Agj/98khIZGsSVjoQ5Mi7ETcSbsH7rqWdHFIu+dbKw8vlnqWnEQYjXA8CIOY3X7\nB6/u5FBS5AdxjO6AR/MuaqpWdDTZwXwgZ9c4wqO4Re4z73sGM1PugUK4dbXIGwVe\nRtzv7cSFOB7OPmj53miJWt2ILACb+Dpnxt9h6TzT1C0JAgMBAAE=\n-----END PUBLIC KEY-----"
    }
   }
  }, 
  "roles": {
   "release": {
    "keyids": [
     "e1ccde549849bb6aa548412e79c407c93f303f3a3ca0ab1e5923ff7d5c4de769"
    ], 
    "threshold": 1
   }, 
   "root": {
    "keyids": [
     "92db731bf30e31c13e775360453f0adc8bfd3107f5000b99431c4bdbcebb31ed"
    ], 
    "threshold": 1
   }, 
   "targets": {
    "keyids": [
     "07f12f6c470e60d49fe6a60cd893dfba870db387083d50fe4fc43c6171a0be59"
    ], 
    "threshold": 1
   }, 
   "timestamp": {
    "keyids": [
     "2d9f41a1b79429e9d950a687fe00da0bb4fd751da98aeade1524b8d28968eb89"
    ], 
    "threshold": 1
   }
  }, 
  "version": 1
 }
}`

	data.SetTUFTypes(
		map[string]string{
			"snapshot": "Release",
		},
	)
	data.SetValidRoles(
		map[string]string{
			"snapshot": "release",
		},
	)

	s := &data.Signed{}
	err := json.Unmarshal([]byte(pypiRoot), s)
	if err != nil {
		t.Fatal(err)
	}
	kdb := keys.NewDB()

	logrus.SetLevel(logrus.DebugLevel)

	// Being able to set the second argument, signer, to nil is a great
	// test as we shouldn't need to instantiate a signer just for reading
	// a repo.
	repo := tuf.NewTufRepo(kdb, nil)
	repo.SetRoot(s)
	remote, err := store.NewHTTPStore(
		"http://mirror1.poly.edu/test-pypi/",
		"metadata",
		"txt",
		"targets",
	)
	cached := store.NewFileCacheStore(remote, "/tmp/tuf")
	if err != nil {
		t.Fatal(err)
	}
	client := Client{
		local:  repo,
		remote: cached,
		keysDB: kdb,
	}

	err = client.Update()
	if err != nil {
		t.Fatal(err)
	}

	testTarget := "packages/2.3/T/TracHTTPAuth/TracHTTPAuth-1.0.1-py2.3.egg"
	expectedHash := "dbcaa6dc0035a636234f9b457d24bf1aeecac0a29b4da97a3b32692f2729f9db"
	expectedSize := int64(8140)
	m := client.TargetMeta(testTarget)
	if m == nil {
		t.Fatal("Failed to find existing target")
	}
	if m.Hashes["sha256"].String() != expectedHash {
		t.Fatal("Target hash incorrect.\nExpected:", expectedHash, "\nReceived:", m.Hashes["sha256"].String())
	}
	if m.Length != expectedSize {
		t.Fatal("Target size incorrect.\nExpected:", expectedSize, "\nReceived:", m.Length)
	}
	err = client.DownloadTarget(ioutil.Discard, testTarget, m)
	if err != nil {
		t.Fatal(err)
	}
}
