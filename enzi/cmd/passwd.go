package cmd

import (
	"crypto/tls"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/docker/orca/enzi/passwords"
	"github.com/docker/orca/enzi/schema"
)

// Passwd is the command for changing a user's password.
var Passwd = cli.Command{
	Name:  "passwd",
	Usage: "Change password of a User Account",
	UsageText: `Change password of a User Account with the provided username and password.

Use the environment variables USERNAME and PASSWORD or set the --interactive flag to prompt for input on the command line.`,
	Action: runPasswd,
}

func init() {
	Passwd.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "interactive, i",
			Usage:       "prompt for username and password on the command line",
			Destination: &interactive,
		},
	}
}

func runPasswd(*cli.Context) error {
	tlsConfig := GetTLSConfig(tls.NoClientCert)

	username, password := getUsernamePassword()

	log.Debug("connecting to db ...")
	dbSession := GetDBSession(tlsConfig)
	defer dbSession.Close()

	passwordHash, err := passwords.HashPassword(password)
	if err != nil {
		log.Fatalf("unable to hash password: %s", err)
	}

	mgr := schema.NewRethinkDBManager(dbSession)

	user, err := mgr.GetUserByName(username)
	if err != nil {
		log.Fatalf("unable to get user: %s", err)
	}

	isActive := true // Always activate the user account.
	updateFields := schema.AccountUpdateFields{
		IsActive:     &isActive,
		PasswordHash: &passwordHash,
	}

	log.Debug("updating user account...")
	if err := mgr.UpdateAccount(user.ID, updateFields); err != nil {
		log.Fatalf("unable to update user: %s", err)
	}

	// Invalidate all sessions for this user since the password has been
	// changed.
	if err := mgr.DeleteSessionsForUser(user.ID, ""); err != nil {
		log.Fatalf("unable to invalidate login sessions for user: %s", err)
	}

	log.Infof("successfully set user account password: %s", user.Name)

	return nil
}
