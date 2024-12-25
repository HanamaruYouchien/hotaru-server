package middleware

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/go-chi/httprate"
)

func WithHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := w.Header()
		header.Add("Access-Control-Allow-Origin", "*")
		header.Add("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		header.Add("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Authorization")

		// header.Add("Content-Type", "application/json; charset=utf-8")

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
					ErrorUnknown(w, http.StatusInternalServerError)
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

func RenderJSON(w http.ResponseWriter, resp any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	raw, _ := json.Marshal(resp)
	w.Write(raw)
}

// Changing the header map after a call to WriteHeader has no effect
func RenderJSONWithStatusCode(w http.ResponseWriter, status int, resp any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	raw, _ := json.Marshal(resp)
	w.Write(raw)
}

func RateLimiter( /*maxRequests int, duration time.Duration*/ ) func(http.Handler) http.Handler {
	maxRequests := 5
	duration := time.Millisecond * 500
	return httprate.Limit(
		maxRequests,
		duration,
		httprate.WithKeyFuncs(httprate.KeyByIP),
		httprate.WithLimitHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ErrorLimitExceeded(w)
		})),
	)
}
