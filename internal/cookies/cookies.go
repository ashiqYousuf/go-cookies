package cookies

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

var (
	ErrValueTooLong = errors.New("cookie value too long")
	ErrInvalidValue = errors.New("invalid cookie value")
)

func WriteEncrypted(w http.ResponseWriter, cookie http.Cookie, secretKey []byte) error {
	// Create a new AES cipher block from the secret key.
	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return err
	}

	// Wrap the cipher block in Galois Counter Mode.
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return err
	}

	plaintext := fmt.Sprintf("%s:%s", cookie.Name, cookie.Value)
	// Encrypt the data using aesGCM.Seal().
	// "{nonce}{encrypted plaintext data}" form of encryped data.
	encrypedValue := aesGCM.Seal(nonce, nonce, []byte(plaintext), nil)

	cookie.Value = string(encrypedValue)

	return Write(w, cookie)
}

func ReadEncryped(r *http.Request, name string, secretKey []byte) (string, error) {
	encrypedValue, err := Read(r, name)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := aesGCM.NonceSize()

	if len(encrypedValue) < nonceSize {
		return "", ErrInvalidValue
	}

	nonce := encrypedValue[:nonceSize]
	cipherText := encrypedValue[nonceSize:]

	// Use aesGCM.Open() to decrypt and authenticate the data.
	plaintext, err := aesGCM.Open(nil, []byte(nonce), []byte(cipherText), nil)
	if err != nil {
		return "", ErrInvalidValue
	}

	expectedName, value, ok := strings.Cut(string(plaintext), ":")
	if !ok {
		return "", ErrInvalidValue
	}

	if expectedName != name {
		return "", ErrInvalidValue
	}

	return value, nil
}

func WriteSigned(w http.ResponseWriter, cookie http.Cookie, secretKey []byte) error {
	// Calculate a HMAC signature of the cookie name and value, using SHA256 and
	// a secret key
	mac := hmac.New(sha256.New, secretKey)
	mac.Write([]byte(cookie.Name))
	mac.Write([]byte(cookie.Value))
	signature := mac.Sum(nil)

	cookie.Value = string(signature) + cookie.Value

	return Write(w, cookie)
}

func ReadSigned(r *http.Request, name string, secretKey []byte) (string, error) {
	signedValue, err := Read(r, name)
	if err != nil {
		return "", err
	}

	// A SHA256 HMAC signature has a fixed length of 32 bytes
	if len(signedValue) < sha256.Size {
		return "", ErrInvalidValue
	}

	signature := signedValue[:sha256.Size]
	value := signedValue[sha256.Size:]

	mac := hmac.New(sha256.New, secretKey)
	mac.Write([]byte(name))
	mac.Write([]byte(value))
	expectedSignature := mac.Sum(nil)

	if !hmac.Equal([]byte(signature), expectedSignature) {
		return "", ErrInvalidValue
	}

	return value, nil
}

func Write(w http.ResponseWriter, cookie http.Cookie) error {
	// base64-encode your cookie values before writing them
	cookie.Value = base64.URLEncoding.EncodeToString([]byte(cookie.Value))

	if len(cookie.String()) > 4096 {
		return ErrValueTooLong
	}

	http.SetCookie(w, &cookie)
	return nil
}

func Read(r *http.Request, name string) (string, error) {
	cookie, err := r.Cookie(name)

	if err != nil {
		return "", err
	}

	// Decode the base64-encoded cookie value
	value, err := base64.URLEncoding.DecodeString(cookie.Value)

	if err != nil {
		return "", ErrInvalidValue
	}

	return string(value), nil
}
