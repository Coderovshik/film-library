package middleware

import (
	"errors"
	"log"
	"net/http"

	"github.com/Coderovshik/film-library/internal/util"
	"github.com/golang-jwt/jwt/v5"
)

func NewAuthMiddleware(key string, adminOnly bool) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("jwt")
			if err != nil {
				if errors.Is(err, http.ErrNoCookie) {
					log.Printf("ERROR: no auth cookie\n")
					util.Unauthorized(w, r)
					return
				}

				log.Printf("ERROR: failed to get auth cookie")
				util.InternalServerError(w, r)
				return
			}

			token := cookie.Value
			uc, err := util.ParseUserClaims(token, key)
			if err != nil {
				if errors.Is(err, jwt.ErrTokenExpired) {
					log.Printf("ERROR: token expired")
					util.Unauthorized(w, r)
					return
				}
				if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
					log.Printf("ERROR: jwt signature is invalid\n")
					util.Unauthorized(w, r)
					return
				}
				if errors.Is(err, util.ErrUnknownClaimsType) {
					log.Printf("ERROR: unknown token claims\n")
					util.Unauthorized(w, r)
					return
				}

				log.Printf("ERROR: failed to parse claims err=%s", err.Error())
				util.InternalServerError(w, r)
				return
			}

			log.Printf("INFO: authenticated user=%s user_id=%d", uc.Username, uc.ID)
			if adminOnly {
				if !uc.IsAdmin {
					log.Printf("ERROR: permission denied")
					util.Forbidden(w, r)
					return
				}

				log.Printf("INFO: admin request")
			}

			next.ServeHTTP(w, r)
		})
	}
}
