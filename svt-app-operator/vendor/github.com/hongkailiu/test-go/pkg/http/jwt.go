package http

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	log "github.com/sirupsen/logrus"
)

func getToken(localID string, key interface{}) (string, error) {
	claims := &jwt.StandardClaims{
		//https://godoc.org/github.com/dgrijalva/jwt-go#pkg-examples
		Id:        localID,
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(key)
}

func getLocalIDFromToken(tokenString string, key interface{}) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		alg := token.Header["alg"]
		log.WithFields(log.Fields{"token.Method": token.Method, "alg": alg}).Debug("found in token")
		if alg != "HS256" {
			return nil, fmt.Errorf("unexpected signing alg: %v", alg)
		}
		return key, nil
	})
	if err != nil {
		return "", fmt.Errorf("found error when jwt.Parse(), %s", err.Error())
	}

	claims, ok := token.Claims.(*jwt.StandardClaims)
	log.WithFields(log.Fields{"ok": ok, "token.Valid": token.Valid, "token.Claims": token.Claims, "claims": claims}).Debug("found claims")

	if ok && token.Valid {
		localID := claims.Id
		exp := claims.ExpiresAt
		log.WithFields(log.Fields{"localID": localID, "exp": exp}).Debug("found in token")
		if err != nil {
			return "", fmt.Errorf("found error when getLocalID(s), %s", err.Error())
		}
		if time.Now().Unix() > exp {
			return "", fmt.Errorf("token is expired, %d", exp)
		}
		return localID, nil
	}
	return "", fmt.Errorf("unexpected error when getLocalIDFromToken: %s", tokenString)
}
