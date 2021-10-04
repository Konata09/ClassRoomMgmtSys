package main

import (
	"github.com/google/uuid"
	"reflect"
	"strings"
)

func getUidByUsernameAndPassword(username string, password string) int {
	passwordMD5 := getPasswordMD5(password)
	stmt, err := db.Prepare("select uid from user where username = ? and password = ?")
	if err != nil {
		return -1
	}
	defer stmt.Close()
	var uid int
	err = stmt.QueryRow(username, passwordMD5).Scan(&uid)
	if err != nil {
		return -1
	}
	return uid
}

func getUidByUsername(username string) int {
	stmt, err := db.Prepare("select uid from user where username = ?")
	if err != nil {
		return -1
	}
	defer stmt.Close()
	var uid int
	err = stmt.QueryRow(username).Scan(&uid)
	if err != nil {
		return -1
	}
	return uid
}

func getPhoneByUid(uid int) int {
	stmt, err := db.Prepare("select phone from user where uid = ?")
	if err != nil {
		return -1
	}
	defer stmt.Close()
	var phone int
	err = stmt.QueryRow(uid).Scan(&phone)
	if err != nil {
		return -1
	}
	return phone
}
func setPhoneByUid(uid int, phone int) bool {
	stmt, err := db.Prepare("update user set phone = ? where uid = ?")
	if err != nil {
		return false
	}
	defer stmt.Close()
	res, err := stmt.Exec(phone, uid)
	if err != nil {
		return false
	}
	affected, _ := res.RowsAffected()
	return affected > 0
}

func getRoleByUid(uid int) *Role {
	stmt, err := db.Prepare("select rolename, isadmin, isstaff from user,role where uid = ? and user.roleid = role.roleid")
	if err != nil {
		return nil
	}
	defer stmt.Close()
	var role Role
	err = stmt.QueryRow(uid).Scan(&role.Rolename, &role.Isadmin, &role.Isstaff)
	if err != nil {
		return nil
	}
	return &role
}

func getRoleidByRolename(rolename string) int {
	stmt, err := db.Prepare("select roleid from role where rolename = ?")
	if err != nil {
		return -1
	}
	var roleid int
	defer stmt.Close()
	err = stmt.QueryRow(rolename).Scan(&roleid)
	if err != nil {
		return -1
	}
	return roleid
}

func getPasswordByUid(uid int) (res string, err error) {
	stmt, err := db.Prepare("select password from user where uid = ?")
	if err != nil {
		return "", err
	}
	defer stmt.Close()
	var passMD5 string
	err = stmt.QueryRow(uid).Scan(&passMD5)
	if err != nil {
		return "", err
	}
	return passMD5, nil
}

func setPasswordByUid(uid int, newPassword string) bool {
	stmt, err := db.Prepare("update user set password = ? where uid = ?")
	if err != nil {
		return false
	}
	defer stmt.Close()
	_, err = stmt.Exec(newPassword, uid)
	if err != nil {
		return false
	}
	return true
}

func getUserByUid(uid int) *User {
	stmt, err := db.Prepare("select username, rolename, isadmin, isstaff from user,role where user.roleid = role.roleid and uid = ?")
	if err != nil {
		return nil
	}
	defer stmt.Close()
	var user User
	err = stmt.QueryRow(uid).Scan(&user.Username, &user.Rolename, &user.Isadmin, &user.Isstaff)
	if err != nil {
		return nil
	}
	return &user
}

func getUsers() []User {
	stmt, err := db.Prepare("select uid, username, rolename, isadmin, isstaff, phone from user,role where user.roleid = role.roleid")
	if err != nil {
		return nil
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		return nil
	}
	var users []User
	for rows.Next() {
		var user User
		rows.Scan(&user.Uid, &user.Username, &user.Rolename, &user.Isadmin, &user.Isstaff, &user.Phone)
		users = append(users, user)
	}
	return users
}

func addUser(username string, password string, phone int, roleid int) bool {
	stmt, err := db.Prepare("insert into user (username, password, phone, roleid) values (?, ?, ?, ?)")
	if err != nil {
		return false
	}
	defer stmt.Close()
	_, err = stmt.Exec(username, password, phone, roleid)
	if err != nil {
		return false
	}
	return true
}

func deleteUser(uid int) bool {
	stmt, err := db.Prepare("delete from user where uid = ?")
	if err != nil {
		return false
	}
	defer stmt.Close()
	_, err = stmt.Exec(uid)
	if err != nil {
		return false
	}
	return true
}

func getUserCommands() []UserCommand {
	stmt, err := db.Prepare("select id, name from command")
	if err != nil {
		return nil
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		return nil
	}
	var commands []UserCommand
	for rows.Next() {
		var command UserCommand
		rows.Scan(&command.CommandId, &command.CommandName)
		commands = append(commands, command)
	}
	return commands
}

func getCommands() []Command {
	stmt, err := db.Prepare("select id, name, value, port from command")
	if err != nil {
		return nil
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		return nil
	}
	var commands []Command
	for rows.Next() {
		var command Command
		rows.Scan(&command.CommandId, &command.CommandName, &command.CommandValue, &command.CommandPort)
		commands = append(commands, command)
	}
	return commands
}

func getCommandById(commandId int) *Command {
	stmt, err := db.Prepare("select name, value, port from command where id = ?")
	if err != nil {
		return nil
	}
	defer stmt.Close()
	var command Command
	err = stmt.QueryRow(commandId).Scan(&command.CommandName, &command.CommandValue, &command.CommandPort)
	if err != nil {
		return nil
	}
	return &command
}

func addCommand(commands []Command) bool {
	stmt, err := db.Prepare("insert into command (name, value, port) values (?, ?, ?)")
	if err != nil {
		return false
	}
	defer stmt.Close()
	for _, cmd := range commands {
		_, err = stmt.Exec(cmd.CommandName, trimCommandToStor(cmd.CommandValue), cmd.CommandPort)
		if err != nil {
			return false
		}
	}
	return true
}

func deleteCommand(commandId int) bool {
	stmt, err := db.Prepare("delete from command where id = ?")
	if err != nil {
		return false
	}
	defer stmt.Close()
	_, err = stmt.Exec(commandId)
	if err != nil {
		return false
	}
	return true
}

func setCommand(commandId int, commandName string, commandValue string, commandPort int) bool {
	stmt, err := db.Prepare("update command set name = ?, value = ?, port = ? where id = ?")
	if err != nil {
		return false
	}
	defer stmt.Close()
	_, err = stmt.Exec(commandName, trimCommandToStor(commandValue), commandPort, commandId)
	if err != nil {
		return false
	}
	return true
}

func getDevices() []Device {
	stmt, err := db.Prepare("select device.id, name, ip, mac, typeid, classid from device, devicetype where device.typeid = devicetype.id")
	if err != nil {
		return nil
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		return nil
	}
	var devices []Device
	for rows.Next() {
		var device Device
		rows.Scan(&device.DeviceId, &device.DeviceName, &device.DeviceIp, &device.DeviceMac, &device.DeviceTypeId, &device.DeviceClassId)
		device.DeviceMac = trimMACtoShow(device.DeviceMac)
		devices = append(devices, device)
	}
	return devices
}

func getUserDevices() []UserDevice {
	stmt, err := db.Prepare("select device.id, name from device, devicetype where typeid = devicetype.id")
	if err != nil {
		return nil
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		return nil
	}
	var devices []UserDevice
	for rows.Next() {
		var device UserDevice
		rows.Scan(&device.DeviceId, &device.DeviceName)
		devices = append(devices, device)
	}
	return devices
}

func getDeviceById(deviceId int) *Device {
	stmt, err := db.Prepare("select name, ip, mac, typeid, classid from device, devicetype where device.typeid = devicetype.id and device.id = ?")
	if err != nil {
		return nil
	}
	defer stmt.Close()
	var device Device
	err = stmt.QueryRow(deviceId).Scan(&device.DeviceName, &device.DeviceIp, &device.DeviceMac, &device.DeviceTypeId, &device.DeviceClassId)
	if err != nil {
		return nil
	}
	device.DeviceMac = trimMACtoShow(device.DeviceMac)
	return &device
}

func getDevicesByClassId(classId int) []Device {
	stmt, err := db.Prepare("select device.id, name, ip, mac, typeid, classid from device, devicetype where device.typeid = devicetype.id and device.classid = ?")
	if err != nil {
		return nil
	}
	defer stmt.Close()
	rows, err := stmt.Query(classId)
	if err != nil {
		return nil
	}
	var devices []Device
	for rows.Next() {
		var device Device
		rows.Scan(&device.DeviceId, &device.DeviceName, &device.DeviceIp, &device.DeviceMac, &device.DeviceTypeId, &device.DeviceClassId)
		device.DeviceMac = trimMACtoShow(device.DeviceMac)
		devices = append(devices, device)
	}
	return devices
}

func addDevice(devices []Device) bool {
	stmt, err := db.Prepare("insert into device (ip, mac, typeid, classid) values (?, ?, ?, ?)")
	if err != nil {
		return false
	}
	defer stmt.Close()
	for _, dev := range devices {
		if reflect.ValueOf(dev).IsZero() {
			continue
		}
		_, err = stmt.Exec(dev.DeviceIp, trimMACtoStor(dev.DeviceMac), dev.DeviceTypeId, dev.DeviceClassId)
		if err != nil {
			return false
		}
	}
	return true
}

func deleteDevice(deviceId int) bool {
	stmt, err := db.Prepare("delete from device where id = ?")
	if err != nil {
		return false
	}
	defer stmt.Close()
	_, err = stmt.Exec(deviceId)
	if err != nil {
		return false
	}
	return true
}

func setDevice(deviceId int, deviceIp string, deviceMac string, deviceTypeId int, deviceClassId int) bool {
	stmt, err := db.Prepare("update device set ip = ?, mac = ?, typeid = ?, classid = ? where id = ?")
	if err != nil {
		return false
	}
	defer stmt.Close()
	res, err := stmt.Exec(deviceIp, trimMACtoStor(deviceMac), deviceTypeId, deviceClassId, deviceId)
	if err != nil {
		return false
	}
	affected, _ := res.RowsAffected()
	return affected > 0
}

func getClassrooms() []Classroom {
	stmt, err := db.Prepare("select classroom.id, classroom.name, classroomgroup.id, classroomgroup.name from classroom,classroomgroup where classroom.groupid = classroomgroup.id")
	if err != nil {
		return nil
	}
	rows, err := stmt.Query()
	defer stmt.Close()
	if err != nil {
		return nil
	}
	var classrooms []Classroom
	for rows.Next() {
		var classroom Classroom
		rows.Scan(&classroom.Id, &classroom.Name, &classroom.GroupId, &classroom.GroupName)
		classrooms = append(classrooms, classroom)
	}
	return classrooms
}

func getClassroom(id int) *ClassroomDetail {
	stmt, err := db.Prepare("select classroom.id, classroom.name, classroomgroup.id, classroomgroup.name from classroom,classroomgroup where classroom.groupid = classroomgroup.id and classroom.id = ?")
	if err != nil {
		return nil
	}
	defer stmt.Close()
	var classroomDetail ClassroomDetail
	err = stmt.QueryRow(id).Scan(&classroomDetail.Id, &classroomDetail.Name, &classroomDetail.GroupId, &classroomDetail.GroupName)
	if err != nil {
		return nil
	}
	classroomDetail.Devices = getDevicesByClassId(id)
	return &classroomDetail
}

func getClassroomControllers() []Device {
	stmt, err := db.Prepare("select id, ip, mac, classid, typeid from device where typeid = 1")
	if err != nil {
		return nil
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		return nil
	}
	var devices []Device
	for rows.Next() {
		var device Device
		rows.Scan(&device.DeviceId, &device.DeviceIp, &device.DeviceMac, &device.DeviceClassId, &device.DeviceTypeId)
		devices = append(devices, device)
	}
	return devices
}

func getClassroomLindges() []Device {
	stmt, err := db.Prepare("select id, ip, mac, classid, typeid from device where typeid = 2")
	if err != nil {
		return nil
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		return nil
	}
	var devices []Device
	for rows.Next() {
		var device Device
		rows.Scan(&device.DeviceId, &device.DeviceIp, &device.DeviceMac, &device.DeviceClassId, &device.DeviceTypeId)
		devices = append(devices, device)
	}
	return devices
}

func setClassroom(id int, name string, groupid int) bool {
	stmt, err := db.Prepare("update classroom set name = ?, groupid = ? where id = ?")
	if err != nil {
		return false
	}
	defer stmt.Close()
	_, err = stmt.Exec(name, groupid, id)
	if err != nil {
		return false
	}
	return true
}

func addTicket(title string, detail string, severity int, place string, createUser int, dutyUser1 int, dutyUser2 int, dutyUser3 int, createTime string, startTime string) bool {
	stmt, err := db.Prepare("insert into ticket(title, detail, severity, place, createuser, dutyuser1, dutyuser2, dutyuser3, createtime, starttime, completetime, completeuser,completedetail) values (?,?,?,?,?,?,?,?,?,?,'',0,'')")
	if err != nil {
		return false
	}
	defer stmt.Close()
	_, err = stmt.Exec(title, detail, severity, place, createUser, dutyUser1, dutyUser2, dutyUser3, createTime, startTime)
	if err != nil {
		return false
	}
	return true
}

func getTickets(limit int) []TicketOverview {
	stmt, err := db.Prepare("select id, title, severity, status,place, createuser, username from ticket, user where createuser = uid order by id desc limit ?")
	if err != nil {
		return nil
	}
	defer stmt.Close()
	rows, err := stmt.Query(limit)
	if err != nil {
		return nil
	}
	var ticketOverviews []TicketOverview
	for rows.Next() {
		var ticketOverview TicketOverview
		rows.Scan(&ticketOverview.Id, &ticketOverview.Title, &ticketOverview.Severity, &ticketOverview.Status, &ticketOverview.Place, &ticketOverview.CreateUser, &ticketOverview.CreateUserName)
		ticketOverviews = append(ticketOverviews, ticketOverview)
	}
	return ticketOverviews
}

func getTicket(id int) *Ticket {
	stmt, err := db.Prepare("select ticket.id, title, detail, severity, status, place, createuser, dutyuser1, dutyuser2, dutyuser3, completeuser, createtime, starttime, completetime, completedetail from ticket where ticket.id = ? ")
	if err != nil {
		return nil
	}
	var ticket Ticket
	defer stmt.Close()
	err = stmt.QueryRow(id).Scan(&ticket.Id, &ticket.Title, &ticket.Detail, &ticket.Severity, &ticket.Status, &ticket.Place, &ticket.CreateUser, &ticket.DutyUser1, &ticket.DutyUser2, &ticket.DutyUser3, &ticket.CompleteUser, &ticket.CreateTime, &ticket.StartTime, &ticket.CompleteTime, &ticket.CompleteDetail)
	if err != nil {
		return nil
	}
	ticket.CreateUserName = getUserByUid(ticket.CreateUser).Username
	if ticket.DutyUser1 != 0 {
		ticket.DutyUser1Name = getUserByUid(ticket.DutyUser1).Username
	}
	if ticket.DutyUser2 != 0 {
		ticket.DutyUser2Name = getUserByUid(ticket.DutyUser2).Username
	}
	if ticket.DutyUser3 != 0 {
		ticket.DutyUser3Name = getUserByUid(ticket.DutyUser3).Username
	}
	if ticket.CompleteUser != 0 {
		ticket.CompleteUserName = getUserByUid(ticket.CompleteUser).Username
	}
	return &ticket
}

func getUserDutyTicketOverview(id int) []TicketOverview {
	stmt, err := db.Prepare("select ticket.id, title, severity, status,place, createuser, username from ticket, user where (dutyuser1 = ? or dutyuser2 = ? or dutyuser3 = ?) and createuser = user.uid")
	if err != nil {
		return nil
	}
	defer stmt.Close()
	rows, err := stmt.Query(id, id, id)
	if err != nil {
		return nil
	}
	var ticketOverviews []TicketOverview
	for rows.Next() {
		var ticketOverview TicketOverview
		rows.Scan(&ticketOverview.Id, &ticketOverview.Title, &ticketOverview.Severity, &ticketOverview.Status, &ticketOverview.Place, &ticketOverview.CreateUser, &ticketOverview.CreateUserName)
		ticketOverviews = append(ticketOverviews, ticketOverview)
	}
	return ticketOverviews
}

func setTicketStatus(id int, status int) bool {
	stmt, err := db.Prepare("update ticket set status = ? where id = ?")
	if err != nil {
		return false
	}
	defer stmt.Close()
	_, err = stmt.Exec(status, id)
	if err != nil {
		return false
	}
	return true
}

func setTicketDutyUser(id int, dutyUser1 int, dutyUser2 int, dutyUser3 int) bool {
	stmt, err := db.Prepare("update ticket set dutyuser1 = ?, dutyuser2 = ?, dutyuser3 = ? where id = ?")
	if err != nil {
		return false
	}
	defer stmt.Close()
	_, err = stmt.Exec(dutyUser1, dutyUser2, dutyUser3, id)
	if err != nil {
		return false
	}
	return true
}

func deleteTicket(id int) bool {
	stmt, err := db.Prepare("delete from ticket where id = ?")
	if err != nil {
		return false
	}
	defer stmt.Close()
	_, err = stmt.Exec(id)
	if err != nil {
		return false
	}
	return true
}

func getDutyCalender() *DutyCalender {
	stmt, err := db.Prepare("select id, user.uid, day, username from duty_v3, user where duty_v3.uid = user.uid")
	if err != nil {
		return nil
	}
	defer stmt.Close()
	var dutyCalender DutyCalender
	rows, err := stmt.Query()
	if err != nil {
		return nil
	}
	for rows.Next() {
		var dutyUser DutyUser
		var day string
		rows.Scan(&dutyUser.Id, &dutyUser.Uid, &day, &dutyUser.Username)
		switch day {
		case "Monday":
			dutyCalender.Monday = append(dutyCalender.Monday, dutyUser)
		case "Tuesday":
			dutyCalender.Tuesday = append(dutyCalender.Tuesday, dutyUser)
		case "Wednesday":
			dutyCalender.Wednesday = append(dutyCalender.Wednesday, dutyUser)
		case "Thursday":
			dutyCalender.Thursday = append(dutyCalender.Thursday, dutyUser)
		case "Friday":
			dutyCalender.Friday = append(dutyCalender.Friday, dutyUser)
		case "Saturday":
			dutyCalender.Saturday = append(dutyCalender.Saturday, dutyUser)
		case "Sunday":
			dutyCalender.Sunday = append(dutyCalender.Sunday, dutyUser)
		}
	}
	return &dutyCalender
}

func addDutyCalender(day string, uid int) bool {
	stmt, err := db.Prepare("insert into duty_v3 (id, day, uid) values (?, ?, ?)")
	if err != nil {
		return false
	}
	defer stmt.Close()
	_, err = stmt.Exec(strings.Replace(uuid.NewString(), "-", "", -1), day, uid)
	if err != nil {
		return false
	}
	return true
}

func updateDutyCalender(day string, uid int, id string) bool {
	stmt, err := db.Prepare("update duty_v3 set day = ?, uid = ? where id = ?")
	if err != nil {
		return false
	}
	defer stmt.Close()
	_, err = stmt.Exec(day, uid, id)
	if err != nil {
		return false
	}
	return true
}

func delDutyCalender(id string) bool {
	stmt, err := db.Prepare("delete from duty_v3 where id = ?")
	if err != nil {
		return false
	}
	_, err = stmt.Exec(id)
	if err != nil {
		return false
	}
	return true
}

func getPreference(name string, value *string) bool {
	stmt, err := db.Prepare("select value from preferences where name = ?")
	if err != nil {
		return false
	}
	defer stmt.Close()
	err = stmt.QueryRow(name).Scan(value)
	if err != nil {
		return false
	}
	return true
}
