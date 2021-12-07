package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-redis/redis/v8"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"log/syslog"
	"net/http"
	"time"
)

var jwtKey = []byte("sd*ust#konata&2O20")
var db *sql.DB
var sysLog *syslog.Writer
var rdb *redis.Client

func initDBConn() {
	var err error
	db, err = sql.Open("sqlite3", "db.db?cache=shared&mode=wrc")
	if err != nil {
		log.Fatal(err)
	}
	if db == nil {
		log.Fatal("Error open database.")
	}
}

func initSyslog() {
	var serverAddr string
	ret := getPreference("syslog_server", &serverAddr)
	if !ret {
		log.Fatal("Syslog server is not configured!")
		return
	}
	fmt.Printf("Syslog server is %s\n", serverAddr)
	var err error
	sysLog, err = syslog.Dial("udp", serverAddr, syslog.LOG_NOTICE|syslog.LOG_USER, "ClassroomMgmtSys")
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println("ClassRoomMgmtSys server was restarted.")
	sysLog.Warning("ClassRoomMgmtSys server was restarted.")
}

func initRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr:        "127.0.0.1:6379",
		Password:    "", // no password set
		DB:          0,  // use default DB
		DialTimeout: 2 * time.Second,
		ReadTimeout: 2 * time.Second,
	})

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		logBoth("Connect to Redis Failed %s", err)
		panic(err)
	}
	// 向 Redis 中存入教室基本信息
	classrooms := getClassrooms()
	var redisStatuses []ClassroomRedisStatus
	for _, classroom := range classrooms {
		var redisStatus ClassroomRedisStatus
		redisStatus.ClassroomId = classroom.Id
		redisStatus.ClassroomName = classroom.Name
		redisStatuses = append(redisStatuses, redisStatus)
	}
	SetMultiClassroomStatusToRedis(redisStatuses)
}

func test() {
}

func main() {
	initDBConn()
	go initRedis()
	initSyslog()
	updateCron()
	//test()
	//dyLogin()
	mux := http.NewServeMux()
	mux.Handle("/api/v2/login", http.HandlerFunc(Login))
	mux.Handle("/api/v2/refresh", http.HandlerFunc(RefreshToken))
	mux.Handle("/api/v2/logout", VerifyHeader(http.HandlerFunc(Logout)))
	//mux.Handle("/api/v2/getCommand", VerifyHeader(http.HandlerFunc(GetCommand)))
	//mux.Handle("/api/v2/getDevice", VerifyHeader(http.HandlerFunc(GetDevice)))
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
	mux.Handle("/api/v2/sendCmd", VerifyHeader(http.HandlerFunc(HandleCmd)))
	// 增加 Ticket
	mux.Handle("/api/v2/addTicket", VerifyHeader(http.HandlerFunc(AddTicket)))
	// 返回我的工单
	mux.Handle("/api/v2/getMyTicket", VerifyHeader(http.HandlerFunc(GetUserDutyTicket)))
	// 返回全部工单 动态确定获取条数
	mux.Handle("/api/v2/getTickets", VerifyHeader(http.HandlerFunc(GetAllTicket)))
	// 返回工单详情
	mux.Handle("/api/v2/getTicketDetail", VerifyHeader(http.HandlerFunc(GetTicketDetail)))
	mux.Handle("/api/v2/admin/setTicketDutyUser", VerifyHeader(VerifyAdmin(http.HandlerFunc(SetTicketDutyUser))))
	mux.Handle("/api/v2/admin/deleteTicket", VerifyHeader(VerifyAdmin(http.HandlerFunc(DeleteTicket))))
	mux.Handle("/api/v2/setTicketStatus", VerifyHeader(http.HandlerFunc(SetTicketDone)))
	// 查看值班表
	mux.Handle("/api/v2/getDuty", VerifyHeader(http.HandlerFunc(GetDutyCalender)))
	// 修改值班表
	mux.Handle("/api/v2/setDuty", VerifyHeader(VerifyAdmin(http.HandlerFunc(SetDutyCalender))))
	logBoth("[INFO] Server listen on 63112")
	log.Panic(http.ListenAndServe(":63112", mux))
}
