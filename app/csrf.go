package app

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"math"
	"math/big"
	"net/http"
	"strings"
)

func (c *Config) CreateCsrfToken(sessionID string) (string, error) {
	// https://cheatsheetseries.owasp.org/cheatsheets/Cross-Site_Request_Forgery_Prevention_Cheat_Sheet.html#hmac-csrf-token

	val, err := rand.Int(rand.Reader, big.NewInt(int64(math.MaxInt64)))
	if err != nil {
		return "", err
	}

	message := []byte(fmt.Sprintf("%s!%v", sessionID, val.Int64()))
	hmac := hmac.New(sha256.New, []byte(*c.CsrfSecret))
	hmac.Write([]byte(message))
	hmacData := hmac.Sum(nil)
	hmacStr := hex.EncodeToString(hmacData)

	csrfToken := fmt.Sprintf("%s.%s", hmacStr, message)

	log.Printf("csrfToken %s\n", csrfToken)

	return csrfToken, nil
}

func (c *Config) ValidCsrfToken(r *http.Request) bool {
	cookie, err := r.Cookie("session")
	if err != nil {
		log.Print("session cookie missing")
		return false
	}
	val := cookie.Value
	splitVal := strings.Split(val, ".")
	if len(splitVal) != 3 {
		log.Print("csrf token wrong length")
		return false
	}

	// check if form token contains session token
	requestToken := strings.Join(splitVal[1:], ".")
	log.Printf("requestToken: %s\n", requestToken)

	if requestToken != r.Header.Get("X-CSRF-Token") {
		log.Print("csrf token header missing")
		return false
	}

	// check if csrf token is valid with server secret
	requestHmac := splitVal[1]

	requestMessage := splitVal[2]
	log.Print("requestMessage", requestMessage)
	hmac := hmac.New(sha256.New, []byte(*c.CsrfSecret))
	hmac.Write([]byte(requestMessage))
	hmacData := hmac.Sum(nil)
	hmacStr := hex.EncodeToString(hmacData)

	return hmacStr == requestHmac
}

func (c *Config) GetCsrfToken(r *http.Request) (string, error) {
	cookie, err := r.Cookie("session")
	if err != nil {
		log.Print("session cookie missing")
		return "", nil
	}
	val := cookie.Value
	splitVal := strings.Split(val, ".")
	if len(splitVal) != 3 {
		log.Print("csrf token wrong length")
		return "", nil
	}

	return strings.Join(splitVal[1:], "."), nil
}
