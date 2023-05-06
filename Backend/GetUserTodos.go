package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type todo struct {
	TodoId          int    `json:"todoid"`
	UserId          int    `json:"userid"`
	TodoName        string `json:"todoname"`
	TodoDescription string `json:"tododescription"`
	CreatedAt       string `json:"createdat"`
	UpdatedAt       string `json:"updatedat"`
	TodoState       string `json:"todostate"`
}

func GetUserTodos(writer http.ResponseWriter, request *http.Request, claims *Claims) {
	writer.Header().Set("Content-Type", "application/json")
	database := ConnectToDatabase()

	defer database.Close()

	userId := claims.UserCredentials.UserId

	rows, err := database.Query(fmt.Sprintf("Select * from `testing`.`todos` where user_id='%d';", userId))
	if err != nil {
		fmt.Println(err)
		writer.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(writer).Encode("Error: Database Error!!")
		return
	}

	todos := make([]todo, 0)

	for rows.Next() {
		var todoItem todo
		var description []byte
		var createdAt []uint8
		var updatedAt []uint8

		if err := rows.Scan(&todoItem.TodoId, &todoItem.UserId, &todoItem.TodoName, &description, &createdAt, &updatedAt, &todoItem.TodoState); err != nil {
			fmt.Println(err)
			writer.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(writer).Encode("Error: Database Error!!")
			return
		}

		todoItem.CreatedAt = string(createdAt)
		todoItem.UpdatedAt = string(updatedAt)
		todoItem.TodoDescription = string(description)

		todos = append(todos, todoItem)
	}

	defer rows.Close()

	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(todos)
}
