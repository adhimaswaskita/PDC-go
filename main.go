package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"github.com/sparrc/go-ping"
)

type pingStruct struct {
	time         string
	deviceStatus string
	ip           string
	deviceName   string
}

func main() {
	var host string
	var deviceList = []string{"www.jeager.io", "www.ketitik.com"}
	fmt.Print("Do you want to start ping ? (y / n) ")
	fmt.Scanln(&host)

	if host != "y" {
		fmt.Println("Ping app canceled")
	}

	for range time.Tick(time.Second * 5) {
		for i := 0; i < len(deviceList); i++ {
			pinger, err := ping.NewPinger(deviceList[i])
			if err != nil {
				panic(err)
			}

			pinger.SetPrivileged(true)

			pinger.Count = 1
			pinger.Run()
			stats := pinger.Statistics()

			var deviceName = stats.Addr
			var ip = stats.IPAddr.String()
			var deviceStatus string

			for i := 0; i < len(stats.Rtts); i++ {
				var upTime = stats.Rtts[i].String()
				if stats.Rtts[i] != time.Duration(0)*time.Second {
					deviceStatus = "Aktif"
					insert(deviceName, ip, deviceStatus, upTime)
				} else {
					deviceStatus = "Tidak aktif"
					insert(deviceName, ip, deviceStatus, upTime)
				}
			}
		}
	}
}

func connect() (*sql.DB, error) {
	connStr := "postgres://postgres:postgres@127.0.0.1/pingReport"
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		return nil, err
	}
	return db, nil
}

func query() {
	db, err := connect()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer db.Close()

	rows, err := db.Query("SELECT * FROM ping")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer rows.Close()

	var result []pingStruct

	for rows.Next() {
		var each = pingStruct{}
		var err = rows.Scan(&each.time, &each.deviceName, &each.ip, &each.deviceStatus)
		fmt.Println(err)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		result = append(result, each)
	}

	if err = rows.Err(); err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(result)
}

func insert(deviceName string, ip string, deviceStatus string, upTime string) {
	db, err := connect()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer db.Close()

	_, err = db.Exec(`INSERT INTO ping VALUES ($1, $2, $3, $4)`, deviceName, ip, deviceStatus, upTime)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("report from", deviceName, "inserted !")
}
