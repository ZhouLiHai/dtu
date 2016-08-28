package main

import (
	"analog"
	"arch"
	"communication"
	"config"
	"fmt"
	"lcd"
	"remote"
	"task"
	"time"
)

func main() {
	fmt.Println("New start.")
	arch.Init()
	config.Init("/opt/core.cfg")
	task.Init()
	lcd.Init("/opt/12")
	analog.Init()
	remote.Init(config.ReadReConfig())

	config.GetQD()
	config.GetContrast()
	config.GetGPIO()
	config.GetPress()
	config.GetDcCoe1()
	config.GetDcCoe2()
	config.GetRebackTime()
	config.GetRebackSwitch()
	config.GetBattHuoHuaSwitch()
	config.GetNetPort()
	config.GetTelnIp()
	config.GetSntpIp()

	communication.Start()

	for {
		time.Sleep(5 * 1e8)
		task.RunningTask()
		lcd.MenuDrawer()
		remote.RefreshSoe()
		analog.RefreshAnValue()
		analog.KalmanFilter()
		analog.Protect()
	}
}
