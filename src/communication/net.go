package communication

import (
	"analog"
	"arch"
	"config"
	"fmt"
	"io"
	"net"
	"os"
	"remote"
	"time"
)

func server() {
	listener, err := net.Listen("tcp", "0.0.0.0:2404")
	if err != nil {
		fmt.Println("[错误]网络监听失败", err)
		os.Exit(1)
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		fmt.Println("[信息]已连接：", conn.LocalAddr(), conn.RemoteAddr())
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	var str104 IEC104Struct
	str104.conn = conn
	str104.MasterAddr = config.GetIECMasterAddr()
	str104.Addr = config.GetIECAddr()
	str104.TranCauseL = config.GetTranCauseL()
	str104.InfoAddrL = config.GetInfoAddrL()
	str104.CommonAddrL = config.GetCommonAddrL()

	for {
		conn.SetDeadline(time.Now().Add(2e8))
		buf := make([]byte, 1024)
		length, err := conn.Read(buf)
		if err == io.EOF {
			break
		}
		if length == 0 {
			continue
		}

		fmt.Println("[信息]接收到信息:长度", length, "内容", buf[:length], err)
		res := str104.lib104Verify(buf, length)

		if IEC101FRAMEERR == res {
			continue
		}
		if IEC101FRAME_U == res {
			str104.lib104Cmd(str104.lib104TypeU(buf))
		}
		if IEC101FRAME_I == res {
			kind := str104.lib104TypeI(buf, length)
			switch byte(kind) {
			case IEC101_C_IC_NA_1:
				str104.IEC104_AP_RQALL_CON()
				str104.IEC104_AP_SP_NA(remote.BuildSingleSoe())
				str104.IEC104_AP_DP_NA(remote.BuildDoubleSoe())
				str104.IEC104_AP_ME_ND(analog.BuildFormValue())
				str104.IEC104_AP_RQALL_END()
			case IEC101_C_CS_NA_1:
				arch.SetFpgaTime(str104.IEC104_AP_CS(buf))
			case IEC101_C_DC_NA_1:

			}
		}
	}
}

func Start() {
	go server()
}
