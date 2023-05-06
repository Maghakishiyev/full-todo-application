package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type UpdateTodoRequest struct {
	TodoId          int    `json:"todoId"`
	TodoName        string `json:"todoname"`
	TodoDescription string `json:"tododescription"`
	TodoState       string `json:"todostate"`
}

func UpdateTodo(writer http.ResponseWriter, request *http.Request, claims *Claims) {
	writer.Header().Set("Content-Type", "application/json")
	database := ConnectToDatabase()

	defer database.Close()

	tx, err := database.Begin()
	if err != nil {
		fmt.Println(err)
		writer.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(writer).Encode("Error: Database error!!")
		return
	}

	UserCredentials := claims.UserCredentials
	var updateRequestBody UpdateTodoRequest

	err = json.NewDecoder(request.Body).Decode(&updateRequestBody)
	if err != nil {
		fmt.Println(err)
		writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writer).Encode("Error: Bad request!!")
		return
	}

	if updateRequestBody.TodoName == "" {
		writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writer).Encode("Error: Very bad request, todoname can not be empty")
		return
	}
	if updateRequestBody.TodoState != "" && updateRequestBody.TodoState != "done" {
		updateRequestBody.TodoState = "pending"
	}

	updateQuery := "Update `testing`.`todos` set Todo_Name=?, Todo_Description=?, Todo_State=?, Updated_At=? where Id_Todo=? and User_Id=?;"

	updateTodo, err := tx.Prepare(updateQuery)
	if err != nil {
		fmt.Println(err)
		writer.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(writer).Encode("Error: while preparing db!!")
		return
	}

	updatedAtVal := time.Now().Local()

	updateResponse, err := updateTodo.Exec(updateRequestBody.TodoName, updateRequestBody.TodoDescription, updateRequestBody.TodoState, updatedAtVal, updateRequestBody.TodoId, UserCredentials.UserId)
	if err != nil {
		println(err)
		writer.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(writer).Encode("Error: while updating todo!!")
		return
	}

	updateTodo.Close()

	err = tx.Commit()
	if err != nil {
		println(err)
		writer.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(writer).Encode("Error: while updating todo!!")
		return
	}

	affectedRows, err := updateResponse.RowsAffected()
	if err != nil {
		println(err)
		writer.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(writer).Encode("Error: while updating todo!!")
		return
	}

	if affectedRows > 0 {
		writer.WriteHeader(http.StatusAccepted)
		json.NewEncoder(writer).Encode("Congratulations: Todo item was successfully updated!!")
		return
	} else {
		fmt.Println(err)
		writer.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(writer).Encode("Error: Can not update todo!!")
		return
	}
}
