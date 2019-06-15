package main

import (
	"time"
	"log"
)

func DbTestChecker() {
	var err error

	for {
		err = dbTest.Ping()
		stats.Dbtest_check = b2i[err == nil]
		if err != nil { log.Println(time.Now(), "нет ответа от DbTest") }

		time.Sleep(cDbTestPause)
	}
}