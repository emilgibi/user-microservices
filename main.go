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

	HOST := os.Getenv("DATABASE_HOST")
	USER := os.Getenv("DATABASE_USER")
	PASS := os.Getenv("DATABASE_PASS")
	NAME := os.Getenv("DATABASE_NAME")
	PORT := os.Getenv("DATABASE_PORT")

	handlerObj.Connect(HOST, USER, PASS, NAME, PORT)
	handlerObj.ORDERSVC = os.Getenv("ORDER_SVC")
	handlerObj.PRODUCTSVC = os.Getenv("PRODUCT_SVC")
	handlerObj.STOCKSVC = os.Getenv("STOCK_SVC")

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
