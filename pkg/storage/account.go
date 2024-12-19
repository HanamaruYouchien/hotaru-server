package storage

import "time"

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

func (db *Storage) CreateAccount(localpart, password string, accountType AccountType) {}

func (db *Storage) Ping() error {
	return db.engine.Ping()
}
