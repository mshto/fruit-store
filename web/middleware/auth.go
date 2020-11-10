package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/mshto/fruit-store/authentication"
	"github.com/mshto/fruit-store/web/common/response"
)

type contextKey int

// sadassa
const (
	UserUUID contextKey = iota
	AccessUUID
)

// AuthMiddleware AuthMiddleware
func AuthMiddleware(auth authentication.Auth) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			bearToken := r.Header.Get("Authorization")
			strArr := strings.Split(bearToken, " ")
			if len(strArr) != 2 {
				response.RenderResponse(w, http.StatusUnauthorized, response.EmptyResp{})
				return
			}

			accessDetails, err := auth.ValidateToken(strArr[1])
			if err != nil {
				response.RenderResponse(w, http.StatusUnauthorized, response.EmptyResp{})
				return
			}

			userUUID, err := auth.GetUserUUID(accessDetails.AccessUUID)
			if err != nil || userUUID != accessDetails.UserUUID {
				response.RenderResponse(w, http.StatusUnauthorized, response.EmptyResp{})
				return
			}

			ctx := context.WithValue(r.Context(), UserUUID, userUUID)
			ctx = context.WithValue(ctx, AccessUUID, accessDetails.AccessUUID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
