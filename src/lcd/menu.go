package lcd

import (
	"arch"
	"container/list"
	"ctrl"
	"fmt"
	"net"
	"time"
)

//菜单函数站,使用出栈和进栈的方式完成菜单的上一级和下一级切换
var menuStack = list.New()

func Init(path string) {
	readHZBuf(path)
	menuStack.PushBack(level0_main)
}

//菜单绘制函数完成按键采集和图形绘制
func MenuDrawer() {
	if key := GetKey(); key != NONE_VALUE {
		go drawMenu(int(key))
	}
}

func drawMenu(key int) {
	menu := menuStack.Back()
	res, function := (menu.Value).(func(int) (int, interface{}))(key)
	if key == ESC_VALUE && menuStack.Len() != 1 {
		menuStack.Remove(menu)
	}
	switch res {
	case 1:
		menuStack.PushBack(function)
	}

	textShow(fmt.Sprintf("UTC:%s", time.Now())[:23], 12, 2, false)
	lcdHLine(0, 15, 160)
	lcdHLine(0, 145, 160)
	draw()
	clearRect(0, 160, 0, 160)
}

var level0_main_offset uint8 = 1

func level0_main(key int) (int, interface{}) {
	textShow("配网自动化系统终端", 30, 20, false)
	textShow("1.系统综合信息", 42, 60, level0_main_offset&0x01 != 0)
	textShow("2.参数配置", 42, 76, level0_main_offset&0x02 != 0)
	textShow("3.系统调试选项", 42, 92, level0_main_offset&0x04 != 0)

	switch key {
	case ACQ_VALUE:
		switch level0_main_offset {
		case 0x01:
			return 1, level1_1
		case 0x02:
			return 1, level1_2
		case 0x04:
			return 1, level1_3

		}
	case UPS_VALUE:
		if level0_main_offset > 0x01 {
			level0_main_offset >>= 1
		}
	case DOW_VALUE:
		if level0_main_offset < 0x04 {
			level0_main_offset <<= 1
		}
	}
	return 0, nil
}

var level1_1_offset uint8 = 1

func level1_1(key int) (int, interface{}) {
	textShow("配网自动化系统终端", 30, 20, false)
	textShow("1.遥测", 42, 60, level1_1_offset&0x01 != 0)
	textShow("2.系统信息", 42, 76, level1_1_offset&0x02 != 0)
	textShow("3.状态自检", 42, 92, level1_1_offset&0x04 != 0)

	switch key {
	case ACQ_VALUE:
		switch level1_1_offset {
		case 0x01:
			return 1, level1_2_1
		case 0x02:
			return 1, level1_2_2
		case 0x04:
			return 1, level1_2_3
		}
	case UPS_VALUE:
		if level1_1_offset > 0x01 {
			level1_1_offset >>= 1
		}
	case DOW_VALUE:
		if level1_1_offset < 0x04 {
			level1_1_offset <<= 1
		}
	}
	return 0, nil
}

var level1_2_1_offset uint8 = 1

func level1_2_1(key int) (int, interface{}) {
	return 0, nil
}

func level1_2_2(key int) (int, interface{}) {
	textShow("软件版本:V0.56", 12, 30, false)
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println("[异常]获取网络地址出错", err)
	}
	textShow(fmt.Sprintf("网络 : %s", addrs[1]), 12, 45, false)
	textShow(fmt.Sprintf("网络 : %s", addrs[2]), 12, 60, false)
	textShow("直流量:", 12, 75, false)
	textShow("直流量:", 12, 90, false)

	switch key {
	}
	return 0, nil
}

func level1_2_3(key int) (int, interface{}) {
	textShow("主板-电源:正常", 25, 45, false)
	textShow("遥测:", 25, 60, false)
	for i := 8; i < 11; i++ {
		if arch.GetBoardState(uint16(i)) == true {
			textShow("正常", 55+(i-8)*30, 60, false)
		} else {
			textShow("异常", 55+(i-8)*30, 60, true)
		}
	}
	textShow("遥信:", 25, 75, false)
	for i := 0; i < 6; i++ {
		if arch.GetBoardState(uint16(i)) == true {
			textShow("正常", 55+(i%3)*30, 75+(i/3)*15, false)
		} else {
			textShow("异常", 55+(i%3)*30, 75+(i/3)*15, true)
		}
	}
	return 0, nil
}

var level1_2_offset uint8 = 1

func level1_2(key int) (int, interface{}) {
	textShow("配网自动化系统终端", 30, 20, false)
	textShow("1.保护定值", 42, 60, level1_2_offset&0x01 != 0)

	switch key {
	case ACQ_VALUE:
		switch level1_2_offset {
		case 0x01:
			return 10, nil
		}
	}
	return 0, nil
}

var level1_3_offset uint8 = 1

func level1_3(key int) (int, interface{}) {
	textShow("配网自动化系统终端", 30, 20, false)
	textShow("1.遥控测试", 42, 60, level1_3_offset&0x01 != 0)
	textShow("2.全遥信状态", 42, 76, level1_3_offset&0x02 != 0)

	switch key {
	case ACQ_VALUE:
		switch level1_3_offset {
		case 0x01:
			return 1, level1_3_1
		case 0x02:
			return 1, level1_3_2
		}
	case UPS_VALUE:
		if level1_3_offset > 0x01 {
			level1_3_offset >>= 1
		}
	case DOW_VALUE:
		if level1_3_offset < 0x04 {
			level1_3_offset <<= 1
		}
	}
	return 0, nil
}

var level1_3_1offset uint16 = 1

func level1_3_1(key int) (int, interface{}) {
	textShow("遥控功能测试:", 45, 20, false)
	textShow("第一路分", 35, 44, level1_3_1offset&0x01 != 0)
	textShow("合", 95, 44, level1_3_1offset&0x02 != 0)
	textShow("第二路分", 35, 60, level1_3_1offset&0x04 != 0)
	textShow("合", 95, 60, level1_3_1offset&0x08 != 0)
	textShow("第三路分", 35, 76, level1_3_1offset&0x10 != 0)
	textShow("合", 95, 76, level1_3_1offset&0x20 != 0)
	textShow("第四路分", 35, 92, level1_3_1offset&0x40 != 0)
	textShow("合", 95, 92, level1_3_1offset&0x80 != 0)
	textShow("第五路分", 35, 108, level1_3_1offset&0x100 != 0)
	textShow("合", 95, 108, level1_3_1offset&0x200 != 0)
	textShow("第六路分", 35, 124, level1_3_1offset&0x400 != 0)
	textShow("合", 95, 124, level1_3_1offset&0x800 != 0)

	switch key {
	case ACQ_VALUE:
		var line = 0
		for level1_3_1offset != 0x01 {
			line++
			level1_3_1offset >>= 1
		}
		fmt.Println("[信息]液晶遥控动作 路:", line/2, "动:", line%2)
		ctrl.JustDo(line/2, line%2)
	case UPS_VALUE:
		if level1_3_1offset > 0x01 {
			level1_3_1offset >>= 1
		}
	case DOW_VALUE:
		if level1_3_1offset < 0x400 {
			level1_3_1offset <<= 1
		}
	}
	return 0, nil
}

func level1_3_2(key int) (int, interface{}) {
	for i := 0; i < 80; i++ {
		if arch.GetSignal(uint16(i)) == 1 {
			textShow("O", 26+(i/10)*15, 30+(i%10)*10, false)
		} else {
			textShow("X", 26+(i/10)*15, 30+(i%10)*10, false)
		}
	}
	return 0, nil
}
