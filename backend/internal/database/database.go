package database

import (
	"github.com/ibeckermayer/teleport-interview/backend/internal/model"
	"github.com/jmoiron/sqlx"
	"github.com/pborman/uuid"
)

// Config is a database config object.
// Env determines whether the production or development database is created/used;
// if \"dev\", the app will seed the database with sample data for manual testing.
type Config struct {
	Env string
}

// Database is a handle to the database layer
// TODO: This should eventually take a database.Config object to swap out
// different drivers/settings/security features
type Database struct {
	cfg Config
	db  *sqlx.DB
}

// New creates a new *Database and initializes it's schema.
// Set filldb to true to fill the database with some fake data (for development purposes).
func New(cfg Config) (*Database, error) {
	dbfile := "./teleport-interview-" + cfg.Env + ".db"
	sqlxdb, err := sqlx.Open("sqlite3", dbfile)
	if err != nil {
		return nil, err
	}

	// force a connection and test that it worked
	if err = sqlxdb.Ping(); err != nil {
		return nil, err
	}

	db := &Database{cfg, sqlxdb}
	if err = db.init(); err != nil {
		return nil, err
	}

	return db, nil
}

func (db *Database) init() error {

	if _, err := db.db.Exec(model.AccountTableSQL); err != nil {
		return err
	}

	if db.cfg.Env == "dev" {
		// Fill the db with data for development purposes
		// Ignore error so the server doesn't panic every time it's recompiled due to the email column failing the UNIQUE check
		db.CreateAccount(uuid.New(), "dev@goteleport.com", "dev@goteleport.com")
	}

	return nil
}
