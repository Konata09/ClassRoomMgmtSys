package main

import (
	"database/sql"
	"fmt"
	"log"
	"log/syslog"
	"net/http"
)

var jwtKey = []byte("sd*ust#konata&2O20")
var db *sql.DB
var sysLog *syslog.Writer

func initDBConn() {
	var err error
	db, err = sql.Open("sqlite3", "db.db?cache=shared&mode=wrc")
	if err != nil {
		log.Fatal(err)
	}
	if db == nil {
		fmt.Println("Error open database.")
	}
}

func initSyslog() {
	var serverAddr string
	ret := getPreference("syslog_server", &serverAddr)
	if !ret {
		log.Fatal("syslog server not configured!")
	}
	fmt.Printf("syslog server: %s\n", serverAddr)
	var err error
	sysLog, err = syslog.Dial("udp", serverAddr, syslog.LOG_NOTICE|syslog.LOG_USER, "ClassroomMgmtSys")
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println("syslog established")
	sysLog.Info("syslog established")
}

func main() {
	initDBConn()
	initSyslog()
	mux := http.NewServeMux()
	mux.Handle("/api/v2/login", http.HandlerFunc(Login))
	mux.Handle("/api/v2/refresh", http.HandlerFunc(RefreshToken))
	mux.Handle("/api/v2/logout", VerifyHeader(http.HandlerFunc(Logout)))
	//mux.Handle("/api/v2/getCommand", VerifyHeader(http.HandlerFunc(GetCommand)))
	//mux.Handle("/api/v2/getDevice", VerifyHeader(http.HandlerFunc(GetDevice)))
	//mux.Handle("/api/v2/sendUDP", VerifyHeader(http.HandlerFunc(SendUDP)))
	//mux.Handle("/api/v2/sendWOL", VerifyHeader(http.HandlerFunc(SendWOL)))
	mux.Handle("/api/v2/user/changePassword", VerifyHeader(http.HandlerFunc(UserChangePassword)))
	mux.Handle("/api/v2/admin/changePassword", VerifyHeader(VerifyAdmin(http.HandlerFunc(AdminChangePassword))))
	mux.Handle("/api/v2/admin/setUser", VerifyHeader(VerifyAdmin(http.HandlerFunc(SetUser))))
	//mux.Handle("/api/v2/admin/setCommand", VerifyHeader(VerifyAdmin(http.HandlerFunc(SetCommand))))
	//mux.Handle("/api/v2/admin/setDevice", VerifyHeader(VerifyAdmin(http.HandlerFunc(SetDevice))))
	mux.Handle("/api/v2/user/ChangePhone", VerifyHeader(http.HandlerFunc(ChangePhone)))
	// 返回全部教室和分组 含基本状态
	mux.Handle("/api/v2/getRooms", VerifyHeader(http.HandlerFunc(GetClassrooms)))
	// 返回教室详情
	mux.Handle("/api/v2/getRoomDetail", VerifyHeader(http.HandlerFunc(GetClassroomDetail)))
	// 返回教室 ping 结果
	mux.Handle("/api/v2/getRoomStatus", VerifyHeader(http.HandlerFunc(GetClassroomStatus)))
	// post 修改教室的名称 组
	mux.Handle("/api/v2/admin/setRoom", VerifyHeader(VerifyAdmin(http.HandlerFunc(SetClassroom))))
	// post 修改设备ip地址
	mux.Handle("/api/v2/admin/setDevice", VerifyHeader(VerifyAdmin(http.HandlerFunc(SetDevice))))
	// 发送教室控制命令
	//mux.Handle("/api/v2/sendCmd", VerifyHeader(http.HandlerFunc()))
	// 返回我的工单
	//mux.Handle("/api/v2/getMyTicket", VerifyHeader(http.HandlerFunc()))
	// 返回全部工单 动态确定获取条数
	//mux.Handle("/api/v2/getTicket", VerifyHeader(http.HandlerFunc()))
	// 返回工单详情
	//mux.Handle("/api/v2/getTicketDetail", VerifyHeader(http.HandlerFunc()))
	//mux.Handle("/api/v2/admin/setTicketDutyUser", VerifyHeader(VerifyAdmin(http.HandlerFunc())))
	//mux.Handle("/api/v2/deleteTicket", VerifyHeader(VerifyAdmin(http.HandlerFunc())))
	//mux.Handle("/api/v2/setTicketDone", VerifyHeader(http.HandlerFunc()))
	// 修改值班表
	//mux.Handle("/api/v2/setDuty", VerifyHeader(VerifyAdmin(http.HandlerFunc())))
	// 获得当日值班
	//mux.Handle("/api/v2/getDutyUser", VerifyHeader(http.HandlerFunc()))

	//fmt.Println(getSubnetBroadcast("172.16.0.254",32))
	fmt.Println("UDPServer listen on 63112")
	log.Panic(http.ListenAndServe(":63112", mux))
}
