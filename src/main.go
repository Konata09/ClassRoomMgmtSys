package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
)

var jwtKey = []byte("sd*ust#konata&2O20")
var db *sql.DB

func initDBConn() {
	var err error
	db, err = sql.Open("sqlite3", "db.db?cache=shared&mode=wrc")
	if err != nil {
		log.Fatal(err)
	}
	if db == nil {
		fmt.Printf("Error open database.\n")
	}
}

func main() {
	initDBConn()
	mux := http.NewServeMux()
	mux.Handle("/api/v1/login", http.HandlerFunc(Login))
	mux.Handle("/api/v1/refresh", http.HandlerFunc(RefreshToken))
	mux.Handle("/api/v1/getCommand", VerifyHeader(http.HandlerFunc(GetCommand)))
	mux.Handle("/api/v1/getDevice", VerifyHeader(http.HandlerFunc(GetDevice)))
	mux.Handle("/api/v1/sendUDP", VerifyHeader(http.HandlerFunc(SendUDP)))
	mux.Handle("/api/v1/sendWOL", VerifyHeader(http.HandlerFunc(SendWOL)))
	mux.Handle("/api/v1/user/changePassword", VerifyHeader(http.HandlerFunc(UserChangePassword)))
	mux.Handle("/api/v1/admin/changePassword", VerifyHeader(VerifyAdmin(http.HandlerFunc(AdminChangePassword))))
	mux.Handle("/api/v1/admin/setUser", VerifyHeader(VerifyAdmin(http.HandlerFunc(SetUser))))
	mux.Handle("/api/v1/admin/setCommand", VerifyHeader(VerifyAdmin(http.HandlerFunc(SetCommand))))
	mux.Handle("/api/v1/admin/setDevice", VerifyHeader(VerifyAdmin(http.HandlerFunc(SetDevice))))

	//fmt.Println(getSubnetBroadcast("172.16.0.254",32))
	fmt.Println("UDPServer listen on 63112")
	log.Panic(http.ListenAndServe(":63112", mux))
}
