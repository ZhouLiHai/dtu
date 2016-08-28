package analog

import (
	"arch"
	"ctrl"
	"fmt"
	"mtype"
	"remote"
	"time"
)

//全局遥测配置
var procfg []mtype.ProConfig

func Protect() {
	// for i := 0; i < len(procfg); i++ {
	for i := 0; i < 1; i++ {
		go protect(procfg[i], i)
	}
}

func maxI(a, b, c float32) float32 {
	var max float32
	if max < a {
		max = a
	}
	if max < b {
		max = b
	}
	if max < c {
		max = c
	}
	return max
}

//quickCheck 函数使用的状态变量
var quickFlter [12][3]int
var quickTimer [12]mtype.ProFilter

func quickCheck(now float32, cfg mtype.ProConfig, line int) int {
	res := 0
	//I 段判断
	switch quickFlter[line][0] {
	case 0:
		if now > float32(cfg.V1) {
			quickFlter[line][0] = 1
			quickTimer[line].T1 = time.Now()
		}
	case 1:
		if now > float32(cfg.V1) {
			if int32(time.Now().Sub(quickTimer[line].T1)/1e6) > cfg.T1 {
				res = 1
				quickFlter[line][0] = 0
			}
		} else {
			quickFlter[line][0] = 0
		}
	}
	//II 段判断
	switch quickFlter[line][1] {
	case 0:
		if now > float32(cfg.V2) {
			quickFlter[line][1] = 1
			quickTimer[line].T2 = time.Now()
		}
	case 1:
		if now > float32(cfg.V2) {
			if int32(time.Now().Sub(quickTimer[line].T2)/1e6) > cfg.T2 {
				if res == 0 {
					res = 2
				}
				quickFlter[line][1] = 0
			}
		} else {
			quickFlter[line][1] = 0
		}
	}
	//负荷判断
	switch quickFlter[line][2] {
	case 0:
		if now > float32(cfg.V3) {
			quickFlter[line][2] = 1
			quickTimer[line].T3 = time.Now()
		}
	case 1:
		if now > float32(cfg.V3) {
			if int32(time.Now().Sub(quickTimer[line].T3)/1e6) > cfg.T3 {
				if res == 0 {
					res = 3
				}
				quickFlter[line][2] = 0
			}
		} else {
			quickFlter[line][2] = 0
		}
	}
	if res != 0 {
		fmt.Println(res)
	}
	return res
}

func proDuty(res int, cfg mtype.ProConfig, line int) {
	remote.Insert(uint16(cfg.Id), remote.TOUCH, remote.SINGLE, time.Now().Unix())
	switch res {
	case 1:
		arch.SetLed(arch.LED_WAR, arch.LED_ON)
	case 2:
		arch.SetLed(arch.LED_TER, arch.LED_ON)
		ctrl.JustDo(int(cfg.Ctrl), int(cfg.How))
	case 3:
		arch.SetLed(arch.LED_WAR, arch.LED_ON)
		arch.SetLed(arch.LED_TER, arch.LED_ON)
		ctrl.JustDo(int(cfg.Ctrl), int(cfg.How))
	}
}

func checkSwitch(res int, cfg mtype.ProConfig) int {
	state := 0
	if res == 1 && cfg.S1 == 1 {
		state = 1
	}
	if res == 2 && cfg.S2 == 1 {
		state = 1
	}
	if res == 3 && cfg.S3 == 1 {
		state = 1
	}
	return state
}

var proFlter [12]int

const PRO_PEACE = 0
const PRO_WAIT_REBACK = 1
const PRO_WAIT_FCOVER = 2
const PRO_BEFORE_END = 3
const PRO_END = 4

func protect(cfg mtype.ProConfig, line int) {
	now := maxI(GetLine(line))
	switch proFlter[line] {
	case PRO_PEACE:
		if res := quickCheck(now, cfg, line); res != 0 {
			if checkSwitch(res, cfg) == 1 {
				proDuty(res, cfg, line)
				if cfg.Sr == 1 {
					proFlter[line] = PRO_WAIT_REBACK
				} else {
					proFlter[line] = PRO_END
				}
				quickTimer[line].Tr = time.Now()
				quickTimer[line].Tf = time.Now()
			}
		}
	case PRO_WAIT_REBACK:
		if int32(time.Now().Sub(quickTimer[line].Tr)/1e6) > cfg.Tr {
			proFlter[line] = PRO_WAIT_FCOVER
		}
	case PRO_WAIT_FCOVER:
		if cfg.Sf == 0 {
			proFlter[line] = PRO_END
		}
		if int32(time.Now().Sub(quickTimer[line].Tf)/1e6) < cfg.Tf {
			if res := quickCheck(now, cfg, line); res != 0 {
				proDuty(res, cfg, line)
				proFlter[line] = PRO_BEFORE_END
			}
		} else {
			proFlter[line] = PRO_PEACE
		}
	case PRO_BEFORE_END:
		remote.Insert(uint16(cfg.Id), remote.UNTOUCH, remote.SINGLE, time.Now().Unix())
		proFlter[line] = PRO_END
	case PRO_END:
	}
}
