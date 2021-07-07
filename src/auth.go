package main

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"regexp"
	"time"
)

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	Uid      int    `json:"uid"`
	Username string `json:"username"`
	Rolename string `json:"rolename"`
	Isadmin  bool   `json:"isadmin"`
	jwt.StandardClaims
}

type ResponseToken struct {
	Token string `json:"token"`
}

func Login(w http.ResponseWriter, r *http.Request) {
	var creds Credentials

	if r.Method != "POST" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	uid := getUidByUsernameAndPassword(creds.Username, creds.Password)
	if uid < 0 {
		ApiErrMsg(w, "Wrong username or password")
		return
	}
	role := getRoleByUid(uid)
	if role == nil {
		ApiErrMsg(w, "Wrong username or password")
		return
	}
	expirationTime := time.Now().Add(7 * 24 * time.Hour)
	claims := &Claims{
		Uid:      uid,
		Username: creds.Username,
		Rolename: role.Rolename,
		Isadmin:  role.Isadmin,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Printf("%s Login success: %d %s\n", time.Now().Format(time.UnixDate), uid, creds.Username)
	json.NewEncoder(w).Encode(&ApiReturn{
		Retcode: 0,
		Message: "OK",
		Data: &ResponseToken{
			Token: tokenString,
		},
	})
}

func VerifyHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		re := regexp.MustCompile(`Bearer\s(.*)$`)

		headerAuth := r.Header.Get("Authorization")
		if len(headerAuth) == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		tknStr := re.FindStringSubmatch(headerAuth)
		if len(tknStr) == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		claims := &Claims{}

		// Parse the JWT string and store the result in `claims`.
		// Note that we are passing the key in this method as well. This method will return an error
		// if the token is invalid (if it has expired according to the expiry time we set on sign in),
		// or if the signature does not match
		tkn, err := jwt.ParseWithClaims(tknStr[1], claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil || !tkn.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Type", "application/json;charset=utf-8")
		next.ServeHTTP(w, r)
	})
}

func VerifyAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := getUserInfoFromJWT(r)
		if user.Isadmin == true {
			w.Header().Set("Content-Type", "application/json;charset=utf-8")
			next.ServeHTTP(w, r)
		} else {
			ApiErrMsg(w, "权限不足")
		}
	})
}

func getUserInfoFromJWT(r *http.Request) *User {
	re := regexp.MustCompile(`Bearer\s(.*)$`)
	headerAuth := r.Header.Get("Authorization")
	tknStr := re.FindStringSubmatch(headerAuth)
	claims := &Claims{}

	jwt.ParseWithClaims(tknStr[1], claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	return &User{
		Uid:      claims.Uid,
		Username: claims.Username,
		Rolename: claims.Rolename,
		Isadmin:  claims.Isadmin,
	}
}

func RefreshToken(w http.ResponseWriter, r *http.Request) {
	re := regexp.MustCompile(`Bearer\s(.*)$`)
	headerAuth := r.Header.Get("Authorization")
	if len(headerAuth) == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	tknStr := re.FindStringSubmatch(headerAuth)
	if len(tknStr) == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(tknStr[1], claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 24*time.Hour {
		ApiErrMsg(w, "Token not expires in one day")
		return
	}
	expirationTime := time.Now().Add(7 * 24 * time.Hour)
	claims.ExpiresAt = expirationTime.Unix()
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := newToken.SignedString(jwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(&ApiReturn{
		Retcode: 0,
		Message: "OK",
		Data: &ResponseToken{
			Token: tokenString,
		},
	})
}
