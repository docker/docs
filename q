[1mdiff --git a/client/client.go b/client/client.go[m
[1mindex 6916daf..8029996 100644[m
[1m--- a/client/client.go[m
[1m+++ b/client/client.go[m
[36m@@ -618,7 +618,7 @@[m [mfunc (r *NotaryRepository) ListRootKeys() []string {[m
 func (r *NotaryRepository) GenRootKey(passphrase string) (string, error) {[m
 	privKey, err := trustmanager.GenerateRSAKey(rand.Reader, rsaRootKeySize)[m
 	if err != nil {[m
[31m-		return "", fmt.Errorf("failed to convert private key: ", err)[m
[32m+[m		[32mreturn "", fmt.Errorf("failed to convert private key: %v", err)[m
 	}[m
 [m
 	r.rootKeyStore.AddEncryptedKey(privKey.ID(), privKey, passphrase)[m
[1mdiff --git a/cmd/notary/tuf.go b/cmd/notary/tuf.go[m
[1mindex af21933..7825170 100644[m
[1m--- a/cmd/notary/tuf.go[m
[1m+++ b/cmd/notary/tuf.go[m
[36m@@ -2,14 +2,13 @@[m [mpackage main[m
 [m
 import ([m
 	"crypto/sha256"[m
[32m+[m	[32m"errors"[m
 	"fmt"[m
 	"io/ioutil"[m
 	"os"[m
 [m
 	"github.com/Sirupsen/logrus"[m
 	notaryclient "github.com/docker/notary/client"[m
[31m-	"github.com/endophage/gotuf/data"[m
[31m-	"github.com/endophage/gotuf/keys"[m
 	"github.com/spf13/cobra"[m
 	"github.com/spf13/viper"[m
 )[m
[36m@@ -107,13 +106,30 @@[m [mfunc tufInit(cmd *cobra.Command, args []string) {[m
 		fatalf(err.Error())[m
 	}[m
 [m
[31m-	// TODO(diogo): We don't want to generate a new root every time. Ask the user[m
[31m-	// which key she wants to use if there > 0 root keys available.[m
[31m-	rootKeyID, err := nRepo.GenRootKey("passphrase")[m
[31m-	if err != nil {[m
[31m-		fatalf(err.Error())[m
[32m+[m	[32mkeysList := nRepo.ListRootKeys()[m
[32m+[m	[32mvar passphrase string[m
[32m+[m	[32mvar rootKeyID string[m
[32m+[m	[32mif len(keysList) < 1 {[m
[32m+[m		[32mfmt.Println("No root keys found. Generating a new root key...")[m
[32m+[m		[32mpassphrase, err = passphraseRetriever()[m
[32m+[m		[32mif err != nil {[m
[32m+[m			[32mfatalf(err.Error())[m
[32m+[m		[32m}[m
[32m+[m		[32mrootKeyID, err = nRepo.GenRootKey(passphrase)[m
[32m+[m		[32mif err != nil {[m
[32m+[m			[32mfatalf(err.Error())[m
[32m+[m		[32m}[m
[32m+[m	[32m} else {[m
[32m+[m		[32mrootKeyID = keysList[0][m
[32m+[m		[32mfmt.Println("Root key found.")[m
[32m+[m		[32mfmt.Printf("Enter passphrase for: %s (%d)\n", rootKeyID, len(rootKeyID))[m
[32m+[m		[32mpassphrase, err = passphraseRetriever()[m
[32m+[m		[32mif err != nil {[m
[32m+[m			[32mfatalf(err.Error())[m
[32m+[m		[32m}[m
 	}[m
[31m-	rootSigner, err := nRepo.GetRootSigner(rootKeyID, "passphrase")[m
[32m+[m
[32m+[m	[32mrootSigner, err := nRepo.GetRootSigner(rootKeyID, passphrase)[m
 	if err != nil {[m
 		fatalf(err.Error())[m
 	}[m
[36m@@ -185,7 +201,7 @@[m [mfunc tufPublish(cmd *cobra.Command, args []string) {[m
 		fatalf(err.Error())[m
 	}[m
 [m
[31m-	err = repo.Publish(passwordRetriever)[m
[32m+[m	[32merr = repo.Publish(passphraseRetriever)[m
 	if err != nil {[m
 		fatalf(err.Error())[m
 	}[m
[36m@@ -249,76 +265,20 @@[m [mfunc verify(cmd *cobra.Command, args []string) {[m
 	return[m
 }[m
 [m
[31m-//func generateKeys(kdb *keys.KeyDB, signer *signed.Signer, remote store.RemoteStore) (string, string, string, string, error) {[m
[31m-//	rawTSKey, err := remote.GetKey("timestamp")[m
[31m-//	if err != nil {[m
[31m-//		return "", "", "", "", err[m
[31m-//	}[m
[31m-//	fmt.Println("RawKey: ", string(rawTSKey))[m
[31m-//	parsedKey := &data.TUFKey{}[m
[31m-//	err = json.Unmarshal(rawTSKey, parsedKey)[m
[31m-//	if err != nil {[m
[31m-//		return "", "", "", "", err[m
[31m-//	}[m
[31m-//	timestampKey := data.NewPublicKey(parsedKey.Cipher(), parsedKey.Public())[m
[31m-//[m
[31m-//	rootKey, err := signer.Create("root")[m
[31m-//	if err != nil {[m
[31m-//		return "", "", "", "", err[m
[31m-//	}[m
[31m-//	targetsKey, err := signer.Create("targets")[m
[31m-//	if err != nil {[m
[31m-//		return "", "", "", "", err[m
[31m-//	}[m
[31m-//	snapshotKey, err := signer.Create("snapshot")[m
[31m-//	if err != nil {[m
[31m-//		return "", "", "", "", err[m
[31m-//	}[m
[31m-//[m
[31m-//	kdb.AddKey(rootKey)[m
[31m-//	kdb.AddKey(targetsKey)[m
[31m-//	kdb.AddKey(snapshotKey)[m
[31m-//	kdb.AddKey(timestampKey)[m
[31m-//	return rootKey.ID(), targetsKey.ID(), snapshotKey.ID(), timestampKey.ID(), nil[m
[31m-//}[m
[31m-[m
[31m-func generateRoles(kdb *keys.KeyDB, rootKeyID, targetsKeyID, snapshotKeyID, timestampKeyID string) error {[m
[31m-	rootRole, err := data.NewRole("root", 1, []string{rootKeyID}, nil, nil)[m
[31m-	if err != nil {[m
[31m-		return err[m
[31m-	}[m
[31m-	targetsRole, err := data.NewRole("targets", 1, []string{targetsKeyID}, nil, nil)[m
[31m-	if err != nil {[m
[31m-		return err[m
[31m-	}[m
[31m-	snapshotRole, err := data.NewRole("snapshot", 1, []string{snapshotKeyID}, nil, nil)[m
[31m-	if err != nil {[m
[31m-		return err[m
[31m-	}[m
[31m-	timestampRole, err := data.NewRole("timestamp", 1, []string{timestampKeyID}, nil, nil)[m
[31m-	if err != nil {[m
[31m-		return err[m
[31m-	}[m
[32m+[m[32m// func passwordRetriever() (string, error) {[m
[32m+[m[32m// 	return "passphrase", nil[m
[32m+[m[32m// }[m
 [m
[31m-	err = kdb.AddRole(rootRole)[m
[32m+[m[32mfunc passphraseRetriever() (string, error) {[m
[32m+[m	[32mfmt.Println("Please provide a passphrase for this root key: ")[m
[32m+[m	[32mvar passphrase string[m
[32m+[m	[32m_, err := fmt.Scanln(&passphrase)[m
 	if err != nil {[m
[31m-		return err[m
[32m+[m		[32mreturn "", err[m
 	}[m
[31m-	err = kdb.AddRole(targetsRole)[m
[31m-	if err != nil {[m
[31m-		return err[m
[31m-	}[m
[31m-	err = kdb.AddRole(snapshotRole)[m
[31m-	if err != nil {[m
[31m-		return err[m
[32m+[m	[32mif len(passphrase) < 8 {[m
[32m+[m		[32mfmt.Println("Please use a password manager to generate and store a good random passphrase.")[m
[32m+[m		[32mreturn "", errors.New("Passphrase too short")[m
 	}[m
[31m-	err = kdb.AddRole(timestampRole)[m
[31m-	if err != nil {[m
[31m-		return err[m
[31m-	}[m
[31m-	return nil[m
[31m-}[m
[31m-[m
[31m-func passwordRetriever() (string, error) {[m
[31m-	return "passphrase", nil[m
[32m+[m	[32mreturn passphrase, nil[m
 }[m
[1mdiff --git a/trustmanager/keyfilestore.go b/trustmanager/keyfilestore.go[m
[1mindex 6418139..f076c79 100644[m
[1m--- a/trustmanager/keyfilestore.go[m
[1m+++ b/trustmanager/keyfilestore.go[m
[36m@@ -1,6 +1,11 @@[m
 package trustmanager[m
 [m
[31m-import "github.com/endophage/gotuf/data"[m
[32m+[m[32mimport ([m
[32m+[m	[32m"path/filepath"[m
[32m+[m	[32m"strings"[m
[32m+[m
[32m+[m	[32m"github.com/endophage/gotuf/data"[m
[32m+[m[32m)[m
 [m
 const ([m
 	keyExtension = "key"[m
[36m@@ -79,5 +84,10 @@[m [mfunc (s *KeyFileStore) GetDecryptedKey(name string, passphrase string) (*data.Pr[m
 // There might be symlinks associating Certificate IDs to Public Keys, so this[m
 // method only returns the IDs that aren't symlinks[m
 func (s *KeyFileStore) ListKeys() []string {[m
[31m-	return s.ListFiles(false)[m
[32m+[m	[32mvar keyIDList []string[m
[32m+[m	[32mfor _, f := range s.ListFiles(false) {[m
[32m+[m		[32mkeyID := strings.TrimSpace(strings.TrimSuffix(filepath.Base(f), filepath.Ext(f)))[m
[32m+[m		[32mkeyIDList = append(keyIDList, keyID)[m
[32m+[m	[32m}[m
[32m+[m	[32mreturn keyIDList[m
 }[m
