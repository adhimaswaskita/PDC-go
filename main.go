package main

import (
	"fmt"
	"time"

	"github.com/adhimaswaskita/ping-data/src"
	_ "github.com/lib/pq"
	"github.com/sparrc/go-ping"
)

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
					src.Insert(deviceName, ip, deviceStatus, upTime)
				} else {
					deviceStatus = "Tidak aktif"
					src.Insert(deviceName, ip, deviceStatus, upTime)
				}
			}
		}
	}
}
