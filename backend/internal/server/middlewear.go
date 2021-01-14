package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/ibeckermayer/teleport-interview/backend/internal/auth"
	"github.com/ibeckermayer/teleport-interview/backend/internal/handlers"
)

func getAPIkey(r *http.Request) (auth.Key, error) {
	s, err := auth.GetBearerToken(r)
	return auth.Key(s), err
}

// getAccountIDfromBody attempts to get the "account_id" field from the *http.Request body without
// mutating the *http.Request (see https://stackoverflow.com/a/47295689/6277051). Handles http error
// responses and returns an error if "account_id" field doesn't exist or other error occurs. Caller
// should return immediately on non-nil error
func getAccountIDfromBody(w http.ResponseWriter, r *http.Request) (string, error) {
	// Check that request body exists
	if r.Body == nil {
		err := errors.New("request body was nil")
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
		return "", err
	}

	r.Body = http.MaxBytesReader(w, r.Body, handlers.BodyMaxSize) // limit request body size

	// Read request body into buffer
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return "", err
	}

	// A map container to decode the JSON body into
	b := make(map[string]json.RawMessage)

	// Unmarshal json
	if err := json.Unmarshal(bodyBytes, &b); err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return "", err
	}

	// Extract account_id
	accountIDraw, ok := b["account_id"]
	if !ok {
		err := errors.New("request body missing required field \"account_id\"")
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
		return "", err
	}
	accountID := string(accountIDraw)

	// Now that accountID is extracted, restore request body to its original state so that next can use it
	r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	return accountID, nil
}

// WithAPIkeyAuth is a middlewear function for protecting handlers for routes that require an API key.
// APIkey protected requests require that the sender send an API key in the Authorization header, as well
// as its corresponding account_id field in the request's body.
func (srv *Server) WithAPIkeyAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the api key from the Authorization header
		key, err := getAPIkey(r)
		if err != nil {
			// Could not get APIkey, return 401
			log.Println(err)
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		// Get accoutID from the request body
		accountID, err := getAccountIDfromBody(w, r)
		if err != nil {
			return
		}

		// Get corresponding row from the apikey table in the database
		apikey, err := srv.db.GetAPIkey(accountID)
		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		// Check that key hashes match
		if !auth.CheckKeyHash(key, apikey.KeyHash) {
			log.Printf("recieved invalid api key for account %v", accountID)
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		// Request authorized, call next
		next.ServeHTTP(w, r)
	})
}
