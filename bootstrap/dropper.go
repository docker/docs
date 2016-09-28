package bootstrap

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"golang.org/x/net/context"

	"github.com/docker/dhe-deploy"
	"github.com/docker/dhe-deploy/bootstrap/flags"
	"github.com/docker/dhe-deploy/bootstrap/ucpclient"
	"github.com/docker/dhe-deploy/shared/containers"
	"github.com/docker/dhe-deploy/shared/dtrutil"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/strslice"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

type DropperOpts struct {
	Host          string
	AdminUsername string
	AdminPassword string
	Node          string
	NoUCP         bool
	UseBundle     bool
	InsecureTLS   bool
	UCPTLSCA      string
	ExtraEnvs     string
}

func splitOnAnd(str string) []string {
	tmp := strings.Split(str, "&")
	res := []string{}
	accum := ""
	for _, e := range tmp {
		accum += e
		if len(e) == 0 || (len(e) > 0 && e[len(e)-1] != '\\') {
			res = append(res, accum)
			accum = ""
		} else {
			accum = accum[:len(accum)-1] + "&"
		}
	}
	return res
}

func (d DropperOpts) MakeContainerConfig(args []string) *containers.ContainerConfig {
	// TODO: randomly generate phase 2 container name, consider including the command being executed?
	container := containers.ContainerConfig{
		Name:         deploy.BootstrapPhase2ContainerName,
		Image:        deploy.BootstrapRepo.TaggedName(),
		AttachStdout: true,
		AttachStderr: true,
		Environment: map[string]string{
			flags.UCPHostFlag.EnvVar:        d.Host,
			flags.UsernameFlag.EnvVar:       d.AdminUsername,
			flags.PasswordFlag.EnvVar:       d.AdminPassword,
			flags.UCPInsecureTLSFlag.EnvVar: strconv.FormatBool(d.InsecureTLS),
			flags.UCPCAFlag.EnvVar:          d.UCPTLSCA,
			flags.UseBundleFlag.EnvVar:      strconv.FormatBool(d.UseBundle),
			deploy.Phase2EnvVar:             "true",
			flags.UCPNodeFlag.EnvVar:        d.Node,
			flags.NoUCPFlag.EnvVar:          strconv.FormatBool(d.NoUCP),
		},
		Tty: true,
		Volumes: []containers.Volume{
			containers.CAVolume,
			containers.NotaryVolume,
		},
		Entrypoint: strslice.StrSlice(append([]string{"dtr"}, args...)),
	}
	// If the user has already specified a node constraint, start the bootstrapper
	// on that node
	if d.Node != "" {
		container.Environment["constraint:node="] = d.Node
	}
	if d.ExtraEnvs != "" {
		for _, env := range splitOnAnd(d.ExtraEnvs) {
			container.Constraints = append(container.Constraints, env)
		}
	}
	// This is necessary for DaemonVersionHack to work. TODO: remove this in a future version
	// when we stop caring about checking the docker version on UCP 1.1.0
	//if d.NoUCP {
	container.DumbVolumes = append(
		container.DumbVolumes,
		fmt.Sprintf("%s:%s", deploy.DockerSocketPath, deploy.DockerSocketPath),
		fmt.Sprintf("%s:%s", deploy.LogsCertPathInHost, deploy.LogsCertPathInContainer),
	)
	return &container
}

func Phase2Execute(img *containers.ContainerConfig, bs Bootstrap) error {
	if bs.GetReplicaID() == "" {
		return fmt.Errorf("attempt to execute phase 2 dropper without a replica id")
	}

	cResp, err := bs.ContainerCreateFromContainerConfig(*img)
	if err != nil {
		return fmt.Errorf("Problem running container '%s' from image '%s': %s", img.Name, img.Image, err)
	}
	cId := cResp.ID

	aResp, err := bs.GetDockerClient().ContainerAttach(context.Background(), cId, types.ContainerAttachOptions{
		Stream: true,
		Stdin:  false,
		Stdout: true,
		Stderr: false,
		//DetachKeys: "",
	})

	if err != nil {
		// XXX - fixme
		return fmt.Errorf("Couldn't attach to phase2 container: %s", err)
	}

	if err = bs.ContainerStart(cId); err != nil {
		return fmt.Errorf("There was a problem running the dtr-phase2 container: %s", err)
	}

	attachErrorChan := make(chan error, 1)
	go func() {
		defer aResp.Close()
		io.Copy(os.Stdout, aResp.Reader)
		attachErrorChan <- err
	}()

	select {
	case err = <-attachErrorChan:
		if err != nil {
			return err
		}
	}

	rc, err := bs.ContainerWait(cId)
	if err != nil {
		log.Errorf("Error waiting for phase 2 to finish: %s, code: %d", err, rc)
		return err
	}

	err = bs.ContainerRemove(cId, types.ContainerRemoveOptions{})
	if err != nil {
		log.Errorf("Error removing the phase 2 container: %s", err)
		return err
	}

	if rc != 0 {
		return fmt.Errorf("Phase 2 returned non-zero status: %d", rc)
	}
	return nil
}

func ParseAndPromptDropperOpts(c *cli.Context) (*DropperOpts, error) {
	var dropperOpts DropperOpts

	if !flags.NoUCP {
		PromptIfNotSet(c, flags.UCPHostFlag)
	}

	PromptIfNotSet(c, flags.UsernameFlag, flags.PasswordFlag)
	dropperOpts.AdminUsername = flags.Username
	dropperOpts.AdminPassword = flags.Password

	if flags.NoUCP {
		dropperOpts.NoUCP = true
		return &dropperOpts, nil
	}

	dropperOpts.Host = flags.UCPHost()
	dropperOpts.Node = flags.UCPNode
	dropperOpts.ExtraEnvs = flags.ExtraEnvs
	dropperOpts.InsecureTLS = flags.UCPInsecureTLS
	dropperOpts.UCPTLSCA = flags.UCPCA
	dropperOpts.UseBundle = flags.UseBundle
	return &dropperOpts, nil
}

// TODO: make this not prompt for things
func GetBootstrapClient(c *cli.Context, dropperOpts *DropperOpts) (Bootstrap, error) {
	var bs Bootstrap

	hubUser := flags.HubUsername
	hubPassword := flags.HubPassword
	// if a username is provided without a password, prompt for it
	if hubUser != "" && hubPassword == "" {
		hubPassword = PromptPassword("Docker Hub password: ")
		// if we don't set this, it won't be passed on to stage 2
		flags.HubPassword = hubPassword
	}

	if !dropperOpts.NoUCP {
		err := UCPCertTest(flags.UCPHost(), flags.UCPInsecureTLS, flags.UCPCA)
		if err != nil {
			return nil, err
		}

		httpClient, err := dtrutil.HTTPClient(dropperOpts.InsecureTLS, dropperOpts.UCPTLSCA)
		if err != nil {
			return nil, err
		}

		ucp := ucpclient.New(dropperOpts.Host, httpClient)
		if err := ucp.Login(dropperOpts.AdminUsername, dropperOpts.AdminPassword); err != nil {
			if strings.Contains(err.Error(), "x509: certificate signed by unknown authority") {
				return nil, fmt.Errorf(fmt.Sprintf(`Certificate validation for UCP failed. You can get the UCP CA from https://%s/ca and then use it by appending --%s "$(cat ca.pem)" to this command.`, dropperOpts.Host, flags.UCPCAFlag.Name))
			}
			return nil, fmt.Errorf("Couldn't login: %s", err)
		}

		if flags.UseBundle {
			bundle, err := ucp.GetBundle()
			if err != nil {
				return nil, fmt.Errorf("Couldn't get bundle: %s", err)
			}
			defer func() {
				err := ucp.DeleteOwnBundle(string(bundle.CertPub))
				if err != nil {
					log.WithField("error", err).Warn("Failed to delete temporary UCP bundle.")
				}
				err = ucp.Logout()
				if err != nil {
					log.WithField("error", err).Warn("Failed to log out.")
				}
			}()
			err = ucp.LabelOwnBundle(string(bundle.CertPub), "Temporary Docker Trusted Registry installer bundle. (safe to delete)")
			if err != nil {
				log.WithField("error", err).Warn("Failed to re-label temporary UCP bundle.")
			}
			certPath := "/tmp/bundle"
			err = ucpclient.BundleToDisk(bundle, certPath)
			if err != nil {
				log.WithField("error", err).Error("failed to write out bundle to disk")
				return nil, err
			}

			bs, err = NewFromBundle(dropperOpts.Host, dropperOpts.AdminUsername, dropperOpts.AdminPassword, certPath, hubUser, hubPassword)
			if err != nil {
				return nil, fmt.Errorf("Couldn't get bootstrap: %s", err)
			}
		} else {
			bs, err = NewFromJWT(dropperOpts.Host, dropperOpts.AdminUsername, dropperOpts.AdminPassword, ucp.JWT(), httpClient, hubUser, hubPassword)
			if err != nil {
				return nil, fmt.Errorf("Couldn't get bootstrap: %s", err)
			}
			/*
				// XXX - make this as part of Close()
				defer func() {
					err := ucp.Logout()
					if err != nil {
						log.WithField("error", err).Warn("Failed to log out.")
					}
				}()
			*/
		}
		bs.SetUCPClient(ucp)
	} else {
		var err error
		bs, err = NewFromSocket(hubUser, hubPassword)
		if err != nil {
			return nil, fmt.Errorf("Couldn't get bootstrap: %s", err)
		}
	}

	err := bs.VersionCheck()
	if err != nil {
		log.Errorf("Docker version check failed: %s", err)
		return nil, err
	}

	return bs, nil
}
