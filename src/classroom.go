package main

type Classroom struct {
	Id int `json:"id"`
	Name string `json:"name"`
	GroupId int `json:"group_id"`
	GroupName string `json:"group_name"`
}