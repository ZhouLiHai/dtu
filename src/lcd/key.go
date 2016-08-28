package lcd

import (
	"arch"
)

const UPS_VALUE = 0x1F7
const DOW_VALUE = 0x1EF
const LFT_VALUE = 0x1DF
const RGT_VALUE = 0x1FE
const ACQ_VALUE = 0x1FD
const ESC_VALUE = 0x0FF
const RST_VALUE = 0x1FB
const ADD_VALUE = 0x1BF
const SUB_VALUE = 0x17F

const NONE_VALUE = 0x1FF

//全局按键信号,以及滤波过滤器
var Key uint16 = NONE_VALUE
var filter = 0
var engine = 0

func GetKey() int32 {
	//实体按键部分
	if Key == arch.GetKeySignal() {
		filter++
		if filter > 40 {
			filter = 0
			return int32(Key)
		}
	} else {
		filter = 0
		Key = arch.GetKeySignal()
	}
	//液晶驱动器部分,每秒中给一个脉冲信号,驱动液晶更新
	engine++
	if engine > 200 {
		engine = 0
		//返回一个无效信号,驱动液晶显示
		return 0xFFFF
	}
	return NONE_VALUE
}
