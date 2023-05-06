package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type RegisterRequest struct {
	Username        string `json:"username"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmpassword"`
}

func RegisterUser(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	database := ConnectToDatabase()

	defer database.Close()

	pingErr := database.Ping()

	if pingErr != nil {
		panic(pingErr.Error())
	}

	var registerBody RegisterRequest

	err := json.NewDecoder(request.Body).Decode(&registerBody)

	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writer).Encode("Error: missing credentials!!")
		return
	}

	if registerBody.Username == "" || len(registerBody.Username) < 3 || registerBody.Password == "" || len(registerBody.Password) < 6 || registerBody.Password != registerBody.ConfirmPassword {
		writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writer).Encode("Error: wrong credentials!!")
		return
	}

	selectQuery, err := database.Query(fmt.Sprintf("Select * from `testing`.`users` where User_Name='%s'", registerBody.Username))
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err.Error())
		json.NewEncoder(writer).Encode("Error: database error")
		return
	}

	if selectQuery.Next() {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(writer).Encode("Error: User allready exists, please try new username")
		return
	}

	tx, err := database.Begin()
	if err != nil {
		fmt.Println(err)
		writer.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(writer).Encode("Error: Database error!!")
		return
	}

	createQuery := "Insert into `testing`.`users` (`User_Name`, `User_Password`) Values (?, ?)"

	createUser, err := tx.Prepare(createQuery)
	if err != nil {
		fmt.Println(err)
		writer.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(writer).Encode("Error: while preparing db!!")
		return
	}

	createUserResponse, err := createUser.Exec(registerBody.Username, registerBody.Password)
	if err != nil {
		println(err)
		writer.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(writer).Encode("Error: while creating User!!")
		return
	}

	createUser.Close()
	err = tx.Commit()
	if err != nil {
		println(err)
		writer.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(writer).Encode("Error: while creating user!!")
		return
	}

	affectedRows, err := createUserResponse.RowsAffected()
	if err != nil {
		println(err)
		writer.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(writer).Encode("Error: while creating user!!")
		return
	}

	if affectedRows > 0 {
		writer.WriteHeader(http.StatusAccepted)
		json.NewEncoder(writer).Encode("Congratulations: User was successfully created!!")
		return
	} else {
		fmt.Println(err)
		writer.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(writer).Encode("Error: Can not create user!!")
		return
	}
}
