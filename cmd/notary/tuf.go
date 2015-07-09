package main

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/Sirupsen/logrus"
	notaryclient "github.com/docker/notary/client"
	"github.com/endophage/gotuf/data"
	"github.com/endophage/gotuf/keys"
	"github.com/spf13/cobra"
)

var remoteTrustServer string

var cmdTufList = &cobra.Command{
	Use:   "list [ GUN ]",
	Short: "Lists targets for a trusted collection.",
	Long:  "Lists all targets for a trusted collection identified by the Globally Unique Name.",
	Run:   tufList,
}

var cmdTufAdd = &cobra.Command{
	Use:   "add [ GUN ] <target> <file>",
	Short: "adds the file as a target to the trusted collection.",
	Long:  "adds the file as a target to the local trusted collection identified by the Globally Unique Name.",
	Run:   tufAdd,
}

var cmdTufRemove = &cobra.Command{
	Use:   "remove [ GUN ] <target>",
	Short: "Removes a target from a trusted collection.",
	Long:  "removes a target from the local trusted collection identified by the Globally Unique Name.",
	Run:   tufRemove,
}

var cmdTufInit = &cobra.Command{
	Use:   "init [ GUN ]",
	Short: "initializes a local trusted collection.",
	Long:  "initializes a local trusted collection identified by the Globally Unique Name.",
	Run:   tufInit,
}

var cmdTufLookup = &cobra.Command{
	Use:   "lookup [ GUN ] <target>",
	Short: "Looks up a specific target in a trusted collection.",
	Long:  "looks up a specific target in a trusted collection identified by the Globally Unique Name.",
	Run:   tufLookup,
}

var cmdTufPublish = &cobra.Command{
	Use:   "publish [ GUN ]",
	Short: "publishes the local trusted collection.",
	Long:  "publishes the local trusted collection identified by the Globally Unique Name, sending the local changes to a remote trusted server.",
	Run:   tufPublish,
}

var cmdVerify = &cobra.Command{
	Use:   "verify [ GUN ] <target>",
	Short: "verifies if the content is included in the trusted collection",
	Long:  "verifies if the data passed in STDIN is included in the trusted collection identified by the Global Unique Name.",
	Run:   verify,
}

func tufAdd(cmd *cobra.Command, args []string) {
	if len(args) < 3 {
		cmd.Usage()
		fatalf("must specify a GUN, target, and path to target data")
	}

	gun := args[0]
	targetName := args[1]
	targetPath := args[2]

	t := &http.Transport{}
	repo, err := nClient.GetRepository(gun, "", t)
	if err != nil {
		fatalf(err.Error())
	}

	target, err := notaryclient.NewTarget(targetName, targetPath)
	if err != nil {
		fatalf(err.Error())
	}
	err = repo.AddTarget(target)
	if err != nil {
		fatalf(err.Error())
	}
	fmt.Println("Successfully added targets")
}

func tufInit(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		cmd.Usage()
		fatalf("Must specify a GUN")
	}

	t := &http.Transport{}

	// TODO(diogo): We don't want to generate a new root every time. Ask the user
	// which key she wants to use if there > 0 root keys available.
	rootKeyID, err := nClient.GenRootKey("passphrase")
	if err != nil {
		fatalf(err.Error())
	}
	rootSigner, err := nClient.GetRootSigner(rootKeyID, "passphrase")
	if err != nil {
		fatalf(err.Error())
	}

	_, err = nClient.InitRepository(args[0], "", t, rootSigner)
	if err != nil {
		fatalf(err.Error())
	}
}

func tufList(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		cmd.Usage()
		fatalf("must specify a GUN")
	}
	gun := args[0]

	t := &http.Transport{}
	repo, err := nClient.GetRepository(gun, "", t)
	if err != nil {
		fatalf(err.Error())
	}

	// Retreive the remote list of signed targets
	targetList, err := repo.ListTargets()
	if err != nil {
		fatalf(err.Error())
	}

	// Print all the available targets
	for _, t := range targetList {
		fmt.Println(t.Name, " ", t.Hashes["sha256"], " ", t.Length)
	}
}

func tufLookup(cmd *cobra.Command, args []string) {
	if len(args) < 2 {
		cmd.Usage()
		fatalf("must specify a GUN and target")
	}
	gun := args[0]
	targetName := args[1]

	t := &http.Transport{}
	repo, err := nClient.GetRepository(gun, "", t)
	if err != nil {
		fatalf(err.Error())
	}

	// TODO(diogo): Parse Targets and print them
	target, err := repo.GetTargetByName(targetName)
	if err != nil {
		fatalf(err.Error())
	}

	fmt.Println(target.Name, fmt.Sprintf("sha256:%s", target.Hashes["sha256"]), target.Length)
}

func tufPublish(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		cmd.Usage()
		fatalf("Must specify a GUN")
	}

	gun := args[0]

	fmt.Println("Pushing changes to ", gun, ".")

	t := &http.Transport{}
	repo, err := nClient.GetRepository(gun, "", t)
	if err != nil {
		fatalf(err.Error())
	}

	err = repo.Publish()
	if err != nil {
		fatalf(err.Error())
	}
}

func tufRemove(cmd *cobra.Command, args []string) {
	if len(args) < 2 {
		cmd.Usage()
		fatalf("must specify a GUN and target")
	}
	gun := args[0]
	targetName := args[1]

	//c := changelist.NewTufChange(changelist.ActionDelete, "targets", "target", targetName, nil)
	//err := cl.Add(c)
	//if err != nil {
	//	fatalf(err.Error())
	//}

	// TODO(diogo): Implement RemoveTargets in libnotary
	fmt.Println("Removing target ", targetName, " from ", gun)
}

func verify(cmd *cobra.Command, args []string) {
	if len(args) < 2 {
		cmd.Usage()
		fatalf("must specify a GUN and target")
	}

	// Reads all of the data on STDIN
	//TODO (diogo): Change this to do a streaming hash
	payload, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fatalf("error reading content from STDIN: %v", err)
	}

	//TODO (diogo): This code is copy/pasted from lookup.
	gun := args[0]
	targetName := args[1]
	t := &http.Transport{}
	repo, err := nClient.GetRepository(gun, "", t)
	if err != nil {
		fatalf(err.Error())
	}

	// TODO(diogo): Parse Targets and print them
	target, err := repo.GetTargetByName(targetName)
	if err != nil {
		logrus.Error("notary: data not present in the trusted collection.")
		os.Exit(-11)
	}

	// Create hasher and hash data
	stdinHash := fmt.Sprintf("sha256:%x", sha256.Sum256(payload))
	serverHash := fmt.Sprintf("sha256:%s", target.Hashes["sha256"])
	if stdinHash != serverHash {
		logrus.Error("notary: data not present in the trusted collection.")
		os.Exit(1)
	} else {
		_, _ = os.Stdout.Write(payload)
	}
	return
}

//func generateKeys(kdb *keys.KeyDB, signer *signed.Signer, remote store.RemoteStore) (string, string, string, string, error) {
//	rawTSKey, err := remote.GetKey("timestamp")
//	if err != nil {
//		return "", "", "", "", err
//	}
//	fmt.Println("RawKey: ", string(rawTSKey))
//	parsedKey := &data.TUFKey{}
//	err = json.Unmarshal(rawTSKey, parsedKey)
//	if err != nil {
//		return "", "", "", "", err
//	}
//	timestampKey := data.NewPublicKey(parsedKey.Cipher(), parsedKey.Public())
//
//	rootKey, err := signer.Create("root")
//	if err != nil {
//		return "", "", "", "", err
//	}
//	targetsKey, err := signer.Create("targets")
//	if err != nil {
//		return "", "", "", "", err
//	}
//	snapshotKey, err := signer.Create("snapshot")
//	if err != nil {
//		return "", "", "", "", err
//	}
//
//	kdb.AddKey(rootKey)
//	kdb.AddKey(targetsKey)
//	kdb.AddKey(snapshotKey)
//	kdb.AddKey(timestampKey)
//	return rootKey.ID(), targetsKey.ID(), snapshotKey.ID(), timestampKey.ID(), nil
//}

func generateRoles(kdb *keys.KeyDB, rootKeyID, targetsKeyID, snapshotKeyID, timestampKeyID string) error {
	rootRole, err := data.NewRole("root", 1, []string{rootKeyID}, nil, nil)
	if err != nil {
		return err
	}
	targetsRole, err := data.NewRole("targets", 1, []string{targetsKeyID}, nil, nil)
	if err != nil {
		return err
	}
	snapshotRole, err := data.NewRole("snapshot", 1, []string{snapshotKeyID}, nil, nil)
	if err != nil {
		return err
	}
	timestampRole, err := data.NewRole("timestamp", 1, []string{timestampKeyID}, nil, nil)
	if err != nil {
		return err
	}

	err = kdb.AddRole(rootRole)
	if err != nil {
		return err
	}
	err = kdb.AddRole(targetsRole)
	if err != nil {
		return err
	}
	err = kdb.AddRole(snapshotRole)
	if err != nil {
		return err
	}
	err = kdb.AddRole(timestampRole)
	if err != nil {
		return err
	}
	return nil
}
