package crypto

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"

	"golang.org/x/crypto/argon2"
)

const (
	hashTime    uint32 = 2
	hashMemory  uint32 = 64 * 1024
	hashKeyLen  uint32 = 50
	hashThreads uint8  = 8
)

const SaltByteLength = 16

func GeneratePasswordSalt() ([]byte, error) {
	b, err := CryptoRandomBytes(SaltByteLength)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func CryptoRandomBytes(size int) ([]byte, error) {
	buf := make([]byte, size)
	_, err := rand.Read(buf)
	return buf, err
}

func HashWithSalt(password string, salt []byte) string {
	return hex.EncodeToString(argon2.IDKey([]byte(password), salt, hashTime, hashMemory, hashThreads, hashKeyLen))
}

func VerifyPassword(hash, password string, salt []byte) bool {
	return subtle.ConstantTimeCompare([]byte(hash), []byte(HashWithSalt(password, salt))) == 1
}
