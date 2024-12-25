package crypto

import "encoding/base64"

const DeviceIDByteLength = 6
const RoomIDByteLength = 12

func GenerateDeviceID() (string, error) {
	b, err := CryptoRandomBytes(DeviceIDByteLength)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func GenerateRoomID() (string, error) {
	b, err := CryptoRandomBytes(RoomIDByteLength)
	if err != nil {
		return "", nil
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
