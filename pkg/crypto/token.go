package crypto

import (
	"crypto/rsa"
	"time"

	"github.com/lestrrat-go/jwx/v3/jwa"
	"github.com/lestrrat-go/jwx/v3/jwk"
	"github.com/lestrrat-go/jwx/v3/jwt"
)

const issuer = "hotaru"

func GenerateToken(privateKey *rsa.PrivateKey) (token string, err error) {
	privKey, err := jwk.Import(privateKey)
	if err != nil {
		return
	}

	tok, err := jwt.NewBuilder().Issuer(issuer).IssuedAt(time.Now()).Expiration(time.Now().Add(time.Hour * 144)).Build()
	if err != nil {
		return
	}

	signed, err := jwt.Sign(tok, jwt.WithKey(jwa.RS256(), privKey))
	if err != nil {
		return
	}
	return string(signed), nil
}

func VerifyToken(publicKey *rsa.PublicKey, token string) error {
	pubKey, err := jwk.Import(publicKey)
	if err != nil {
		return err
	}

	_, err = jwt.Parse([]byte(token), jwt.WithKey(jwa.RS256(), pubKey))
	if err != nil {
		return err
	}
	return nil
}
