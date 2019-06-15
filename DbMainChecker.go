package main

import (
	"time"
	"log"
)

func DbMainChecker() {
	var err error

	for {
		err = dbMain.Ping()
		stats.Dbmain_check = b2i[err == nil]
		if err != nil { log.Println(time.Now(), "нет ответа от DbMain") }

		time.Sleep(cDbMainPause)
	}
}