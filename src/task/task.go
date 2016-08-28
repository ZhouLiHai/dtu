package task

import (
	"arch"
	"time"
)

//上电后硬件插口状态
var board_state [12]bool

func Init() {
	for i := 0; i < 11; i++ {
		board_state[i] = arch.GetBoardState(uint16(i))
	}
}

func RunningTask() {
	if CheckBoardState() {
		arch.SetLed(arch.LED_CHK, arch.LED_ON)
	} else {
		arch.SetLed(arch.LED_CHK, arch.LED_FF)
	}

	if RunLightTurn() {
		arch.SetLed(arch.LED_RUN, arch.LED_ON)
	} else {
		arch.SetLed(arch.LED_RUN, arch.LED_FF)
	}
}

func CheckBoardState() bool {
	for i := 0; i < 11; i++ {
		if board_state[i] != arch.GetBoardState(uint16(i)) {
			return true
		}
	}
	return false
}

func RunLightTurn() bool {
	if time.Now().Second()%2 == 0 {
		return true
	} else {
		return false
	}
}
