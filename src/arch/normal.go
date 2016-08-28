package arch

/*
#cgo linux CFLAGS: -mno-thumb-interwork
#cgo linux LDFLAGS: -lrt

#include <fcntl.h>
#include <unistd.h>
#include <sys/mman.h>

unsigned short * MFileOpen() {
	int fd = open("/dev/mem", O_RDWR|O_SYNC);
	return (unsigned short *)mmap(0, 4096, PROT_READ | PROT_WRITE, MAP_SHARED, fd, 0x14000000);
}
*/
import "C"
import (
	"fmt"
	"os"
	"syscall"
	"time"
	"unsafe"
)

//内存映射地址范围
const FPGA_SIZE = 2048

//内存映射基地址
var Ghead *[FPGA_SIZE]uint16

//全局 led 灯状态
var led_state uint16 = 0x00

func openMem() *[FPGA_SIZE]uint16 {
	ptr, err := C.MFileOpen()
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return (*[FPGA_SIZE]uint16)(unsafe.Pointer(ptr))
}

func Init() {
	Ghead = openMem()
	if Ghead == nil {
		fmt.Println("初始化 FPGA 地址失败...")
		os.Exit(1)
	}
}

func GetBoardState(id uint16) bool {
	if id < 0 || id > 10 {
		return false
	}
	/*
	 *  0-5 :遥信遥控板
	 *  6   :电源板
	 *  7   :按键板
	 *  8-10:遥测板
	 */
	if id < 8 {
		var state uint16 = Ghead[0x434]
		return (state & (1 << id)) != 0
	} else {
		switch id {
		case 8:
			return Ghead[0x500] == 0xEB90
		case 9:
			return Ghead[0x600] == 0xEB90
		case 10:
			return Ghead[0x700] == 0xEB90
		}
	}
	return false
}

func GetAnState(i int) bool {
	return GetBoardState(uint16(i + 8))
}

func Select(id int, how int) {
	var cmd_open uint16
	var cmd_shut uint16
	if id < 0 || id >= 12 {
		return
	}
	if id%2 == 0 {
		cmd_open = 0x1c00
		cmd_shut = 0x1500
	} else {
		cmd_open = 0x2c00
		cmd_shut = 0x2500
	}
	switch how {
	case 0:
		Ghead[id/2+0x420] = cmd_open

	case 1:
		Ghead[id/2+0x420] = cmd_shut
	}
}

func Action(id int) {
	var cmd uint16
	if id < 0 || id >= 12 {
		return
	}
	if id%2 == 0 {
		cmd = 0x1a00
	} else {
		cmd = 0x2a00
	}
	Ghead[id/2+0x420] = cmd
}

func Cancel(id int) {
	if id < 0 || id >= 12 {
		return
	}
	Ghead[id/2+0x420] = 0x1F00
}

const LED_ON = 1
const LED_FF = 0

const LED_RUN = 5
const LED_CHK = 4
const LED_WAR = 3
const LED_TER = 2
const LED_BK1 = 1
const LED_BK2 = 0

func SetLed(id int, s int) {
	offset := uint16(id + 2)

	if id >= 0 && id < 7 {
		if 0 == s {
			led_state &= ^(1 << offset)
		}
		if 1 == s {
			led_state |= (1 << offset)
		}
	}
	Ghead[0x41C] = led_state
}

func BattOn(hw int) {
	if hw == 1 {
		Ghead[0x43A] = 1
	}
	if hw == 0 {
		Ghead[0x43A] = 0
	}
}

func BattOff(hw int) {
	if hw == 1 {
		Ghead[0x439] = 1
	}
	if hw == 0 {
		Ghead[0x439] = 0
	}
}

func DevTime() int64 {
	var tv int64 = 0
	tv = tv + int64(Ghead[0x446])
	tv = tv << 16
	tv = tv + int64(Ghead[0x447])
	tv = tv << 16
	tv = tv + int64(Ghead[0x448])
	return tv
}

func GetKeySignal() uint16 {
	return Ghead[0x41A]
}

func TransTime(k uint16) (int, int) {
	high := (k >> 8)
	low := k & 0x00ff
	return int((high/16)*10 + high%16), int((low/16)*10 + low%16)
}

func SetCpuTime() {
	year, mon := TransTime(Ghead[0x44A])
	year += 2000
	day, hour := TransTime(Ghead[0x44B])
	mins, sec := TransTime(Ghead[0x44C])
	unix_time := time.Date(year, time.Month(mon), day, hour, mins, sec, 0, time.Local).Unix()
	time := syscall.NsecToTimeval(unix_time * 1e9)
	if err := syscall.Settimeofday(&time); err != nil {
		fmt.Println(err)
	}
}

func SetFpgaTime(unix_time int64) {
	ttime := syscall.NsecToTimeval(unix_time)
	if err := syscall.Settimeofday(&ttime); err != nil {
		fmt.Println(err)
	}
	Ghead[0x44E] = 0x0000 + uint16(time.Now().Second())
	Ghead[0x44E] = 0x0100 + uint16(time.Now().Minute())
	Ghead[0x44E] = 0x0200 + uint16(time.Now().Hour())
	Ghead[0x44E] = 0x0400 + uint16(time.Now().Day())
	Ghead[0x44E] = 0x0500 + uint16(time.Now().Month())
	Ghead[0x44E] = 0x0600 + uint16(time.Now().Year())

}

func ReadAnCoe() [3][24]float32 {
	var Coes [3][24]float32

	fmt.Println("[信息]遥测系数读取>>>>>>>>>>>>>>>>")

	for j := 0; j < 3; j++ {
		if GetAnState(j) != true {
			continue
		}
		fmt.Println("[信息]第一块遥测板:")

		for i := 0; i < 24; i++ {
			if i < 4 {
				Coes[j][i] = 220.0 / float32(Ghead[0x500+j*0x100+i+2])
				fmt.Printf("|%8.6f|", Coes[j][i])
			}
			if i >= 4 && i < 16 {
				Coes[j][i] = 5.000 / float32(Ghead[0x500+j*0x100+i+2])
				fmt.Printf("|%8.6f|", Coes[j][i])
			}
			if i >= 16 {
				Coes[j][i] = float32(Ghead[0x500+j*0x100+i+2])
				fmt.Printf("|%5.0f|", Coes[j][i])
			}
			if i == 3 || i == 15 || i == 23 {
				fmt.Println()
			}
		}
	}
	fmt.Println("[信息]系数读取结束<<<<<<<<<<<<<<<<")

	return Coes
}

func ReadAnValue(Coes [3][24]float32) [3][24]float32 {
	var AnValue [3][24]float32

	for i := 0; i < 3; i++ {
		if GetAnState(i) != true {
			continue
		}

		for j := 0; j < 16; j++ {
			//第一块遥测板基地址0x500+板间偏移i*0x100+板内基地址26+板内偏移
			AnValue[i][j] = Coes[i][j] * float32(Ghead[0x500+(i*0x100)+26+j])
		}
		var offset int = 42
		for j := 16; j < 24; j++ {
			temp := int32(Ghead[0x500+(i*0x100)+26+offset+1])<<16 + int32(Ghead[0x500+(i*0x100)+26+offset])
			AnValue[i][j] = float32(temp) / Coes[i][j]
			offset += 2
		}
	}
	return AnValue
}

func GetSoeHead() int32 {
	return int32(Ghead[0x400])
}

func GetSoeBody(index int32) (id uint16, sw uint16, tv int64) {
	id = Ghead[index] & 0x7FFF
	sw = Ghead[index] >> 15
	tv = int64(Ghead[index+1])<<32 + int64(Ghead[index+2])<<16 + int64(Ghead[index+3])

	return id, sw, tv
}

func GetSignal(index uint16) uint16 {
	if Ghead[0x402+index/12]&(1<<(index%12)) == (1 << (index % 12)) {
		return 1
	} else {
		return 0
	}
}
