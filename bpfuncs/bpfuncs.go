package bpfuncs

import (
	"log"
	"net/url"
	"strings"
	"net/http"
	"io/ioutil"
)

func GetBody(resp *http.Response) string {
	data, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	return string(data)
}

func Login(client *http.Client, login string, password string) bool {
	const (
		siteLogin = "https://companysite/LoginPage"
		sitePage = "https://companysite/Page"
	)

		//Проверка. Не залогинен
	resp, err := (*client).Get(sitePage)
	data := GetBody(resp)
	if strings.Contains(data, "Пожалуйста, войдите в систему") {
		resp, err = (*client).PostForm(siteLogin, url.Values{
			"LoginUser$UserName": {login},
			"LoginUser$Password": {password},
			"LoginUser$LoginButton_input": {"Войти"},
		})
		if err != nil {
			log.Println(err)
			return false
		}
		data = GetBody(resp)
		if strings.Contains(data, "Пожалуйста, войдите в систему") {
			return false
		}
	}
	return true
}