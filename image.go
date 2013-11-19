package docker

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/dotcloud/docker/archive"
	"github.com/dotcloud/docker/graphdriver"
	"github.com/dotcloud/docker/utils"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

type Image struct {
	ID              string    `json:"id"`
	Parent          string    `json:"parent,omitempty"`
	Comment         string    `json:"comment,omitempty"`
	Created         time.Time `json:"created"`
	Container       string    `json:"container,omitempty"`
	ContainerConfig Config    `json:"container_config,omitempty"`
	DockerVersion   string    `json:"docker_version,omitempty"`
	Author          string    `json:"author,omitempty"`
	Config          *Config   `json:"config,omitempty"`
	Architecture    string    `json:"architecture,omitempty"`
	graph           *Graph
	Size            int64
}

func LoadImage(root string) (*Image, error) {
	// Load the json data
	jsonData, err := ioutil.ReadFile(jsonPath(root))
	if err != nil {
		return nil, err
	}
	img := &Image{}

	if err := json.Unmarshal(jsonData, img); err != nil {
		return nil, err
	}
	if err := ValidateID(img.ID); err != nil {
		return nil, err
	}

	if buf, err := ioutil.ReadFile(path.Join(root, "layersize")); err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
	} else {
		size, err := strconv.Atoi(string(buf))
		if err != nil {
			return nil, err
		}
		img.Size = int64(size)
	}

	return img, nil
}

func StoreImage(img *Image, jsonData []byte, layerData archive.Archive, root, rootfs string) error {
	// Store the layer
	layer := rootfs
	if err := os.MkdirAll(layer, 0755); err != nil {
		return err
	}

	// If layerData is not nil, unpack it into the new layer
	if layerData != nil {
		if differ, ok := img.graph.driver.(graphdriver.Differ); ok {
			if err := differ.ApplyDiff(img.ID, layerData); err != nil {
				return err
			}
		} else {
			start := time.Now()
			utils.Debugf("Start untar layer")
			if err := archive.ApplyLayer(layer, layerData); err != nil {
				return err
			}
			utils.Debugf("Untar time: %vs", time.Now().Sub(start).Seconds())
		}
	}

	// If raw json is provided, then use it
	if jsonData != nil {
		return ioutil.WriteFile(jsonPath(root), jsonData, 0600)
	}
	// Otherwise, unmarshal the image
	jsonData, err := json.Marshal(img)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(jsonPath(root), jsonData, 0600); err != nil {
		return err
	}
	// Compute and save the size of the rootfs
	size, err := utils.TreeSize(rootfs)
	if err != nil {
		return fmt.Errorf("Error computing size of rootfs %s: %s", img.ID, err)
	}
	img.Size = size
	if err := img.SaveSize(root); err != nil {
		return err
	}

	return nil
}

// SaveSize stores the current `size` value of `img` in the directory `root`.
func (img *Image) SaveSize(root string) error {
	if err := ioutil.WriteFile(path.Join(root, "layersize"), []byte(strconv.Itoa(int(img.Size))), 0600); err != nil {
		return fmt.Errorf("Error storing image size in %s/layersize: %s", root, err)
	}
	return nil
}

func jsonPath(root string) string {
	return path.Join(root, "json")
}

// TarLayer returns a tar archive of the image's filesystem layer.
func (img *Image) TarLayer(compression archive.Compression) (archive.Archive, error) {
	if img.graph == nil {
		return nil, fmt.Errorf("Can't load storage driver for unregistered image %s", img.ID)
	}
	layerPath, err := img.graph.driver.Get(img.ID)
	if err != nil {
		return nil, err
	}
	return archive.Tar(layerPath, compression)
}

func ValidateID(id string) error {
	if id == "" {
		return fmt.Errorf("Image id can't be empty")
	}
	if strings.Contains(id, ":") {
		return fmt.Errorf("Invalid character in image id: ':'")
	}
	return nil
}

func GenerateID() string {
	id := make([]byte, 32)
	_, err := io.ReadFull(rand.Reader, id)
	if err != nil {
		panic(err) // This shouldn't happen
	}
	return hex.EncodeToString(id)
}

// Image includes convenience proxy functions to its graph
// These functions will return an error if the image is not registered
// (ie. if image.graph == nil)
func (img *Image) History() ([]*Image, error) {
	var parents []*Image
	if err := img.WalkHistory(
		func(img *Image) error {
			parents = append(parents, img)
			return nil
		},
	); err != nil {
		return nil, err
	}
	return parents, nil
}

func (img *Image) WalkHistory(handler func(*Image) error) (err error) {
	currentImg := img
	for currentImg != nil {
		if handler != nil {
			if err := handler(currentImg); err != nil {
				return err
			}
		}
		currentImg, err = currentImg.GetParent()
		if err != nil {
			return fmt.Errorf("Error while getting parent image: %v", err)
		}
	}
	return nil
}

func (img *Image) GetParent() (*Image, error) {
	if img.Parent == "" {
		return nil, nil
	}
	if img.graph == nil {
		return nil, fmt.Errorf("Can't lookup parent of unregistered image")
	}
	return img.graph.Get(img.Parent)
}

func (img *Image) root() (string, error) {
	if img.graph == nil {
		return "", fmt.Errorf("Can't lookup root of unregistered image")
	}
	return img.graph.imageRoot(img.ID), nil
}

func (img *Image) getParentsSize(size int64) int64 {
	parentImage, err := img.GetParent()
	if err != nil || parentImage == nil {
		return size
	}
	size += parentImage.Size
	return parentImage.getParentsSize(size)
}

// Build an Image object from raw json data
func NewImgJSON(src []byte) (*Image, error) {
	ret := &Image{}

	utils.Debugf("Json string: {%s}", src)
	// FIXME: Is there a cleaner way to "purify" the input json?
	if err := json.Unmarshal(src, ret); err != nil {
		return nil, err
	}
	return ret, nil
}
