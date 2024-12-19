package storage

import (
	"encoding/hex"
	"errors"
	"time"

	"hotaru.hana.im/server/pkg/crypto"
)

type Account struct {
	Localpart     string `xorm:"pk"`
	PasswordHash  string
	Salt          string
	CreatedAt     time.Time `xorm:"created"`
	IsDeactivated bool
	Type          AccountType
}

type AccountType int

const (
	AccountTypeUser AccountType = iota
	AccountTypeAdmin
)

var ErrPasswordTooWeak = errors.New("password too weak")

func (db *Storage) CreateAccount(localpart, password string, accountType AccountType) error {
	if !validatePasswordStrength(password) {
		return ErrPasswordTooWeak
	}

	salt, err := crypto.GeneratePasswordSalt()
	if err != nil {
		return err
	}
	passwordHash := crypto.HashWithSalt(password, salt)
	account := &Account{
		Localpart:    localpart,
		PasswordHash: passwordHash,
		Salt:         hex.EncodeToString(salt),
		Type:         AccountTypeUser,
	}

	_, err = db.engine.InsertOne(account)
	return err
}

func (db *Storage) Ping() error {
	return db.engine.Ping()
}

func validatePasswordStrength(password string) bool {
	return len(password) >= 8
}
