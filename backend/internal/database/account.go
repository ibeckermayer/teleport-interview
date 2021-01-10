package database

import (
	"time"

	"github.com/ibeckermayer/teleport-interview/backend/internal/auth"
)

var accountSchema = `CREATE TABLE IF NOT EXISTS account (
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

func (db *Database) insert(a *Account) error {
	_, err := db.db.NamedExec("INSERT INTO account (account_id, plan, email, password_hash, created_at, updated_at) VALUES (:account_id, :plan, :email, :password_hash, :created_at, :updated_at)", a)
	return err
}

// CreateAccount creates a new account and saves it in the database. Returns the Account that was created.
func (db *Database) CreateAccount(accountID string, email string, password string) error {
	passwordHash, err := auth.HashPassword(password)
	if err != nil {
		return err
	}
	account := &Account{accountID, "FREE", email, passwordHash, time.Now(), time.Now()}
	err = db.insert(account)
	return err
}
