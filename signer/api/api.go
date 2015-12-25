package api

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/docker/notary/signer"
	"github.com/docker/notary/signer/keys"
	"github.com/docker/notary/tuf/signed"
	"github.com/gorilla/mux"

	pb "github.com/docker/notary/proto"
)

// Handlers sets up all the handers for the routes, injecting a specific CryptoService object for them to use
func Handlers(cryptoServices signer.CryptoServiceIndex) *mux.Router {
	r := mux.NewRouter()

	r.Methods("GET").Path("/{ID}").Handler(KeyInfo(cryptoServices))
	r.Methods("POST").Path("/new/{Algorithm}").Handler(CreateKey(cryptoServices))
	r.Methods("POST").Path("/delete").Handler(DeleteKey(cryptoServices))
	r.Methods("POST").Path("/sign").Handler(Sign(cryptoServices))
	return r
}

// getCryptoService handles looking up the correct signing service, given the
// algorithm specified in the HTTP request. If the algorithm isn't specified
// or isn't supported, an error is returned to the client and this function
// returns a nil CryptoService
func getCryptoService(algorithm string, cryptoServices signer.CryptoServiceIndex) (signed.CryptoService, error) {
	if algorithm == "" {
		return nil, fmt.Errorf("algorithm not specified")
	}

	if service, ok := cryptoServices[algorithm]; ok {
		return service, nil
	}

	return nil, fmt.Errorf("algorithm " + algorithm + " not supported")
}

// KeyInfo returns a Handler that given a specific Key ID param, returns the public key bits of that key
func KeyInfo(cryptoServices signer.CryptoServiceIndex) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		tufKey, _, err := FindKeyByID(cryptoServices, &pb.KeyID{ID: vars["ID"]})
		if err != nil {
			switch err {
			// If we received an ErrInvalidKeyID, the key doesn't exist, return 404
			case keys.ErrInvalidKeyID:
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(err.Error()))
				return
			// If we received anything else, it is unexpected, and we return a 500
			default:
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}
		}
		key := &pb.PublicKey{
			KeyInfo: &pb.KeyInfo{
				KeyID:     &pb.KeyID{ID: tufKey.ID()},
				Algorithm: &pb.Algorithm{Algorithm: tufKey.Algorithm()},
			},
			PublicKey: tufKey.Public(),
		}
		json.NewEncoder(w).Encode(key)
		return
	})
}

// CreateKey returns a handler that generates a new key using the provided
// algorithm. Only the public component of the key is returned.
func CreateKey(cryptoServices signer.CryptoServiceIndex) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		cryptoService, err := getCryptoService(vars["Algorithm"], cryptoServices)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		tufKey, err := cryptoService.Create("", vars["Algorithm"])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		key := &pb.PublicKey{
			KeyInfo: &pb.KeyInfo{
				KeyID:     &pb.KeyID{ID: tufKey.ID()},
				Algorithm: &pb.Algorithm{Algorithm: tufKey.Algorithm()},
			},
			PublicKey: tufKey.Public(),
		}
		json.NewEncoder(w).Encode(key)
		return
	})
}

// DeleteKey returns a handler that delete a specific KeyID
func DeleteKey(cryptoServices signer.CryptoServiceIndex) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var keyID *pb.KeyID
		err := json.NewDecoder(r.Body).Decode(&keyID)
		defer r.Body.Close()
		if err != nil || keyID.ID == "" {
			w.WriteHeader(http.StatusBadRequest)
			jsonErr, _ := json.Marshal("Malformed request")
			w.Write([]byte(jsonErr))
			return
		}

		_, cryptoService, err := FindKeyByID(cryptoServices, keyID)

		if err != nil {
			switch err {
			// If we received an ErrInvalidKeyID, the key doesn't exist, return 404
			case keys.ErrInvalidKeyID:
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(err.Error()))
				return
			// If we received anything else, it is unexpected, and we return a 500
			default:
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}
		}

		if err = cryptoService.RemoveKey(keyID.ID); err != nil {
			switch err {
			// If we received an ErrInvalidKeyID, the key doesn't exist, return 404
			case keys.ErrInvalidKeyID:
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(err.Error()))
				return
			// If we received anything else, it is unexpected, and we return a 500
			default:
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}
		}
		// In case we successfully delete this key, return 200
		return
	})
}

// Sign returns a handler that is able to perform signatures on a given blob
func Sign(cryptoServices signer.CryptoServiceIndex) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var sigRequest *pb.SignatureRequest
		err := json.NewDecoder(r.Body).Decode(&sigRequest)
		defer r.Body.Close()
		if err != nil || sigRequest.Content == nil ||
			sigRequest.KeyID == nil {
			w.WriteHeader(http.StatusBadRequest)
			jsonErr, _ := json.Marshal("Malformed request")
			w.Write([]byte(jsonErr))
			return
		}

		tufKey, cryptoService, err := FindKeyByID(cryptoServices, sigRequest.KeyID)
		if err == keys.ErrInvalidKeyID {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(err.Error()))
			return
		} else if err != nil {
			// We got an unexpected error
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		privKey, _, err := cryptoService.GetPrivateKey(tufKey.ID())
		if err != nil {
			// We got an unexpected error
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		sig, err := privKey.Sign(rand.Reader, sigRequest.Content, nil)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		signature := &pb.Signature{
			KeyInfo: &pb.KeyInfo{
				KeyID:     &pb.KeyID{ID: tufKey.ID()},
				Algorithm: &pb.Algorithm{Algorithm: tufKey.Algorithm()},
			},
			Algorithm: &pb.Algorithm{Algorithm: privKey.SignatureAlgorithm().String()},
			Content:   sig,
		}

		json.NewEncoder(w).Encode(signature)
		return
	})
}
