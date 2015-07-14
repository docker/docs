package keys

import pb "github.com/docker/rufus/proto"

// KeyDB represents an in-memory key keystore
type KeyDB struct {
	keys map[string]*Key
}

// CreateKey is needed to implement KeyManager. Returns an empty key.
func (db *KeyDB) CreateKey() (*pb.PublicKey, error) {
	k := &pb.PublicKey{}

	return k, nil
}

// AddKey Adds a new key to the database
func (db *KeyDB) AddKey(key *Key) error {
	if _, ok := db.keys[key.ID]; ok {
		return ErrExists
	}
	db.keys[key.ID] = key

	return nil
}

// GetKey returns the private bits of a key
func (db *KeyDB) GetKey(keyInfo *pb.KeyInfo) (*Key, error) {
	if key, ok := db.keys[keyInfo.ID]; ok {
		return key, nil
	}
	return nil, ErrInvalidKeyID
}

// DeleteKey deletes the keyID from the database
func (db *KeyDB) DeleteKey(keyInfo *pb.KeyInfo) (*pb.Void, error) {
	_, err := db.GetKey(keyInfo)
	if err != nil {
		return nil, err
	}
	delete(db.keys, keyInfo.ID)
	return nil, nil
}

// KeyInfo returns the public bits of a key, given a specific keyID
func (db *KeyDB) KeyInfo(keyInfo *pb.KeyInfo) (*pb.PublicKey, error) {
	key, err := db.GetKey(keyInfo)
	if err != nil {
		return nil, err
	}
	return &pb.PublicKey{KeyInfo: &pb.KeyInfo{ID: keyInfo.ID, Algorithm: &pb.Algorithm{Algorithm: key.Algorithm}}, PublicKey: key.Public[:]}, nil
}

// NewKeyDB returns an instance of KeyDB
func NewKeyDB() *KeyDB {
	return &KeyDB{
		keys: make(map[string]*Key),
	}
}
