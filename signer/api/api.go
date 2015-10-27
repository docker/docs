package api

import (
	"encoding/json"
	"net/http"

	"github.com/docker/notary/signer"
	"github.com/docker/notary/signer/keys"
	"github.com/docker/notary/tuf/data"
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
func getCryptoService(w http.ResponseWriter, algorithm string, cryptoServices signer.CryptoServiceIndex) signed.CryptoService {
	if algorithm == "" {
		http.Error(w, "algorithm not specified", http.StatusBadRequest)
		return nil
	}

	service := cryptoServices[data.KeyAlgorithm(algorithm)]

	if service == nil {
		http.Error(w, "algorithm "+algorithm+" not supported", http.StatusBadRequest)
		return nil
	}

	return service
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
				Algorithm: &pb.Algorithm{Algorithm: tufKey.Algorithm().String()},
			},
			PublicKey: tufKey.Public(),
		}
		json.NewEncoder(w).Encode(key)
		return
	})
}

// CreateKey returns a handler that generates a new
func CreateKey(cryptoServices signer.CryptoServiceIndex) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		cryptoService := getCryptoService(w, vars["Algorithm"], cryptoServices)
		if cryptoService == nil {
			// Error handled inside getCryptoService
			return
		}

		tufKey, err := cryptoService.Create("", data.KeyAlgorithm(vars["Algorithm"]))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		key := &pb.PublicKey{
			KeyInfo: &pb.KeyInfo{
				KeyID:     &pb.KeyID{ID: tufKey.ID()},
				Algorithm: &pb.Algorithm{Algorithm: tufKey.Algorithm().String()},
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

		signatures, err := cryptoService.Sign([]string{sigRequest.KeyID.ID}, sigRequest.Content)
		if err != nil || len(signatures) != 1 {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		signature := &pb.Signature{
			KeyInfo: &pb.KeyInfo{
				KeyID:     &pb.KeyID{ID: tufKey.ID()},
				Algorithm: &pb.Algorithm{Algorithm: tufKey.Algorithm().String()},
			},
			Content: signatures[0].Signature,
		}

		json.NewEncoder(w).Encode(signature)
		return
	})
}
