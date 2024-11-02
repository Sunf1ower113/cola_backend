package midlleware

import (
	"bytes"
	"context"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

var secretKey = []byte(os.Getenv("SECRET_KEY"))

type Claims struct {
	UserID int64  `json:"user_id"`
	Role   string `json:"role"` // Добавляем роль в токен
	jwt.StandardClaims
}

func LoggerRequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request: %s %s\n", r.Method, r.URL.Path)
		if r.Body != nil {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				log.Println("Failed to read request body:", err)
			}
			defer r.Body.Close()
			log.Println("Request Body:", string(body))
			r.Body = io.NopCloser(bytes.NewBuffer(body))
		}
		next.ServeHTTP(w, r)
	})
}

func TimeoutMiddleware(next http.Handler) http.Handler {
	log.Println("timeout middleware")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// AuthMiddleware for general user authentication
func AuthMiddleware(next http.Handler) http.Handler {
	log.Println("auth middleware")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, err := parseToken(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		// Adding claims to context for downstream handlers
		ctx := context.WithValue(r.Context(), "userClaims", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// AdminMiddleware ensures that only admin users can access certain routes
func AdminMiddleware(next http.Handler) http.Handler {
	log.Println("admin middleware")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, err := parseToken(r)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		// Check if user role is admin
		if claims.Role != "admin" {
			http.Error(w, "Access denied: admin only", http.StatusForbidden)
			return
		}

		// Adding claims to context for downstream handlers
		ctx := context.WithValue(r.Context(), "userClaims", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// parseToken validates JWT token and returns Claims
func parseToken(r *http.Request) (*Claims, error) {
	cookie, err := r.Cookie("token")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			return nil, errors.New("Token not found")
		}
		return nil, errors.New("Internal server error")
	}

	tokenString := cookie.Value
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("Invalid token")
	}

	return claims, nil
}

func EnableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Разрешаем запросы с фронтенда
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000") // Замените на нужный домен фронтенда
		w.Header().Set("Vary", "Origin")                                       // Указывает, что заголовок зависит от источника

		// Указываем разрешённые методы и заголовки
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true") // Если нужно передавать сессии и куки

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent) // Возвращаем пустой ответ для preflight-запроса
			return
		}

		next.ServeHTTP(w, r)
	})
}
