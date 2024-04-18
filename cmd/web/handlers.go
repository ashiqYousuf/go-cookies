package main

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/ashiqYousuf/go-cookies/internal/cookies"
)

func setCookieHandler(w http.ResponseWriter, r *http.Request) {
	user := User{Name: "Alice", Age: 21}

	var buf bytes.Buffer

	err := gob.NewEncoder(&buf).Encode(&user)
	if err != nil {
		log.Println(err)
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	cookie := http.Cookie{
		Name:     "exampleCookie",
		Value:    buf.String(),
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}

	err = cookies.WriteEncrypted(w, cookie, secretKey)
	// err := cookies.WriteSigned(w, cookie, secretKey)
	if err != nil {
		log.Println(err)
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("<h1>Cookie Set</h1>"))
}

func getCookieHandler(w http.ResponseWriter, r *http.Request) {
	gobEncodedValue, err := cookies.ReadEncryped(r, "exampleCookie", secretKey)
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

	var user User

	reader := strings.NewReader(gobEncodedValue)
	if err := gob.NewDecoder(reader).Decode(&user); err != nil {
		log.Println(err)
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "<h1>Name: %s </h1>\n", user.Name)
	fmt.Fprintf(w, "<h1>Age: %v </h1>\n", user.Age)
}
