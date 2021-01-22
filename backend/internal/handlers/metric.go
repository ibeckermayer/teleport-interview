package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/ibeckermayer/teleport-interview/backend/internal/database"
	"github.com/ibeckermayer/teleport-interview/backend/internal/util"
)

// MetricsPostHandler handles calls to "api/metrics"
type MetricsPostHandler struct {
	db *database.Database
}

// NewMetricsPostHandler creates a new MetricsHandler
func NewMetricsPostHandler(db *database.Database) *MetricsPostHandler {
	return &MetricsPostHandler{db}
}

type metricsPostRequestBody struct {
	AccountID string    `json:"account_id"`
	UserID    string    `json:"user_id"`
	Timestamp time.Time `json:"timestamp"`
}

// Handles "/api/metrics" POST requests. Should be wrapped with WithAPIkeyAuth middlewear
func (mph *MetricsPostHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var body metricsPostRequestBody

	// Decode json
	err := util.DecodeJSONBody(w, r, &body)
	if err != nil {
		util.HandleJSONdecodeError(w, err)
		return
	}

	// Save the metric to the database
	// TODO: should metric be saved regardless of whether CreateUser below fails?
	if err := mph.db.CreateMetric(body.AccountID, body.UserID, body.Timestamp); err != nil {
		log.Println(err)
		util.ErrorJSON(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Check if this user exists
	_, err = mph.db.GetUser(body.UserID)
	if err != nil {
		// If user DNE, save new user to the database
		if err := mph.db.CreateUser(body.UserID, body.AccountID); err != nil {
			log.Println(err)
			util.ErrorJSON(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}
}
