package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var SECRET = []byte("Super-secret-auth-key")

type Claims struct {
	UserCredentials *User
	jwt.StandardClaims
}

func createJwt(userParams User) (string, time.Time, error) {
	expirationTime := time.Now().Add(time.Minute * 5)

	claims := &Claims{
		UserCredentials: &userParams,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(SECRET)

	if err != nil {
		fmt.Println(err)
		return "", expirationTime, err
	}

	return tokenString, expirationTime, nil
}

func validateJwt(next func(writer http.ResponseWriter, request *http.Request)) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		cookie, err := request.Cookie("token")
		if err != nil {
			RemoveCookie(writer)
			if err == http.ErrNoCookie {
				writer.WriteHeader(http.StatusUnauthorized)
				return
			}
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		if cookie != nil {
			tokenStr := cookie.Value

			claims := &Claims{}

			token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
				return SECRET, nil
			})

			if err != nil {
				RemoveCookie(writer)
				if err == jwt.ErrSignatureInvalid {
					writer.WriteHeader(http.StatusUnauthorized)
					return
				}
				writer.WriteHeader(http.StatusBadRequest)
				return
			}

			if !token.Valid {
				RemoveCookie(writer)
				writer.WriteHeader(http.StatusUnauthorized)
				return
			}

			if token.Valid {
				next(writer, request)
			}

		} else {
			RemoveCookie(writer)
			writer.WriteHeader(http.StatusUnauthorized)
			writer.Write([]byte("Not authorized"))
		}
	})
}

func RefreshHandler(writer http.ResponseWriter, request *http.Request) {
	claims := &Claims{}
	expirationTime := time.Now().Add(time.Minute * 5)

	claims.ExpiresAt = expirationTime.Unix()

	Newtoken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := Newtoken.SignedString(SECRET)

	if err != nil {
		RemoveCookie(writer)
		fmt.Println(tokenString)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(
		writer,
		&http.Cookie{
			Name:    "token",
			Value:   tokenString,
			Expires: expirationTime,
		},
	)
}

func RemoveCookie(writer http.ResponseWriter) {
	http.SetCookie(
		writer,
		&http.Cookie{
			Name:    "token",
			Value:   "",
			Expires: time.Now().Add(-time.Hour),
		},
	)
}

func validateJwtAndReturnClaims(next func(writer http.ResponseWriter, request *http.Request, claims *Claims)) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		cookie, err := request.Cookie("token")
		if err != nil {
			RemoveCookie(writer)
			if err == http.ErrNoCookie {
				writer.WriteHeader(http.StatusUnauthorized)
				return
			}
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		if cookie != nil {
			tokenStr := cookie.Value

			claims := &Claims{}

			token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
				return SECRET, nil
			})

			if err != nil {
				RemoveCookie(writer)
				if err == jwt.ErrSignatureInvalid {
					writer.WriteHeader(http.StatusUnauthorized)
					return
				}
				writer.WriteHeader(http.StatusBadRequest)
				return
			}

			if !token.Valid {
				RemoveCookie(writer)
				writer.WriteHeader(http.StatusUnauthorized)
				return
			}

			if token.Valid {
				next(writer, request, claims)
			}

		} else {
			RemoveCookie(writer)
			writer.WriteHeader(http.StatusUnauthorized)
			writer.Write([]byte("Not authorized"))
		}
	})
}
