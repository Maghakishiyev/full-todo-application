package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type DeleteRequestBody struct {
	TodoId int `json:"todoid"`
}

func DeleteTodo(writer http.ResponseWriter, request *http.Request, claims *Claims) {
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

	query := "DELETE from `testing`.`todos` WHERE User_Id = ? AND Id_Todo = ?;"
	deleteTodo, err := tx.Prepare(query)
	if err != nil {
		fmt.Println(err)
		writer.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(writer).Encode("Error: while preparing db!!")
		return
	}

	UserCredentials := claims.UserCredentials

	var deleteRequestBody DeleteRequestBody

	errDecoder := json.NewDecoder(request.Body).Decode(&deleteRequestBody)
	if errDecoder != nil {
		fmt.Println(errDecoder)
		writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writer).Encode("Error: Bad request!!")
		return
	}

	deleteResponse, err := deleteTodo.Exec(UserCredentials.UserId, deleteRequestBody.TodoId)
	if err != nil {
		println(err)
		writer.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(writer).Encode("Error: while deleting todo!!")
		return
	}

	deleteTodo.Close()

	err = tx.Commit()
	if err != nil {
		println(err)
		writer.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(writer).Encode("Error: while deleting todo!!")
		return
	}

	affectedRows, err := deleteResponse.RowsAffected()
	if err != nil {
		println(err)
		writer.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(writer).Encode("Error: while deleting todo!!")
		return
	}

	if affectedRows > 0 {
		writer.WriteHeader(http.StatusOK)
		json.NewEncoder(writer).Encode("Successfully deleted todo!!")
		return
	} else {
		println(err)
		writer.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(writer).Encode("Error: while deleting todo!!(Not authorized)")
		return
	}
}
