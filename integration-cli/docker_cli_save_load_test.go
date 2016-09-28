package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/docker/docker/pkg/integration/checker"
	"github.com/go-check/check"
)

// save a repo using gz compression and try to load it using stdout
func (s *DockerSuite) TestSaveXzAndLoadRepoStdout(c *check.C) {
	testRequires(c, DaemonIsLinux)
	name := "test-save-xz-and-load-repo-stdout"
	dockerCmd(c, "run", "--name", name, "busybox", "true")

	repoName := "foobar-save-load-test-xz-gz"
	out, _ := dockerCmd(c, "commit", name, repoName)

	dockerCmd(c, "inspect", repoName)

	repoTarball, _, err := runCommandPipelineWithOutput(
		exec.Command(dockerBinary, "save", repoName),
		exec.Command("xz", "-c"),
		exec.Command("gzip", "-c"))
	c.Assert(err, checker.IsNil, check.Commentf("failed to save repo: %v %v", out, err))
	deleteImages(repoName)

	loadCmd := exec.Command(dockerBinary, "load")
	loadCmd.Stdin = strings.NewReader(repoTarball)
	out, _, err = runCommandWithOutput(loadCmd)
	c.Assert(err, checker.NotNil, check.Commentf("expected error, but succeeded with no error and output: %v", out))

	after, _, err := dockerCmdWithError("inspect", repoName)
	c.Assert(err, checker.NotNil, check.Commentf("the repo should not exist: %v", after))
}

// save a repo using xz+gz compression and try to load it using stdout
func (s *DockerSuite) TestSaveXzGzAndLoadRepoStdout(c *check.C) {
	testRequires(c, DaemonIsLinux)
	name := "test-save-xz-gz-and-load-repo-stdout"
	dockerCmd(c, "run", "--name", name, "busybox", "true")

	repoName := "foobar-save-load-test-xz-gz"
	dockerCmd(c, "commit", name, repoName)

	dockerCmd(c, "inspect", repoName)

	out, _, err := runCommandPipelineWithOutput(
		exec.Command(dockerBinary, "save", repoName),
		exec.Command("xz", "-c"),
		exec.Command("gzip", "-c"))
	c.Assert(err, checker.IsNil, check.Commentf("failed to save repo: %v %v", out, err))

	deleteImages(repoName)

	loadCmd := exec.Command(dockerBinary, "load")
	loadCmd.Stdin = strings.NewReader(out)
	out, _, err = runCommandWithOutput(loadCmd)
	c.Assert(err, checker.NotNil, check.Commentf("expected error, but succeeded with no error and output: %v", out))

	after, _, err := dockerCmdWithError("inspect", repoName)
	c.Assert(err, checker.NotNil, check.Commentf("the repo should not exist: %v", after))
}

func (s *DockerSuite) TestSaveSingleTag(c *check.C) {
	testRequires(c, DaemonIsLinux)
	repoName := "foobar-save-single-tag-test"
	dockerCmd(c, "tag", "busybox:latest", fmt.Sprintf("%v:latest", repoName))

	out, _ := dockerCmd(c, "images", "-q", "--no-trunc", repoName)
	cleanedImageID := strings.TrimSpace(out)

	out, _, err := runCommandPipelineWithOutput(
		exec.Command(dockerBinary, "save", fmt.Sprintf("%v:latest", repoName)),
		exec.Command("tar", "t"),
		exec.Command("grep", "-E", fmt.Sprintf("(^repositories$|%v)", cleanedImageID)))
	c.Assert(err, checker.IsNil, check.Commentf("failed to save repo with image ID and 'repositories' file: %s, %v", out, err))
}

func (s *DockerSuite) TestSaveCheckTimes(c *check.C) {
	repoName := "busybox:latest"
	out, _ := dockerCmd(c, "inspect", repoName)
	data := []struct {
		ID      string
		Created time.Time
	}{}
	err := json.Unmarshal([]byte(out), &data)
	c.Assert(err, checker.IsNil, check.Commentf("failed to marshal from %q: err %v", repoName, err))
	c.Assert(len(data), checker.Not(checker.Equals), 0, check.Commentf("failed to marshal the data from %q", repoName))
	tarTvTimeFormat := "2006-01-02 15:04"
	out, _, err = runCommandPipelineWithOutput(
		exec.Command(dockerBinary, "save", repoName),
		exec.Command("tar", "tv"),
		exec.Command("grep", "-E", fmt.Sprintf("%s %s", data[0].Created.Format(tarTvTimeFormat), data[0].ID)))
	c.Assert(err, checker.IsNil, check.Commentf("failed to save repo with image ID and 'repositories' file: %s, %v", out, err))
}

func (s *DockerSuite) TestSaveImageId(c *check.C) {
	testRequires(c, DaemonIsLinux)
	repoName := "foobar-save-image-id-test"
	dockerCmd(c, "tag", "emptyfs:latest", fmt.Sprintf("%v:latest", repoName))

	out, _ := dockerCmd(c, "images", "-q", "--no-trunc", repoName)
	cleanedLongImageID := strings.TrimSpace(out)

	out, _ = dockerCmd(c, "images", "-q", repoName)
	cleanedShortImageID := strings.TrimSpace(out)

	// Make sure IDs are not empty
	c.Assert(cleanedLongImageID, checker.Not(check.Equals), "", check.Commentf("Id should not be empty."))
	c.Assert(cleanedShortImageID, checker.Not(check.Equals), "", check.Commentf("Id should not be empty."))

	saveCmd := exec.Command(dockerBinary, "save", cleanedShortImageID)
	tarCmd := exec.Command("tar", "t")

	var err error
	tarCmd.Stdin, err = saveCmd.StdoutPipe()
	c.Assert(err, checker.IsNil, check.Commentf("cannot set stdout pipe for tar: %v", err))
	grepCmd := exec.Command("grep", cleanedLongImageID)
	grepCmd.Stdin, err = tarCmd.StdoutPipe()
	c.Assert(err, checker.IsNil, check.Commentf("cannot set stdout pipe for grep: %v", err))

	c.Assert(tarCmd.Start(), checker.IsNil, check.Commentf("tar failed with error: %v", err))
	c.Assert(saveCmd.Start(), checker.IsNil, check.Commentf("docker save failed with error: %v", err))
	defer saveCmd.Wait()
	defer tarCmd.Wait()

	out, _, err = runCommandWithOutput(grepCmd)

	c.Assert(err, checker.IsNil, check.Commentf("failed to save repo with image ID: %s, %v", out, err))
}

// save a repo and try to load it using flags
func (s *DockerSuite) TestSaveAndLoadRepoFlags(c *check.C) {
	testRequires(c, DaemonIsLinux)
	name := "test-save-and-load-repo-flags"
	dockerCmd(c, "run", "--name", name, "busybox", "true")

	repoName := "foobar-save-load-test"

	deleteImages(repoName)
	dockerCmd(c, "commit", name, repoName)

	before, _ := dockerCmd(c, "inspect", repoName)

	out, _, err := runCommandPipelineWithOutput(
		exec.Command(dockerBinary, "save", repoName),
		exec.Command(dockerBinary, "load"))
	c.Assert(err, checker.IsNil, check.Commentf("failed to save and load repo: %s, %v", out, err))

	after, _ := dockerCmd(c, "inspect", repoName)
	c.Assert(before, checker.Equals, after, check.Commentf("inspect is not the same after a save / load"))
}

func (s *DockerSuite) TestSaveMultipleNames(c *check.C) {
	testRequires(c, DaemonIsLinux)
	repoName := "foobar-save-multi-name-test"

	// Make one image
	dockerCmd(c, "tag", "emptyfs:latest", fmt.Sprintf("%v-one:latest", repoName))

	// Make two images
	dockerCmd(c, "tag", "emptyfs:latest", fmt.Sprintf("%v-two:latest", repoName))

	out, _, err := runCommandPipelineWithOutput(
		exec.Command(dockerBinary, "save", fmt.Sprintf("%v-one", repoName), fmt.Sprintf("%v-two:latest", repoName)),
		exec.Command("tar", "xO", "repositories"),
		exec.Command("grep", "-q", "-E", "(-one|-two)"),
	)
	c.Assert(err, checker.IsNil, check.Commentf("failed to save multiple repos: %s, %v", out, err))
}

func (s *DockerSuite) TestSaveRepoWithMultipleImages(c *check.C) {
	testRequires(c, DaemonIsLinux)
	makeImage := func(from string, tag string) string {
		var (
			out string
		)
		out, _ = dockerCmd(c, "run", "-d", from, "true")
		cleanedContainerID := strings.TrimSpace(out)

		out, _ = dockerCmd(c, "commit", cleanedContainerID, tag)
		imageID := strings.TrimSpace(out)
		return imageID
	}

	repoName := "foobar-save-multi-images-test"
	tagFoo := repoName + ":foo"
	tagBar := repoName + ":bar"

	idFoo := makeImage("busybox:latest", tagFoo)
	idBar := makeImage("busybox:latest", tagBar)

	deleteImages(repoName)

	// create the archive
	out, _, err := runCommandPipelineWithOutput(
		exec.Command(dockerBinary, "save", repoName),
		exec.Command("tar", "t"),
		exec.Command("grep", "VERSION"),
		exec.Command("cut", "-d", "/", "-f1"))
	c.Assert(err, checker.IsNil, check.Commentf("failed to save multiple images: %s, %v", out, err))
	actual := strings.Split(strings.TrimSpace(out), "\n")

	// make the list of expected layers
	out, _ = dockerCmd(c, "history", "-q", "--no-trunc", "busybox:latest")
	expected := append(strings.Split(strings.TrimSpace(out), "\n"), idFoo, idBar)

	sort.Strings(actual)
	sort.Strings(expected)
	c.Assert(actual, checker.DeepEquals, expected, check.Commentf("archive does not contains the right layers: got %v, expected %v", actual, expected))
}

// Issue #6722 #5892 ensure directories are included in changes
func (s *DockerSuite) TestSaveDirectoryPermissions(c *check.C) {
	testRequires(c, DaemonIsLinux)
	layerEntries := []string{"opt/", "opt/a/", "opt/a/b/", "opt/a/b/c"}
	layerEntriesAUFS := []string{"./", ".wh..wh.aufs", ".wh..wh.orph/", ".wh..wh.plnk/", "opt/", "opt/a/", "opt/a/b/", "opt/a/b/c"}

	name := "save-directory-permissions"
	tmpDir, err := ioutil.TempDir("", "save-layers-with-directories")
	c.Assert(err, checker.IsNil, check.Commentf("failed to create temporary directory: %s", err))
	extractionDirectory := filepath.Join(tmpDir, "image-extraction-dir")
	os.Mkdir(extractionDirectory, 0777)

	defer os.RemoveAll(tmpDir)
	_, err = buildImage(name,
		`FROM busybox
	RUN adduser -D user && mkdir -p /opt/a/b && chown -R user:user /opt/a
	RUN touch /opt/a/b/c && chown user:user /opt/a/b/c`,
		true)
	c.Assert(err, checker.IsNil, check.Commentf("%v", err))

	out, _, err := runCommandPipelineWithOutput(
		exec.Command(dockerBinary, "save", name),
		exec.Command("tar", "-xf", "-", "-C", extractionDirectory),
	)
	c.Assert(err, checker.IsNil, check.Commentf("failed to save and extract image: %s", out))

	dirs, err := ioutil.ReadDir(extractionDirectory)
	c.Assert(err, checker.IsNil, check.Commentf("failed to get a listing of the layer directories: %s", err))

	found := false
	for _, entry := range dirs {
		var entriesSansDev []string
		if entry.IsDir() {
			layerPath := filepath.Join(extractionDirectory, entry.Name(), "layer.tar")

			f, err := os.Open(layerPath)
			c.Assert(err, checker.IsNil, check.Commentf("failed to open %s: %s", layerPath, err))

			entries, err := listTar(f)
			for _, e := range entries {
				if !strings.Contains(e, "dev/") {
					entriesSansDev = append(entriesSansDev, e)
				}
			}
			c.Assert(err, checker.IsNil, check.Commentf("encountered error while listing tar entries: %s", err))

			if reflect.DeepEqual(entriesSansDev, layerEntries) || reflect.DeepEqual(entriesSansDev, layerEntriesAUFS) {
				found = true
				break
			}
		}
	}

	c.Assert(found, checker.Equals, true, check.Commentf("failed to find the layer with the right content listing"))

}

/*
 * REBASE NOTE: migrate this test
// TestSaveAndLoadContentAddressable does a repeated save/load and checks that
// the ID changes the first time (because of content addressability), but
// doesn't change with another save/load (again, because of content
// addressability).
func (s *DockerSuite) TestSaveAndLoadContentAddressable(c *check.C) {
	baseImage := "ca-save-load-test-base"
	derivedImage := "ca-save-load-test-derived"

	_, err := buildImage(baseImage, fmt.Sprintf(`
	    FROM busybox
	    ENV FOO bar
	    CMD echo %s
	`, baseImage), true)
	if err != nil {
		c.Fatal(err)
	}

	baseIDBeforeSave, err := inspectField(baseImage, "Id")
	c.Assert(err, check.IsNil)
	baseParentBeforeSave, err := inspectField(baseImage, "Parent")
	c.Assert(err, check.IsNil)

	_, err = buildImage(derivedImage, fmt.Sprintf(`
	    FROM %s
	    CMD echo %s
	`, baseImage, derivedImage), true)
	if err != nil {
		c.Fatal(err)
	}

	derivedIDBeforeSave, err := inspectField(derivedImage, "Id")
	c.Assert(err, check.IsNil)
	derivedParentBeforeSave, err := inspectField(derivedImage, "Parent")
	c.Assert(err, check.IsNil)

	tmpDir, err := ioutil.TempDir("", "contentaddressable")
	if err != nil {
		c.Errorf("failed to create temporary directory: %s", err)
	}
	defer os.RemoveAll(tmpDir)
	baseTar := filepath.Join(tmpDir, "base.tar")
	derivedTar := filepath.Join(tmpDir, "derived.tar")
	dockerCmd(c, "save", "--output", baseTar, baseImage)
	dockerCmd(c, "save", "--output", derivedTar, derivedImage)

	// Delete images and reload
	dockerCmd(c, "rmi", baseImage)
	dockerCmd(c, "rmi", derivedImage)

	dockerCmd(c, "load", "--input", baseTar)
	dockerCmd(c, "load", "--input", derivedTar)

	baseIDAfterLoad, err := inspectField(baseImage, "Id")
	c.Assert(err, check.IsNil)
	baseParentAfterLoad, err := inspectField(baseImage, "Parent")
	c.Assert(err, check.IsNil)
	derivedIDAfterLoad, err := inspectField(derivedImage, "Id")
	c.Assert(err, check.IsNil)
	derivedParentAfterLoad, err := inspectField(derivedImage, "Parent")
	c.Assert(err, check.IsNil)

	// New IDs must be digests
	_, err = digest.ParseDigest(baseIDAfterLoad)
	c.Assert(err, check.IsNil)
	_, err = digest.ParseDigest(baseParentAfterLoad)
	c.Assert(err, check.IsNil)
	_, err = digest.ParseDigest(derivedIDAfterLoad)
	c.Assert(err, check.IsNil)
	_, err = digest.ParseDigest(derivedParentAfterLoad)
	c.Assert(err, check.IsNil)

	// They must NOT match the old IDs
	c.Assert(baseIDAfterLoad, check.Not(check.Equals), baseIDBeforeSave)
	c.Assert(baseParentAfterLoad, check.Not(check.Equals), baseParentBeforeSave)
	c.Assert(derivedIDAfterLoad, check.Not(check.Equals), derivedIDBeforeSave)
	c.Assert(derivedParentAfterLoad, check.Not(check.Equals), derivedParentBeforeSave)

	// Derived image parent should match base image ID
	if derivedParentAfterLoad != baseIDAfterLoad {
		c.Fatal("loaded derived image did not use already-loaded base image as parent")
	}

	tmpDir2, err := ioutil.TempDir("", "contentaddressable2")
	if err != nil {
		c.Errorf("failed to create temporary directory: %s", err)
	}
	defer os.RemoveAll(tmpDir2)
	baseTar = filepath.Join(tmpDir2, "base.tar")
	derivedTar = filepath.Join(tmpDir2, "derived.tar")
	dockerCmd(c, "save", "--output", baseTar, baseImage)
	dockerCmd(c, "save", "--output", derivedTar, derivedImage)

	// Delete images and reload
	dockerCmd(c, "rmi", baseImage)
	dockerCmd(c, "rmi", derivedImage)

	dockerCmd(c, "load", "--input", baseTar)
	dockerCmd(c, "load", "--input", derivedTar)

	baseIDAfterLoad2, err := inspectField(baseImage, "Id")
	c.Assert(err, check.IsNil)
	baseParentAfterLoad2, err := inspectField(baseImage, "Parent")
	c.Assert(err, check.IsNil)
	derivedIDAfterLoad2, err := inspectField(derivedImage, "Id")
	c.Assert(err, check.IsNil)
	derivedParentAfterLoad2, err := inspectField(derivedImage, "Parent")
	c.Assert(err, check.IsNil)

	// New IDs must match IDs from first load
	c.Assert(baseIDAfterLoad2, check.Equals, baseIDAfterLoad)
	c.Assert(baseParentAfterLoad2, check.Equals, baseParentAfterLoad)
	c.Assert(derivedIDAfterLoad2, check.Equals, derivedIDAfterLoad)
	c.Assert(derivedParentAfterLoad2, check.Equals, derivedParentAfterLoad)
}
*/
