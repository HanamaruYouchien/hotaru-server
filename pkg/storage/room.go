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

	if err := db.IsRoomAliasExist(alias); err == nil {
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

	// TODO: Apply events
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

func (db *Storage) IsRoomAliasExist(alias string) error {
	has, err := db.engine.ID(alias).Exist(&RoomAlias{})
	if err != nil {
		return err
	}
	if !has {
		return ErrRoomAliasNotExist
	}
	return nil
}

func (db *Storage) GetRoomAliasesByRoomID(roomID string) []string {
	if roomID == "" {
		return []string{}
	}

	aliases := make([]RoomAlias, 0)
	db.engine.Where("room_id = ?", roomID).Find(&aliases)
	res := make([]string, 0, len(aliases))
	for _, v := range aliases {
		res = append(res, v.Alias)
	}
	return res
}

func (db *Storage) GetRoomIDByRoomAlias(alias string) (string, error) {
	roomAlias := &RoomAlias{}
	has, err := db.engine.ID(alias).Get(roomAlias)
	if err != nil {
		return "", err
	}
	if !has {
		return "", ErrRoomAliasNotExist
	}
	return roomAlias.RoomId, nil
}

func (db *Storage) CreateRoomAlias(alias, roomID string) error {
	roomAlias := &RoomAlias{
		Alias:  alias,
		RoomId: roomID,
	}
	if _, err := db.engine.InsertOne(roomAlias); err != nil {
		return err
	}
	return nil
}

func (db *Storage) DeleteRoomAlias(alias string) error {
	_, err := db.engine.ID(alias).Delete(&RoomAlias{})
	return err
}
