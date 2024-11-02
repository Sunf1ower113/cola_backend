package utils

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

func RenderJSON(w http.ResponseWriter, code int, v interface{}) {
	js, err := json.Marshal(v)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(code)
	n, err := w.Write(js)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(n))
}

func SetCookie(w http.ResponseWriter, v string) {
	cookie := http.Cookie{
		Name:     "token",
		Value:    v,
		Expires:  time.Now().Add(time.Hour * 24),
		SameSite: http.SameSiteNoneMode,
		Secure:   true, // Установите true для HTTPS; false для локальной разработки
		HttpOnly: true,
		Path:     "/",
	}
	http.SetCookie(w, &cookie)
}
