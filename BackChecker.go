package main

import (
	"net/http"
	"log"
	"net/http/cookiejar"
	"golang.org/x/net/publicsuffix"
	"bpfuncs"
	"strings"
	"time"
)

func BackChecker() {
	var (
		err error
		resp *http.Response
		data string
		res bool
	)

	options := cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	}
	jar, err := cookiejar.New(&options)
	if err != nil { log.Fatal(err) }
	client := http.Client{Jar: jar}
	bpfuncs.Login(&client, "login", "password", false)

	for {
		resp, err = client.Get("https://companysite/Page1")
		data = bpfuncs.GetBody(resp)
		res = strings.Contains(data, "...")
		stats.Back_check_1 = b2i[res]
		if !res { log.Println(time.Now(), "нет ответа от Бэка (1)") }

		resp, err = client.Get("https://companysite/Page2")
		data = bpfuncs.GetBody(resp)
		res = strings.Contains(data, "...")
		stats.Back_check_2 = b2i[res]
		if !res { log.Println(time.Now(), "нет ответа от Бэка (2)") }

		time.Sleep(cBackPause)
	}
}