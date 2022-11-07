package main

import (
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
)

func main() {

	handler := http.HandlerFunc(GetImages)

	http.Handle("/images", basicAuth(handler))
	fmt.Println("Server started at port 8080")
	fmt.Println("Available endpoint: `/images`")
	fmt.Println("To access the endpoint, add basic authorization header with username and password")
	fmt.Println("Add query parameters to the request: [x, y, z] (int)")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		return
	}

}
