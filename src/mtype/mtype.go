package mtype

import "time"

//遥信配置文件结构
type ReConfig struct {
	Id   int32
	Type int32
	L0   int32
	L1   int32
}

//遥测配置文件结构
type AnConfig struct {
	Id   int32
	Swit bool
	Line int32
	Form float32
	Larg float32
	Thre float32
	Up   float32
	Down float32
}

//保护配置文件
type ProConfig struct {
	Id   int32
	S1   int32
	V1   int32
	T1   int32
	S2   int32
	V2   int32
	T2   int32
	S3   int32
	V3   int32
	T3   int32
	Sr   int32
	Tr   int32
	Sf   int32
	Tf   int32
	Ctrl int32
	How  int32
}

type ProFilter struct {
	T1 time.Time
	T2 time.Time
	T3 time.Time
	Tr time.Time
	Tf time.Time
}

type SingleSoe struct {
	Id    int
	State int
}

type DoubleSoe struct {
	Id     int
	State1 int
	State2 int
}

type AnForm struct {
	Id    int
	Value float32
}
