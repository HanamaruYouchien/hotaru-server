package storage

import (
	"encoding/hex"
	"errors"
	"regexp"
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

var (
	ErrPasswordTooWeak    = errors.New("password too weak")
	ErrAccountNotExist    = errors.New("account not exist")
	ErrPasswordNotCorrect = errors.New("username or password not correct")
	ErrUserInUse          = errors.New("user id is already taken")
	ErrInvalidUsername    = errors.New("invalid username")
)

func (db *Storage) CreateAccount(localpart, password string, accountType AccountType) error {
	if !validatePasswordStrength(password) {
		return ErrPasswordTooWeak
	}

	if err := db.ValidateLocalpart(localpart); err != nil {
		return err
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

func (db *Storage) GetAccount(localpart string) (*Account, error) {
	account := &Account{}
	has, err := db.engine.ID(localpart).Get(account)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, ErrAccountNotExist
	}
	return account, nil
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
	if err := validateLocalpart(localpart); err != nil {
		return err
	}

	has, err := db.engine.ID(localpart).Exist(&Account{})
	if err != nil {
		return err
	}

	if has {
		return ErrUserInUse
	}
	return nil
}

var localpartValidator = regexp.MustCompile(`^[0-9a-z_\-+=./]+$`)

const MaxLocalpartLength = 255

func validateLocalpart(localpart string) error {
	if len(localpart) > MaxLocalpartLength || !localpartValidator.MatchString(localpart) {
		return ErrInvalidUsername
	}
	return nil
}

func (db *Storage) Ping() error {
	return db.engine.Ping()
}

func validatePasswordStrength(password string) bool {
	return len(password) >= 8
}
