package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"

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
func AuthMiddleware(auth authentication.Auth, log *logrus.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			bearToken := r.Header.Get("Authorization")
			strArr := strings.Split(bearToken, " ")
			if len(strArr) != 2 {
				log.Errorf("failed to split token, value: %v", bearToken)
				response.RenderResponse(w, http.StatusUnauthorized, response.EmptyResp{})
				return
			}

			accessDetails, err := auth.ValidateToken(strArr[1])
			if err != nil {
				log.Errorf("failed to validate token, error: %v", err)
				response.RenderResponse(w, http.StatusUnauthorized, response.EmptyResp{})
				return
			}

			userUUID, err := auth.GetUserUUID(accessDetails.AccessUUID)
			if err != nil || userUUID != accessDetails.UserUUID {
				log.Errorf("failed to get user uuid, userUUID: %v, error: %v", accessDetails.UserUUID, err)
				response.RenderResponse(w, http.StatusUnauthorized, response.EmptyResp{})
				return
			}

			ctx := context.WithValue(r.Context(), UserUUID, userUUID)
			ctx = context.WithValue(ctx, AccessUUID, accessDetails.AccessUUID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
