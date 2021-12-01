package main

import (
	"bytes"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
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
	Isstaff  bool   `json:"isstaff"`
	Phone    int    `json:"phone"`
	jwt.StandardClaims
}

type ResponseToken struct {
	Token string `json:"token"`
}

func Login(w http.ResponseWriter, r *http.Request) {
	var creds Credentials

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
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
		Isstaff:  role.Isstaff,
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
	logBoth("[LOGIN]Login success: %d %s from %s", uid, creds.Username, r.RemoteAddr)
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
		w.Header().Set("Content-Type", "application/json;charset=utf-8")
		if strings.Contains(r.RemoteAddr, "127.0.0.1") {
			next.ServeHTTP(w, r)
			return
		}
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
		buf, _ := ioutil.ReadAll(r.Body)
		logBoth("[%s]%s %d %s %s %s", r.Method, r.RemoteAddr, claims.Uid, claims.Username, r.URL.Path, string(buf))
		reader := ioutil.NopCloser(bytes.NewBuffer(buf))
		r.Body = reader
		next.ServeHTTP(w, r)
	})
}

func VerifyAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := GetUserInfoFromJWT(r)
		if user.Isadmin == true {
			next.ServeHTTP(w, r)
		} else {
			ApiErrMsg(w, "权限不足")
		}
	})
}

func GetUserInfoFromJWT(r *http.Request) *User {
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
		Isstaff:  claims.Isstaff,
		Phone:    claims.Phone,
	}
}

func RefreshToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
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

	var body struct {
		Username string `json:"username"`
		Uid      int    `json:"uid"`
	}
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		ApiErr(w)
		return
	}

	// token 刷新限制
	//if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 96*time.Hour {
	//	ApiErrMsg(w, "Token not expires in 4 day")
	//	return
	//}

	expirationTime := time.Now().Add(7 * 24 * time.Hour)
	claims.ExpiresAt = expirationTime.Unix()
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := newToken.SignedString(jwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	logBoth("[LOGIN]Token refresh success: uid: %d username: %s from %s", body.Uid, body.Username, r.RemoteAddr)

	json.NewEncoder(w).Encode(&ApiReturn{
		Retcode: 0,
		Message: "OK",
		Data: &ResponseToken{
			Token: tokenString,
		},
	})
}

func Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
	// TODO: JWT 在服务端不好实现无效化
	ApiOk(w)
}
