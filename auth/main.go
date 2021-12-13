package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type TokenRequest struct {
	Token string `json:"token"`
}

type Response struct {
	Message string `json:"message"`
}

func authenticate(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	var credentials AuthRequest
	json.Unmarshal(body, &credentials)

	w.Header().Set("Content-Type", "application/json")
	if credentials.Username == "username" && credentials.Password == "password" {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Response{
			Message: "TOKEN: SECRET AUTH TOKEN",
		})
	} else {
		w.WriteHeader(http.StatusUnprocessableEntity)
	}
}

func validate(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	var credentials TokenRequest
	err = json.Unmarshal(body, &credentials)

	w.Header().Set("Content-Type", "application/json")
	if credentials.Token == "TOKEN: SECRET AUTH TOKEN" {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Response{
			Message: "VALID TOKEN",
		})
	} else {
		w.WriteHeader(http.StatusUnauthorized)
	}
}

func main() {
	http.HandleFunc("/authenticate", authenticate)
	http.HandleFunc("/validate", validate)

	log.Println("Listening Authentication service on port :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
