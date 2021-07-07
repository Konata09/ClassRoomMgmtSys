package main

import "reflect"

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

func getRoleByUid(uid int) *Role {
	stmt, err := db.Prepare("select rolename, isadmin from user,role where uid = ? and user.roleid = role.roleid")
	if err != nil {
		return nil
	}
	defer stmt.Close()
	var role Role
	err = stmt.QueryRow(uid).Scan(&role.Rolename, &role.Isadmin)
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
	stmt, err := db.Prepare("select username, rolename, isadmin from user,role where user.roleid = role.roleid and uid = ?")
	if err != nil {
		return nil
	}
	defer stmt.Close()
	var user User
	err = stmt.QueryRow(uid).Scan(&user.Username, &user.Rolename, &user.Isadmin)
	if err != nil {
		return nil
	}
	return &user
}

func getUsers() []User {
	stmt, err := db.Prepare("select uid, username, rolename, isadmin from user,role where user.roleid = role.roleid")
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
		rows.Scan(&user.Uid, &user.Username, &user.Rolename, &user.Isadmin)
		users = append(users, user)
	}
	return users
}

func addUser(username string, password string, roleid int) bool {
	stmt, err := db.Prepare("insert into user (username, password, roleid) values (?, ?, ?)")
	if err != nil {
		return false
	}
	defer stmt.Close()
	_, err = stmt.Exec(username, password, roleid)
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
	_, err = stmt.Exec(commandName, trimCommandToStor(commandValue), commandPort, commandId)
	if err != nil {
		return false
	}
	return true
}

func getDevices() []Device {
	stmt, err := db.Prepare("select id, name, ip, mac, udp, wol, submask from device")
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
		rows.Scan(&device.DeviceId, &device.DeviceName, &device.DeviceIp, &device.DeviceMac, &device.DeviceUdp, &device.DeviceWol, &device.DeviceSubmask)
		device.DeviceMac = trimMACtoShow(device.DeviceMac)
		devices = append(devices, device)
	}
	return devices
}

func getUserDevices() []UserDevice {
	stmt, err := db.Prepare("select id, name, udp, wol from device")
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
		rows.Scan(&device.DeviceId, &device.DeviceName, &device.DeviceUdp, &device.DeviceWol)
		devices = append(devices, device)
	}
	return devices
}

func getDeviceById(deviceId int) *Device {
	stmt, err := db.Prepare("select name, ip, mac, udp, wol, submask from device where id = ?")
	if err != nil {
		return nil
	}
	defer stmt.Close()
	var device Device
	err = stmt.QueryRow(deviceId).Scan(&device.DeviceName, &device.DeviceIp, &device.DeviceMac, &device.DeviceUdp, &device.DeviceWol, &device.DeviceSubmask)
	if err != nil {
		return nil
	}
	device.DeviceMac = trimMACtoShow(device.DeviceMac)
	return &device
}

func addDevice(devices []Device) bool {
	stmt, err := db.Prepare("insert into device (name, ip, mac, udp, wol, submask) values (?, ?, ?, ?, ?)")
	if err != nil {
		return false
	}
	defer stmt.Close()
	for _, dev := range devices {
		if reflect.ValueOf(dev).IsZero() {
			continue
		}
		_, err = stmt.Exec(dev.DeviceName, dev.DeviceIp, trimMACtoStor(dev.DeviceMac), &dev.DeviceUdp, &dev.DeviceWol, &dev.DeviceSubmask)
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

func setDevice(deviceId int, deviceName string, deviceIp string, deviceMac string, deviceUdp bool, deviceWol bool, deviceSubmask int) bool {
	if deviceSubmask == 0 {
		deviceSubmask = 24
	}
	stmt, err := db.Prepare("update device set name = ?, ip = ?, mac = ?, udp = ?, wol = ?, submask = ? where id = ?")
	if err != nil {
		return false
	}
	_, err = stmt.Exec(deviceName, deviceIp, trimMACtoStor(deviceMac), deviceUdp, deviceWol, deviceSubmask, deviceId)
	if err != nil {
		return false
	}
	return true
}
