package storage

import (
	"errors"
	"time"

	"hotaru.hana.im/server/pkg/crypto"
	"xorm.io/xorm/schemas"
)

type Device struct {
	Localpart   string `xorm:"pk"`
	DeviceID    string `xorm:"pk 'device_id'"`
	DisplayName string
	AccessToken string    `xorm:"notnull unique index"`
	CreatedAt   time.Time `xorm:"created"`
}

var ErrDeviceNotExist = errors.New("device not exist")
var ErrEmptyAccessToken = errors.New("empty access token")
var ErrEmptyLocalpart = errors.New("empty localpart")

func (db *Storage) CreateDevice(localpart, deviceID, displayName string) (accessToken string, err error) {
	if accessToken, err = crypto.GenerateAccessToken(); err != nil {
		return
	}
	_, err = db.engine.InsertOne(&Device{
		Localpart:   localpart,
		DeviceID:    deviceID,
		DisplayName: displayName,
		AccessToken: accessToken,
	})
	if err != nil {
		return "", err
	}
	return
}

func (db *Storage) UpdateAccessToken(localpart, deviceID, displayName string) (accessToken string, err error) {
	if accessToken, err = crypto.GenerateAccessToken(); err != nil {
		return
	}
	if _, err = db.engine.ID(schemas.PK{localpart, deviceID}).Update(&Device{AccessToken: accessToken}); err != nil {
		return "", err
	}
	return
}

func (db *Storage) IsDeviceExist(localpart, deviceID string) error {
	if localpart == "" || deviceID == "" {
		return ErrDeviceNotExist
	}
	has, err := db.engine.ID(schemas.PK{localpart, deviceID}).Exist(&Device{})
	if err != nil {
		return err
	}
	if !has {
		return ErrDeviceNotExist
	}
	return nil
}

func (db *Storage) GetDeviceByAccessToken(accessToken string) (*Device, error) {
	if accessToken == "" {
		return nil, ErrEmptyAccessToken
	}

	dev := &Device{
		AccessToken: accessToken,
	}
	has, err := db.engine.Get(dev)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, ErrDeviceNotExist
	}

	return dev, nil
}

func (db *Storage) DeleteDevice(localpart, deviceID string) error {
	_, err := db.engine.ID(schemas.PK{localpart, deviceID}).Delete(&Device{})
	return err
}

func (db *Storage) DeleteDeviceByLocalpart(localpart string) error {
	if localpart == "" {
		return ErrEmptyLocalpart
	}
	_, err := db.engine.Delete(&Device{Localpart: localpart})
	return err
}
