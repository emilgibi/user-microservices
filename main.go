package main

import (
	"log"
	"net/http"
	"os"

	"github.com/emilgibi/user-microservices/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {

	var handlerObj handlers.Handler

	godotenv.Load()

	HOST := os.Getenv("DB_HOST")
	USER := os.Getenv("USER_NAME")
	PASS := os.Getenv("PASS")

	handlerObj.Connect(HOST, USER, PASS, "postgres", "5432")

	dbinstance, _ := handlerObj.DB.DB()
	defer dbinstance.Close()

	router := mux.NewRouter()

	router.HandleFunc("/user", handlerObj.GetUser).Methods("GET")
	router.HandleFunc("/user/{user_id}", handlerObj.GetUserid).Methods("GET")
	router.HandleFunc("/user", handlerObj.UpdateUser).Methods("POST")
	router.HandleFunc("/products", handlerObj.GetProduct).Methods("GET")
	router.HandleFunc("/product/{product_id}", handlerObj.GetProductById).Methods("GET")
	router.HandleFunc("/stock/check", handlerObj.CheckStock).Methods("GET")
	http.Handle("/", router)

	//start and listen to requests
	log.Fatal(http.ListenAndServe(":8083", nil))
}
