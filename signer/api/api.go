package api

import (
	"encoding/json"
	"net/http"

	"github.com/docker/rufus"
	"github.com/docker/rufus/keys"
	"github.com/gorilla/mux"

	pb "github.com/docker/rufus/proto"
)

// Handlers sets up all the handers for the routes, injecting a specific SigningService object for them to use
func Handlers(sigServices rufus.SigningServiceIndex) *mux.Router {
	r := mux.NewRouter()

	r.Methods("GET").Path("/{Algorithm}/{ID}").Handler(KeyInfo(sigServices))
	r.Methods("POST").Path("/new/{Algorithm}").Handler(CreateKey(sigServices))
	r.Methods("POST").Path("/delete").Handler(DeleteKey(sigServices))
	r.Methods("POST").Path("/sign").Handler(Sign(sigServices))
	return r
}

// getSigningService handles looking up the correct signing service, given the
// algorithm specified in the HTTP request. If the algorithm isn't specified
// or isn't supported, an error is returned to the client and this function
// returns a nil SigningService
func getSigningService(w http.ResponseWriter, algorithm string, sigServices rufus.SigningServiceIndex) rufus.SigningService {
	if algorithm == "" {
		http.Error(w, "algorithm not specified", http.StatusBadRequest)
		return nil
	}

	service := sigServices[algorithm]

	if service == nil {
		http.Error(w, "algorithm "+algorithm+" not supported", http.StatusBadRequest)
		return nil
	}

	return service
}

// KeyInfo returns a Handler that given a specific Key ID param, returns the public key bits of that key
func KeyInfo(sigServices rufus.SigningServiceIndex) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		algorithm := vars["Algorithm"]

		sigService := getSigningService(w, algorithm, sigServices)
		if sigService == nil {
			// Error handled inside getSigningService
			return
		}

		keyInfo := &pb.KeyInfo{ID: vars["ID"], Algorithm: &pb.Algorithm{Algorithm: algorithm}}
		key, err := sigService.KeyInfo(keyInfo)
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
func CreateKey(sigServices rufus.SigningServiceIndex) http.Handler {
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
func DeleteKey(sigServices rufus.SigningServiceIndex) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var keyInfo *pb.KeyInfo
		err := json.NewDecoder(r.Body).Decode(&keyInfo)
		defer r.Body.Close()
		if err != nil || keyInfo.ID == "" || keyInfo.Algorithm == nil || keyInfo.Algorithm.Algorithm == "" {
			w.WriteHeader(http.StatusBadRequest)
			jsonErr, _ := json.Marshal("Malformed request")
			w.Write([]byte(jsonErr))
			return
		}

		sigService := getSigningService(w, keyInfo.Algorithm.Algorithm, sigServices)
		if sigService == nil {
			// Error handled inside getSigningService
			return
		}

		_, err = sigService.DeleteKey(keyInfo)

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
func Sign(sigServices rufus.SigningServiceIndex) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var sigRequest *pb.SignatureRequest
		err := json.NewDecoder(r.Body).Decode(&sigRequest)
		defer r.Body.Close()
		if err != nil || sigRequest.Content == nil ||
			sigRequest.KeyInfo == nil || sigRequest.KeyInfo.Algorithm == nil {
			w.WriteHeader(http.StatusBadRequest)
			jsonErr, _ := json.Marshal("Malformed request")
			w.Write([]byte(jsonErr))
			return
		}

		sigService := getSigningService(w, sigRequest.KeyInfo.Algorithm.Algorithm, sigServices)
		if sigService == nil {
			// Error handled inside getSigningService
			return
		}

		signer, err := sigService.Signer(sigRequest.KeyInfo)
		if err == keys.ErrInvalidKeyID {
			w.WriteHeader(http.StatusNotFound)
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
