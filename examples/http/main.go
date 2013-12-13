package main

import (
	"github.com/mvrilo/mstat"
	"net/http"
)

func simple() {
	http.ListenAndServe(":8000", mstat.New())
}

func modular() {
	m := http.NewServeMux()
	m.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hey"))
	})

	s := mstat.New()
	s.Next = m

	http.ListenAndServe(":8000", s)
}

func main() {
	// simple()
	modular()
}
