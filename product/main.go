package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Product struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type Response struct {
	Message string `json:"message"`
}

var products = []Product{
	{Id: 1, Name: "Product 1"},
	{Id: 2, Name: "Product 2"},
	{Id: 3, Name: "Product 3"},
}

func list(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(products)
}

func order(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	authReq, _ := json.Marshal(map[string]string{
		"token": token,
	})

	authRes, _ := http.Post("http://localhost:8081/validate", "application/json", bytes.NewBuffer(authReq))
	defer authRes.Body.Close()

	if authRes.StatusCode == http.StatusUnauthorized {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(Response{
			Message: fmt.Sprintf("Cannot order product: access denied"),
		})
		return
	}

	body, _ := ioutil.ReadAll(r.Body)

	var product Product
	json.Unmarshal(body, &product)

	w.Header().Set("Content-Type", "application/json")
	for _, p := range products {
		if p.Id == product.Id {
			if pay(p, token) {
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(Response{
					Message: fmt.Sprintf("Successfully ordered product %d", product.Id),
				})
			} else {
				w.WriteHeader(http.StatusUnprocessableEntity)
				json.NewEncoder(w).Encode(Response{
					Message: fmt.Sprintf("Cannot order product %d: payment failed", product.Id),
				})
			}
			return
		}
	}

	w.WriteHeader(http.StatusUnprocessableEntity)
	json.NewEncoder(w).Encode(Response{
		Message: fmt.Sprintf("Product not found %d", product.Id),
	})
}

func pay(product Product, token string) bool {
	payReq, _ := json.Marshal(product)
	req, _ := http.NewRequest("POST", "http://localhost:8082/pay", bytes.NewBuffer(payReq))
	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, _ := client.Do(req)
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

func main() {
	http.HandleFunc("/list", list)
	http.HandleFunc("/order", order)

	log.Println("Listening Product service on port :8083")
	log.Fatal(http.ListenAndServe(":8083", nil))
}
