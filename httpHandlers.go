package main

import (
	"net/http"
)

func spam(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/spam.html")
}

func listen(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/listen.html")
}

func squawk(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/squawk.html")
}
