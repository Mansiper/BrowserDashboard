package main

import (
	"log"
	"strings"
	"time"
	"database/sql"
	"strconv"
)

func DbStat() {
	const (
		mc int64 = 2592000
		dc int64 = 86400
		hc int64 = 3600
		nc int64 = 60
	)
	var (
		err error
		rows *sql.Rows
		key string
		val string
		oldReceived, oldSent, oldQueries, oldDelete, oldInsert, oldSelect, oldUpdate int64
		speed, uptimeStr string
		secs, m, d, h, n, s int64
	)

	for {
		rows, err = dbMain.Query("show global status")
		if err != nil {
			log.Println("Global Status error")
			log.Println(err)
		}
		defer rows.Close()

		for rows.Next() {
			if err := rows.Scan(&key, &val); err != nil {
				log.Println("Global Status read error")
				log.Println(err)
			}
			switch key {
				case "Bytes_received":	//Байт получено (записано)
					speed = ByteToHumanString(oldReceived, val, int64(cDbStatPause.Seconds()))
					oldReceived = GetInt64(val)
					stats.DB_Bytes_received = speed
				case "Bytes_sent":			//Байт отправлено (прочитано)
					speed = ByteToHumanString(oldSent, val, int64(cDbStatPause.Seconds()))
					oldSent = GetInt64(val)
					stats.DB_Bytes_sent = speed
				//case "Com_commit":			//Количество коммитов	? Handler_commit
				//	speed = ByteToHumanString(oldCommit, val.(int64), int64(cDbStatPause.Seconds()))
				//	oldCommit = GetInt64(val)
				//	stats.DB_Com_commit = speed
				case "Com_delete":			//Количество delete ? Handler_delete
					speed = CountToHumanString(oldDelete, val, int64(cDbStatPause.Seconds()), "deletes")
					oldDelete = GetInt64(val)
					stats.DB_Com_delete = speed
				case "Com_insert":			//Количество insert ? Handler_write
					speed = CountToHumanString(oldInsert, val, int64(cDbStatPause.Seconds()), "inserts")
					oldInsert = GetInt64(val)
					stats.DB_Com_insert = speed
				case "Com_select":			//Количество select ? Handler_read_rnd_next - количество по строкам
					speed = CountToHumanString(oldSelect, val, int64(cDbStatPause.Seconds()), "selects")
					oldSelect = GetInt64(val)
					stats.DB_Com_select = speed
				case "Com_update":			//Количество update ? Handler_update
					speed = CountToHumanString(oldUpdate, val, int64(cDbStatPause.Seconds()), "updates")
					oldUpdate = GetInt64(val)
					stats.DB_Com_update = speed
				case "Max_used_connections":			//Пик подключеий
					stats.DB_Max_used_connections = val
				case "Max_used_connections_time":	//Время пика
					stats.DB_Max_used_connections_time = val
				case "Queries":					//Количество запросов
					speed = CountToHumanString(oldQueries, val, int64(cDbStatPause.Seconds()), "queries")
					oldQueries = GetInt64(val)
					stats.DB_Queries = speed
				case "Threads_connected":	//Количество активных соединений
					stats.DB_Threads_connected = val
				case "Threads_running":	//Количество открытых не спящих соединений
					stats.DB_Threads_running = val
				case "Uptime":					//Время работы
				 	secs = GetInt64(val)
					m = int64(secs / mc)
					secs = secs - m * mc
					d = int64(secs / dc)
					secs = secs - d * dc
					h = int64(secs / hc)
					secs = secs - h * hc
					n = int64(secs / nc)
					secs = secs - n * nc
					s = secs
					uptimeStr = ""
					//if y > 0 { uptimeStr += WordNumberCase(y, " год", " года", " лет", true) + " " }
					if m > 0 { uptimeStr += WordNumberCase(int(m), " месяц", " месяца", " месяцев", true) + " " }
					if d > 0 { uptimeStr += WordNumberCase(int(d), " день", " дня", " дней", true) + " " }
					if h > 0 { uptimeStr += WordNumberCase(int(h), " час", " часа", " часов", true) + " " }
					if n > 0 { uptimeStr += WordNumberCase(int(n), " минута", " минуты", " минут", true) + " " }
					if s > 0 { uptimeStr += WordNumberCase(int(s), " секунда", " секунды", " секунд", true) + " " }

					stats.DB_Uptime = strings.Trim(uptimeStr, " ")
			}
		}

		time.Sleep(cDbStatPause)
	}
}

func ByteToHumanString(from int64, toS string, sec int64) string {
	if toS == "" { return "" }
	diff := GetInt64(toS) - from
	if diff > 107374182 * sec /*0.1Gb*/ {
		return strconv.FormatFloat(float64(float64(diff) / 107374182.0 / float64(sec)), 'f', 2, 32) + " GB/s"
	} else if diff > 104857 * sec /*0.1Mb*/ {
		return strconv.FormatFloat(float64(float64(diff) / 104857.0 / float64(sec)), 'f', 2, 32) + " MB/s"
	} else if diff > 102 * sec /*0.1Kb*/ {
		return strconv.FormatFloat(float64(float64(diff) / 102.0 / float64(sec)), 'f', 2, 32) + " KB/s"
	} else {
		return strconv.FormatFloat(float64(diff / sec), 'f', 2, 32) + " B/s"
	}
}

func CountToHumanString(from int64, toS string, sec int64, nums string) string {
	if toS == "" { return "" }
	diff := GetInt64(toS) - from
	return strconv.FormatFloat(float64(diff / sec), 'f', 0, 32) + " " + nums + "/s"
}

func WordNumberCase(number int, ifOne, ifTwo, ifFive string, addNumber bool) string {
	var result string
	switch number % 10 {
		case 1:			result = ifOne
		case 2,3,4:	result = ifTwo
		default:		result = ifFive
	}
	m := number % 100
	if m >=11 && m <= 14 {
		result = ifFive
	}
	if addNumber {
		result = strconv.Itoa(number) + result
	}
	return result
}

func GetInt64(str string) int64 {
	res, err := strconv.ParseInt(str, 10, 64)
	if err != nil { return 0 } else { return res }
}