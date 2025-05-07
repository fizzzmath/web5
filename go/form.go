package main

import (
	"fmt"
	"net/http"
	"text/template"
)

func formHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("form.html")

	if err != nil {
		fmt.Fprintf(w, "Template error: %v", err)
		return
	}

	tmpl.Execute(w, nil)
}