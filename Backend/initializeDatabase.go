package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func ConnectToDatabase() *sql.DB {
	database, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/testing?parseTime=true")

	if err != nil {
		fmt.Println(err)
		panic(err.Error())
	}

	fmt.Println("Connected to db")

	database.SetConnMaxLifetime(time.Minute * 3)
	database.SetMaxOpenConns(10)
	database.SetMaxIdleConns(10)

	return database

	// // insert, err := database.Query("INSERT INTO `testing`.`users` (`User_Name`, `User_Password`) VALUES (\"Season\", \"season2801\")")
	// // if err != nil {
	// // 	panic(err.Error())
	// // }

	// // defer insert.Close()

	// selectUsers, err := database.Query("Select * from `testing`.`users`")

	// if err != nil {
	// 	panic(err.Error())
	// }

	// for selectUsers.Next() {
	// 	var userId int
	// 	var userName string
	// 	var UserPassword string
	// 	createdAt := make([]byte, 0)
	// 	updatedAt := make([]byte, 0)
	// 	err := selectUsers.Scan(&userId, &userName, &UserPassword, &createdAt, &updatedAt)

	// 	if err != nil {
	// 		panic(err.Error())
	// 	}

	// 	fmt.Println(userId, userName, UserPassword, string(createdAt), string((updatedAt)))
	// }

	// defer selectUsers.Close()
	// fmt.Println("Successfully insered into users table")
}
