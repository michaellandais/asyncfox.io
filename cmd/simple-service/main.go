package main

import (
	"errors"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
)

type Page struct {
	Name string
}

var templates = template.Must(template.ParseFiles("index.html"))

func getRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got / request\n")
	io.WriteString(w, "Hello !")
}

func getOS(w http.ResponseWriter, r *http.Request) {
	info, _ := GetOSInfo()
	p := Page{Name: info.Name}
	err := templates.ExecuteTemplate(w, "index.html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", getRoot)
	mux.HandleFunc("/os", getOS)

	err := http.ListenAndServe(":3000", mux)

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
