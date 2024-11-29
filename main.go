package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Response struct {
	Message string `json:"message"`
}

func longRunningTask(done chan string) {

	time.Sleep(6 * time.Second)

	done <- "Task done successfully"
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	// response obj

	response := Response{
		Message: "Hello this is mc",
	}

	//create a channel

	done := make(chan string)

	go longRunningTask(done)

	//set the header

	w.Header().Set("Content-Type", "application/json")

	//wait for the task to complete

	//select case
	select {

	case result := <-done:

		response.Message = fmt.Sprintf("%s - %s", response.Message, result)

		w.WriteHeader(http.StatusOK) // send HTTP 200 OK

		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Println("Error encoding response", err)
		}

	case <-time.After(5 * time.Second):

		w.WriteHeader(http.StatusRequestTimeout) // send HTTP 408 Request Timeout

		errorResponse := Response{
			Message: "Request TImed out successfully",
		}

		if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
			log.Println("Error encoding response", err)
		}

	}

}

func main() {

	router := mux.NewRouter()

	router.HandleFunc("/api", helloHandler).Methods("GET")

	fmt.Println("Starting server on the port 8080...")
	log.Fatal(http.ListenAndServe(":8080", router))
}
