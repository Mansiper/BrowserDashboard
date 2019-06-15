package main

import (
	"net/http"
	"log"
	"net/http/cookiejar"
	"golang.org/x/net/publicsuffix"
	bp "bpfuncs"
	"strings"
	"time"
)

func Api2Checker() {
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
		resp, err = client.Get("https://companysite/api2/test")
		data = bp.GetBody(resp)
		res = strings.Contains(data, `{"Code":200,"Status":"success","Data":`)
		stats.Api2_check = b2i[res]
		if !res { log.Println(time.Now(), "нет ответа от API2") }

		time.Sleep(cApi2Pause)
	}
}