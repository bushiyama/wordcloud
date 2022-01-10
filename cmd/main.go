package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, r.URL.RawQuery)
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("http.ListenAndServe:", err)
	}
}
