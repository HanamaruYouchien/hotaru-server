package middleware

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
)

func WithHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := w.Header()
		header.Add("Access-Control-Allow-Origin", "*")
		header.Add("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		header.Add("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Authorization")

		header.Add("Content-Type", "application/json; charset=utf-8")

		next.ServeHTTP(w, r)
	})
}

func MaxBodyLength(length int64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rr := new(http.Request)
			*rr = *r
			rr.Body = http.MaxBytesReader(w, rr.Body, length)
			next.ServeHTTP(w, rr)
		})
	}
}

const CtxKeyObject = "object"

func Bind[T any]() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			obj := new(T)
			raw, err := io.ReadAll(r.Body)
			if err != nil {
				if _, ok := err.(*http.MaxBytesError); ok {
					ErrorTooLarge(w)
				} else {
					ErrorUnknown(w)
				}
				return
			}

			err = json.Unmarshal(raw, obj)
			if err != nil {
				ErrorNotJson(w) // ErrCodeBadJson should be tested manually
				return
			}

			rr := r.WithContext(context.WithValue(r.Context(), CtxKeyObject, obj))
			next.ServeHTTP(w, rr)
		})
	}
}

func GetObject(r *http.Request) any {
	return r.Context().Value(CtxKeyObject)
}
