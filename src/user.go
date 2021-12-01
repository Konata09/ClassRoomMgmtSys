package main

import (
	"encoding/json"
	"net/http"
)

type User struct {
	Uid      int    `json:"uid"`
	Username string `json:"username"`
	Rolename string `json:"rolename"`
	Isadmin  bool   `json:"isadmin"`
	Isstaff  bool   `json:"isstaff"`
	Phone    int    `json:"phone"`
}

type Role struct {
	Rolename string `json:"rolename"`
	Isadmin  bool   `json:"isadmin"`
	Isstaff  bool   `json:"isstaff"`
}

type AllUsers struct {
	Count int    `json:"count"`
	Users []User `json:"users"`
}

type ChangePasswordBody struct {
	Uid     int    `json:"uid"`
	OldPass string `json:"old_pass"`
	NewPass string `json:"new_pass"`
}

type PutUserBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Rolename string `json:"rolename"`
	Phone    int    `json:"phone"`
}

func UserChangePassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	user := GetUserInfoFromJWT(r)
	var body ChangePasswordBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		ApiErr(w)
		return
	}
	oldPass, err := getPasswordByUid(user.Uid)
	if err != nil {
		ApiErr(w)
		return
	}
	if oldPass != getPasswordMD5(body.OldPass) {
		json.NewEncoder(w).Encode(&ApiReturn{
			Retcode: -1,
			Message: "Wrong password",
		})
		return
	}
	if setPasswordByUid(user.Uid, getPasswordMD5(body.NewPass)) {
		ApiOk(w)
	} else {
		ApiErrMsg(w, "修改失败")
	}
}

func AdminChangePassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var body ChangePasswordBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		ApiErr(w)
		return
	}
	if getUserByUid(body.Uid) == nil {
		ApiErrMsg(w, "用户不存在")
		return
	}
	if setPasswordByUid(body.Uid, getPasswordMD5(body.NewPass)) {
		ApiOk(w)
	} else {
		ApiErrMsg(w, "修改失败")
	}
}

func SetUser(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		users := getUsers()
		json.NewEncoder(w).Encode(&ApiReturn{
			Retcode: 0,
			Message: "OK",
			Data: &AllUsers{
				Count: len(users),
				Users: users,
			},
		})
	case "PUT":
		var body PutUserBody
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			ApiErr(w)
			return
		}
		roleid := getRoleidByRolename(body.Rolename)
		if roleid < 0 {
			ApiErrMsg(w, "用户组不存在")
			return
		}
		if getUidByUsername(body.Username) > 0 {
			ApiErrMsg(w, "用户名已占用")
			return
		}
		ok := addUser(body.Username, getPasswordMD5(body.Password), body.Phone, roleid)
		if ok {
			ApiOk(w)
		} else {
			ApiErr(w)
		}
	case "DELETE":
		var body User
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			ApiErr(w)
			return
		}
		if getUserByUid(body.Uid) == nil {
			ApiErrMsg(w, "用户不存在")
			return
		}
		ok := deleteUser(body.Uid)
		if ok {
			ApiOk(w)
		} else {
			ApiErr(w)
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func ChangePhone(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var body struct {
		Uid   int
		Phone int
	}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		ApiErr(w)
		return
	}
	if setPhoneByUid(body.Uid, body.Phone) {
		ApiOk(w)
	} else {
		ApiErrMsg(w, "修改失败")
	}
	return
}
