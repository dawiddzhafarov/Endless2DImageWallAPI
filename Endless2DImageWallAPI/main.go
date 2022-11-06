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

	if err := http.ListenAndServe(":8080", nil); err != nil {
		return
	}

}
