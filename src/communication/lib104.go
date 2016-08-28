package communication

import (
	"fmt"
	"mtype"
	"net"
	"time"
)

type IEC104Struct struct {
	usRecvNum                uint16 //已经接收到的帧
	usSendNum                uint16 //已经发送出的帧
	usServRecvNum            uint16 //服务器接收到的帧
	usServSendNum            uint16 //服务器已发送的帧
	ucSendCountOverturn_Flag uint8  //发送计数翻转标志
	ucRecvCountOverturn_Flag uint8  //接收计数翻转标志
	usAckNum                 uint16 //已经认可的帧
	ucTimeOut_t0             uint8  //连接建立超时值,单位s，默认30s
	ucTimeOut_t1             uint8  //APDU的发送或测试的超时时间,默认：15s
	ucTimeOut_t2             uint8  //无数据报文t2t1情况下发送S-帧的超时时间,默认：10s
	ucTimeOut_t3             uint8  //长时间空闲状态下发送测试帧的超时   默认:20s
	k                        uint8  //发送I格式应用规约数据单元的未认可帧数
	w                        uint8  //接收I格式应用规约数据单元的帧数
	ucMax_k                  uint8  //发送状态变量的最大不同的接收序号
	ucMax_w                  uint8  //接收w个I格式APDUs之后的最后的认可
	ucDataSend_Flag          uint8  //是否允许发送标志
	ucIdleCount_Flag         uint8  //是否允许t2计数
	ucIdleCount              uint8  //t2计数
	ucWaitServConCount_Flag  uint8  //是否t1计数
	ucWaitServConCount       uint8  //t1计数
	ucNoRecvCount            uint8  //t0计数
	ucNoRecvT3               uint8  //t3计时器
	ucYK_Limit_Time          uint8
	ucYK_TimeCount_Flag      uint8
	ucYK_TimeCount           uint8
	conn                     net.Conn //发送套接字
	MasterAddr               uint16
	Addr                     uint16
	TranCauseL               uint16
	InfoAddrL                uint16
	CommonAddrL              uint16
}

const IEC101FRAMEOK uint16 = 0
const IEC101FRAMEERR uint16 = 1
const IEC101FRAMESEQERR uint16 = 2
const IECFRAMENEEDACK uint16 = 3
const IEC101FRAMELOGIN uint16 = 4
const IEC101FRAMERESET uint16 = 5
const IEC101FRAMECHERR uint16 = 6

const IEC101FRAME_I uint16 = 0x10
const IEC101FRAME_S uint16 = 0x11
const IEC101FRAME_U uint16 = 0x13

//ASDU的类型标识
const IEC101_M_SP_NA_1 byte = 1
const IEC101_M_DP_NA_1 byte = 3
const IEC101_M_ME_NA_1 byte = 9
const IEC101_M_ME_NB_1 byte = 11
const IEC101_M_ME_NC_1 byte = 11
const IEC101_M_IT_NA_1 byte = 15
const IEC101_M_ME_ND_1 byte = 21
const IEC101_M_SP_TB_1 byte = 30
const IEC101_M_DP_TB_1 byte = 31
const IEC101_M_ST_TB_1 byte = 32
const IEC101_M_ME_TD_1 byte = 34
const IEC101_M_ME_TE_1 byte = 35
const IEC101_M_ME_TF_1 byte = 36
const IEC101_M_IT_TB_1 byte = 37
const IEC101_M_EP_TD_1 byte = 38
const IEC101_M_EP_TE_1 byte = 39
const IEC101_M_EP_TF_1 byte = 40
const IEC101_C_DC_NA_1 byte = 45
const IEC101_C_RC_NA_1 byte = 47
const IEC101_C_SC_TA_1 byte = 58
const IEC101_M_EI_NA_1 byte = 70
const IEC101_C_IC_NA_1 byte = 100
const IEC101_C_CI_NA_1 byte = 101
const IEC101_C_RD byte = 102
const IEC101_C_CS_NA_1 byte = 103
const IEC101_C_RP_NA_1 byte = 105

//ASDU的传送原因
const IEC101_CAUSE_N_BIT byte = 0x40
const IEC101_CAUSE_SPONT byte = 3            //突发（自发）
const IEC101_CAUSE_INIT byte = 4             //请求或者被请求
const IEC101_CAUSE_REQ byte = 5              //请求或者被请求
const IEC101_CAUSE_ACT byte = 6              //激活
const IEC101_CAUSE_ACTCON byte = 7           //激活确认
const IEC101_CAUSE_DEACT byte = 8            //停止激活
const IEC101_CAUSE_DEACTCON byte = 9         //停止激活确认
const IEC101_CAUSE_ACTTERM byte = 10         //激活终止
const IEC101_CAUSE_INTROGEN byte = 20        //响应站召唤
const IEC101_CAUSE_COUNTGEN byte = 37        //响应计数量站总召唤
const IEC101_CAUSE_UNKNOWNTYPE byte = 44     //未知的类型标识
const IEC101_CAUSE_UNKNOWNCAUSE byte = 45    //未知的传送原因
const IEC101_CAUSE_UNKNOWNCOMMADDR byte = 46 //未知的应用服务数据单元公共地址
const IEC101_CAUSE_UNKNOWNINFOADDR byte = 47 //未知的信息对象地址

func (this *IEC104Struct) lib104TypeU(dealBuf []byte) byte {
	return dealBuf[2]
}

func (this *IEC104Struct) lib104TypeI(dealBuf []byte, bufLen int) uint16 {
	ServSendNum := uint16(dealBuf[2])*256 + uint16(dealBuf[3])
	ServRecvNum := uint16(dealBuf[4])*256 + uint16(dealBuf[5])

	if !(ServSendNum&0x0001 == 0x0001) {
		ServRecvNum >>= 1
		this.usServRecvNum = ServRecvNum

		if ServSendNum == 0x7fff {
			this.usRecvNum = 0
		} else {
			this.usRecvNum += 1
		}
		return uint16(dealBuf[6])
	} else {
		return IEC101FRAMEERR
	}
}

func (this *IEC104Struct) lib104Verify(dealBuf []byte, bufLen int) uint16 {
	//检查长度
	if (bufLen > 255) || (bufLen < 6) {
		return IEC101FRAMEERR
	}
	//检查开头
	if dealBuf[0] != 0x68 || uint8(dealBuf[1]) != uint8(bufLen-2) {
		return IEC101FRAMEERR
	}

	logMsg := []byte{0x68, 0x04, 0x07, 0x00, 0x00, 0x00}
	if string(dealBuf[:6]) == string(logMsg) || bufLen == 6 {
		return IEC101FRAME_U
	} else {
		return IEC101FRAME_I
	}
}
func (this *IEC104Struct) lib104Cmd(kind byte) {
	buf := []byte{0x68, 0x04, 0x00, 0x00, 0x00, 0x00}

	switch kind {
	case 0x07:
		buf[2] = 0x0b
	case 0x10:
		buf[2] = 0x23
	case 0x43:
		buf[2] = 0x83
	}
	this.conn.Write(buf)
}

func (this *IEC104Struct) IEC104_AP_RQALL_CON() {
	sendBuf := make([]byte, 0, 256)
	sendBuf = append(sendBuf, IEC101_C_IC_NA_1)
	sendBuf = append(sendBuf, 0x01)
	sendBuf = append(sendBuf, IEC101_CAUSE_ACTCON)
	sendBuf = append(sendBuf, byte(this.MasterAddr))
	sendBuf = append(sendBuf, byte(this.Addr&0xff))
	sendBuf = append(sendBuf, byte((this.Addr&0xff00)>>8))
	sendBuf = append(sendBuf, 0x00)
	sendBuf = append(sendBuf, 0x00)
	sendBuf = append(sendBuf, 0x00)
	sendBuf = append(sendBuf, IEC101_CAUSE_INTROGEN)

	this.lib104Send(sendBuf, len(sendBuf))
}

func (this *IEC104Struct) IEC104_AP_RQALL_END() {
	sendBuf := make([]byte, 0, 256)
	sendBuf = append(sendBuf, IEC101_C_IC_NA_1)
	sendBuf = append(sendBuf, 0x01)
	sendBuf = append(sendBuf, IEC101_CAUSE_ACTTERM)
	sendBuf = append(sendBuf, byte(this.MasterAddr))
	sendBuf = append(sendBuf, byte(this.Addr&0xff))
	sendBuf = append(sendBuf, byte((this.Addr&0xff00)>>8))
	sendBuf = append(sendBuf, 0x00)
	sendBuf = append(sendBuf, 0x00)
	sendBuf = append(sendBuf, 0x00)
	sendBuf = append(sendBuf, IEC101_CAUSE_INTROGEN)

	this.lib104Send(sendBuf, len(sendBuf))
}

func (this *IEC104Struct) lib104Send(buf []byte, length int) {
	sendBuf := make([]byte, 0, 256)

	sendBuf = append(sendBuf, 0x68)
	sendBuf = append(sendBuf, byte(length+4))

	sendnum := this.usSendNum << 1
	sendBuf = append(sendBuf, byte(sendnum&0xff))
	sendBuf = append(sendBuf, byte((sendnum>>8)&0xff))

	recvnum := this.usRecvNum << 1
	sendBuf = append(sendBuf, byte(recvnum&0xff))
	sendBuf = append(sendBuf, byte((recvnum>>8)&0xff))

	sendBuf = append(sendBuf, buf...)

	// sendBuf[1] = byte(len(sendBuf))
	this.conn.Write(sendBuf)

	fmt.Println("[消息]发送", sendBuf)

	//判断发送序列是否需要反转
	if this.usSendNum == 0x7fff {
		this.usSendNum = 0
		this.ucSendCountOverturn_Flag = 0x01
	} else {
		this.usSendNum++
	}
	this.k++
	this.w = 0
}

//不带时标的单点遥信，响应总召测，类型标识1
func (this *IEC104Struct) IEC104_AP_SP_NA(soe []mtype.SingleSoe) {
	sendBuf := make([]byte, 0, 256)
	sendBuf = append(sendBuf, IEC101_M_SP_NA_1)
	sendBuf = append(sendBuf, 0x00)
	sendBuf = append(sendBuf, IEC101_CAUSE_INTROGEN)

	sendBuf = append(sendBuf, byte(this.MasterAddr))
	sendBuf = append(sendBuf, byte(this.Addr&0xff))
	sendBuf = append(sendBuf, byte((this.Addr>>8)&0xff))

	for i := 0; i < len(soe); i++ {
		sendBuf = append(sendBuf, byte(soe[i].Id&0xff))
		sendBuf = append(sendBuf, byte((soe[i].Id>>8)&0xff))
		sendBuf = append(sendBuf, 0x00)
		sendBuf = append(sendBuf, byte(soe[i].State))

		if (i+1)%25 == 0 {
			sendBuf[1] = byte(i)
			this.lib104Send(sendBuf, len(sendBuf))
			sendBuf = sendBuf[:6]
		}
	}
	sendBuf[1] = byte(len(soe) % 25)
	this.lib104Send(sendBuf, len(sendBuf))
}

//不带时标的双点遥信，响应总召测，类型标识2
func (this *IEC104Struct) IEC104_AP_DP_NA(soe []mtype.DoubleSoe) {
	sendBuf := make([]byte, 0, 256)
	sendBuf = append(sendBuf, IEC101_M_SP_NA_1)
	sendBuf = append(sendBuf, 0x00)
	sendBuf = append(sendBuf, IEC101_CAUSE_INTROGEN)

	sendBuf = append(sendBuf, byte(this.MasterAddr))
	sendBuf = append(sendBuf, byte(this.Addr&0xff))
	sendBuf = append(sendBuf, byte((this.Addr>>8)&0xff))

	for i := 0; i < len(soe); i++ {
		sendBuf = append(sendBuf, byte(soe[i].Id&0xff))
		sendBuf = append(sendBuf, byte((soe[i].Id>>8)&0xff))
		sendBuf = append(sendBuf, 0x00)
		sendBuf = append(sendBuf, byte((soe[i].State1<<1)+soe[i].State2))

		if (i+1)%25 == 0 {
			sendBuf[1] = byte(i)
			this.lib104Send(sendBuf, len(sendBuf))
			sendBuf = sendBuf[:6]
		}
	}
	sendBuf[1] = byte(len(soe) % 25)
	this.lib104Send(sendBuf, len(sendBuf))
}

func (this *IEC104Struct) IEC104_AP_ME_ND(formVal []mtype.AnForm) {
	sendBuf := make([]byte, 0, 256)
	sendBuf = append(sendBuf, IEC101_M_ME_ND_1)
	sendBuf = append(sendBuf, 0x00)
	sendBuf = append(sendBuf, IEC101_CAUSE_SPONT)

	sendBuf = append(sendBuf, byte(this.MasterAddr))
	sendBuf = append(sendBuf, byte(this.Addr&0xff))
	sendBuf = append(sendBuf, byte((this.Addr>>8)&0xff))

	for i := 0; i < len(formVal); i++ {
		sendBuf = append(sendBuf, byte(formVal[i].Id&0xff))
		sendBuf = append(sendBuf, byte((formVal[i].Id>>8)&0xff))
		sendBuf = append(sendBuf, 0x00)
		sendBuf = append(sendBuf, byte(uint32(formVal[i].Value)&0xff))
		sendBuf = append(sendBuf, byte((uint32(formVal[i].Value)>>8)&0xff))

		if (i+1)%25 == 0 {
			sendBuf[1] = byte(i)
			this.lib104Send(sendBuf, len(sendBuf))
			sendBuf = sendBuf[:6]
		}
	}
	sendBuf[1] = byte(len(formVal) % 25)
	this.lib104Send(sendBuf, len(sendBuf))
}

//对时，类型标识103
func (this *IEC104Struct) IEC104_AP_CS(buf []byte) int64 {
	usMSec := uint16(buf[16])*256 + uint16(buf[15])
	ucSec := usMSec / 1000
	ucMin := buf[17] & 0x3f
	ucHour := buf[18] & 0x1f
	ucDay := buf[19] & 0x1f
	ucMonth := buf[20] & 0x0f
	ucYear := uint16(buf[21] & 0x7f)

	if ucSec > 59 || ucMin > 59 || ucHour > 23 || ucYear > 99 || ucMonth > 12 {
		return 0
	}

	if (ucMonth == 0) || (ucDay == 0) {
		return 0
	}
	if ucMonth == 2 {
		if ((ucYear+2000)%4 == 0 && (ucYear+2000)%100 != 0) || (ucYear+2000)%400 == 0 {
			if ucDay > 29 {
				return 0
			}
		} else {
			if ucDay > 28 {
				return 0
			}
		}
	}
	if ucMonth <= 7 && ucMonth%2 == 0 && ucDay > 30 {
		return 0
	}
	if ucMonth >= 8 && ucMonth%2 == 1 && ucDay > 30 {
		return 0
	}

	unix_time := time.Date(int(ucYear), time.Month(int(ucMonth)), int(ucDay), int(ucHour), int(ucMin), int(ucSec), 0, time.Local).Unix()

	sendBuf := make([]byte, 0, 256)
	sendBuf = append(sendBuf, buf[6])
	sendBuf = append(sendBuf, 0x01)
	sendBuf = append(sendBuf, IEC101_CAUSE_ACTCON)

	sendBuf = append(sendBuf, byte(this.MasterAddr))
	sendBuf = append(sendBuf, byte(this.Addr&0xff))
	sendBuf = append(sendBuf, byte((this.Addr>>8)&0xff))
	sendBuf = append(sendBuf, buf[12:22]...)

	this.lib104Send(sendBuf, len(sendBuf))
	return unix_time
}
