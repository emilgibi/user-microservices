package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/emilgibi/user-microservices/models"
	resty "github.com/go-resty/resty/v2"
	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Handler struct {
	DB         *gorm.DB
	ORDERSVC   string
	PRODUCTSVC string
	STOCKSVC   string
}

func (handler *Handler) Connect(host, user, pass, dbName, port string) {
	var err error
	dsn := "host=" + host + " user=" + user + " password=" + pass + " dbname=" + dbName + " port=" + port + " sslmode=disable TimeZone=Asia/Shanghai"
	handler.DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	handler.DB.AutoMigrate(models.User{})
	if err != nil {
		panic(err)
	}
}

func (handler *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	var users []models.User
	handler.DB.Find(&users)
	json.NewEncoder(w).Encode(users)
}

func (handler *Handler) GetUserid(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var users models.User
	handler.DB.First(&users, params["id"])
	json.NewEncoder(w).Encode(users)
}

func (handler *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var update models.User
	_ = json.NewDecoder(r.Body).Decode(&update)

	handler.DB.Save(&update)
	json.NewEncoder(w).Encode(update)
}

func (handler *Handler) GetProduct(w http.ResponseWriter, r *http.Request) {

	client := resty.New()
	resp, _ := client.R().
		SetHeader("Content-Type", "application/json").
		Get("http://" + handler.PRODUCTSVC + ":8081/products")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp.Body())
}

func (handler *Handler) GetProductById(w http.ResponseWriter, r *http.Request) {

	client := resty.New()
	resp, _ := client.R().
		SetHeader("Content-Type", "application/json").
		Get("http://" + handler.PRODUCTSVC + ":8081/products{product_id}")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp.Body())
}

func (handler *Handler) CheckStock(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	client := resty.New()
	resp, err := client.R().SetResult(map[string]interface{}{}).Get("http://" + handler.STOCKSVC + ":8082/stock/check/" + params["id"] + params["product_quantity"])
	if err != nil {
		log.Fatal(err)
		return
	}
	type response struct {
		Message string `json:"message"`
		Status  int    `json:"status"`
	}

	result := resp.Result().(*map[string]interface{})
	message, ok := (*result)["message"].(bool)
	if !ok {
		res := response{Message: "Item not available", Status: 0}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(res)
		return
	}
	if message {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(true)
		client := resty.New()
		resp, err := client.R().
			SetHeader("Content-Type", "application/json").
			SetBody(`{"ID": 12345, "ProductName": "asdf", "OrderQuantity": 5}`).
			Post("http://" + handler.ORDERSVC + ":8084/order")

		if err != nil {
			log.Fatal(err)
			return
		}
		fmt.Println(resp)
	} else {
		type Response struct {
			Success bool   `json:"success"`
			Message string `json:"message"`
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Some technical issues have arisen",
		})
	}
}
