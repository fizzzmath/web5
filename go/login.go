package main

import (
	"crypto/sha256"
	"database/sql"
	"fmt"
	"net/http"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

type (
	loginResponse struct {
		Login string
		Password string
		Message string
		Type string
	}

	Application struct {
		Fio string
		Phone string
		Email string
		Birthdate string
		Gender string
		Bio string
		Langs []string
	}

	Errors struct {
		Fio string
		Phone string
		Email string
		Birthdate string
		Gender string
		Bio string
		Langs string
	}

	formResponse struct {
		ID string
		Application Application
		Errors Errors
		Message string
	}
)

func dataIsCorrect(login string, password string) (bool, error) {
	db, err := sql.Open("mysql", "u68867:6788851@/u68867")

	if err != nil {
		return false, err
	}

	defer db.Close();

	p := ""

	sel, err := db.Query(`
		SELECT PASSWORD
		FROM USER
		WHERE LOGIN = ?;
	`, login)
	
	if err != nil {
		return false, err
	}

	defer sel.Close();

	for sel.Next() {
		err := sel.Scan(&p)

		if err != nil {
			return false, err
		}
	}

	return fmt.Sprintf("%x", sha256.Sum256([]byte(password))) == p, nil
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("login.html")

	if err != nil {
		fmt.Fprintf(w, "Template error: %v", err)
		return
	}

	response := loginResponse{}
	
	if r.Method == http.MethodPost {
		login := r.FormValue("login")
		password := r.FormValue("password")

		valid, err := dataIsCorrect(login, password)

		if err != nil {
			fmt.Fprintf(w, "MySQL error: %v", err)
			return
		}

		if !valid {
			response.Type = "error"
			response.Message = "Неверные логин или пароль"
			tmpl.Execute(w, response)
			return
		}

		tmpl, err = template.ParseFiles("form.html")

		if err != nil {
			fmt.Fprintf(w, "Template error: %v", err)
			return
		}

		response := loginResponse{Login: login}

		tmpl.Execute(w, response)
		return
	}

	tmpl.Execute(w, response)
}