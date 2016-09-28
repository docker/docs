package cmd

import (
	"crypto/tls"
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/docker/orca/enzi/api/forms"
	"github.com/docker/orca/enzi/passwords"
	"github.com/docker/orca/enzi/schema"
	"golang.org/x/crypto/ssh/terminal"
)

// CreateAdmin is the command for creating an initial admin user.
var CreateAdmin = cli.Command{
	Name:  "create-admin",
	Usage: "Create an Admin User Account",
	Description: `
	Create an Admin User Account with the provided username and password.

	Use the environment variables USERNAME and PASSWORD or set the
	--interactive flag to prompt for input on the command line.`,
	Action: runCreateAdmin,
}

var interactive bool

func init() {
	CreateAdmin.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "interactive, i",
			Usage:       "prompt for username and password on the command line",
			Destination: &interactive,
		},
	}
}

func runCreateAdmin(*cli.Context) error {
	tlsConfig := GetTLSConfig(tls.NoClientCert)

	username, password := getUsernamePassword()

	fmt.Println("connecting to db ...")
	dbSession := GetDBSession(tlsConfig)
	defer dbSession.Close()

	passwordHash, err := passwords.HashPassword(password)
	if err != nil {
		log.Fatalf("unable to hash password: %s", err)
	}

	mgr := schema.NewRethinkDBManager(dbSession)

	user := &schema.Account{
		Name:         username,
		IsAdmin:      true,
		IsActive:     true,
		PasswordHash: passwordHash,
	}

	fmt.Println("creating admin user account...")
	if err := mgr.CreateAccount(user); err != nil {
		log.Fatalf("unable to create admin user: %s", err)
	}

	fmt.Fprintf(os.Stderr, "successfully created admin user account: %s\n", user.Name)

	return nil
}

func getUsernamePassword() (username, password string) {
	if interactive {
		username = promptUsername()
		password = promptPassword(true, "")

		return username, password
	}

	// Get the username and password from the environment.
	username = os.Getenv("USERNAME")
	password = os.Getenv("PASSWORD")

	if username == "" || password == "" {
		log.Fatalln("must specify both USERNAME and PASSWORD environment variables or use --interactive flag")
	}

	if validationErr := forms.ValidateAccountName(&username, "username"); validationErr != nil {
		log.Fatalf("invalid: %s\n", validationErr.Detail)
	}

	if validationErr := forms.ValidatePassword(password, "password"); validationErr != nil {
		log.Fatalf("invalid: %s\n", validationErr.Detail)
	}

	return username, password
}

func promptUsername() string {
	var username string
	for {
		fmt.Print("Username: ")
		if _, err := fmt.Scanln(&username); err != nil {
			log.Fatalf("unable to scan admin username: %s", err)
		}

		validationErr := forms.ValidateAccountName(&username, "username")
		if validationErr == nil {
			break
		}

		fmt.Fprintf(os.Stderr, "invalid: %s\n", validationErr.Detail)
	}

	return username
}

func promptPassword(validate bool, prompt string) string {
	if prompt == "" {
		prompt = "Password"
	}

	var password string
	for {
		fmt.Printf("%s: ", prompt)
		passwordBytes, err := terminal.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			log.Fatalf("unable to scan password: %s", err)
		}
		fmt.Println()

		password = string(passwordBytes)

		if !validate {
			return password
		}

		if validationErr := forms.ValidatePassword(password, "password"); validationErr != nil {
			fmt.Fprintf(os.Stderr, "invalid: %s", validationErr.Detail)
			continue
		}

		fmt.Print("Verify Password: ")
		passwordBytes, err = terminal.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			log.Fatalf("unable to scan password: %s", err)
		}
		fmt.Println()

		if password == string(passwordBytes) {
			break
		}

		fmt.Fprintln(os.Stderr, "invalid: passwords do not match")
	}

	return password
}
