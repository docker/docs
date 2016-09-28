package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
	cfssllog "github.com/cloudflare/cfssl/log"
	"github.com/codegangsta/cli"
	"github.com/howeyc/gopass"

	orcaconfig "github.com/docker/orca/config"
)

// Handle the common args that all commands use
func HandleGlobalArgs(c *cli.Context) {
	// Quiet down the cfssl logging
	cfssllog.Level = cfssllog.LevelWarning
	if c.Bool("debug") {
		log.SetLevel(log.DebugLevel)
	}
	if c.Bool("jsonlog") {
		log.SetFormatter(&log.JSONFormatter{})
	}
	if c.String("host-address") != "" {
		OrcaLocalName = c.String("host-address")
		os.Setenv("UCP_HOST_ADDRESS", OrcaLocalName)
	}
	OrcaSANs = c.StringSlice("san")
	OrcaInstanceID = os.Getenv("UCP_BOOTSTRAP_PHASE2")
	InPhase2 = OrcaInstanceID != ""
	if OrcaInstanceKey == "" {
		OrcaInstanceKey = os.Getenv("UCP_INSTANCE_KEY")
	}

	// Optional DNS settings
	DNS = c.StringSlice("dns")
	DNSOpt = c.StringSlice("dns-opt")
	DNSSearch = c.StringSlice("dns-search")
	pullBehavior := c.String("pull")
	if pullBehavior != "" {
		switch pullBehavior {
		case "always", "missing", "never":
			PullBehavior = pullBehavior
		default:
			log.Fatalf("Invalid '--pull' argument %s - must be 'always', 'missing', or 'never'", pullBehavior)
		}
	}

	ver := c.String("image-version")
	if ver != "" {
		orcaconfig.ImageVersion = ver
	}

	// Kinda lame that codegangsta doesn't automatically set the env too, but we need that in some cases
	os.Setenv("UCP_ADMIN_USER", c.String("admin-username"))
	os.Setenv("UCP_ADMIN_PASSWORD", c.String("admin-password"))
	os.Setenv("REGISTRY_USERNAME", c.String("registry-username"))
	os.Setenv("REGISTRY_PASSWORD", c.String("registry-password"))
	ucpURL := c.String("url")
	if ucpURL != "" {
		os.Setenv("UCP_URL", ucpURL)
	}
	fingerprint := c.String("fingerprint")
	if fingerprint != "" {
		os.Setenv("UCP_FINGERPRINT", fingerprint)
	}
}

// Gather the specified environment setting interactively
func InteractivePrompt(envVar string) error {
	s := InteractiveArgs[envVar]
	defValue := os.Getenv(envVar)
	defPrompt := ""
	if defValue != "" {
		if strings.Contains(envVar, "PASSWORD") {
			defPrompt = " (" + strings.Repeat("*", len(defValue)) + ")"
		} else {
			defPrompt = " (" + defValue + ")"
		}
	}
	fmt.Printf("%s%s: ", s.Prompt, defPrompt)

	if s.Echo {
		reader := bufio.NewReader(os.Stdin)
		value, err := reader.ReadString('\n')
		if err != nil {
			log.Debugf("Failed to read input: %s", err)
			return err
		}
		if strings.TrimSpace(value) == "" && defValue != "" {
			value = defValue
		}
		os.Setenv(envVar, strings.TrimSpace(value))
	} else {
		os.Setenv(envVar, string(gopass.GetPasswd()))
	}
	return nil
}
