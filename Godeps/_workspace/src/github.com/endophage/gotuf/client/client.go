package client

import (
	"encoding/json"
	"fmt"
	"io"
	"path"
	"path/filepath"
	"strings"

	"github.com/Sirupsen/logrus"
	tuf "github.com/endophage/gotuf"
	"github.com/endophage/gotuf/data"
	"github.com/endophage/gotuf/keys"
	"github.com/endophage/gotuf/signed"
	"github.com/endophage/gotuf/store"
	"github.com/endophage/gotuf/utils"
)

type Client struct {
	local  *tuf.TufRepo
	remote store.RemoteStore
	keysDB *keys.KeyDB
}

func NewClient(local *tuf.TufRepo, remote store.RemoteStore, keysDB *keys.KeyDB) *Client {
	return &Client{
		local:  local,
		remote: remote,
		keysDB: keysDB,
	}
}

// Update an in memory copy of the TUF Repo. If an error is returned, the
// Client instance should be considered corrupted and discarded as it may
// be left in a partially updated state
func (c *Client) Update() error {
	err := c.update()
	if err != nil {
		switch err.(type) {
		case tuf.ErrSigVerifyFail:
		case tuf.ErrMetaExpired:
		case tuf.ErrLocalRootExpired:
			if err := c.downloadRoot(); err != nil {
				logrus.Errorf("Client Update (Root):", err)
				return err
			}
		default:
			return err
		}
	}
	// If we error again, we now have the latest root and just want to fail
	// out as there's no expectation the problem can be resolved automatically
	return c.update()
}

func (c *Client) update() error {
	err := c.downloadTimestamp()
	if err != nil {
		logrus.Errorf("Client Update (Timestamp):", err)
		return err
	}
	err = c.downloadSnapshot()
	if err != nil {
		logrus.Errorf("Client Update (Snapshot):", err)
		return err
	}
	err = c.checkRoot()
	if err != nil {
		return err
	}
	err = c.downloadTargets("targets")
	if err != nil {
		logrus.Errorf("Client Update (Targets):", err)
		return err
	}
	return nil
}

// checkRoot determines if the hash, and size are still those reported
// in the snapshot file. It will also check the expiry, however, if the
// hash and size in snapshot are unchanged but the root file has expired,
// there is little expectation that the situation can be remedied.
func (c Client) checkRoot() error {
	return nil
}

// downloadRoot is responsible for downloading the root.json
func (c *Client) downloadRoot() error {
	role := data.RoleName("root")
	size := c.local.Snapshot.Signed.Meta[role].Length

	raw, err := c.remote.GetMeta(role, size)
	if err != nil {
		return err
	}
	s := &data.Signed{}
	err = json.Unmarshal(raw, s)
	if err != nil {
		return err
	}
	err = signed.Verify(s, role, 0, c.keysDB)
	if err != nil {
		return err
	}
	c.local.SetRoot(s)
	return nil
}

// downloadTimestamp is responsible for downloading the timestamp.json
func (c *Client) downloadTimestamp() error {
	role := data.RoleName("timestamp")
	raw, err := c.remote.GetMeta(role, 5<<20)
	if err != nil {
		return err
	}
	s := &data.Signed{}
	err = json.Unmarshal(raw, s)
	if err != nil {
		return err
	}
	err = signed.Verify(s, role, 0, c.keysDB)
	if err != nil {
		return err
	}
	c.local.SetTimestamp(s)
	return nil
}

// downloadSnapshot is responsible for downloading the snapshot.json
func (c *Client) downloadSnapshot() error {
	role := data.RoleName("snapshot")
	size := c.local.Timestamp.Signed.Meta[role+".txt"].Length
	raw, err := c.remote.GetMeta(role, size)
	if err != nil {
		return err
	}
	s := &data.Signed{}
	err = json.Unmarshal(raw, s)
	if err != nil {
		return err
	}
	err = signed.Verify(s, role, 0, c.keysDB)
	if err != nil {
		return err
	}
	c.local.SetSnapshot(s)
	return nil
}

// downloadTargets is responsible for downloading any targets file
// including delegates roles. It will download the whole tree of
// delegated roles below the given one
func (c *Client) downloadTargets(role string) error {
	role = data.RoleName(role) // this will really only do something for base targets role
	snap := c.local.Snapshot.Signed
	root := c.local.Root.Signed
	r := c.keysDB.GetRole(role)
	if r == nil {
		return fmt.Errorf("Invalid role: %s", role)
	}
	keyIDs := r.KeyIDs
	s, err := c.GetTargetsFile(role, keyIDs, snap.Meta, root.ConsistentSnapshot, r.Threshold)
	if err != nil {
		logrus.Error("Error getting targets file:", err)
		return err
	}
	err = c.local.SetTargets(role, s)
	if err != nil {
		return err
	}
	t := c.local.Targets[role].Signed
	for _, r := range t.Delegations.Roles {
		err := c.downloadTargets(r.Name)
		if err != nil {
			logrus.Error("Failed to download ", role, err)
			return err
		}
	}
	return nil
}

func (c Client) GetTargetsFile(roleName string, keyIDs []string, snapshotMeta data.Files, consistent bool, threshold int) (*data.Signed, error) {
	rolePath, err := c.RoleTargetsPath(roleName, snapshotMeta, consistent)
	if err != nil {
		return nil, err
	}
	r, err := c.remote.GetMeta(rolePath, snapshotMeta[roleName+".txt"].Length)
	if err != nil {
		return nil, err
	}
	s := &data.Signed{}
	err = json.Unmarshal(r, s)
	if err != nil {
		logrus.Error("Error unmarshalling targets file:", err)
		return nil, err
	}
	err = signed.Verify(s, roleName, 0, c.keysDB)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (c Client) RoleTargetsPath(roleName string, snapshotMeta data.Files, consistent bool) (string, error) {
	if consistent {
		roleMeta, ok := snapshotMeta[roleName]
		if !ok {
			return "", fmt.Errorf("Consistent Snapshots Enabled but no meta found for target role")
		}
		if _, ok := roleMeta.Hashes["sha256"]; !ok {
			return "", fmt.Errorf("Consistent Snapshots Enabled and sha256 not found for targets file in snapshot meta")
		}
		dir := filepath.Dir(roleName)
		if strings.Contains(roleName, "/") {
			lastSlashIdx := strings.LastIndex(roleName, "/")
			roleName = roleName[lastSlashIdx+1:]
		}
		roleName = path.Join(
			dir,
			fmt.Sprintf("%s.%s.json", roleMeta.Hashes["sha256"], roleName),
		)
	}
	return roleName, nil
}

func (c Client) TargetMeta(path string) *data.FileMeta {
	return c.local.FindTarget(path)
}

func (c Client) DownloadTarget(dst io.Writer, path string, meta *data.FileMeta) error {
	reader, err := c.remote.GetTarget(path)
	if err != nil {
		return err
	}
	defer reader.Close()
	r := io.TeeReader(
		io.LimitReader(reader, meta.Length),
		dst,
	)
	err = utils.ValidateTarget(r, meta)
	return err
}
