package main

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/Sirupsen/logrus"
	notaryclient "github.com/docker/notary/client"
	"github.com/endophage/gotuf/store"
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
	repo.AddTarget(target)
	fmt.Println("Successfully added targets")
}

func tufInit(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		cmd.Usage()
		fatalf("Must specify a GUN")
	}

	t := &http.Transport{}
	repo, err := nClient.GetRepository(args[0], "", t)
	if err != nil {
		fatalf(err.Error())
	}

	// TODO(diogo): We don't want to generate a new root every time. Ask the user
	// which key she wants to use if there > 0 root keys available.
	newRootKey, err := nClient.GenRootKey("passphrase")
	if err != nil {
		fatalf(err.Error())
	}
	repo.Initialize(newRootKey)
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

	// TODO(diogo): Parse Targets and print them
	_, _ = repo.ListTargets()
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

	repo.Publish()
}

func tufRemove(cmd *cobra.Command, args []string) {
	if len(args) < 2 {
		cmd.Usage()
		fatalf("must specify a GUN and target")
	}
	gun := args[0]
	targetName := args[1]

	t := &http.Transport{}
	_, err := nClient.GetRepository(gun, "", t)
	if err != nil {
		fatalf(err.Error())
	}

	// TODO(diogo): Implement RemoveTargets in libnotary
	fmt.Println("Removing target ", targetName, " from ", gun)
	// repo.RemoveTargets("targets", targetName)
	// if err != nil {
	// 	fatalf(err.Error())
	// }
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

// Use this to initialize remote HTTPStores from the config settings
func getRemoteStore(gun string) (store.RemoteStore, error) {
	return store.NewHTTPStore(
		"https://notary:4443/v2/"+gun+"/_trust/tuf/",
		"",
		"json",
		"",
		"key",
	)
}
