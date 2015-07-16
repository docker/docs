package api

import (
	"encoding/json"
	"net/http"

	"github.com/docker/notary/signer"
	"github.com/docker/notary/signer/keys"
	"github.com/endophage/gotuf/data"
	"github.com/gorilla/mux"

	pb "github.com/docker/notary/proto"
)

// Handlers sets up all the handers for the routes, injecting a specific SigningService object for them to use
func Handlers(sigServices signer.SigningServiceIndex) *mux.Router {
	r := mux.NewRouter()

	r.Methods("GET").Path("/{ID}").Handler(KeyInfo(sigServices))
	r.Methods("POST").Path("/new/{Algorithm}").Handler(CreateKey(sigServices))
	r.Methods("POST").Path("/delete").Handler(DeleteKey(sigServices))
	r.Methods("POST").Path("/sign").Handler(Sign(sigServices))
	return r
}

// getSigningService handles looking up the correct signing service, given the
// algorithm specified in the HTTP request. If the algorithm isn't specified
// or isn't supported, an error is returned to the client and this function
// returns a nil SigningService
func getSigningService(w http.ResponseWriter, algorithm string, sigServices signer.SigningServiceIndex) signer.SigningService {
	if algorithm == "" {
		http.Error(w, "algorithm not specified", http.StatusBadRequest)
		return nil
	}

	service := sigServices[data.KeyAlgorithm(algorithm)]

	if service == nil {
		http.Error(w, "algorithm "+algorithm+" not supported", http.StatusBadRequest)
		return nil
	}

	return service
}

// KeyInfo returns a Handler that given a specific Key ID param, returns the public key bits of that key
func KeyInfo(sigServices signer.SigningServiceIndex) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		key, _, err := FindKeyByID(sigServices, &pb.KeyID{ID: vars["ID"]})
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
		json.NewEncoder(w).Encode(key)
		return
	})
}

// CreateKey returns a handler that generates a new
func CreateKey(sigServices signer.SigningServiceIndex) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		sigService := getSigningService(w, vars["Algorithm"], sigServices)
		if sigService == nil {
			// Error handled inside getSigningService
			return
		}

		key, err := sigService.CreateKey()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		json.NewEncoder(w).Encode(key)
		return
	})
}

// DeleteKey returns a handler that delete a specific KeyID
func DeleteKey(sigServices signer.SigningServiceIndex) http.Handler {
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

		_, sigService, err := FindKeyByID(sigServices, keyID)

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

		_, err = sigService.DeleteKey(keyID)

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
		// In case we successfully delete this key, return 200
		return
	})
}

// Sign returns a handler that is able to perform signatures on a given blob
func Sign(sigServices signer.SigningServiceIndex) http.Handler {
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

		_, sigService, err := FindKeyByID(sigServices, sigRequest.KeyID)
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

		signer, err := sigService.Signer(sigRequest.KeyID)
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

		signature, err := signer.Sign(sigRequest)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		json.NewEncoder(w).Encode(signature)
		return
	})
}
