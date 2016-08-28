package config

/*
#cgo CFLAGS: -I/Volumes/arm-x/libconfig/include
#cgo LDFLAGS: -L/Volumes/arm-x/libconfig/lib -lconfig
#include "libconfig.h"
*/
import "C"
import (
	"fmt"
	"mtype"
	"os"
)

//全局配置文件结构体
var Gconfig C.config_t

func Init(path string) {
	if err := C.config_read_file(&Gconfig, C.CString(path)); err == 0 {
		//err_line := C.config_error_text(unsafe.Pointer(&Gconfig))
		fmt.Println("[错误]配置文件读取错误...", Gconfig.error_line, C.GoString(Gconfig.error_text))
		os.Exit(1)
	}
}

func GetQD() C.int {
	misc := C.config_lookup(&Gconfig, C.CString("go.misc"))

	var ret C.int
	if err := C.config_setting_lookup_int(misc, C.CString("delay"), &ret); err == 0 {
		fmt.Println("[警告]遥信去抖参数读取失败...")
	}
	fmt.Println("[配置]遥信去抖时间", ret, "ms")
	return ret
}

func GetContrast() C.int {
	misc := C.config_lookup(&Gconfig, C.CString("go.misc"))

	var ret C.int
	if err := C.config_setting_lookup_int(misc, C.CString("contrast"), &ret); err == 0 {
		fmt.Println("[警告]液晶对比度参数读取失败...")
	}
	fmt.Println("[配置]液晶对比度", ret, "")
	return ret
}

func GetGPIO() []string {
	init := C.config_lookup(&Gconfig, C.CString("go.misc.init"))

	var ret []string
	for i := 0; i < int(C.config_setting_length(init)); i++ {
		ret = append(ret, C.GoString(C.config_setting_get_string_elem(init, C.int(i))))
	}

	fmt.Println("[配置]GPIO读取数量", C.config_setting_length(init), "条")
	return ret
}

func GetPress() C.int {
	misc := C.config_lookup(&Gconfig, C.CString("go.misc"))

	var ret C.int
	if err := C.config_setting_lookup_int(misc, C.CString("pressure"), &ret); err == 0 {
		fmt.Println("[警告]有压鉴别值参数读取失败...")
	}
	fmt.Println("[配置]有压鉴别值", ret, "V")
	return ret
}

func GetActivate() C.int {
	misc := C.config_lookup(&Gconfig, C.CString("go.misc"))

	var ret C.int
	if err := C.config_setting_lookup_int(misc, C.CString("activate"), &ret); err == 0 {
		fmt.Println("[警告]活化周期参数读取失败...")
	}
	fmt.Println("[配置]活化周期", ret, "秒")
	return ret
}

func GetOnBattery() C.int {
	misc := C.config_lookup(&Gconfig, C.CString("go.misc"))

	var ret C.int
	if err := C.config_setting_lookup_int(misc, C.CString("onbattery"), &ret); err == 0 {
		fmt.Println("[警告]活化时间参数读取失败...")
	}
	fmt.Println("[配置]活化时间", ret, "V")
	return ret
}

func GetDcCoe1() C.double {
	var ret C.double
	if err := C.config_lookup_float(&Gconfig, C.CString("go.misc.dcCoe1"), &ret); err == 0 {
		fmt.Println("[警告]直流量系数1参数读取失败...")
	}
	fmt.Println("[配置]直流量系数1", ret, "")
	return ret
}

func GetDcCoe2() C.double {
	var ret C.double
	if err := C.config_lookup_float(&Gconfig, C.CString("go.misc.dcCoe2"), &ret); err == 0 {
		fmt.Println("[警告]直流量系数2参数读取失败...")
	}
	fmt.Println("[配置]直流量系数2", ret, "")
	return ret
}

func GetRebackTime() C.int {
	misc := C.config_lookup(&Gconfig, C.CString("go.misc"))

	var ret C.int
	if err := C.config_setting_lookup_int(misc, C.CString("rebacktime"), &ret); err == 0 {
		fmt.Println("[警告]复归时间参数读取失败...")
	}
	fmt.Println("[配置]复归时间", ret, "ms")
	return ret
}

func GetRebackSwitch() C.int {
	misc := C.config_lookup(&Gconfig, C.CString("go.misc"))

	var ret C.int
	if err := C.config_setting_lookup_int(misc, C.CString("reback"), &ret); err == 0 {
		fmt.Println("[警告]复归开关参数读取失败...")
	}
	fmt.Println("[配置]复归开关", ret, "")
	return ret
}

func GetBattHuoHuaSwitch() C.int {
	misc := C.config_lookup(&Gconfig, C.CString("go.misc"))

	var ret C.int
	if err := C.config_setting_lookup_int(misc, C.CString("loop"), &ret); err == 0 {
		fmt.Println("[警告]活化开关参数读取失败...")
	}
	fmt.Println("[配置]活化开关", ret, "")
	return ret
}

func GetIECAddr() uint16 {
	var ret C.int
	if err := C.config_lookup_int(&Gconfig, C.CString("go.protocol.Addr"), &ret); err == 0 {
		fmt.Println("[警告]规约主站地址参数读取失败...")
	}
	fmt.Println("[配置]规约主站地址", ret, "")
	return uint16(ret)
}

func GetIECMasterAddr() uint16 {
	var ret C.int
	if err := C.config_lookup_int(&Gconfig, C.CString("go.protocol.MasterAddr"), &ret); err == 0 {
		fmt.Println("[警告]规约公共地址参数读取失败...")
	}
	fmt.Println("[配置]规约公共地址", ret, "")
	return uint16(ret)
}

func GetTranCauseL() uint16 {
	var ret C.int
	if err := C.config_lookup_int(&Gconfig, C.CString("go.protocol.iTranCauseL"), &ret); err == 0 {
		fmt.Println("[警告]iTranCauseL参数读取失败...")
	}
	fmt.Println("[配置]iTranCauseL", ret, "")
	return uint16(ret)
}
func GetInfoAddrL() uint16 {
	var ret C.int
	if err := C.config_lookup_int(&Gconfig, C.CString("go.protocol.iInfoAddrL"), &ret); err == 0 {
		fmt.Println("[警告]iInfoAddrL参数读取失败...")
	}
	fmt.Println("[配置]iInfoAddrL", ret, "")
	return uint16(ret)
}

func GetCommonAddrL() uint16 {
	var ret C.int
	if err := C.config_lookup_int(&Gconfig, C.CString("go.protocol.iCommonAddrL"), &ret); err == 0 {
		fmt.Println("[警告]iCommonAddrL参数读取失败...")
	}
	fmt.Println("[配置]iCommonAddrL", ret, "")
	return uint16(ret)
}

func GetNetPort() C.int {
	var ret C.int
	if err := C.config_lookup_int(&Gconfig, C.CString("go.net.port"), &ret); err == 0 {
		fmt.Println("[警告]监听端口参数读取失败...")
	}
	fmt.Println("[配置]监听端口", ret, "")
	return ret
}

func GetTelnIp() string {
	var ret *C.char
	if err := C.config_lookup_string(&Gconfig, C.CString("go.net.telnIP"), &ret); err == 0 {
		fmt.Println("[警告]备自投设备IP参数读取失败...")
	}
	fmt.Println("[配置]备自投设备IP", C.GoString(ret), "")
	return C.GoString(ret)
}

func GetSntpIp() string {
	var ret *C.char
	if err := C.config_lookup_string(&Gconfig, C.CString("go.net.SntpIP"), &ret); err == 0 {
		fmt.Println("[警告]备自投设备IP参数读取失败...")
	}
	fmt.Println("[配置]备自投设备IP", C.GoString(ret), "")
	return C.GoString(ret)
}

func ReadAnConfig() []mtype.AnConfig {
	//结构体内部定义为 C 语言类型,方便 C 函数库的交互(局部使用)
	type AC struct {
		Id      C.int
		Switchs C.int
		Shift   C.int
		Form    C.double
		Larg    C.double
		Thre    C.double
		Up      C.double
		Down    C.double
	}
	var res []mtype.AnConfig

	pack := C.config_lookup(&Gconfig, C.CString("go.analog"))

	for i := 0; i < int(C.config_setting_length(pack)); i++ {
		hooker := new(AC)

		r := C.config_setting_get_elem(pack, C.uint(i))
		C.config_setting_lookup_int(r, C.CString("n"), &hooker.Id)
		C.config_setting_lookup_int(r, C.CString("s"), &hooker.Switchs)
		C.config_setting_lookup_int(r, C.CString("l"), &hooker.Shift)
		C.config_setting_lookup_float(r, C.CString("f"), &hooker.Form)
		C.config_setting_lookup_float(r, C.CString("la"), &hooker.Larg)
		C.config_setting_lookup_float(r, C.CString("th"), &hooker.Thre)
		C.config_setting_lookup_float(r, C.CString("down"), &hooker.Down)
		C.config_setting_lookup_float(r, C.CString("up"), &hooker.Up)

		transformer := mtype.AnConfig{int32(hooker.Id), bool(int(hooker.Switchs) == 1),
			int32(hooker.Shift), float32(hooker.Form), float32(hooker.Larg),
			float32(hooker.Thre), float32(hooker.Down), float32(hooker.Up)}
		res = append(res, transformer)
	}

	for i := 0; i < len(res); i++ {
		fmt.Println(res[i])
	}
	return res
}

func ReadReConfig() []mtype.ReConfig {
	//结构体内部定义为 C 语言类型,方便 C 函数库的交互(局部使用)
	type RC struct {
		Id   C.int
		Type C.int
		L0   C.int
		L1   C.int
	}

	var res []mtype.ReConfig

	pack := C.config_lookup(&Gconfig, C.CString("go.remote"))

	for i := 0; i < int(C.config_setting_length(pack)); i++ {
		hooker := new(RC)

		r := C.config_setting_get_elem(pack, C.uint(i))
		C.config_setting_lookup_int(r, C.CString("n"), &hooker.Id)
		C.config_setting_lookup_int(r, C.CString("t"), &hooker.Type)
		C.config_setting_lookup_int(r, C.CString("l0"), &hooker.L0)
		C.config_setting_lookup_int(r, C.CString("l1"), &hooker.L1)

		transformer := mtype.ReConfig{int32(hooker.Id), int32(hooker.Type), int32(hooker.L0),
			int32(hooker.L1)}

		res = append(res, transformer)
	}

	for i := 0; i < len(res); i++ {
		fmt.Println(res[i])
	}
	return res
}

func ReadProConfig() []mtype.ProConfig {
	//结构体内部定义为 C 语言类型,方便 C 函数库的交互(局部使用)
	type PC struct {
		Id   C.int
		S1   C.int
		V1   C.int
		T1   C.int
		S2   C.int
		V2   C.int
		T2   C.int
		S3   C.int
		V3   C.int
		T3   C.int
		Sr   C.int
		Tr   C.int
		Sf   C.int
		Tf   C.int
		Ctrl C.int
		How  C.int
	}

	var res []mtype.ProConfig

	pack := C.config_lookup(&Gconfig, C.CString("go.protect"))

	for i := 0; i < int(C.config_setting_length(pack)); i++ {
		hooker := new(PC)

		r := C.config_setting_get_elem(pack, C.uint(i))
		C.config_setting_lookup_int(r, C.CString("n"), &hooker.Id)
		C.config_setting_lookup_int(r, C.CString("s1"), &hooker.S1)
		C.config_setting_lookup_int(r, C.CString("v1"), &hooker.V1)
		C.config_setting_lookup_int(r, C.CString("t1"), &hooker.T1)
		C.config_setting_lookup_int(r, C.CString("s2"), &hooker.S2)
		C.config_setting_lookup_int(r, C.CString("v2"), &hooker.V2)
		C.config_setting_lookup_int(r, C.CString("t2"), &hooker.T2)
		C.config_setting_lookup_int(r, C.CString("s3"), &hooker.S3)
		C.config_setting_lookup_int(r, C.CString("v3"), &hooker.V3)
		C.config_setting_lookup_int(r, C.CString("t3"), &hooker.T3)
		C.config_setting_lookup_int(r, C.CString("sr"), &hooker.Sr)
		C.config_setting_lookup_int(r, C.CString("tr"), &hooker.Tr)
		C.config_setting_lookup_int(r, C.CString("sf"), &hooker.Sf)
		C.config_setting_lookup_int(r, C.CString("tf"), &hooker.Tf)
		C.config_setting_lookup_int(r, C.CString("c"), &hooker.Ctrl)
		C.config_setting_lookup_int(r, C.CString("w"), &hooker.How)

		transformer := mtype.ProConfig{int32(hooker.Id),
			int32(hooker.S1), int32(hooker.V1 / 1000), int32(hooker.T1),
			int32(hooker.S2), int32(hooker.V2 / 1000), int32(hooker.T2),
			int32(hooker.S3), int32(hooker.V3 / 1000), int32(hooker.T3),
			int32(hooker.Sr), int32(hooker.Tr),
			int32(hooker.Sf), int32(hooker.Tf * 1000),
			int32(hooker.Ctrl), int32(hooker.How)}

		res = append(res, transformer)
	}

	for i := 0; i < len(res); i++ {
		fmt.Println(res[i])
	}
	return res
}
