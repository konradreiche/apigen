package server

import (
	"context"
	"net/http"
	"strings"
)

func metadataMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(context.WithValue(r.Context(), "userAgent", r.UserAgent()))
		r = r.WithContext(context.WithValue(r.Context(), "ipAddress", getIP(r)))
		next.ServeHTTP(w, r)
	})
}

func getIP(req *http.Request) string {
	xff := req.Header.Get("X-Forwarded-For")
	if xff != "" {
		return strings.Split(xff, ",")[0]
	}
	remoteAddr := req.RemoteAddr
	return strings.Split(remoteAddr, ":")[0]
}
