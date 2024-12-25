package storage

import (
	"errors"
	"time"

	"hotaru.hana.im/server/pkg/crypto"
)

type Room struct {
	RoomId      string `xorm:"pk"`
	Name        string
	Topic       string `xorm:"text"`
	RoomVersion string
	Visibility  string
	Creator     string
	Type        string
	CreatedAt   time.Time `xorm:"created"`
}

type RoomAlias struct {
	Alias  string `xorm:"pk"`
	RoomId string
}

type AccountRoom struct {
	Localpart  string `xorm:"pk"`
	RoomId     string `xorm:"pk"`
	Membership MembershipType
	PowerLevel int
}

type MembershipType int

const (
	MembershipTypeUnrelated MembershipType = iota
	MembershipTypeKnocking
	MembershipTypeInvited
	MembershipTypeJoined
	MembershipTypeBanned
)

var ErrRoomAliasNotExist = errors.New("alias not exist")
var ErrRoomAliasInUsed = errors.New("alias in used")

func (db *Storage) CreateRoom(name, topic, visibility, creator, alias string) (string, error) {
	roomID, err := crypto.GenerateRoomID()
	if err != nil {
		return "", err
	}

	if err := db.IsAliasExist(alias); err == nil {
		return "", ErrRoomAliasInUsed
	} else if !errors.Is(err, ErrRoomAliasNotExist) {
		return "", err
	}

	room := &Room{
		RoomId:      roomID,
		Name:        name,
		Topic:       topic,
		RoomVersion: "v10",
		Visibility:  visibility,
		Creator:     creator,
		Type:        "",
	}
	roomAlias := &RoomAlias{
		Alias:  alias,
		RoomId: roomID,
	}
	accountRoom := &AccountRoom{
		Localpart:  creator,
		RoomId:     roomID,
		Membership: MembershipTypeJoined,
		PowerLevel: 100,
	}

	session := db.engine.NewSession()
	defer session.Close()

	if err := session.Begin(); err != nil {
		return "", nil
	}
	if _, err := session.Insert(room, accountRoom); err != nil {
		return "", err
	}
	if alias != "" {
		if _, err := session.Insert(roomAlias); err != nil {
			return "", err
		}
	}
	if err := session.Commit(); err != nil {
		return "", nil
	}

	return roomID, nil
}

func (db *Storage) IsAliasExist(alias string) error {
	has, err := db.engine.ID(alias).Exist(&RoomAlias{})
	if err != nil {
		return err
	}
	if !has {
		return ErrRoomAliasNotExist
	}
	return nil
}
