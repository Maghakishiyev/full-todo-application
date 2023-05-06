package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type CreateTodoRequest struct {
	TodoName        string `json:"todoname"`
	TodoDescription string `json:"tododescription"`
}

func CreateTodo(writer http.ResponseWriter, request *http.Request, claims *Claims) {
	writer.Header().Set("Content-Type", "application/json")
	database := ConnectToDatabase()

	defer database.Close()

	UserCredentials := claims.UserCredentials

	var createRequestBody CreateTodoRequest

	err := json.NewDecoder(request.Body).Decode(&createRequestBody)
	if err != nil {
		fmt.Println(err)
		writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writer).Encode("Error: Bad request!!")
		return
	}

	tx, err := database.Begin()
	if err != nil {
		fmt.Println(err)
		writer.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(writer).Encode("Error: Database error!!")
		return
	}

	insertQuery := "Insert into `testing`.`todos` (User_Id, Todo_Name, Todo_Description) values (?, ?, ?)"

	insertTodo, err := tx.Prepare(insertQuery)
	if err != nil {
		fmt.Println(err)
		writer.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(writer).Encode("Error: while preparing db!!")
		return
	}

	insertResponse, err := insertTodo.Exec(UserCredentials.UserId, createRequestBody.TodoName, createRequestBody.TodoDescription)
	if err != nil {
		println(err)
		writer.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(writer).Encode("Error: while deleting todo!!")
		return
	}

	insertTodo.Close()

	err = tx.Commit()
	if err != nil {
		println(err)
		writer.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(writer).Encode("Error: while creating todo!!")
		return
	}

	affectedRows, err := insertResponse.RowsAffected()
	if err != nil {
		fmt.Println(err)
		writer.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(writer).Encode("Error: Can not create todo!!")
	}

	if affectedRows > 0 {
		writer.WriteHeader(http.StatusAccepted)
		json.NewEncoder(writer).Encode("Congratulations: Todo item created successfully!!")
		return
	} else {
		fmt.Println(err)
		writer.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(writer).Encode("Error: Can not create todo!!")
		return
	}
}
