package activitypub

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"strings"
)

type Pair struct {
	K string
	V string
}

func Sign(privKey string, keyId string, target string, headers []Pair) (string, error) {
	var signedString string
	for _, h := range headers {
		if strings.ToLower(h.K) == "digest" {
			signedString += fmt.Sprintf("%s: SHA-256=%s\n", strings.ToLower(h.K), base64.StdEncoding.EncodeToString(Hash(h.V)))
		} else {
			signedString += fmt.Sprintf("%s: %s\n", strings.ToLower(h.K), h.V)
		}
	}
	signedString += fmt.Sprintf("(request-target): %s", target)

	block, _ := pem.Decode([]byte(privKey))
	if block == nil {
		return "", errors.New("failed to parse private key PEM")
	}

	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}

	signature, err := rsa.SignPKCS1v15(rand.Reader, priv, crypto.SHA256, Hash(signedString))
	if err != nil {
		return "", err
	}

	header := fmt.Sprintf("keyId=\"%s\",algorithm=\"rsa-sha256\",headers=\"", keyId)
	for _, h := range headers {
		header += fmt.Sprintf("%s ", strings.ToLower(h.K))
	}
	header += fmt.Sprintf("(request-target)\",signature=\"%s\"", base64.StdEncoding.EncodeToString(signature))

	return header, nil
}

func Verify(pubKey string, digest string, signature string) (bool, error) {

	// err := rsa.VerifyPKCS1v15(&priv.PublicKey, crypto.SHA256, digest, signature)
	// if err != nil {
	// 	return false, err
	// }

	return true, nil
}

func Hash(signedString string) []byte {
	sh := sha256.New()
	sh.Write([]byte(signedString))
	return sh.Sum(nil)
}
