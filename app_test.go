package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

var a App

func TestMain(m *testing.M) {
	err := a.Initialise(DB_USER, DB_USER_PASSWORD, DB_HOST, "test")
	if err != nil {
		log.Fatal("error occured zhile initialising the database")
	}
	createTable()
	m.Run()

}

func createTable() {
	createTableQuery := `CREATE TABLE IF NOT EXISTS products (
		id INT NOT NULL AUTO_INCREMENT, 
		name VARCHAR(255) NOT NULL, 
		quantity INT, 
		price DECIMAL(10, 7),
		PRIMARY KEY(id)
	);`
	_, err := a.DB.Exec(createTableQuery)
	if err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	a.DB.Exec("DELETE from products")
	a.DB.Exec("ALTER table products AUTO_INCREMENT=1")
	log.Println("clear table")
}

func addProduct(name string, quantity int, price float64) {
	query := fmt.Sprintf("insert into products(name,quantity,price) values('%v', %v, %v)", name, quantity, price)
	_, err := a.DB.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}
func TestGetProduct(t *testing.T) {
	clearTable()

	addProduct("keyboard", 100, 22.3)
	request, _ := http.NewRequest("GET", "/products/1", nil)
	response := sendRequest(request)
	checkStatusCode(t, http.StatusOK, response.Code)

}

func sendRequest(request *http.Request) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	a.Router.ServeHTTP(recorder, request)
	return recorder
}

func checkStatusCode(t *testing.T, expected int, actual int) {
	if expected != actual {
		t.Errorf("Expected status: %v, received %v", expected, actual)
	}

}

func TestCreateProduct(t *testing.T) {
	clearTable()
	product := []byte(`{"name":"chair", "quantity":1, "price":100}`)
	req, _ := http.NewRequest("POST", "/products", bytes.NewBuffer(product))
	req.Header.Set("Content-Type", "application/json")
	response := sendRequest(req)
	checkStatusCode(t, http.StatusCreated, response.Code)
	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["name"] != "chair" {
		t.Errorf("Expected name: %v, get %v", "chair", m["name"])

	}
	// convert number to float64
	if m["quantity"] != 1.0 {
		t.Errorf("Expected quantity: %v, get %v", 1, m["quantity"])

	}
	if m["price"] != 100.0 {
		t.Errorf("Expected price: %v, get %v", 100, m["price"])

	}
}

func TestDeleteProduct(t *testing.T) {
	clearTable()
	// addProduct("connector", 10, 10)
	// request, _ := http.NewRequest("GET", "/products/1", nil)
	// response := sendRequest(request)
	// checkStatusCode(t, http.StatusOK, response.Code)

	// request, _ = http.NewRequest("DELETE", "/products/1", nil)
	// response = sendRequest(request)
	// checkStatusCode(t, http.StatusOK, response.Code)

	request, _ := http.NewRequest("GET", "/products/1", nil)
	response := sendRequest(request)
	checkStatusCode(t, http.StatusNotFound, response.Code)

}

func TestUpdateProduct(t *testing.T) {
	clearTable()
	addProduct("connector", 10, 10)
	request, _ := http.NewRequest("GET", "/products/1", nil)
	response := sendRequest(request)
	checkStatusCode(t, http.StatusOK, response.Code)

	var oldValue map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &oldValue)

	product := []byte(`{"name":"connector", "quantity":1, "price":10}`)
	request, _ = http.NewRequest("PUT", "/products/1", bytes.NewBuffer(product))
	response = sendRequest(request)
	checkStatusCode(t, http.StatusOK, response.Code)

	var newValue map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &newValue)

	if newValue["id"] != oldValue["id"] {
		t.Errorf("Expected price: %v, get %v", oldValue["id"], newValue["id"])

	}

	if newValue["name"] != oldValue["name"] {
		t.Errorf("Expected price: %v, get %v", oldValue["name"], newValue["name"])

	}

	if newValue["price"] != oldValue["price"] {
		t.Errorf("Expected price: %v, get %v", oldValue["price"], newValue["price"])

	}

	if newValue["quantity"] == oldValue["quantity"] {
		t.Errorf("Expected price: %v, get %v", oldValue["quantity"], newValue["quantity"])

	}
}
