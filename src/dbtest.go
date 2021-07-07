package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

func a() {
	db, err := sql.Open("sqlite3", "./db.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	//sqlStmt := `
	//create table foo (id integer not null primary key, name text);
	//delete from foo;
	//`
	//_, err = db.Exec(sqlStmt)
	//if err != nil {
	//	log.Printf("%q: %s\n", err, sqlStmt)
	//	return
	//}

	//tx, err := db.Begin()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//stmt, err := tx.Prepare("insert into user(uid, name) values(?, ?)")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer stmt.Close()
	//for i := 0; i < 100; i++ {
	//	_, err = stmt.Exec(i, fmt.Sprintf("んにちわ世界%03d", i))
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//}
	//tx.Commit()

	rows, err := db.Query("select uid, username from user")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var name string
		err = rows.Scan(&id, &name)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(id, name)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := db.Prepare("select uid, password, isadmin from user,role where username = ? and user.roleid = role.roleid")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	var uid int
	var password string
	var isadmin bool
	err = stmt.QueryRow("demo").Scan(&uid, &password, &isadmin)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("uid: %d pass: %s isadmin %t", uid, password, isadmin)

	//_, err = db.Exec("delete from foo")
	//if err != nil {
	//	log.Fatal(err)
	//}

	//_, err = db.Exec("insert into foo(id, name) values(1, 'foo'), (2, 'bar'), (3, 'baz')")
	//if err != nil {
	//	log.Fatal(err)
	//}

}
