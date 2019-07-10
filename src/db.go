package src

import (
	"database/sql"
	"fmt"
)

type pingStruct struct {
	time         string
	deviceStatus string
	ip           string
	deviceName   string
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

func Insert(deviceName string, ip string, deviceStatus string, upTime string) {
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
