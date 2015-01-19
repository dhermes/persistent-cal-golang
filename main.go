package main

import (
	"fmt"
	"net/http"
)

func init() {
	http.HandleFunc("/", handler)
}

func handler(writer http.ResponseWriter, unusedReq *http.Request) {
	fmt.Fprint(writer, "Hello, world!")
}
