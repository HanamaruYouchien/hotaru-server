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
	IsLocked      bool
	Type          AccountType
}

type AccountType int

const (
	AccountTypeUser AccountType = iota
	AccountTypeAdmin
)

var ErrPasswordTooWeak = errors.New("password too weak")
var ErrAccountNotExist = errors.New("account not exist")
var ErrPasswordNotCorrect = errors.New("username or password not correct")
var ErrUserInUse = errors.New("user id is already taken")

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

func (db *Storage) VerifyAccount(localpart, password string) error {
	account := &Account{}
	isFind, err := db.engine.ID(localpart).Get(account)
	if err != nil {
		return err
	}
	if !isFind {
		return ErrAccountNotExist
	}

	saltBytes, err := hex.DecodeString(account.Salt)
	if err != nil {
		return err
	}
	if !crypto.VerifyPassword(account.PasswordHash, password, saltBytes) {
		return ErrPasswordNotCorrect
	}
	return nil
}

func (db *Storage) ValidateLocalpart(localpart string) error {
	// TODO: validate localpart format(https://spec.matrix.org/v1.12/appendices/#common-identifier-format)
	has, err := db.engine.ID(localpart).Exist(&Account{})
	if err != nil {
		return err
	}

	if has {
		return ErrUserInUse
	}
	return nil
}

func (db *Storage) Ping() error {
	return db.engine.Ping()
}

func validatePasswordStrength(password string) bool {
	return len(password) >= 8
}
