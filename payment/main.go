package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type Product struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

func pay(w http.ResponseWriter, r *http.Request) {
	authReq, _ := json.Marshal(map[string]string{
		"token": r.Header.Get("Authorization"),
	})

	authRes, _ := http.Post("http://auth:8081/validate", "application/json", bytes.NewBuffer(authReq))
	defer authRes.Body.Close()

	if authRes.StatusCode == http.StatusUnauthorized {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	body, _ := ioutil.ReadAll(r.Body)

	var product Product
	json.Unmarshal(body, &product)
	body, _ = ioutil.ReadAll(r.Body)

	w.WriteHeader(http.StatusOK)
}

func main() {
	http.HandleFunc("/pay", pay)

	log.Println("Listening Payment service on port :8082")
	log.Fatal(http.ListenAndServe(":8082", nil))
}
