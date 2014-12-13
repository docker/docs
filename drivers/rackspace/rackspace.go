package rackspace

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/docker/machine/drivers"
	"github.com/docker/machine/drivers/openstack"
)

// Driver is a machine driver for Rackspace. It's a specialization of the generic OpenStack one.
type Driver struct {
	*openstack.Driver

	APIKey string
}

// CreateFlags stores the command-line arguments given to "machine create".
type CreateFlags struct {
	Username     *string
	APIKey       *string
	Region       *string
	MachineName  *string
	EndpointType *string
	ImageID      *string
	FlavorID     *string
	SSHUser      *string
	SSHPort      *int
}

func init() {
	drivers.Register("rackspace", &drivers.RegisteredDriver{
		New:            NewDriver,
		GetCreateFlags: GetCreateFlags,
	})
}

// GetCreateFlags registers the "machine create" flags recognized by this driver, including
// their help text and defaults.
func GetCreateFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			EnvVar: "OS_USERNAME",
			Name:   "rackspace-username",
			Usage:  "Rackspace account username",
			Value:  "",
		},
		cli.StringFlag{
			EnvVar: "OS_API_KEY",
			Name:   "rackspace-api-key",
			Usage:  "Rackspace API key",
			Value:  "",
		},
		cli.StringFlag{
			EnvVar: "OS_REGION_NAME",
			Name:   "rackspace-region",
			Usage:  "Rackspace region name",
			Value:  "",
		},
		cli.StringFlag{
			EnvVar: "OS_ENDPOINT_TYPE",
			Name:   "rackspace-endpoint-type",
			Usage:  "Rackspace endpoint type (adminURL, internalURL or the default publicURL)",
			Value:  "publicURL",
		},
		cli.StringFlag{
			Name:  "rackspace-image-id",
			Usage: "Rackspace image ID. Default: Ubuntu 14.10 (Utopic Unicorn) (PVHVM)",
			Value: "",
		},
		cli.StringFlag{
			Name:  "rackspace-flavor-id",
			Usage: "Rackspace flavor ID. Default: General Purpose 1GB",
			Value: "general1-1",
		},
		cli.StringFlag{
			Name:  "rackspace-ssh-user",
			Usage: "SSH user for the newly booted machine. Set to root by default",
			Value: "root",
		},
		cli.IntFlag{
			Name:  "rackspace-ssh-port",
			Usage: "SSH port for the newly booted machine. Set to 22 by default",
			Value: 22,
		},
	}
}

// NewDriver instantiates a Rackspace driver.
func NewDriver(storePath string) (drivers.Driver, error) {
	log.WithFields(log.Fields{
		"storePath": storePath,
	}).Debug("Instantiating Rackspace driver.")

	client := &Client{}
	inner, err := openstack.NewDerivedDriver(storePath, client)
	if err != nil {
		return nil, err
	}

	driver := &Driver{Driver: inner}
	client.driver = driver
	return driver, nil
}

// DriverName is the user-visible name of this driver.
func (d *Driver) DriverName() string {
	return "rackspace"
}

// SetConfigFromFlags assigns and verifies the command-line arguments presented to the driver.
func (d *Driver) SetConfigFromFlags(flags drivers.DriverOptions) error {
	d.Username = flags.String("rackspace-username")
	d.APIKey = flags.String("rackspace-api-key")
	d.Region = flags.String("rackspace-region")
	d.EndpointType = flags.String("rackspace-endpoint-type")
	d.ImageId = flags.String("rackspace-image-id")
	d.FlavorId = flags.String("rackspace-flavor-id")
	d.SSHUser = flags.String("rackspace-ssh-user")
	d.SSHPort = flags.Int("rackspace-ssh-port")
	return nil
}

func missingEnvOrOption(setting, envVar, opt string) error {
	return fmt.Errorf(
		"%s must be specified either using the environment variable %s or the CLI option %s",
		setting,
		envVar,
		opt,
	)
}

func (d *Driver) checkConfig() error {
	if d.Username == "" {
		return missingEnvOrOption("Username", "OS_USERNAME", "--rackspace-username")
	}
	if d.APIKey == "" {
		return missingEnvOrOption("API key", "OS_API_KEY", "--rackspace-api-key")
	}

	if d.ImageId == "" {
		// Default to the Ubuntu 14.10 image.
		// This is done here, rather than in the option registration, to keep the default value
		// from making "machine create --help" ugly.
		d.ImageId = "0766e5df-d60a-4100-ae8c-07f27ec0148f"
	}

	if d.EndpointType != "publicURL" && d.EndpointType != "adminURL" && d.EndpointType != "internalURL" {
		return fmt.Errorf(`Invalid endpoint type "%s". Endpoint type must be publicURL, adminURL or internalURL.`, d.EndpointType)
	}

	return nil
}
