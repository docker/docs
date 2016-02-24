package main

import (
	"fmt"
	"io/ioutil"

	notaryclient "github.com/docker/notary/client"
	"github.com/docker/notary/passphrase"
	"github.com/docker/notary/trustmanager"
	"github.com/docker/notary/tuf/data"
	"github.com/docker/notary/tuf/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cmdDelegationTemplate = usageTemplate{
	Use:   "delegation",
	Short: "Operates on delegations.",
	Long:  `Operations on TUF delegations.`,
}

var cmdDelegationListTemplate = usageTemplate{
	Use:   "list [ GUN ]",
	Short: "Lists delegations for the Global Unique Name.",
	Long:  "Lists all delegations known to notary for a specific Global Unique Name.",
}

var cmdDelegationRemoveTemplate = usageTemplate{
	Use:   "remove [ GUN ] [ Role ] <KeyID 1> ...",
	Short: "Remove KeyID(s) from the specified Role delegation.",
	Long:  "Remove KeyID(s) from the specified Role delegation in a specific Global Unique Name.",
}

var cmdDelegationAddTemplate = usageTemplate{
	Use:   "add [ GUN ] [ Role ] <X509 file path 1> ...",
	Short: "Add a keys to delegation using the provided public key X509 certificates.",
	Long:  "Add a keys to delegation using the provided public key PEM encoded X509 certificates in a specific Global Unique Name.",
}

type delegationCommander struct {
	// these need to be set
	configGetter func() (*viper.Viper, error)
	retriever    passphrase.Retriever

	paths                         []string
	allPaths, removeAll, forceYes bool
}

func (d *delegationCommander) GetCommand() *cobra.Command {
	cmd := cmdDelegationTemplate.ToCommand(nil)
	cmd.AddCommand(cmdDelegationListTemplate.ToCommand(d.delegationsList))

	cmdRemDelg := cmdDelegationRemoveTemplate.ToCommand(d.delegationRemove)
	cmdRemDelg.Flags().StringSliceVar(&d.paths, "paths", nil, "List of paths to remove")
	cmdRemDelg.Flags().BoolVarP(&d.forceYes, "yes", "y", false, "Answer yes to the removal question (no confirmation)")
	cmdRemDelg.Flags().BoolVar(&d.allPaths, "all-paths", false, "Remove all paths from this delegation")
	cmd.AddCommand(cmdRemDelg)

	cmdAddDelg := cmdDelegationAddTemplate.ToCommand(d.delegationAdd)
	cmdAddDelg.Flags().StringSliceVar(&d.paths, "paths", nil, "List of paths to add")
	cmdAddDelg.Flags().BoolVar(&d.allPaths, "all-paths", false, "Add all paths to this delegation")
	cmd.AddCommand(cmdAddDelg)
	return cmd
}

// delegationsList lists all the delegations for a particular GUN
func (d *delegationCommander) delegationsList(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		cmd.Usage()
		return fmt.Errorf(
			"Please provide a Global Unique Name as an argument to list")
	}

	config, err := d.configGetter()
	if err != nil {
		return err
	}

	gun := args[0]

	rt, err := getTransport(config, gun, true)
	if err != nil {
		return err
	}

	// initialize repo with transport to get latest state of the world before listing delegations
	nRepo, err := notaryclient.NewNotaryRepository(
		config.GetString("trust_dir"), gun, getRemoteTrustServer(config), rt, d.retriever)
	if err != nil {
		return err
	}

	delegationRoles, err := nRepo.GetDelegationRoles()
	if err != nil {
		return fmt.Errorf("Error retrieving delegation roles for repository %s: %v", gun, err)
	}

	cmd.Println("")
	prettyPrintRoles(delegationRoles, cmd.Out(), "delegations")
	cmd.Println("")
	return nil
}

// delegationRemove removes a public key from a specific role in a GUN
func (d *delegationCommander) delegationRemove(cmd *cobra.Command, args []string) error {
	if len(args) < 2 {
		cmd.Usage()
		return fmt.Errorf("must specify the Global Unique Name and the role of the delegation along with optional keyIDs and/or a list of paths to remove")
	}

	config, err := d.configGetter()
	if err != nil {
		return err
	}

	gun := args[0]
	role := args[1]

	// Check if role is valid delegation name before requiring any user input
	if !data.IsDelegation(role) {
		return fmt.Errorf("invalid delegation name %s", role)
	}

	// If we're only given the gun and the role, attempt to remove all data for this delegation
	if len(args) == 2 && d.paths == nil && !d.allPaths {
		d.removeAll = true
	}

	keyIDs := []string{}

	if len(args) > 2 {
		keyIDs = args[2:]
	}

	// If the user passes --all-paths, don't use any of the passed in --paths
	if d.allPaths {
		d.paths = nil
	}

	// no online operations are performed by add so the transport argument
	// should be nil
	nRepo, err := notaryclient.NewNotaryRepository(
		config.GetString("trust_dir"), gun, getRemoteTrustServer(config), nil, d.retriever)
	if err != nil {
		return err
	}

	if d.removeAll {
		cmd.Println("\nAre you sure you want to remove all data for this delegation? (yes/no)")
		// Ask for confirmation before force removing delegation
		if !d.forceYes {
			confirmed := askConfirm()
			if !confirmed {
				fatalf("Aborting action.")
			}
		} else {
			cmd.Println("Confirmed `yes` from flag")
		}
		// Delete the entire delegation
		err = nRepo.RemoveDelegationRole(role)
		if err != nil {
			return fmt.Errorf("failed to remove delegation: %v", err)
		}
	} else {
		if d.allPaths {
			err = nRepo.ClearDelegationPaths(role)
			if err != nil {
				return fmt.Errorf("failed to remove delegation: %v", err)
			}
		}
		// Remove any keys or paths that we passed in
		err = nRepo.RemoveDelegationKeysAndPaths(role, keyIDs, d.paths)
		if err != nil {
			return fmt.Errorf("failed to remove delegation: %v", err)
		}
	}

	cmd.Println("")
	if d.removeAll {
		cmd.Printf("Forced removal (including all keys and paths) of delegation role %s to repository \"%s\" staged for next publish.\n", role, gun)
	} else {
		removingItems := ""
		if len(keyIDs) > 0 {
			removingItems = removingItems + fmt.Sprintf("with keys %s, ", keyIDs)
		}
		if d.allPaths {
			removingItems = removingItems + "with all paths, "
		}
		if d.paths != nil {
			removingItems = removingItems + fmt.Sprintf("with paths [%s], ", prettyPrintPaths(d.paths))
		}
		cmd.Printf("Removal of delegation role %s %sto repository \"%s\" staged for next publish.\n", role, removingItems, gun)
	}
	cmd.Println("")

	return nil
}

// delegationAdd creates a new delegation by adding a public key from a certificate to a specific role in a GUN
func (d *delegationCommander) delegationAdd(cmd *cobra.Command, args []string) error {
	// We must have at least the gun and role name, and at least one key or path (or the --all-paths flag) to add
	if len(args) < 2 || len(args) < 3 && d.paths == nil && !d.allPaths {
		cmd.Usage()
		return fmt.Errorf("must specify the Global Unique Name and the role of the delegation along with the public key certificate paths and/or a list of paths to add")
	}

	config, err := d.configGetter()
	if err != nil {
		return err
	}

	gun := args[0]
	role := args[1]

	pubKeys := []data.PublicKey{}
	if len(args) > 2 {
		pubKeyPaths := args[2:]
		for _, pubKeyPath := range pubKeyPaths {
			// Read public key bytes from PEM file
			pubKeyBytes, err := ioutil.ReadFile(pubKeyPath)
			if err != nil {
				return fmt.Errorf("unable to read public key from file: %s", pubKeyPath)
			}

			// Parse PEM bytes into type PublicKey
			pubKey, err := trustmanager.ParsePEMPublicKey(pubKeyBytes)
			if err != nil {
				return fmt.Errorf("unable to parse valid public key certificate from PEM file %s: %v", pubKeyPath, err)
			}
			pubKeys = append(pubKeys, pubKey)
		}
	}

	for _, path := range d.paths {
		if path == "" {
			d.allPaths = true
			break
		}
	}

	// If the user passes --all-paths (or gave the "" path in --paths), give the "" path
	if d.allPaths {
		d.paths = []string{""}
	}

	// no online operations are performed by add so the transport argument
	// should be nil
	nRepo, err := notaryclient.NewNotaryRepository(
		config.GetString("trust_dir"), gun, getRemoteTrustServer(config), nil, d.retriever)
	if err != nil {
		return err
	}

	// Add the delegation to the repository
	err = nRepo.AddDelegation(role, pubKeys, d.paths)
	if err != nil {
		return fmt.Errorf("failed to create delegation: %v", err)
	}

	// Make keyID slice for better CLI print
	pubKeyIDs := []string{}
	for _, pubKey := range pubKeys {
		pubKeyID, err := utils.CanonicalKeyID(pubKey)
		if err != nil {
			return err
		}
		pubKeyIDs = append(pubKeyIDs, pubKeyID)
	}

	cmd.Println("")
	addingItems := ""
	if len(pubKeyIDs) > 0 {
		addingItems = addingItems + fmt.Sprintf("with keys %s, ", pubKeyIDs)
	}
	if d.paths != nil || d.allPaths {
		addingItems = addingItems + fmt.Sprintf("with paths [%s], ", prettyPrintPaths(d.paths))
	}
	cmd.Printf(
		"Addition of delegation role %s %sto repository \"%s\" staged for next publish.\n",
		role, addingItems, gun)
	cmd.Println("")
	return nil
}
