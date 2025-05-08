package main

import (
	"fmt"
	"net/http"
	"regexp"
	"text/template"

	"github.com/golang-jwt/jwt/v5"
)

func validate(appl Application) Errors {
	var re *regexp.Regexp

	var errors Errors

	pattern := `^([А-ЯA-Z][а-яa-z]+ ){2}[А-ЯA-Z][а-яa-z]+$`
	re = regexp.MustCompile(pattern)

	if appl.Fio == "" {
		errors.Fio = "Поле должно быть заполнено"
	} else if !re.MatchString(appl.Fio) {
		errors.Fio = "Поле должно быть заполнено в формате: Фамилия Имя Отчество"
	}

	pattern = `^(\+7|8)9\d{9}$`
	re = regexp.MustCompile(pattern)

	if appl.Phone == "" {
		errors.Phone = "Поле должно быть заполнено"
	} else if !re.MatchString(appl.Phone) {
		errors.Phone = "Поле должно быть заполнено в формате: +79XXXXXXXXX или 89XXXXXXXXX"
	}

	pattern = `^[A-Za-z][\w\.-_]+@\w+(\.[a-z]{2,})+$`
	re = regexp.MustCompile(pattern)

	if appl.Email == "" {
		errors.Email = "Поле должно быть заполнено"
	} else if !re.MatchString(appl.Email) {
		errors.Email = "Поле должно быть заполнено в формате: имя@домен"
	}

	if appl.Birthdate == "" {
		errors.Birthdate = "Поле должно быть заполнено"
	}

	if appl.Gender == "" {
		errors.Gender = "Поле должно быть заполнено"
	}

	if appl.Bio == "" {
  		errors.Bio = "Поле должно быть заполнено"
	}

	if len(appl.Langs) == 0 {
		errors.Langs = "Поле должно быть заполнено"
	}

	return errors
}

func (e Errors) Count() int {
	count := 0

	if e.Fio != "" { count++ }
	if e.Phone != "" { count++ }
	if e.Email != "" { count++ }
	if e.Birthdate != "" { count++ }
	if e.Gender != "" { count++ }
	if e.Bio != "" { count++ }
	if e.Langs != "" { count++ }

	return count
}

func isAuthorized(r *http.Request) bool {
	cookie, err := r.Cookie("accessToken")

	if err != nil {
		return false
	}

	token, err := jwt.ParseWithClaims(cookie.Value, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("access-token-secret-key"), nil
	})

	return err == nil && token.Valid
}

func formHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("form.html")

	if err != nil {
		fmt.Fprintf(w, "Template error: %v", err)
		return
	}

	response := FormResponse{}

	if r.Method == http.MethodPost {
		err := r.ParseForm()

		if err != nil {
			fmt.Fprintf(w, "Form parsing error: %v", err)
			return
		}

		id := r.URL.Query().Get("id")
		appl := Application{
			Fio: r.FormValue("fio"),
			Phone: r.FormValue("phone"),
			Email: r.FormValue("email"),
			Birthdate: r.FormValue("birthdate"),
			Gender: r.FormValue("gender"),
			Bio: r.FormValue("bio"),
			Langs: r.PostForm["langs[]"],
		}
		errors := validate(appl)

		response = FormResponse{
			ID: id,
			Application: appl,
			Errors: errors,
		}

		if errors.Count() != 0 {
			tmpl.Execute(w, response)
			return
		}

		if id != "" && isAuthorized(r) {
			// err := updateApplication(appl)

			// if err != nil {
			// 	fmt.Fprintf(w, "MySQL error: %v", err)
			// 	return
			// }

			response.Message = "Ваши данные успешно изменены!"

			tmpl.Execute(w, response)
			return
		}

		response := LoginResponse{}
		
		if id != "" {
			response.Type = "warning_red"
			response.Message = "Для выполнения этого действия необходимо авторизоваться"
		} else {
			//insertApplication(appl)
			//login, password, err := generateLAP()

			// if err != nil {
			// 	fmt.Fprintf(w, "MySQL error: %v", err)
			// 	return
			// }

			response.Login = "u0000002";
			response.Password = "qwertyman2345678";
			response.Type = "warning_green"
			response.Message = "Вы успешно зарегистрировались. Перед нажатием на кнопку Войти сохраните ваш логин и пароль!"
		}

		tmpl, err := template.ParseFiles("login.html")

		if err != nil {
			fmt.Fprintf(w, "Template error: %v", err)
			return
		}

		tmpl.Execute(w, response)
		return
	}

	tmpl.Execute(w, response)
}