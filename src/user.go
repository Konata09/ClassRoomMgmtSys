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
}

type Role struct {
	Rolename string
	Isadmin  bool
}

type AllUsers struct {
	Count int    `json:"count"`
	Users []User `json:"users"`
}

type ChangePasswordBody struct {
	Uid     int    `json:"uid"`
	OldPass string `json:"oldpass"`
	NewPass string `json:"newpass"`
}

type PutUserBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Rolename string `json:"rolename"`
}

func UserChangePassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	user := getUserInfoFromJWT(r)
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
	ok := setPasswordByUid(user.Uid, getPasswordMD5(body.NewPass))
	if ok {
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
	ok := setPasswordByUid(body.Uid, getPasswordMD5(body.NewPass))
	if ok {
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
		ok := addUser(body.Username, getPasswordMD5(body.Password), roleid)
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
