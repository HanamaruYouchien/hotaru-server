package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"hotaru.hana.im/server/pkg/storage"
)

const CtxKeyAuth = "authentication"
const CtxKeyDevice = "device"
const CtxKeyAccount = "account"

func Auth(db *storage.Storage, required bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" { // no token
				if required {
					ErrorMissingToken(w)
				} else {
					rr := r.WithContext(context.WithValue(r.Context(), CtxKeyAuth, false))
					next.ServeHTTP(w, rr)
				}
				return
			}
			if !strings.HasPrefix(authHeader, "Bearer ") {
				ErrorUnknownToken(w)
				return
			}

			accessToken := strings.TrimPrefix(authHeader, "Bearer ")
			accessToken = strings.TrimSpace(accessToken)

			dev, err := db.GetDeviceByAccessToken(accessToken)
			if err != nil {
				switch {
				case errors.Is(err, storage.ErrEmptyAccessToken):
					fallthrough
				case errors.Is(err, storage.ErrDeviceNotExist):
					ErrorUnknownToken(w)
				default:
					ErrorUnknownMsg(w, "unknown error")
				}
				return
			}

			account, err := db.GetAccount(dev.Localpart)
			if err != nil {
				switch {
				case errors.Is(err, storage.ErrAccountNotExist):
					ErrorUnknownToken(w)
				default:
					ErrorUnknownMsg(w, "unknown error")
				}
				return
			}

			rr := r.WithContext(context.WithValue(r.Context(), CtxKeyAuth, true))
			rr = r.WithContext(context.WithValue(rr.Context(), CtxKeyAccount, account))
			rr = r.WithContext(context.WithValue(rr.Context(), CtxKeyDevice, dev))

			next.ServeHTTP(w, rr)
		})
	}
}

func AuthRequired(db *storage.Storage) func(http.Handler) http.Handler {
	return Auth(db, true)
}

func AuthOptional(db *storage.Storage) func(http.Handler) http.Handler {
	return Auth(db, false)
}

func IsAuthenticated(r *http.Request) bool {
	return r.Context().Value(CtxKeyAuth).(bool)
}

func GetDevice(r *http.Request) *storage.Device {
	dev, _ := r.Context().Value(CtxKeyDevice).(*storage.Device)
	return dev
}

func GetAccount(r *http.Request) *storage.Account {
	acc, _ := r.Context().Value(CtxKeyAccount).(*storage.Account)
	return acc
}
