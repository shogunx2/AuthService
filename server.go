package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	DAO "github.com/shogunx2/AuthService/DAO"
	Services "github.com/shogunx2/AuthService/services"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, HTTP")
}

func main() {
	/*
		amd := AuthMapDatastore{}
		amd.Init()
	*/

	apgd := DAO.AuthPGDatastore{}
	apgd.Init()

	as := Services.AuthService{}
	as.Init(&apgd)

	http.HandleFunc("/", helloHandler)
	http.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("add function start")
		var authReq Services.AuthRequest
		fmt.Println("Body: ", r.Body)
		err := json.NewDecoder(r.Body).Decode(&authReq)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			fmt.Println("err1 called")
			return
		}
		fmt.Println(authReq)

		authResp, err := as.Add(&authReq)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			fmt.Println("err2 called")
			return
		}

		json.NewEncoder(w).Encode(authResp)
		fmt.Println("add function called")
		fmt.Fprintf(w, "add function called")
	})

	http.HandleFunc("/remove", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("add function start")
		var authReq Services.AuthRequest
		fmt.Println("Body: ", r.Body)
		err := json.NewDecoder(r.Body).Decode(&authReq)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			fmt.Println("err1 called")
			return
		}
		fmt.Println(authReq)

		authResp, err := as.Remove(&authReq)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			fmt.Println("err2 called")
			return
		}

		json.NewEncoder(w).Encode(authResp)
		fmt.Println("remove function called")
		fmt.Fprintf(w, "remove function called")
	})

	http.HandleFunc("/authenticate", func(w http.ResponseWriter, r *http.Request) {
		var authReq Services.AuthRequest
		err := json.NewDecoder(r.Body).Decode(&authReq)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Println("-------------------------")
		fmt.Println(authReq)

		authResult, err := as.Authenticate(&authReq)
		fmt.Println(authResult)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		json.NewEncoder(w).Encode(map[string]bool{"authenticated": authResult})
		fmt.Println("authenticate function called")
		fmt.Fprintf(w, "authenticate function called")
		fmt.Println("*******************************")
	})

	fmt.Println("Starting server on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}
