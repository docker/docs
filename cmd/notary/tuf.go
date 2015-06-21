package main

import (
	"bytes"
	"crypto/sha256"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/Sirupsen/logrus"
	"github.com/endophage/gotuf"
	"github.com/endophage/gotuf/client"
	"github.com/endophage/gotuf/data"
	"github.com/endophage/gotuf/keys"
	"github.com/endophage/gotuf/signed"
	"github.com/endophage/gotuf/store"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	kdb := keys.NewDB()
	signer := signed.NewSigner(NewCryptoService(gun))
	repo := tuf.NewTufRepo(kdb, signer)

	b, err := ioutil.ReadFile(targetPath)
	if err != nil {
		fatalf(err.Error())
	}

	filestore := bootstrapRepo(gun, repo)

	fmt.Println("Generating metadata for target")
	meta, err := data.NewFileMeta(bytes.NewBuffer(b))
	if err != nil {
		fatalf(err.Error())
	}

	fmt.Printf("Adding target \"%s\" with sha256 \"%s\" and size %d bytes.\n", targetName, meta.Hashes["sha256"], meta.Length)
	_, err = repo.AddTargets("targets", data.Files{targetName: meta})
	if err != nil {
		fatalf(err.Error())
	}

	saveRepo(repo, filestore)
}

func tufInit(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		cmd.Usage()
		fatalf("Must specify a GUN")
	}

	gun := args[0]
	kdb := keys.NewDB()
	signer := signed.NewSigner(NewCryptoService(gun))

	rootKey, err := signer.Create("root")
	targetsKey, err := signer.Create("targets")
	snapshotKey, err := signer.Create("snapshot")
	timestampKey, err := signer.Create("timestamp")

	kdb.AddKey(rootKey)
	kdb.AddKey(targetsKey)
	kdb.AddKey(snapshotKey)
	kdb.AddKey(timestampKey)

	rootRole, err := data.NewRole("root", 1, []string{rootKey.ID()}, nil, nil)
	targetsRole, err := data.NewRole("targets", 1, []string{targetsKey.ID()}, nil, nil)
	snapshotRole, err := data.NewRole("snapshot", 1, []string{snapshotKey.ID()}, nil, nil)
	timestampRole, err := data.NewRole("timestamp", 1, []string{timestampKey.ID()}, nil, nil)

	kdb.AddRole(rootRole)
	kdb.AddRole(targetsRole)
	kdb.AddRole(snapshotRole)
	kdb.AddRole(timestampRole)

	repo := tuf.NewTufRepo(kdb, signer)

	filestore, err := store.NewFilesystemStore(
		path.Join(viper.GetString("tufDir"), gun), // TODO: base trust dir from config
		"metadata",
		"json",
		"targets",
	)
	if err != nil {
		fatalf(err.Error())
	}

	err = repo.InitRepo(false)
	if err != nil {
		fatalf(err.Error())
	}
	saveRepo(repo, filestore)
}

func tufList(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		cmd.Usage()
		fatalf("must specify a GUN")
	}
	gun := args[0]
	kdb := keys.NewDB()
	repo := tuf.NewTufRepo(kdb, nil)

	remote, err := store.NewHTTPStore(
		"https://notary:4443/v2/"+gun+"/_trust/tuf/",
		"",
		"json",
		"",
	)
	c, err := bootstrapClient(gun, remote, repo, kdb)
	if err != nil {
		return
	}
	err = c.Update()
	if err != nil {
		logrus.Error("Error updating client: ", err.Error())
		return
	}

	if rawOutput {
		for name, meta := range repo.Targets["targets"].Signed.Targets {
			fmt.Println(name, " ", meta.Hashes["sha256"], " ", meta.Length)
		}
	} else {
		for name, meta := range repo.Targets["targets"].Signed.Targets {
			fmt.Println(name, " ", meta.Hashes["sha256"], " ", meta.Length)
		}
	}
}

func tufLookup(cmd *cobra.Command, args []string) {
	if len(args) < 2 {
		cmd.Usage()
		fatalf("must specify a GUN and target")
	}
	gun := args[0]
	targetName := args[1]
	kdb := keys.NewDB()
	repo := tuf.NewTufRepo(kdb, nil)

	remote, err := store.NewHTTPStore(
		"https://notary:4443/v2/"+gun+"/_trust/tuf/",
		"",
		"json",
		"",
	)
	c, err := bootstrapClient(gun, remote, repo, kdb)
	if err != nil {
		return
	}
	err = c.Update()
	if err != nil {
		logrus.Error("Error updating client: ", err.Error())
		return
	}
	meta := c.TargetMeta(targetName)
	if meta == nil {
		logrus.Infof("Target %s not found in %s.", targetName, gun)
		return
	}
	if rawOutput {
		fmt.Println(targetName, fmt.Sprintf("sha256:%s", meta.Hashes["sha256"]), meta.Length)
	} else {
		fmt.Println(targetName, fmt.Sprintf("sha256:%s", meta.Hashes["sha256"]), meta.Length)
	}
}

func tufPublish(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		cmd.Usage()
		fatalf("Must specify a GUN")
	}

	gun := args[0]
	fmt.Println("Pushing changes to ", gun, ".")

	remote, err := store.NewHTTPStore(
		"https://notary:4443/v2/"+gun+"/_trust/tuf/",
		"",
		"json",
		"",
	)
	filestore, err := store.NewFilesystemStore(
		path.Join(viper.GetString("tufDir"), gun),
		"metadata",
		"json",
		"targets",
	)
	if err != nil {
		fatalf(err.Error())
	}

	root, err := filestore.GetMeta("root", 0)
	if err != nil {
		fatalf(err.Error())
	}
	targets, err := filestore.GetMeta("targets", 0)
	if err != nil {
		fatalf(err.Error())
	}
	snapshot, err := filestore.GetMeta("snapshot", 0)
	if err != nil {
		fatalf(err.Error())
	}
	timestamp, err := filestore.GetMeta("timestamp", 0)
	if err != nil {
		fatalf(err.Error())
	}

	err = remote.SetMeta("root", root)
	if err != nil {
		fatalf(err.Error())
	}
	err = remote.SetMeta("targets", targets)
	if err != nil {
		fatalf(err.Error())
	}
	err = remote.SetMeta("snapshot", snapshot)
	if err != nil {
		fatalf(err.Error())
	}
	err = remote.SetMeta("timestamp", timestamp)
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
	kdb := keys.NewDB()
	signer := signed.NewSigner(NewCryptoService(gun))
	repo := tuf.NewTufRepo(kdb, signer)

	fmt.Println("Removing target ", targetName, " from ", gun)

	filestore := bootstrapRepo(gun, repo)

	err := repo.RemoveTargets("targets", targetName)
	if err != nil {
		fatalf(err.Error())
	}

	saveRepo(repo, filestore)
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
	kdb := keys.NewDB()
	repo := tuf.NewTufRepo(kdb, nil)

	remote, err := store.NewHTTPStore(
		"https://notary:4443/v2/"+gun+"/_trust/tuf/",
		"",
		"json",
		"",
	)
	rootJSON, err := remote.GetMeta("root", 5<<20)
	if err != nil {
		fmt.Println("Couldn't get initial root")
		fatalf(err.Error())
	}
	root := &data.Signed{}
	err = json.Unmarshal(rootJSON, root)
	if err != nil {
		fmt.Println("Couldn't parse initial root")
		fatalf(err.Error())
	}
	// TODO: Validate the root file against the key store
	err = repo.SetRoot(root)
	if err != nil {
		fatalf("Error setting root %v", err)
	}

	c := client.NewClient(
		repo,
		remote,
		kdb,
	)

	err = c.Update()
	if err != nil {
		fmt.Println("Update failed")
		fatalf(err.Error())
	}
	meta := c.TargetMeta(targetName)
	if meta == nil {
		return
	}

	// Create hasher and hash data
	stdinHash := fmt.Sprintf("sha256:%x", sha256.Sum256(payload))
	serverHash := fmt.Sprintf("sha256:%s", meta.Hashes["sha256"])
	if stdinHash != serverHash {
		_, _ = os.Stderr.Write([]byte("Data not present in the trusted collection.\n"))
		os.Exit(1)
	} else {
		_, _ = os.Stdout.Write(payload)
	}
	return
}

func saveRepo(repo *tuf.TufRepo, filestore store.MetadataStore) error {
	fmt.Println("Saving changes to Trusted Collection.")
	signedRoot, err := repo.SignRoot(data.DefaultExpires("root"))
	if err != nil {
		return err
	}
	rootJSON, _ := json.Marshal(signedRoot)
	filestore.SetMeta("root", rootJSON)

	for r, _ := range repo.Targets {
		signedTargets, err := repo.SignTargets(r, data.DefaultExpires("targets"))
		if err != nil {
			return err
		}
		targetsJSON, _ := json.Marshal(signedTargets)
		parentDir := filepath.Dir(r)
		os.MkdirAll(parentDir, 0755)
		filestore.SetMeta(r, targetsJSON)
	}

	signedSnapshot, err := repo.SignSnapshot(data.DefaultExpires("snapshot"))
	if err != nil {
		return err
	}
	snapshotJSON, _ := json.Marshal(signedSnapshot)
	filestore.SetMeta("snapshot", snapshotJSON)

	signedTimestamp, err := repo.SignTimestamp(data.DefaultExpires("timestamp"))
	if err != nil {
		return err
	}
	timestampJSON, _ := json.Marshal(signedTimestamp)
	filestore.SetMeta("timestamp", timestampJSON)
	return nil
}

func bootstrapClient(gun string, remote store.RemoteStore, repo *tuf.TufRepo, kdb *keys.KeyDB) (*client.Client, error) {
	rootJSON, err := remote.GetMeta("root", 5<<20)
	if err != nil {
		return nil, err
	}
	root := &data.Signed{}
	err = json.Unmarshal(rootJSON, root)
	if err != nil {
		return nil, err
	}
	err = validateRoot(gun, root)
	if err != nil {
		return nil, err
	}
	err = repo.SetRoot(root)
	if err != nil {
		return nil, err
	}
	return client.NewClient(
		repo,
		remote,
		kdb,
	), nil
}

func validateRoot(gun string, root *data.Signed) error {
	rootSigned := &data.Root{}
	err := json.Unmarshal(root.Signed, rootSigned)
	if err != nil {
		return err
	}
	certs := make(map[string]*data.PublicKey)
	for _, kID := range rootSigned.Roles["root"].KeyIDs {
		// TODO: currently assuming only one cert contained in
		// public key entry
		k, _ := pem.Decode([]byte(rootSigned.Keys["kid"].Public()))
		rootCert, err := x509.ParseCertificates(k.Bytes)
		if err != nil {
			continue
		}
		err = caStore.Verify(gun, rootCert[0])
		if err != nil {
			continue
		}
		certs[kID] = rootSigned.Keys[kID]
	}
	_, err = signed.VerifyRoot(root, 0, certs, 1)
	if err != nil {
		// failed to validate the signatures against the certificates
		return err
	}
	return nil
}

func bootstrapRepo(gun string, repo *tuf.TufRepo) store.MetadataStore {
	filestore, err := store.NewFilesystemStore(
		path.Join(viper.GetString("tufDir"), gun),
		"metadata",
		"json",
		"targets",
	)
	if err != nil {
		fatalf(err.Error())
	}

	fmt.Println("Loading trusted collection.")
	rootJSON, err := filestore.GetMeta("root", 0)
	if err != nil {
		fatalf(err.Error())
	}
	root := &data.Signed{}
	err = json.Unmarshal(rootJSON, root)
	if err != nil {
		fatalf(err.Error())
	}
	repo.SetRoot(root)
	targetsJSON, err := filestore.GetMeta("targets", 0)
	if err != nil {
		fatalf(err.Error())
	}
	targets := &data.Signed{}
	err = json.Unmarshal(targetsJSON, targets)
	if err != nil {
		fatalf(err.Error())
	}
	repo.SetTargets("targets", targets)
	snapshotJSON, err := filestore.GetMeta("snapshot", 0)
	if err != nil {
		fatalf(err.Error())
	}
	snapshot := &data.Signed{}
	err = json.Unmarshal(snapshotJSON, snapshot)
	if err != nil {
		fatalf(err.Error())
	}
	repo.SetSnapshot(snapshot)
	timestampJSON, err := filestore.GetMeta("timestamp", 0)
	if err != nil {
		fatalf(err.Error())
	}
	timestamp := &data.Signed{}
	err = json.Unmarshal(timestampJSON, timestamp)
	if err != nil {
		fatalf(err.Error())
	}
	repo.SetTimestamp(timestamp)
	return filestore
}
