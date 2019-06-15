package main

import (
	"net/http"
	"net/http/cookiejar"
	"golang.org/x/net/publicsuffix"
	"log"
	"time"
	"bpfuncs"
)

func AsChecker() {
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

	for {
		resp, err = client.Get("https://companyurl/api/calls")
		data = bpfuncs.GetBody(resp)
		res = data == `...`
		stats.As_check = b2i[res]
		if !res { log.Println(time.Now(), "нет ответа от сервиса звонков") }

		time.Sleep(cAsPause)
	}
}