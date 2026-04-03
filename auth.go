package main

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	log "github.com/sirupsen/logrus"
)

//	based on
//	https://yandex.cloud/ru/docs/iam/operations/iam-token/create-for-sa#go_1

// Формирование JWT.
func signedToken() string {
	claims := jwt.RegisteredClaims{
		Issuer:    auth.ServiceAccountID,
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(1 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		NotBefore: jwt.NewNumericDate(time.Now().UTC()),
		Audience:  []string{"https://iam.api.cloud.yandex.net/iam/v1/tokens"},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodPS256, claims)
	token.Header["kid"] = auth.KID

	privateKey := loadPrivateKey()
	signed, err := token.SignedString(privateKey)
	if err != nil {
		panic(err)
	}
	return signed
}

type keyFileStruct struct {
	PrivateKey string `json:"private_key"`
}

func getPrivateData() []byte {
	if conf.keyFile != "" {
		data, err := os.ReadFile(conf.keyFile)
		if err != nil {
			panic(err)
		}
		return data
	} else {
		log.Println(conf.keyData)
		data, err := base64.StdEncoding.DecodeString(conf.keyData)
		if err != nil {
			panic(err)
		}
		return data
	}
	return nil
}

func loadPrivateKey() *rsa.PrivateKey {
	rsaPrivateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(auth.PrivateKey))
	if err != nil {
		panic(err)
	}
	return rsaPrivateKey
}

func getIAMToken() string {
	jot := signedToken()
	log.Debugf("JWT token: %s...", jot[:50])
	resp, err := http.Post(
		"https://iam.api.cloud.yandex.net/iam/v1/tokens",
		"application/json",
		strings.NewReader(fmt.Sprintf(`{"jwt":"%s"}`, jot)),
	)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		panic(fmt.Sprintf("%s: %s", resp.Status, body))
	}
	var data struct {
		IAMToken string `json:"iamToken"`
	}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		panic(err)
	}
	return data.IAMToken
}
