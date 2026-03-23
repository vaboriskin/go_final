package api

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type SignInReq struct {
	Password string `json:"password"`
}

type SignInResp struct {
	Token string `json:"token"`
	Error string `json:"error,omitempty"`
}

func signInHandler(w http.ResponseWriter, r *http.Request) {
	var req SignInReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeJSON(w, SignInResp{Error: err.Error()})
		return
	}

	pass := os.Getenv("TODO_PASSWORD")
	if pass == "" || pass != req.Password {
		writeJSON(w, SignInResp{Error: "неверный пароль"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"hash": hashPassword(pass),
		"exp":  time.Now().Add(8 * time.Hour).Unix(),
	})
	tokenString, err := token.SignedString([]byte(pass))
	if err != nil {
		writeJSON(w, SignInResp{Error: "Ошибка генерации токена"})
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   8 * 3600,
	})

	writeJSON(w, SignInResp{Token: tokenString})
}

func hashPassword(pass string) string {
	var sum int
	for _, c := range pass {
		sum += int(c)
	}
	return string(rune(sum))
}

func auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pass := os.Getenv("TODO_PASSWORD")
		if pass == "" {
			next(w, r)
			return
		}

		var tokenString string

		if cookie, err := r.Cookie("token"); err == nil {
			tokenString = cookie.Value
		}

		if tokenString == "" {
			authHeader := r.Header.Get("Authorization")
			if strings.HasPrefix(authHeader, "Bearer ") {
				tokenString = strings.TrimPrefix(authHeader, "Bearer ")
			}
		}

		if tokenString == "" {
			http.Error(w, "Authentification required", http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(pass), nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Authentification required", http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if claims["hash"] != hashPassword(pass) {
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}
		} else {
			http.Error(w, "Authentification required", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}
