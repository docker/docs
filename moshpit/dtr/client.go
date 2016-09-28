package dtr

import (
	"archive/tar"
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"hash/crc64"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/docker/dhe-deploy/pkg/moshpit-framework"
	"github.com/docker/dhe-deploy/pkg/moshpit-framework/commands"
	"github.com/docker/dhe-deploy/pkg/moshpit-framework/dockerutil"
	"github.com/docker/dhe-deploy/pkg/moshpit-framework/util"
	"github.com/docker/dhe-deploy/shared/dtrutil"
	"github.com/docker/distribution"
	"github.com/docker/distribution/digest"
	"github.com/docker/distribution/manifest/schema1"
	_ "github.com/docker/distribution/manifest/schema2"
	"github.com/docker/distribution/reference"
	rc "github.com/docker/distribution/registry/client"
	"github.com/docker/distribution/registry/client/auth"
	"github.com/docker/distribution/registry/client/transport"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/engine-api/client"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/container"
	"github.com/docker/engine-api/types/network"
	"github.com/docker/go-connections/nat"
	"github.com/docker/libtrust"
	"golang.org/x/net/context"
	"gopkg.in/yaml.v2"
)

type PushConfig struct {
	Method    string
	Namespace string
	RepoName  string

	TagName               string
	TagPattern            string
	PushesPerBatch        int
	Batches               int
	LayerMinBytes         int
	LayerMaxBytes         int
	LayerSizeScalePattern string
	ImageMinLayers        int
	ImageMaxLayers        int
	LayersScalePattern    string
	Retries               int
}

type PullConfig struct {
	Method     string
	Namespace  string
	RepoName   string
	TagName    string
	TagPattern string
}

// The mass-population config is a superset of the push config:
// It has a user/repo range that it will randomly populate within, instead of using a single namepsace/repository
// Mass-population will only run if UsersToPopulate > 0 and ReposToPopulate > 0,
// (note that for both of them being 1, you can just use a regular push and specify which namespace/repo you want to populate)
type MassPopulateConfig struct {
	Push *PushConfig

	// Requirement is that these users and repos exist and they need to follow a similar pattern to
	// the the creation of such in the setup phase (although they can have been generated differently,
	// for example the user pattern is identical to our big LDAP test server)
	// Mass-populate will randomly choose UsersToPopulate users and ReposToPopulate per user repos to populate with the
	// pushconfig settings (replacing the namespace/repoName in it)
	TotalNumUsers        int
	TotalNumReposPerUser int
	UsersToPopulate      int
	ReposToPopulate      int
}

type ClientConfig struct {
	DTRURL         string
	DTRCA          string
	DTRInsecureTLS bool
	Username       string
	Password       string
	RefreshToken   string
	Seed           int
	Push           *PushConfig
	Pull           *PullConfig
	MassPopulate   *MassPopulateConfig
}

type job struct {
	ctx        context.Context
	log        logrus.FieldLogger
	clientName string
	config     ClientConfig
	state      commands.JobState
	err        error
}

func (j *job) State() (commands.JobState, error) {
	return j.state, j.err
}

type credentialStore struct {
	username      string
	password      string
	refreshTokens map[string]string
}

func (tcs *credentialStore) Basic(url *url.URL) (string, string) {
	return tcs.username, tcs.password
}

// refresh tokens are the long lived tokens that can be used instead of a password
func (tcs *credentialStore) RefreshToken(u *url.URL, service string) string {
	fmt.Printf("request for refresh token at %s", service)
	return tcs.refreshTokens[service]
}

func (tcs *credentialStore) SetRefreshToken(u *url.URL, service string, token string) {
	if tcs.refreshTokens != nil {
		tcs.refreshTokens[service] = token
	}
}

func newRepoClient(ctx context.Context, registry, namespace, repo, tag, user, password, refreshToken string, httpClient *http.Client, insecure bool, ca string) (distribution.Repository, error) {
	// XXX: this is pretty aweful. we should handle the creation of a transport in one place
	caTransport, err := util.HTTPTransport(insecure, ca)
	if err != nil {
		return nil, err
	}

	challengeManager := auth.NewSimpleChallengeManager()

	req, err := http.NewRequest("GET", fmt.Sprintf("https://%s/v2/", registry), nil)
	if err != nil {
		return nil, err
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	challengeManager.AddResponse(resp)
	if err != nil {
		return nil, err
	}

	scope := fmt.Sprintf("%s/%s", namespace, repo)
	creds := &credentialStore{}
	if refreshToken != "" {
		creds.refreshTokens = map[string]string{strings.Split(registry, ":")[0]: refreshToken}
	} else {
		creds.username = user
		creds.password = password
	}
	th := auth.NewTokenHandler(caTransport, creds, scope, "pull", "push")
	trans := transport.NewTransport(caTransport, auth.NewAuthorizer(challengeManager, th))
	baseURL := fmt.Sprintf("https://%s/", registry)
	ref, err := reference.WithName(path.Join(namespace, repo))
	if err != nil {
		return nil, err
	}
	repository, err := rc.NewRepository(ctx, ref, baseURL, trans)
	if err != nil {
		return nil, err
	}
	return repository, nil
}

var dummyPrivateKey = []byte(`
-----BEGIN RSA PRIVATE KEY-----
MIIEpQIBAAKCAQEA4CUp7erNTs3B4db3nG0ot1g6vO1jYZYMZx9tPUg/ssryUb+v
fxjyekxCewKxc1oN7aaIHn3nSLLY2c4kZpoKfkJpharPcXPtFXs0+/mvUgccEHSu
RRiW16mNosx0JBsgIvs4xgZHulZF5mmBXGQ7efDCCLKagchsj2gIbQG+7xlsaJNs
EXcl7BcBBP3P0brFs21CHk6USpgbiIyazExcAHfyMAI/3OYwSbDFxxJcxY/NmSPv
uLouoikWzi280C7lBkIjPgRg7E/hCNNfxdFwQPXhvZsJ/S7sjlXcmGqq176BPpID
goJozuEcELM5kJlPey0kX3B8rbJIrcbHMAMN+wIDAQABAoIBAG95zDqpdmZk6rI5
SXigyYk19jCUF8Mm7xAyjw/VaOixCochwFSDwcSVPNMU6dAmz5qMIlKX53k+iZ75
aR0mK2XT/csewoD4WMdAOX/AFDPFmW2NukZfDlY/21NGP9TdMMR1ES1bXj0MP0Ny
4YVjzDi/RqEwwqsdVPCVmusr3RvoKLpHxVIJnJ/5B2VWWryDwCXnBPN/r+2pbA1J
IORZeqA6PGrS6URe7OB9j6R6VkmAzYF6EgLZWeRUxf2Waa+1rqc3KHpvgID5sWP7
SzsvQOFYvIHegebPrq+wCtWoNMpRjkKKJT5mIzV224PByDsVaJXbsiqNlIgye8si
ei3zkHECgYEA7J04HadVBllbkY5lVBME5gbDNHtyK22vXJFKMJSu4JwE8+eGLdh5
xkyyRgfD6bHBz9SMx/uxRlNOgUUzNpsH3joOwAR+uK7GDdktgi33U4r0dyIECckj
TXmBuvWbs+qXb7bVUsmsBCTrbXpyYKu0xcmRIAQWlWR2u4MyHQtEjsUCgYEA8oJr
KMNVhBw7CjAxR02vQVojiRGdsL/up5IKdcxm/2Bz1iPLidzGIXOOy6JjLbnZh28r
4vGtvuc0WuyVHc79aYCU+H09wJ24LXcFhTK7rzf22IfOE4euJQQoAp09UAk8jZQs
yYl7QY8I8tHHK4N40/AARJYqXHz/5BV/drHy9b8CgYEAm8NX+LV6VIaosCaEcBdh
JyiWgsstOoenZJHEvDx07ynmXMYyX6XdbHx6830TLJm3U+DBmLkJV5lp2dG7SBxA
zrt7kE6lOWDcwqsQuV0XLykPjAmZjPObSNpPW8tp58PsUz+SKUDX+5ZuYZC1EQyY
IYhzABeQ4mHTg9d3OwV66V0CgYEA8F7i1lme5r6Qqo2QGqvefXlZ5Z/HXI0xgXjY
02AR6yjwSB3cvj5NSJTgweioQ4eGHJ7Nsjl4zNMgass7Fnu3ZJ5lilOhJM1v4+io
WRkrPQbMrl0VnvgKXXhcLBMs1asCERcAuZaCzD15Ui0qLHA5EGE/8ruhK2FexfWl
DMJfHsMCgYEA29gjRBDBHdeC4RMldg4sxeNCpPlSVGSGFmqS0uLD7jHibbOKbKOd
gCPaW1zf/R9+Q6cnM6gYJHELOMKl/Tbpv0Hp3uYfmDsmbTw6CIVe/VDZ1yHiqGNW
82OwSD2TzFb91OSjnIHH0OCwSLGccgcJuE0XU0kdqV6kLINKspKezyY=
-----END RSA PRIVATE KEY-----`)

// these 4 types were copied out of registry to build the smallest possible valid configJSON
type diffID digest.Digest
type imageRootFS struct {
	Type      string   `json:"type"`
	DiffIDs   []diffID `json:"diff_ids,omitempty"`
	BaseLayer string   `json:"base_layer,omitempty"`
}
type imageHistory struct {
	Created    time.Time `json:"created"`
	Author     string    `json:"author,omitempty"`
	CreatedBy  string    `json:"created_by,omitempty"`
	Comment    string    `json:"comment,omitempty"`
	EmptyLayer bool      `json:"empty_layer,omitempty"`
}
type imageConfig struct {
	RootFS       *imageRootFS   `json:"rootfs,omitempty"`
	History      []imageHistory `json:"history,omitempty"`
	Architecture string         `json:"architecture,omitempty"`
}
type layer struct {
	BlobSum diffID `json:"blobSum"`
}
type pullImageConfig struct {
	FSLayers []layer `json:"fsLayers"`
}

func dummyConfigJSON(numLayers int) []byte {
	dummyConfig := imageConfig{History: make([]imageHistory, numLayers), RootFS: &imageRootFS{DiffIDs: make([]diffID, numLayers)}}
	bytes, _ := json.Marshal(&dummyConfig)
	return bytes
}

// returns the connection string for the daemon
func (j *job) BuildDaemon(hostClient client.APIClient, name, registry string, port int) (string, error) {
	// XXX: we are exposing an unsecured docker daemon to the internet... make this less terrible by using custom networks?

	internalPort := "2375"

	resp, err := hostClient.ContainerCreate(j.ctx, &container.Config{
		Cmd:          []string{"--debug", "--insecure-registry", registry},
		Image:        "docker:dind",
		ExposedPorts: map[nat.Port]struct{}{nat.Port(internalPort): {}},
	}, &container.HostConfig{
		Privileged: true,
		PortBindings: nat.PortMap{
			nat.Port(internalPort): []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: fmt.Sprintf("%d", port)}},
		},
	}, &network.NetworkingConfig{}, name)
	if err != nil {
		return "", err
	}
	id := resp.ID

	err = hostClient.ContainerStart(j.ctx, id)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("172.17.0.1:%d", port), nil
}

func (j *job) PushImage(client client.APIClient, registry, namespace, repo, tag, user, password, refreshToken string, httpClient *http.Client, localRand *rand.Rand, layerSize int64, imageNumLayers int, insecure bool, ca string) error {
	if j.config.Push.Method == "dind-daemon" {
		// TODO: move all this docker client setup outside of PushImage because PushImage is potentially called in a loop
		err := j.PushImageWithDockerClient(client, registry, namespace, repo, tag, user, password, refreshToken, httpClient, localRand, layerSize, imageNumLayers, insecure, ca)
		if err != nil {
			return err
		}
	} else if j.config.Pull.Method == "" || j.config.Pull.Method == "registry" {
		err := j.PushImageWithRegistryClient(registry, namespace, repo, tag, user, password, refreshToken, httpClient, localRand, layerSize, imageNumLayers, insecure, ca)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("unknown push method")
	}
	return nil
}

func (j *job) PushImageWithDockerClient(client client.APIClient, registry, namespace, repo, tag, user, password, refreshToken string, httpClient *http.Client, localRand *rand.Rand, layerSize int64, imageNumLayers int, insecure bool, ca string) error {
	imageName := fmt.Sprintf("%s/%s/%s:%s", registry, namespace, repo, tag)
	dockerfile := "FROM alpine:latest\n"
	for i := 0; i < imageNumLayers; i++ {
		// the size is an approximation because dd is awkward and I'm lazy
		dockerfile += fmt.Sprintf("RUN dd bs=512 count=%d if=/dev/urandom of=/%d\n", layerSize/512+1, i)
	}
	buf := new(bytes.Buffer)
	tarWriter := tar.NewWriter(buf)
	tarWriter.WriteHeader(&tar.Header{
		Name:     "Dockerfile",
		Typeflag: tar.TypeReg,
		Mode:     0600,
		Size:     int64(len(dockerfile)),
	})
	_, err := tarWriter.Write([]byte(dockerfile))
	if err != nil {
		return err
	}
	img, err := client.ImageBuild(j.ctx, types.ImageBuildOptions{
		Tags:        []string{imageName},
		Context:     bytes.NewReader(buf.Bytes()),
		ForceRemove: true,
		NoCache:     true,
	})
	if err != nil {
		return err
	}
	j.log.Info("Docker build output:")
	decoder := json.NewDecoder(img.Body)
	for err == nil {
		response := jsonmessage.JSONMessage{}
		err = decoder.Decode(&response)
		if response.Error == nil {
			j.log.Debugf("%s (%s)", response.Stream, response.Status)
		} else {
			return fmt.Errorf("Daemon error: %s", response.Error)
		}
	}
	reader, err := client.ImagePush(j.ctx, types.ImagePushOptions{
		// because of ambiguities in the registry api we can put the tag in the imageID
		ImageID:      imageName,
		RegistryAuth: dockerutil.MakeRegistryAuth(user, password, refreshToken),
	}, nil)

	j.log.Info("Docker push output:")
	decoder = json.NewDecoder(reader)
	for err == nil {
		response := jsonmessage.JSONMessage{}
		err = decoder.Decode(&response)
		if response.Error == nil {
			j.log.Debugf("%s (%s)", response.Stream, response.Status)
		} else {
			return fmt.Errorf("Daemon error: %s", response.Error)
		}
	}
	if err != io.EOF {
		return fmt.Errorf("Failed to parse response: %s", err)
	}

	return nil
}

func (j *job) PushImageWithRegistryClient(registry, namespace, repo, tag, user, password, refreshToken string, httpClient *http.Client, localRand *rand.Rand, layerSize int64, imageNumLayers int, insecure bool, ca string) error {
	repoClient, err := newRepoClient(j.ctx, registry, namespace, repo, tag, user, password, refreshToken, httpClient, insecure, ca)
	if err != nil {
		return err
	}
	blobs := repoClient.Blobs(j.ctx)

	descriptors := []distribution.Descriptor{}
	// TODO: parallelize
	for i := 0; i < imageNumLayers; i++ {
		// ** start layer upload
		// This pipe contraption looks like this:
		// localRand -> limited reader -> tar -> bufferedWriter -> multiwriter  -> blobWriter => actual descriptor
		//                                                             |
		//                                                             v
		//                                                       digest reader => expected descriptor
		//
		// The multiwriter part is because: "Depending on the implementation, written data may be validated
		// against the provisional descriptor fields." - registry comment
		blobWriter, err := blobs.Create(j.ctx)
		if err != nil {
			return err
		}

		digester := digest.Canonical.New()
		writer := io.MultiWriter(digester.Hash(), blobWriter)

		// The buffered writer is necessary because otherwise the multi-writer has no mechanism to handle
		// short writes, which causes the error from short writes to propagate all the way back up and
		// stop io.Copy from writing out the stream
		bufferedWriter := bufio.NewWriterSize(writer, 5000000)

		tw := tar.NewWriter(bufferedWriter)

		// TODO: maybe check if the blob it exists first? but that would require either checkpointing the rng
		// or storing the layer in memory

		header := &tar.Header{
			Name:     fmt.Sprintf("/filler"),
			Typeflag: tar.TypeReg,
			Mode:     0600,
			Size:     layerSize,
		}
		if err := tw.WriteHeader(header); err != nil {
			_ = tw.Close()
			j.log.WithField("err", err).Warn("failed to write header into tar writer")
			return err
		}
		n, err := io.Copy(tw, &io.LimitedReader{R: localRand, N: layerSize})
		if err != nil {
			_ = tw.Close()
			j.log.WithField("err", err).WithField("n", n).Warn("failed to io.Copy localRand into tar writer")
			return err
		}
		err = tw.Close()
		if err != nil {
			j.log.WithField("err", err).Warn("failed to close tar writer")
			return err
		}

		descriptor, err := blobWriter.Commit(j.ctx, distribution.Descriptor{Digest: digester.Digest()})
		if err != nil {
			return fmt.Errorf("error pushing blob: %s", err)
		}

		descriptors = append(descriptors, descriptor)
		j.log.WithField("descriptor", descriptor).Info("pushed blob")
	}

	// ** start manifest upload
	manifests, err := repoClient.Manifests(j.ctx)
	if err != nil {
		return err
	}

	configJSON := dummyConfigJSON(imageNumLayers)
	ref, err := reference.WithName(path.Join(namespace, repo))
	if err != nil {
		return err
	}
	taggedRef, err := reference.WithTag(ref, tag)
	if err != nil {
		return err
	}
	privateKey, err := libtrust.UnmarshalPrivateKeyPEM(dummyPrivateKey)
	if err != nil {
		return err
	}
	mb := schema1.NewConfigManifestBuilder(blobs, privateKey, taggedRef, configJSON)
	// DTR doesn't support schema v2 yet :(
	//mb := schema2.NewManifestBuilder(blobs, configJSON)
	for _, descriptor := range descriptors {
		err = mb.AppendReference(descriptor)
		if err != nil {
			return err
		}
	}
	manifest, err := mb.Build(j.ctx)
	if err != nil {
		return err
	}
	digest, err := manifests.Put(j.ctx, manifest, distribution.WithTag(tag))
	if err != nil {
		return err
	}
	j.log.WithField("digest", digest).Info("pushed manifest")
	return nil
}

// TODO: make this more of a Puller interface with multiple implementations :/
func (j *job) PullImage(httpClient *http.Client, registry, namespace, repo, tag, user, password, refreshToken string, insecure bool, ca string) error {
	if j.config.Pull.Method == "hostdaemon" {
		// TODO: instead of using the docker socket directly, spawn a dind to use for pushing and pulling (with the right --insecure-registry)
		client, err := dockerutil.MakeDockerClient("unix:///var/run/docker.sock", "", nil, j.config.Username, j.config.Password, j.config.RefreshToken)
		if err != nil {
			return err
		}
		err = j.PullImageWithDockerClient(client, registry, namespace, repo, tag, user, password, refreshToken)
		if err != nil {
			return err
		}
	} else if j.config.Pull.Method == "" || j.config.Pull.Method == "registry" {
		err := j.PullImageWithRegistryClient(httpClient, registry, namespace, repo, tag, user, password, refreshToken, insecure, ca)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("unknown pull method")
	}
	return nil
}

func (j *job) PullImageWithRegistryClient(httpClient *http.Client, registry, namespace, repo, tag, user, password, refreshToken string, insecure bool, ca string) error {
	repoClient, err := newRepoClient(j.ctx, registry, namespace, repo, tag, user, password, refreshToken, httpClient, insecure, ca)
	if err != nil {
		return err
	}
	tags := repoClient.Tags(j.ctx)
	descriptor, err := tags.Get(j.ctx, tag)
	if err != nil {
		return fmt.Errorf("error getting tag %s", err)
	}
	manifests, err := repoClient.Manifests(j.ctx)
	if err != nil {
		return err
	}
	manifest, err := manifests.Get(j.ctx, descriptor.Digest)
	if err != nil {
		return fmt.Errorf("error getting manifest %s", err)
	}
	// We assume v1 TODO: support v2
	_, payload, err := manifest.Payload()
	imgCfg := pullImageConfig{}
	err = json.Unmarshal(payload, &imgCfg)
	if err != nil {
		return err
	}
	blobs := repoClient.Blobs(j.ctx)
	// TODO: parallel pull
	for _, layer := range imgCfg.FSLayers {
		digest := digest.Digest(layer.BlobSum)
		reader, err := blobs.Open(j.ctx, digest)
		if err != nil {
			return err
		}
		// we discard the data after fetching it
		_, err = io.Copy(ioutil.Discard, reader)
		if err != nil {
			return err
		}
	}
	return nil
}

func (j *job) PullImageWithDockerClient(dc client.APIClient, registry, namespace, repo, tag, user, password, refreshToken string) error {
	imageID := ""
	if registry != "" {
		imageID = fmt.Sprintf("%s/%s/%s:%s", registry, namespace, repo, tag)
	} else {
		if namespace != "" {
			imageID = fmt.Sprintf("%s/%s:%s", namespace, repo, tag)
		} else {
			imageID = fmt.Sprintf("%s:%s", repo, tag)
		}
	}
	j.log.WithField("imageID", imageID).Info("pulling image")
	options := types.ImagePullOptions{
		// because of ambiguities in the registry api we can put the tag in the imageID
		ImageID:      imageID,
		RegistryAuth: dockerutil.MakeRegistryAuth(user, password, refreshToken),
	}

	buf, err := dc.ImagePull(j.ctx, options, nil)
	if err != nil {
		return fmt.Errorf("Failed to pull image: %s", err)
	}
	defer buf.Close()
	decoder := json.NewDecoder(buf)
	for err == nil {
		response := jsonmessage.JSONMessage{}
		err = decoder.Decode(&response)
		if response.Error == nil {
			j.log.WithField("status", response.Status).Debug("pulling status")
		} else {
			return fmt.Errorf("Daemon error pulling: %s", response.Error)
		}
	}
	if err != io.EOF {
		return fmt.Errorf("Failed to parse response: %s", err)
	}

	j.log.Infof("pulled %s", imageID)
	return nil
}

func (j *job) executePush(config *PushConfig, hostClient *client.Client, rootRand *rand.Rand) error {
	tag := config.TagName
	if config.TagPattern == "random" {
		tag = fmt.Sprintf("%s-%010d", tag, rootRand.Int63n(1000000))
	}
	// we create a new source for each thread based on a seed
	// that way the layers created from each thread are the same every run
	parallelRands := make([]*rand.Rand, config.PushesPerBatch)
	for n := 0; n < config.PushesPerBatch; n++ {
		parallelRands[n] = rand.New(rand.NewSource(rootRand.Int63()))
	}

	parallelDinds := make([]client.APIClient, config.PushesPerBatch)
	if config.Method == "dind-daemon" {
		for n := 0; n < config.PushesPerBatch; n++ {
			daemonPort := parallelRands[n].Intn(10000) + 30000
			name := fmt.Sprintf("moshpit-docker-daemon-%d", daemonPort)

			_ = hostClient.ContainerRemove(j.ctx, types.ContainerRemoveOptions{ContainerID: name, Force: true, RemoveVolumes: true})
			connStr, err := j.BuildDaemon(hostClient, name, j.config.DTRURL, daemonPort)
			if err != nil {
				return err
			}
			defer hostClient.ContainerRemove(j.ctx, types.ContainerRemoveOptions{ContainerID: name, Force: true, RemoveVolumes: true})

			var client client.APIClient
			err = dtrutil.Poll(time.Second, 10, func() error {
				var err error
				client, err = dockerutil.MakeDockerClient(connStr, "", nil, j.config.Username, j.config.Password, j.config.RefreshToken)
				if err != nil {
					return err
				}
				_, err = client.ContainerList(j.ctx, types.ContainerListOptions{})
				if err != nil {
					return err
				}
				return nil
			})
			if err != nil {
				return err
			}
			parallelDinds[n] = client
		}
	}

	wg := sync.WaitGroup{}
	wg.Add(config.PushesPerBatch)
	warns := []error{}
	errs := []error{}
	for n := 0; n < config.PushesPerBatch; n++ {
		go func(n int) {
			for m := 0; m < config.Batches; m++ {
				// if we don't create individual http clients for each thread, we get timeout errors way more
				httpClient, err := util.HTTPClient(j.config.DTRInsecureTLS, j.config.DTRCA)
				if err != nil {
					errs = append(errs, err)
				} else {

					currTag := fmt.Sprintf("%s-n%d-m%d", tag, n, m)

					j.log.Infof("pushing %s/%s/%s:%s", j.config.DTRURL, config.Namespace, config.RepoName, currTag)
					// keep adding warnings while retrying, but not errors until it gives up
					m := 0
					for ; m < config.Retries; m++ {
						err = j.PushImage(parallelDinds[n], j.config.DTRURL, config.Namespace, config.RepoName, currTag, j.config.Username, j.config.Password, j.config.RefreshToken, httpClient, parallelRands[n], int64(config.LayerMinBytes), config.ImageMinLayers, j.config.DTRInsecureTLS, j.config.DTRCA)
						if err != nil {
							warns = append(warns, err)
						} else {
							break
						}
					}
					if m == config.Retries && err != nil {
						if err != nil {
							errs = append(errs, err)
						}
					}
				}
			}
			wg.Done()
		}(n)
	}
	wg.Wait()
	if len(errs) > 0 {
		return fmt.Errorf("%d errors and %d warnings pushing: %s; warnings: %s", len(errs), len(warns), errs, warns)
	}

	return nil
}

// this performs the heavy lifting of the load test
func (j *job) exec() error {
	table := crc64.MakeTable(0x1b)
	nameSeed := int64(crc64.Checksum([]byte(j.clientName), table))

	seed := int64(j.config.Seed)
	// if there's no seed provided, use the client name as the seed
	if seed == 0 {
		seed = nameSeed
	}
	rootRand := rand.New(rand.NewSource(seed))

	httpClient, err := util.HTTPClient(j.config.DTRInsecureTLS, j.config.DTRCA)
	if err != nil {
		return err
	}
	hostClient, err := dockerutil.MakeDockerClient("unix:///var/run/docker.sock", "", nil, j.config.Username, j.config.Password, j.config.RefreshToken)
	if err != nil {
		return err
	}

	if j.config.Push.Method == "dind-daemon" {
		err := j.PullImageWithDockerClient(hostClient, "", "", "docker", "dind", "", "", "")
		if err != nil {
			return err
		}
	}

	if j.config.Push != nil {
		if err := j.executePush(j.config.Push, hostClient, rootRand); err != nil {
			return err
		}
	}

	if j.config.Pull != nil {
		seed := int64(j.config.Seed)
		// if there's no seed provided, use the client name as the seed
		if seed == 0 {
			seed = nameSeed
		}
		rootRand := rand.New(rand.NewSource(seed))
		// we create a new source for each thread based on a seed
		// that way the layers created from each thread are the same every run
		thread1Rand := rand.New(rand.NewSource(rootRand.Int63()))
		tag := j.config.Pull.TagName

		if j.config.Pull.TagPattern == "random" {
			tag = fmt.Sprintf("%s-%010d", tag, thread1Rand.Int63n(1000000))
		}

		j.log.Infof("pulling %s/%s/%s:%s", j.config.DTRURL, j.config.Pull.Namespace, j.config.Pull.RepoName, tag)
		err = j.PullImage(httpClient, j.config.DTRURL, j.config.Pull.Namespace, j.config.Pull.RepoName, tag, j.config.Username, j.config.Password, j.config.RefreshToken, j.config.DTRInsecureTLS, j.config.DTRCA)
		if err != nil {
			return err
		}
	}

	if j.config.MassPopulate != nil {
		pushConfig := *j.config.MassPopulate.Push
		for i := 0; i < j.config.MassPopulate.UsersToPopulate; i++ {
			pushConfig.Namespace = fmt.Sprintf(fmt.Sprintf("testuser%06d", (rootRand.Int63()%int64(j.config.MassPopulate.TotalNumUsers))+1))

			for k := 0; k < j.config.MassPopulate.ReposToPopulate; k++ {
				pushConfig.RepoName = fmt.Sprintf(fmt.Sprintf("repo%d", (rootRand.Int63() % int64(j.config.MassPopulate.TotalNumReposPerUser))))

				if err := j.executePush(&pushConfig, hostClient, rootRand); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (j *job) Exec() {
	err := j.exec()
	if err != nil {
		j.log.Error("job failure")
		j.state = commands.JobState_FAILURE
		j.err = err
	} else {
		j.log.Info("job success")
		j.state = commands.JobState_SUCCESS
	}
}

func ClientRun(ctx context.Context, name, clientConfig string) (moshpit.Job, error) {
	log := moshpit.LoggerFromCtx(ctx)
	config := ClientConfig{}
	err := yaml.Unmarshal([]byte(clientConfig), &config)
	if err != nil {
		return nil, err
	}
	log.WithField("config", config).Debug("received client config")

	// TODO: kick off the job asynchronously
	j := &job{
		ctx:        ctx,
		clientName: name,
		log:        log,
		config:     config,
		state:      commands.JobState_RUNNING,
	}
	go j.Exec()
	return j, nil
}
