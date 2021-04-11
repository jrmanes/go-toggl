package middleware

import (
	"context"
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/jrmanes/go-toggl/pkg/claim"
	"github.com/jrmanes/go-toggl/pkg/response"
)

func Authorizator(next http.Handler) http.Handler {
	signingString := os.Getenv("JWT_SIGNING_STRING")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("Authorization")
		tokenString, err := tokenFromAuthorization(authorization)
		if err != nil {
			response.HTTPError(w, r, http.StatusUnauthorized, err.Error())
			return
		}

		// call package claim an parse the token received
		c, err := claim.GetFromToken(tokenString, signingString)
		if err != nil {
			response.HTTPError(w, r, http.StatusUnauthorized, err.Error())
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, "id", c.ID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func tokenFromAuthorization(authorization string) (string, error) {
	if authorization == "" {
		return "", errors.New("auth is required")
	}

	if !strings.HasPrefix(authorization, "Bearer") {
		return "", errors.New("invalid autorization format")
	}

	l := strings.Split(authorization, " ")
	if len(l) != 2 {
		return "", errors.New("invalid autorization format")
	}

	return l[1], nil
}
