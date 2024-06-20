package httpserver

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

func (h *handler) middlewareRequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Started %s %s", r.Method, r.URL.Path)

		// Call the next handler
		next.ServeHTTP(w, r)

		log.Printf("Completed %s in %v", r.URL.Path, time.Since(start))
	})
}

func (h *handler) middlewareAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Started %s %s", r.Method, r.URL.Path)

		pattern := fmt.Sprintf("%s %s", r.Method, r.URL.Path)
		rc, ok := h.routes[pattern]
		if !ok {
			log.Printf("%s %s - not configured", r.Method, r.URL.Path)
			h.writePlainResponse(w, http.StatusNotFound, "")
			return
		}

		var authToken string
		if !rc.authUser {
			next.ServeHTTP(w, r)
			return
		}

		authToken = r.Header.Get("Authorization")
		split := strings.Split(authToken, " ")
		if len(split) < 2 {
			respCode := http.StatusUnauthorized
			h.writePlainResponse(w, respCode, "invalid auth token")
			log.Printf("Completed (unauthenticated) %s with response code (%d) in %v", r.URL.Path, respCode, time.Since(start))
			return
		}
		authToken = split[1]

		userID, err := h.dateService.AuthenticateUserToken(authToken)
		if err != nil {
			respCode := http.StatusUnauthorized
			h.logger.Error("authenticate user token", err)
			h.writePlainResponse(w, respCode, "invalid auth token")
			log.Printf("Completed (unauthenticated) %s with response code (%d) in %v", r.URL.Path, respCode, time.Since(start))
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), ctxKeySessionUserID, userID))

		log.Printf("Completed (authenticated) %s in %v", r.URL.Path, time.Since(start))
		next.ServeHTTP(w, r)
	})
}
