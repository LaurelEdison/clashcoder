package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/LaurelEdison/clashcoder/backend/handlers"
	users "github.com/LaurelEdison/clashcoder/backend/handlers/user"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

var jwtSecret []byte

func InitJWT(secret string) {
	jwtSecret = []byte(secret)
}

func generateJWT(userID uuid.UUID, role string) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID.String(),
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})

	return token.SignedString(jwtSecret)

}

func Login(h *handlers.Handlers) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type parameters struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		params := parameters{}

		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			h.RespondWithError(w, http.StatusBadRequest, "Invalid JSON")
			return
		}

		user, err := h.DB.GetUserByEmail(r.Context(), params.Email)
		if err != nil {
			h.RespondWithError(w, http.StatusUnauthorized, "Invalid email or password")
			return
		}
		if err := users.CheckPasswordHash(user.PasswordHash, params.Password); !err {
			h.RespondWithError(w, http.StatusUnauthorized, "Invalid email or password")
			return
		}

		token, err := generateJWT(user.ID, user.Role)
		if err != nil {
			h.RespondWithError(w, http.StatusInternalServerError, "Could not generate token")
			return
		}

		h.RespondWithJSON(w, 200, map[string]string{"token": token})

	}
}

func JWTAuthMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Missing or invalid token", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["Alg"])
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			uidStr, idOk := claims["user_id"].(string)
			role, roleOk := claims["role"].(string)
			if !idOk {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}
			if !roleOk {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}
			userID, err := uuid.Parse(uidStr)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}
			ctx := r.Context()
			ctx = context.WithValue(ctx, users.UserIDKey, userID)
			ctx = context.WithValue(ctx, users.UserRoleKey, role)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

	})
}
