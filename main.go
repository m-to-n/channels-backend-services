package main

import (
	"fmt"
	dapr "github.com/dapr/go-sdk/client"
	"github.com/gorilla/mux"
	"net/http"
)

const DAPR_APP_PORT = "6002"

func bindingHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received Message!")
	/*w.Header().Set("Content-Type", "application/json")
	var orderId int
	err := json.NewDecoder(r.Body).Decode(&orderId)
	fmt.Println("Received Message: ", orderId)
	if err != nil {
		fmt.Println("error parsing checkout input binding payload: %s", err)
		w.WriteHeader(http.StatusOK)
		return
	}*/
}

// https://docs.dapr.io/developing-applications/building-blocks/bindings/howto-triggers/
func main() {
	client, err := dapr.NewClient()
	if err != nil {
		panic(err)
	}
	defer client.Close()

	fmt.Println("Serious DAPR stuff here...!")

	r := mux.NewRouter()
	r.HandleFunc("/channels-backend-services-sqs-wa-twilio", bindingHandler).Methods("POST", "OPTIONS")
	http.ListenAndServe(":"+DAPR_APP_PORT, r)
}
