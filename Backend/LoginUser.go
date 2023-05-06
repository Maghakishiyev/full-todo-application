package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	UserId    int    `json:"userid"`
	Username  string `json:"username"`
	CreatedAt string `json:"createdat"`
	UpdatedAt string `json:"updatedat"`
}

func LoginUser(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	database := ConnectToDatabase()

	defer database.Close()

	pingErr := database.Ping()

	if pingErr != nil {
		fmt.Println(pingErr.Error())
	}

	var LoginBody LoginRequest

	err := json.NewDecoder(request.Body).Decode(&LoginBody)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writer).Encode("Bad Request, missing credentials !!!")
		return
	}

	if LoginBody.Username == "" || len(LoginBody.Username) < 3 || LoginBody.Password == "" || len(LoginBody.Password) < 6 {
		writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writer).Encode("Error: bad credentials!!")
		return
	}

	selectRow, err := database.Query(fmt.Sprintf("Select Id_User, User_Name, Created_At, Updated_At from `testing`.`users` where User_Name='%s' and User_Password='%s'", LoginBody.Username, LoginBody.Password))

	if err != nil {
		fmt.Println()
		writer.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(writer).Encode("Error: Database error!!")
		return
	}

	defer selectRow.Close()

	if selectRow.Next() {
		var selectedUser User
		var createdAt []uint8
		var updatedAt []uint8

		errScan := selectRow.Scan(
			&selectedUser.UserId,
			&selectedUser.Username,
			&createdAt,
			&updatedAt,
		)

		if errScan != nil {
			fmt.Println(errScan)
			writer.WriteHeader(http.StatusNotFound)
			json.NewEncoder(writer).Encode("Error: Database Error!!")
			return
		}

		selectedUser.CreatedAt = string(createdAt)
		selectedUser.UpdatedAt = string(updatedAt)

		tokenString, expirationTime, err := createJwt(selectedUser)

		if err != nil {
			fmt.Println(err)
			writer.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(writer).Encode("Error: can not create token")
		}

		http.SetCookie(
			writer,
			&http.Cookie{
				Name:    "token",
				Value:   tokenString,
				Expires: expirationTime,
			},
		)

		writer.WriteHeader(http.StatusFound)
		json.NewEncoder(writer).Encode(selectedUser)
		return
	} else {
		writer.WriteHeader(http.StatusNotFound)
		json.NewEncoder(writer).Encode("Error: user not found!!")
		return
	}
}

func LogOutUser(writer http.ResponseWriter, request *http.Request) {
	RemoveCookie(writer)
}
