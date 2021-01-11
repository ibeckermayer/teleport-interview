package model

import (
	"time"
)

// AccountTableSQL is the SQL statement for creating a table corresponding to the Account model
var AccountTableSQL = `CREATE TABLE IF NOT EXISTS account (
	account_id CHARACTER(36) PRIMARY KEY,
	plan VARCHAR(50) NOT NULL,
	email VARCHAR(320) UNIQUE NOT NULL,
	password_hash CHARACTER(60) NOT NULL,
	created_at DATETIME NOT NULL,
	updated_at DATETIME NOT NULL);`

// Account represents a row in the "account" table.
type Account struct {
	AccountID    string    `db:"account_id"`
	Plan         string    `db:"plan"` // One of "FREE" or "ENTERPRISE"
	Email        string    `db:"email"`
	PasswordHash string    `db:"password_hash"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}
