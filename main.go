package main

import (
	"fmt"
	dapr "github.com/dapr/go-sdk/client"
)

func main() {
	client, err := dapr.NewClient()
	if err != nil {
		panic(err)
	}
	defer client.Close()
	fmt.Println("Serious DAPR stuff here...!")
}
