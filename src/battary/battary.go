package battary

import (
	"arch"
	"fmt"
	"sync"
	"time"
)

var lock *sync.Mutex = &sync.Mutex{}

func Start() {
	fmt.Println("电池活化投入", time.RFC1123)
	go start()
}

func End() {
	fmt.Println("电池活化退出", time.RFC1123)
	go end()
}

func start() {
	lock.Lock()
	arch.BattOn(1)
	time.Sleep(1e9)
	arch.BattOn(0)
	lock.Unlock()
}

func end() {
	lock.Lock()
	arch.BattOff(1)
	time.Sleep(1e9)
	arch.BattOff(0)
	lock.Unlock()
}
