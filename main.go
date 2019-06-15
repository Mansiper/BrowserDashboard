package main

import (
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"log"
	"time"
	"os"
	"github.com/go-martini/martini"
	"net/http"
	"encoding/json"
)

type Data struct {
	Back_check_1 int										`json:"back_check_1"`
	Back_check_2 int										`json:"back_check_2"`
	BackTest_check_1 int								`json:"backtest_check_1"`
	BackTest_check_2 int								`json:"backtest_check_2"`
	Api_check int												`json:"api_check"`
	ApiTest_check int										`json:"apitest_check"`
	Api2_check int											`json:"api2_check"`
	Api2Test_check int									`json:"api2test_check"`
	Gps_check int												`json:"gps_check"`
	As_check int												`json:"as_check"`
	Asev_check int											`json:"asev_check"`
	Dbmain_check int										`json:"dbmain_check"`
	Dbtest_check int										`json:"dbtest_check"`
	DB_Bytes_sent string								`json:"DB_Bytes_sent"`
	DB_Bytes_received string						`json:"DB_Bytes_received"`
	DB_Queries string										`json:"DB_Queries"`
	DB_Com_select string								`json:"DB_Com_select"`
	DB_Com_insert string								`json:"DB_Com_insert"`
	DB_Com_update string								`json:"DB_Com_update"`
	DB_Com_delete string								`json:"DB_Com_delete"`
	DB_Threads_connected string					`json:"DB_Threads_connected"`
	DB_Threads_running string						`json:"DB_Threads_running"`
	DB_Max_used_connections_time string	`json:"DB_Max_used_connections_time"`
	DB_Max_used_connections string			`json:"DB_Max_used_connections"`
	DB_Uptime string										`json:"DB_Uptime"`
}

const (
	cBackPause = 5 * time.Second
	cBackTestPause = 5 * time.Second
	cApiPause = 5 * time.Second
	cApi2Pause = 5 * time.Second
	cGpsPause = 30 * time.Second
	cAsPause = 5 * time.Second
	cAsEvPause = 5 * time.Minute
	cDbMainPause = 3 * time.Second
	cDbTestPause = 3 * time.Second
	cDbStatPause = 2 * time.Second
)
var (
	dbMain, dbTest *sql.DB
	b2i = map[bool]int{false: 0, true: 1}
	stats Data
)

//----------------------------------------------------------------------------------------------------------------------

func main() {
	var (
		err error
		lf *os.File
	)

	fileName := `logs\` + time.Now().Format("2006.01.02_15.04.05")+".log"
	lf, err = os.Create(fileName)
	if err != nil { panic(err) }
	log.SetOutput(lf)

	//Подключение к боевой базе данных
	dbMain, err = sql.Open("mysql",
		"login:password@tcp(companydburl)/database")
	defer dbMain.Close()
	if err != nil {
		log.Fatal(err)
	} else if err = dbMain.Ping(); err != nil {
		log.Fatal(err)
	}
	//Подключение к тестовой базе данных
	dbTest, err = sql.Open("mysql",
		"login:password@tcp(companytestdburl)/database")
	defer dbMain.Close()
	if err != nil {
		log.Fatal(err)
	} else if err = dbMain.Ping(); err != nil {
		log.Fatal(err)
	}

	go BackChecker()
	go BackTestChecker()
	go ApiChecker()
	go ApiTestChecker()
	go Api2Checker()
	go Api2TestChecker()
	go GpsChecker()
	go AsChecker()
	go AsEvChecker()
	go DbMainChecker()
	go DbTestChecker()
	go DbStat()

	//Подключение сервера
	m := martini.Classic()
	m.Options("/dbdata", optsQuery)
	m.Get("/dbdata", GetDashboardData)
	m.Use(martini.Static(`C:\www\BpDashboard\files`))
	http.Handle("/", m)
	m.RunOnAddr(":123")
}

//----------------------------------------------------------------------------------------------------------------------

func GetDashboardData(res http.ResponseWriter, req *http.Request) {
	defer func() {
		if x := recover(); x != nil {
			log.Printf("Data request caught panic: %v", x)
			http.Error(res, "Something wrong has happend. Try later.", 500)
		}
	}()

	var retNum int

	str, err := json.Marshal(stats)
	if err != nil {
		retNum = 400
		log.Printf(err.Error())
	} else {
		retNum = 200
	}

	res.Header().Add("Access-Control-Allow-Origin", "*")
	res.Header().Add("Access-Control-Allow-Methods", "GET,HEAD,OPTIONS,TRACE")
	res.Header().Add("Access-Control-Allow-Headers", "Accept, X-Requested-With, Content-type")
	res.WriteHeader(retNum)
	res.Write(str)
}
func optsQuery(res http.ResponseWriter, req *http.Request, params martini.Params) {
	res.Header().Add("Access-Control-Allow-Origin", "*")
	res.Header().Add("Access-Control-Allow-Methods", "GET,HEAD,OPTIONS,TRACE")
	res.Header().Add("Access-Control-Allow-Headers", "Accept, X-Requested-With, Content-type")
	res.WriteHeader(200)
}