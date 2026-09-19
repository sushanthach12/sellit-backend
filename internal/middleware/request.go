package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type ctxKey int

const (
	requestIdKey ctxKey = iota
)

const (
	requestId = "x-request-id"
)

func RequestId(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get(requestId)
		if id == "" {
			id = uuid.NewString()
		}

		w.Header().Add(requestId, id)
		// here we are using custom type requestIdKey instead of the "requestId" , because to avoid
		// any overlap of the same key
		ctx := context.WithValue(r.Context(), requestIdKey, id)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Never make the requestIdKey public, make it private and have a method to retrieve it in the same package
func GetRequestIdFromContext(ctx context.Context) string {
	requestId := ctx.Value(requestIdKey).(string) // convert the any -> string using .(string)
	return requestId
}
