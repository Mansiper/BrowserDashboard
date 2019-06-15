package main

import (
	"log"
	"time"
	"database/sql"
)

func GpsChecker() {
	var (
		err error
		row *sql.Row
		dt, oldDt string
		res bool
	)

	for {
		row = dbMain.QueryRow("select time_field from gps_table order by time_field desc limit 1")
		if row == nil { log.Println("GPS row is nil") }
		err = row.Scan(&dt)
		if err != nil { log.Println(err.Error()) }
		res = dt != oldDt
		stats.Gps_check = b2i[res]
		oldDt = dt
		if !res { log.Println(time.Now(), "нет ответа от GPS") }

		time.Sleep(cGpsPause)
	}
}