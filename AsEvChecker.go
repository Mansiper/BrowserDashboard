package main

import (
	"log"
	"time"
	"database/sql"
)

func AsEvChecker() {
	var (
		err error
		row *sql.Row
		dt, oldDt string
		res bool
	)

	for {
		row = dbMain.QueryRow("select time_field from calls_table order by time_field desc limit 1")
		if row == nil { log.Println("PhoneCall row is nil") }
		err = row.Scan(&dt)
		if err != nil { log.Println(err.Error()) }
		res = dt != oldDt
		stats.Asev_check = b2i[res]
		oldDt = dt
		if !res { log.Println(time.Now(), "нет ответа от PhoneCall") }

		time.Sleep(cAsEvPause)
	}
}