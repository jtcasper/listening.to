package main

import (
	"fmt"
	"log"
	"net/http"
)

type indexHandler struct{}

func (h *indexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

type callbackHandler struct{}

func (h *callbackHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%+v\n", r)
}

func main() {
	http.Handle("/", &indexHandler{})
	http.Handle("/callback", &callbackHandler{})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
