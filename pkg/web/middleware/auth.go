package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"hotaru.hana.im/server/pkg/storage"
)

const CtxKeyDevice = "device"
const CtxKeyAccount = "account"

func Auth(db *storage.Storage) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				ErrorMissingToken(w)
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

			rr := r.WithContext(context.WithValue(r.Context(), CtxKeyAccount, account))
			rr = r.WithContext(context.WithValue(rr.Context(), CtxKeyDevice, dev))

			next.ServeHTTP(w, rr)
		})
	}
}

func GetDevice(r *http.Request) *storage.Device {
	return r.Context().Value(CtxKeyDevice).(*storage.Device)
}

func GetAccount(r *http.Request) *storage.Account {
	return r.Context().Value(CtxKeyAccount).(*storage.Account)
}
