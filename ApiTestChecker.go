package main

import (
	"net/http"
	"log"
	"net/http/cookiejar"
	"golang.org/x/net/publicsuffix"
	"strings"
	"time"
	"bpfuncs"
)

func ApiTestChecker() {
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
		resp, err = client.Get("https://companytestsite/api/test")
		data = bpfuncs.GetBody(resp)
		res = strings.Contains(data, `{"Code":200,"Status":"success","Data":`)
		stats.ApiTest_check = b2i[res]
		if !res { log.Println(time.Now(), "нет ответа от APITest") }

		time.Sleep(cApiPause)
	}
}