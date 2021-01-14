package database

import (
	"time"

	"github.com/ibeckermayer/teleport-interview/backend/internal/auth"
	"github.com/ibeckermayer/teleport-interview/backend/internal/model"
)

func (db *Database) insertAccount(a *model.Account) error {
	_, err := db.db.NamedExec("INSERT INTO account (account_id, plan, email, password_hash, created_at, updated_at) VALUES (:account_id, :plan, :email, :password_hash, :created_at, :updated_at)", a)
	return err
}

// GetAccountByEmail retrieves an Account from the database by email address
func (db *Database) GetAccountByEmail(email string) (model.Account, error) {
	a := model.Account{}
	err := db.db.Get(&a, "SELECT * FROM account WHERE email=$1", email)
	return a, err
}

// CreateAccount creates a new account and saves it in the database
func (db *Database) CreateAccount(accountID, email, password string) error {
	passwordHash, err := auth.HashPassword(password)
	if err != nil {
		return err
	}
	account := &model.Account{
		AccountID:    accountID,
		Plan:         model.FREE,
		Email:        email,
		PasswordHash: passwordHash,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	err = db.insertAccount(account)
	return err
}
