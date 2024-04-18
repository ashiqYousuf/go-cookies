package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/ashiqYousuf/go-cookies/internal/cookies"
)

func setCookieHandler(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{
		Name:     "exampleCookie",
		Value:    "Hi Zoë!!",
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}

	err := cookies.WriteEncrypted(w, cookie, secretKey)
	// err := cookies.WriteSigned(w, cookie, secretKey)
	if err != nil {
		log.Println(err)
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("<h1>Cookie Set</h1>"))
}

func getCookieHandler(w http.ResponseWriter, r *http.Request) {
	value, err := cookies.ReadEncryped(r, "exampleCookie", secretKey)
	// value, err := cookies.ReadSigned(r, "exampleCookie", secretKey)
	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			http.Error(w, "<p>cookie not found</p>", http.StatusBadRequest)
		case errors.Is(err, cookies.ErrInvalidValue):
			http.Error(w, "invalid cookie", http.StatusBadRequest)
		default:
			log.Println(err)
			http.Error(w, "server error", http.StatusInternalServerError)
		}
		return
	}

	w.Write([]byte(fmt.Sprintf("<h1>Cookie Value: %s</h1>", value)))
}
