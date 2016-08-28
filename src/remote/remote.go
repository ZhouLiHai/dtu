package remote

import (
	"arch"
	"fmt"
	"mtype"
	"sync"
)

//Soe 记录体结构
type record struct {
	Id   uint16
	Sw   uint16
	Type uint16
	Tv   int64
}

const TOUCH = 1
const UNTOUCH = 0

const SINGLE = 1
const DOUBLE = 2

//全局遥信配置
var configs []mtype.ReConfig

//全局 Soe 在 FPGA 中的头指针
var soeHead int32

//全局 Soe 记录指针,结构体以及锁[循环记录]
var recordHead int32
var records [255]record
var lock *sync.Mutex = &sync.Mutex{}

func Init(cfg []mtype.ReConfig) {
	configs = cfg
	soeHead = arch.GetSoeHead()
	recordHead = 0
}

func RefreshSoe() {
	if soeHead == arch.GetSoeHead() {
		return
	} else {
		go refresh()
	}
}

func Insert(id uint16, sw uint16, tp uint16, tv int64) {
	lock.Lock()
	defer lock.Unlock()

	records[recordHead].Id = id
	records[recordHead].Sw = sw
	records[recordHead].Type = tp
	records[recordHead].Tv = tv
	fmt.Println("[信息]添加遥信", records[recordHead])
}

func refresh() {
	id, sw, tv := arch.GetSoeBody(soeHead)
	for i := 0; i < len(configs); i++ {
		if int32(id) == configs[i].L0 && configs[i].Type == 1 {
			Insert(id, sw, 1, tv)
		}
		if (int32(id) == configs[i].L0 || int32(id) == configs[i].L1) && configs[i].Type == 2 {
			Insert(id, arch.GetSignal(uint16(configs[i].L0))<<1+arch.GetSignal(uint16(configs[i].L1)), 2, tv)
		}
	}
	soeHead = (soeHead + 0x04) % 0x400
}

func BuildSingleSoe() []mtype.SingleSoe {
	soe := make([]mtype.SingleSoe, 0, 128)
	for i := 0; i < len(configs); i++ {
		if configs[i].Type == SINGLE {
			soe = append(soe, mtype.SingleSoe{Id: int(configs[i].Id), State: int(arch.GetSignal(uint16(configs[i].L0)))})
		}
	}
	fmt.Println(soe)
	return soe
}

func BuildDoubleSoe() []mtype.DoubleSoe {
	soe := make([]mtype.DoubleSoe, 0, 64)
	for i := 0; i < len(configs); i++ {
		if configs[i].Type == DOUBLE {
			soe = append(soe, mtype.DoubleSoe{Id: int(configs[i].Id),
				State1: int(arch.GetSignal(uint16(configs[i].L0))),
				State2: int(arch.GetSignal(uint16(configs[i].L1)))})
		}
	}
	return soe
}
