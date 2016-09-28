package dumpcerts

import (
	"archive/tar"
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"golang.org/x/net/context"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/docker/dhe-deploy"
	"github.com/docker/dhe-deploy/bootstrap"
	"github.com/docker/dhe-deploy/bootstrap/flags"
	"github.com/docker/dhe-deploy/shared/containers"
	"github.com/docker/dhe-deploy/shared/dtrutil"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/engine-api/types"
)

func dump(c *cli.Context) (int, error) {
	log.Debug("dump certs")
	if !bootstrap.IsPhase2() {
		return phase1(c)
	} else {
		return phase2(c)
	}

	return 0, nil
}

func phase1(c *cli.Context) (int, error) {
	buf, err := DumpCerts(c, "")
	if err != nil {
		return 1, err
	}
	os.Stdout.Write(*buf)
	return 0, nil
}

func DumpCerts(c *cli.Context, replicaID string) (*[]byte, error) {
	log.Debug("phase 1 starting")

	dropperOpts, err := bootstrap.ParseAndPromptDropperOpts(c)
	if err != nil {
		return nil, err
	}

	bs, err := bootstrap.GetBootstrapClient(c, dropperOpts)
	if err != nil {
		return nil, err
	}
	log.Debugf("bootstrap = %q", bs)

	// if the replicaID isn't set, try to get the ID from the command line,
	// and if that doesn't work, prompt for it
	if replicaID == "" {
		replica, err := bs.ExistingReplicaPicker(flags.ExistingReplicaID, "Choose a replica to dump the certs from", true)
		replicaID = replica.ReplicaID
		if err != nil {
			return nil, err
		}
	}
	bs.SetReplicaID(replicaID)

	// XXX - make this more generic
	img := containers.ContainerConfig{
		Name:         deploy.BootstrapHelperContainerName,
		Image:        deploy.BootstrapRepo.TaggedName(),
		AttachStdout: true,
		AttachStderr: false,
		Environment: map[string]string{
			deploy.Phase2EnvVar: "true",
			// XXX - add node constraint
			"affinity:container=": containers.Etcd.ReplicaName(replicaID),
		},
		Tty: false,
		Volumes: []containers.Volume{
			containers.CAVolume,
		},
		Entrypoint: []string{"dtr", "dumpcerts"},
	}

	log.Debugf("image = %q", img)

	cResp, err := bs.ContainerCreateFromContainerConfig(img)
	if err != nil {
		return nil, fmt.Errorf("Problem running container '%s' from image '%s': %s", img.Name, img.Image, err)
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
		return nil, fmt.Errorf("Couldn't attach to phase2 container: %s", err)
	}

	if err = bs.ContainerStart(cId); err != nil {
		return nil, fmt.Errorf("There was a problem running the dtr-phase2 container: %s", err)
	}

	buf := new(bytes.Buffer)

	//attachErrorChan := make(chan error)
	go func() {
		defer aResp.Close()
		stdcopy.StdCopy(buf, os.Stderr, aResp.Reader)
		//attachErrorChan <- err
	}()

	rc, err := bs.ContainerWait(cId)
	if err != nil {
		log.Errorf("Error waiting for phase 2 to finish: %s, code: %d", err, rc)
		return nil, err
	}

	strBuf := buf.String()
	output, err := base64.StdEncoding.DecodeString(strBuf)
	if err != nil {
		log.Errorf("Couldn't decode phase2 data: %s; data: %s", err, strBuf)
		return nil, err
	}

	err = bs.ContainerRemove(cId, types.ContainerRemoveOptions{})
	if err != nil {
		log.Errorf("Error removing the phase 2 container: %s", err)
		return nil, err
	}

	return &output, nil
}

func phase2(c *cli.Context) (int, error) {
	log.Debug("phase 2 starting")

	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)

	// XXX - dynamically generate the file names
	files := []string{
		containers.EtcdCACertStore.CertPath(),
		containers.EtcdCACertStore.KeyPath(),
		containers.EtcdCertStore.CertPath(),
		containers.EtcdCertStore.KeyPath(),
		containers.RethinkCACertStore.CertPath(),
		containers.RethinkCACertStore.KeyPath(),
		containers.RethinkCertStore.CertPath(),
		containers.RethinkCertStore.KeyPath(),
	}

	for _, file := range files {
		body, err := ioutil.ReadFile(file)
		if err != nil {
			log.Errorf("Couldn't read file '%s': %s", file, err)
			return 1, err
		}
		file := strings.TrimPrefix(file, "/")
		if err := dtrutil.AddBytesToTar(tw, body, file); err != nil {
			log.Errorf("Couldn't write %s tar: %s", file, err)
			return 1, err
		}
	}
	if err := tw.Close(); err != nil {
		log.Errorf("Couldn't close tar file: %s", err)
	}

	os.Stdout.Write([]byte(base64.StdEncoding.EncodeToString(buf.Bytes())))

	return 0, nil
}

func Run(c *cli.Context) {
	bootstrap.ConfigureLogging()
	rc, err := dump(c)
	if err != nil {
		log.Error(err)
	}
	os.Exit(rc)
}
