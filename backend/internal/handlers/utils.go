package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/ibeckermayer/teleport-interview/backend/internal/auth"
)

type malformedRequest struct {
	status int    // HTTP status
	logMsg string // Detailed message for logging, not to be returned to clients
}

// Generic HTTP status message for sending to clients
func (mr malformedRequest) Error() string {
	return mr.logMsg
}

// Modified version of https://www.alexedwards.net/blog/how-to-properly-parse-a-json-request-body
// Attempts to decode a json request body into dst, returning an error if it fails
func decodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	if r.Header.Get("Content-Type") != "" {
		if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
			msg := "Content-Type header is not application/json"
			return &malformedRequest{status: http.StatusUnsupportedMediaType, logMsg: msg}
		}
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1048576) // 1MB

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&dst)
	if err != nil {
		return &malformedRequest{status: http.StatusBadRequest, logMsg: err.Error()}
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		msg := "Request body must only contain a single JSON object"
		return &malformedRequest{status: http.StatusBadRequest, logMsg: msg}
	}

	return nil
}

// Helper function for responding to requests that cause decodeJSONBody to throw an error.
// Logs a detailed error message and responds to the client with a generic HTTP error message.
// The caller should ensure no further writes are done to w.
func handleJSONdecodeError(w http.ResponseWriter, err error) {
	var mr *malformedRequest
	if errors.As(err, &mr) {
		log.Println(mr)
		http.Error(w, http.StatusText(mr.status), mr.status)
	} else {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

var (
	errAuthHeaderNotFound     = errors.New("Authorization header expected but not found")
	errAuthHeaderNotFormatted = errors.New("Authorization header was improperly formatted")
)

// Helper function to retreive a token sent in standard "Bearer" format from a request
// (https://tools.ietf.org/html/rfc6750#page-5). If the request doesn't contain an Authorization
// header or the Authorization header is improperly formatted, getBearerToken returns "".
// Handlers generally shouldn't call this function, and should instead call getSessionID or
// getApiKey (TODO) to specify which type of token they are expecting.
func getBearerToken(r *http.Request) (string, error) {
	reqToken := r.Header.Get("Authorization")
	if reqToken == "" {
		// Request did not contain an Authorization header
		return "", errAuthHeaderNotFound
	}
	splitToken := strings.Split(reqToken, "Bearer ")
	if len(splitToken) == 1 {
		// Split failed, request may have been improperly formatted
		return "", errAuthHeaderNotFormatted
	}
	return splitToken[1], nil
}

func getSessionID(r *http.Request) (auth.SessionID, error) {
	s, err := getBearerToken(r)
	return auth.SessionID(s), err
}

// handlerWithSession behaves like an ordinary http.Hander (https://golang.org/pkg/net/http/#Handler)
// but expects an auth.Session
type handlerWithSession interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request, session auth.Session)
	GetSessionManager() *auth.SessionManager
}

// WithSessionAuth is middlewear for protecting HandlerWithSession's with sessionID auth
func WithSessionAuth(hws handlerWithSession) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get sessionID from the request
		sessionID, err := getSessionID(r)
		if err != nil {
			// Could not get sessionID, return 401
			log.Println(err)
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		// Attempt to get the corresponding session from the SessionManager
		sm := hws.GetSessionManager()
		session, err := sm.GetSession(sessionID)
		if err != nil {
			// Session does not exist or timed out
			log.Println(err)
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		// User is authorized, call the HandlerWithSession
		hws.ServeHTTP(w, r, session)
	})
}
