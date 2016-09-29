package commands

import (
	"encoding/json"
	"errors"
	"io/ioutil"

	"github.com/Sirupsen/logrus"
)

type metaData struct {
	Version      string  `json:"version"`
	Build        string  `json:"build"`
	Sha1         string  `json:"sha1"`
	HumanVersion string  `json:"humanVersion"`
	LastUpdate   string  `json:"lastUpdate"`
	Size         string  `json:"size"`
	Main         string  `json:"main"`
	Assets       []asset `json:"assets"`
	Channel      string  `json:"channel"`
	Arch         string  `json:"arch"`
}

type asset struct {
	Signature string `json:"signature,omitempty"`
	Symbols   string `json:"symbols,omitempty"`
	Notes     string `json:"notes,omitempty"`
}

func genMetaData(release metaData, tmpMetadata string) error {

	jsonByte, err := json.Marshal(release)
	if err != nil {
		logrus.Error("Could not generate Metadata")
		return errors.New("Could not gen Metadata")
	}
	logrus.Debug("Metadata: ", string(jsonByte))
	err = ioutil.WriteFile(tmpMetadata, jsonByte, 0644)
	if err != nil {
		logrus.Error("Could not write Metadata to tmp")
		return errors.New("Could not write Metadata in tmp: " + tmpMetadata)
	}
	return nil
}
