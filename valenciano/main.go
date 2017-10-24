package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

func main() {

	_, err := os.Open("public_html/index.html")
	if err != nil {
		t, err := template.ParseFiles("template/index.html")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		data := struct {
			PrecoSolteiro string
			PrecoCasal    string
		}{
			PrecoSolteiro: "65",
			PrecoCasal:    "110",
		}
		f, err := os.Create("public_html/index.html")
		if err != nil {
			fmt.Println(err)
		}
		t.Execute(f, data)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("In static file handler", r.URL.Path)
		http.ServeFile(w, r, "public_html/"+r.URL.Path[1:])
	})

	log.Fatal(http.ListenAndServe(":8080", nil))

}
