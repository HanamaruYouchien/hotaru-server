package storage

import "errors"

type Profile struct {
	Localpart   string `xorm:"pk"`
	DisplayName string
	AvatarUrl   string
}

var ErrProfileNotExist = errors.New("profile not exist")

func (db *Storage) GetProfile(localpart string) (*Profile, error) {
	profile := &Profile{}
	has, err := db.engine.ID(localpart).Get(profile)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, ErrProfileNotExist
	}
	return profile, nil
}

func (db *Storage) CreateProfile(localpart, displayName, avatarUrl string) error {
	profile := &Profile{
		Localpart:   localpart,
		DisplayName: displayName,
		AvatarUrl:   avatarUrl,
	}
	_, err := db.engine.InsertOne(profile)
	return err
}

func (db *Storage) IsProfileExist(localpart string) error {
	has, err := db.engine.ID(localpart).Exist(&Profile{})
	if err != nil {
		return err
	}
	if !has {
		return ErrProfileNotExist
	}
	return nil
}

func (db *Storage) UpdateProfileDisplayName(localpart, displayName string) error {
	_, err := db.engine.ID(localpart).Cols("display_name").Update(&Profile{DisplayName: displayName})
	return err
}

func (db *Storage) UpdateProfileAvatarUrl(localpart, avatarUrl string) error {
	_, err := db.engine.ID(localpart).Cols("avatar_url").Update(&Profile{AvatarUrl: avatarUrl})
	return err
}
