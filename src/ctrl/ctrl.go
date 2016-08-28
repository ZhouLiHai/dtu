package ctrl

import (
	"arch"
	"fmt"
	"sync"
	"time"
)

var state int = 0

var lock *sync.Mutex = &sync.Mutex{}

//使用 go 语句将操作封装
func Select(id int, hw int) {
	fmt.Println("[信息]遥控预选", "路:", id, "动:", hw, "遥控状态", state, time.Now())
	go choice(id, hw)
}

func Action(id int) {
	fmt.Println("[信息]遥控执行", "路:", id, "遥控状态", state, time.Now())
	go action(id)
}

func Cancel(id int) {
	fmt.Println("[信息]遥控取消", "路:", id, "遥控状态", state, time.Now())
	go cancel(id)
}

func JustDo(id int, hw int) {
	fmt.Println("[信息]遥控速动", "路:", id, "动:", hw, "遥控状态", state, time.Now())
	go justdo(id, hw)
}

func choice(id int, hw int) bool {
	lock.Lock()
	defer lock.Unlock()
	if state != 0 {
		fmt.Println("遥控执行中 --", "已经预选:", state, time.Now())
		return false
	}
	arch.Select(id, hw)
	return true
}

func action(id int) bool {
	lock.Lock()
	defer lock.Unlock()
	if state != id {
		fmt.Println("遥控未预选 --", "预选状态:", state, "欲执行:", id, time.Now())
		return false
	}
	arch.Action(id)
	time.Sleep(1e9)
	arch.Cancel(id)
	return true
}
func cancel(id int) bool {
	lock.Lock()
	defer lock.Unlock()
	if state != id {
		fmt.Println("遥控未预选 --", "预选状态:", state, "欲取消:", id, time.Now())
		return false
	}
	arch.Cancel(id)
	return true
}

func justdo(id int, hw int) bool {
	lock.Lock()
	defer lock.Unlock()
	if state != 0 {
		fmt.Println("[速动]遥控执行中 --", "已经预选:", state, time.Now())
		return false
	}
	arch.Select(id, hw)
	time.Sleep(1e6)
	arch.Action(id)
	time.Sleep(1e9)
	arch.Cancel(id)
	return true
}
