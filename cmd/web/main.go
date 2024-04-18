package main

import (
	"encoding/hex"
	"flag"
	"log"
	"net/http"
)

var secretKey []byte

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	secret := flag.String("secret", "13d6b4dff8f84a10851021ec8608f814570d562c92fe6b5ec4c9f595bcb3234b", "Hex string")

	flag.Parse()

	var err error
	secretKey, err = hex.DecodeString(*secret)
	if err != nil {
		log.Fatal(err)
	}

	mux := routes()

	srv := &http.Server{
		Addr:    *addr,
		Handler: mux,
	}

	log.Println("Server listening on PORT", *addr)

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
