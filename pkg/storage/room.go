package storage

import (
	"errors"
	"time"

	"hotaru.hana.im/server/pkg/crypto"
	"xorm.io/xorm/schemas"
)

type Room struct {
	RoomId            string `xorm:"pk"`
	Name              string
	Topic             string `xorm:"text"`
	RoomVersion       string
	Visibility        string
	JoinRules         string
	HistoryVisibility string
	GuestAccess       string
	Creator           string
	Type              string
	CreatedAt         time.Time `xorm:"created"`
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
var ErrAccountRoomNotExist = errors.New("account room not exist")
var ErrRoomAliasInUsed = errors.New("alias in used")
var ErrEmptyRoomID = errors.New("empty roomID")
var ErrNoPermission = errors.New("no permission")

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

func (db *Storage) GetRoomAliasesByRoomID(roomID string) ([]string, error) {
	if roomID == "" {
		return nil, ErrEmptyRoomID
	}

	aliases := make([]string, 0)
	if err := db.engine.Table(&RoomAlias{}).Where("room_id = ?", roomID).Cols("alias").Find(&aliases); err != nil {
		return nil, err
	}
	return aliases, nil
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

func (db *Storage) GetJoinedRooms(localpart string) (roomIDs []string, err error) {
	if localpart == "" {
		return nil, ErrEmptyLocalpart
	}

	roomIDs = make([]string, 0)
	if err := db.engine.Table(&AccountRoom{}).Where("localpart = ?", localpart).And("membership = ?", MembershipTypeJoined).Cols("room_id").Find(&roomIDs); err != nil {
		return nil, err
	}

	return roomIDs, nil
}

func (db *Storage) IsAccountRoomExist(roomID, localpart string) error {
	has, err := db.engine.ID(schemas.PK{localpart, roomID}).Exist(&AccountRoom{})
	if err != nil {
		return err
	}
	if !has {
		return ErrAccountRoomNotExist
	}
	return nil
}

func (db *Storage) GetAccountRoom(roomID, localpart string) (*AccountRoom, error) {
	accountRoom := &AccountRoom{}
	has, err := db.engine.ID(schemas.PK{localpart, roomID}).Get(accountRoom)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, ErrAccountRoomNotExist
	}
	return accountRoom, nil
}

func (db *Storage) InviteRoom(roomID, from, to string) error {
	if to == "" {
		return ErrEmptyLocalpart
	}
	if roomID == "" {
		return ErrEmptyRoomID
	}

	if accountRoomFrom, err := db.GetAccountRoom(roomID, from); err != nil {
		if errors.Is(err, ErrAccountRoomNotExist) {
			return ErrNoPermission
		} else {
			return err
		}
	} else {
		if accountRoomFrom.Membership != MembershipTypeJoined {
			return ErrNoPermission
		}
		// TODO: check permission
	}

	id := schemas.PK{to, roomID}
	if accountRoom, err := db.GetAccountRoom(roomID, to); err == nil {
		// TODO: use different error types
		switch accountRoom.Membership {
		case MembershipTypeBanned, MembershipTypeJoined:
			return ErrNoPermission
		}
		if _, err := db.engine.ID(id).Update(&AccountRoom{Membership: MembershipTypeInvited}); err != nil {
			return err
		}
	} else {
		if !errors.Is(err, ErrAccountRoomNotExist) {
			return err
		}
		accountRoom := &AccountRoom{
			Localpart:  to,
			RoomId:     roomID,
			Membership: MembershipTypeInvited,
		}
		if _, err := db.engine.InsertOne(accountRoom); err != nil {
			return err
		}
	}
	return nil
}

func (db *Storage) JoinRoom(roomID, localpart string) error {
	if localpart == "" {
		return ErrEmptyLocalpart
	}
	if roomID == "" {
		return ErrEmptyRoomID
	}

	id := schemas.PK{localpart, roomID}
	if accountRoom, err := db.GetAccountRoom(roomID, localpart); err == nil {
		// TODO: use different error types
		switch accountRoom.Membership {
		case MembershipTypeBanned, MembershipTypeJoined, MembershipTypeKnocking:
			return ErrNoPermission
		}
		// TODO: check permission
		if _, err := db.engine.ID(id).Update(&AccountRoom{Membership: MembershipTypeJoined, PowerLevel: 0}); err != nil {
			return err
		}
	} else {
		if !errors.Is(err, ErrAccountRoomNotExist) {
			return err
		}
		// TODO: check permission
		accountRoom := &AccountRoom{
			Localpart:  localpart,
			RoomId:     roomID,
			Membership: MembershipTypeJoined,
			PowerLevel: 0,
		}
		if _, err := db.engine.InsertOne(accountRoom); err != nil {
			return err
		}
	}
	return nil
}
