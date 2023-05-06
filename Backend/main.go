package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func initializeRouter() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/api/register", RegisterUser).Methods("POST")
	router.HandleFunc("/api/login", LoginUser).Methods("POST")
	router.HandleFunc("/api/login", LogOutUser).Methods("POST")
	router.Handle("/api/refreshtoken", validateJwt(RefreshHandler))
	router.Handle("/api/todos", validateJwtAndReturnClaims(GetUserTodos)).Methods("GET")
	router.Handle("/api/todos/create", validateJwtAndReturnClaims(CreateTodo)).Methods("POST")
	router.Handle("/api/todos/update", validateJwtAndReturnClaims(UpdateTodo)).Methods("POST")
	router.Handle("/api/todos/delete", validateJwtAndReturnClaims(DeleteTodo)).Methods("POST")

	return router
}

func main() {
	mainRouter := initializeRouter()

	println("Server Started\nListening to port :8000")
	log.Fatal(http.ListenAndServe(":8000", mainRouter))
}
