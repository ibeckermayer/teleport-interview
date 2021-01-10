package database

import (
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/pborman/uuid"
)

// Config is a database config object
type Config struct {
	env string
}

// NewConfig returns a new Config for the database
func NewConfig(env string) Config {
	return Config{env}
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
	dbfile := "./teleport-interview-" + cfg.env + ".db"
	sqlxdb, err := sqlx.Open("sqlite3", dbfile)
	if err != nil {
		return &Database{}, err
	}

	// force a connection and test that it worked
	if err = sqlxdb.Ping(); err != nil {
		return &Database{}, err
	}

	db := &Database{cfg, sqlxdb}
	if err = db.init(); err != nil {
		return &Database{}, err
	}

	return db, nil
}

func (db *Database) init() error {

	if _, err := db.db.Exec(accountSchema); err != nil {
		return err
	}

	if db.cfg.env == "dev" {
		// Fill the db with data for development purposes
		// Ignore error so the server doesn't panic every time it's recompiled due to the email column failing the UNIQUE check
		err := db.CreateAccount(uuid.New(), "dev@goteleport.com", "dev@goteleport.com")
		log.Println(err)
	}

	return nil
}
