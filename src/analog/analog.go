package analog

import (
	"arch"
	"config"
	// "fmt"
	"math"
	"mtype"
)

//全局遥测系数
var GolobalAnCoe [3][24]float32

//全局遥测值[未滤波]
var GolobalAnVal [3][24]float32

//全局遥测值[滤波后]
var GolobalAnBal [3][24]float32

//全局遥测配置
var ancfg []mtype.AnConfig

func Init() {
	//初始化系数
	RefreshAnCoef()
	//初始化遥测参数
	ancfg = config.ReadAnConfig()
	//初始化保护参数
	procfg = config.ReadProConfig()
	//初始化卡尔曼滤波函数的估计偏差
	for i := 0; i < 3; i++ {
		for j := 0; j < 24; j++ {
			filter[i][j] = 1
		}
	}
}

func RefreshAnCoef() {
	GolobalAnCoe = arch.ReadAnCoe()
}

func RefreshAnValue() {
	GolobalAnVal = arch.ReadAnValue(GolobalAnCoe)
}

//用于卡尔曼滤波的估计误差
var filter [3][24]float32

func KalmanFilter() {
	for i := 0; i < 3; i++ {
		if arch.GetAnState(i) != true {
			continue
		}
		for j := 0; j < 24; j++ {
			var diff float32 = float32(GolobalAnBal[i][j] - GolobalAnVal[i][j])
			var siff float32 = float32(math.Pow(float64(diff), 2) + math.Pow(float64(filter[i][j]), 2))
			kg := siff / (siff + GolobalAnVal[i][j]*0.2)
			GolobalAnBal[i][j] -= float32(kg * diff)
		}
	}
	// fmt.Println(GolobalAnBal[0][4], GolobalAnVal[0][4])
}

func Test(a int) func(a int) int {
	con := 1
	proto := func(a int) int {
		con = con + 1
		return a + con
	}
	return proto
}

func Test2(send []byte) {
	send[2] = 1
}

//根据线路给出三相电流的值
func GetLine(line int) (a, b, c float32) {
	m := line / 4
	n := (line%4)*3 + 4
	return GolobalAnVal[m][n], GolobalAnVal[m][n+1], GolobalAnVal[m][n+2]
}

func BuildFormValue() []mtype.AnForm {
	formVal := make([]mtype.AnForm, 0, 64)
	for i := 0; i < len(ancfg); i++ {
		l := ancfg[i].Line
		formVal = append(formVal, mtype.AnForm{Id: int(ancfg[i].Id),
			Value: float32(GolobalAnBal[l/24][l%24] / ancfg[i].Form)})
	}
	return formVal
}
