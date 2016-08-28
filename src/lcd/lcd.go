package lcd

import (
	"arch"
	"fmt"
	"github.com/axgle/mahonia"
	"io"
	"os"
)

//液晶点坐标
type Point struct {
	X int32
	Y int32
}

//汉字字库
var HZbuf []byte

//液晶点阵
var LcdBuf [160 * 160]uint8

func readHZBuf(path string) {
	fi, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer fi.Close()

	HZbuf = make([]byte, 198576, 198576)
	num, _ := io.ReadFull(fi, HZbuf)

	if num != 198576 {
		fmt.Println("[错误]汉字库读取失败,文件大小有误.", num)
	}
}

func draw() {
	var buf [160][10]uint16

	//我在这里留下这么多角度,就是为了证明一次遍历就能完成矩阵的对换和形变
	//不要试图理解我的这段代码,在它面前颤抖吧!
	//[90度]
	// for i := 0; i < 160; i++ {
	// 	for j := 0; j < 160; j++ {
	// 		if LcdBuf[(i*10+j/16)*16+j%16] == 0xff {
	// 			buf[i][j/16] |= 1 << uint32(j%16)
	// 		}
	// 	}
	// }

	//[180度]
	// for i := 0; i < 160; i++ {
	// 	for j := 0; j < 160; j++ {
	// 		if LcdBuf[(i*10+j/16)*16+j%16] == 0xff {
	// 			buf[j][i/16] |= 1 << uint32(i%16)
	// 		}
	// 	}
	// }

	//[270度]
	// for i := 0; i < 160; i++ {
	// 	for j := 0; j < 160; j++ {
	// 		if LcdBuf[(i*10+j/16)*16+j%16] == 0xff {
	// 			buf[159-i][(159-j)/16] |= 1 << uint32((159-j)%16)
	// 		}
	// 	}
	// }

	//[360度]
	for i := 0; i < 160; i++ {
		for j := 0; j < 160; j++ {
			if LcdBuf[(i*10+j/16)*16+j%16] == 0xff {
				buf[j][(159-i)/16] |= 1 << uint32((159-i)%16)
			}
		}
	}

	for l := 0; l < 1600; l++ {
		arch.Ghead[0x460] = uint16(l)
		arch.Ghead[0x461] = buf[l/10][l%10]
	}
}

func clearRect(a, b, c, d int) {
	for j := a; j < b; j++ {
		for i := c; i < d; i++ {
			LcdBuf[j*160+i] = 0
		}
	}
}

func pixelColor(x int, y int, color uint16) {
	if color == 0 {
		LcdBuf[y*160+x] = 0x00
	} else {
		LcdBuf[y*160+x] = 0xff
	}
}

func lcdHLine(s_x, s_y, end int) {
	for i := s_x; i < end; i++ {
		pixelColor(i, s_y, 1)
	}
}

func textShow(str string, y int, x int, rev_flg bool) {
	//将传入的UTF-8字符转换为双字节的 GBK 字符.
	enc := mahonia.NewEncoder("GBK")
	str = enc.ConvertString(str)

	yp := y
	for i := 0; i < len(str); i++ {
		if str[i] < 0x80 {
			b1 := str[i]
			b1 -= 0x20
			rec_offset := uint32(b1) * 24
			HBuf := HZbuf[rec_offset : rec_offset+24]
			xp := x
			for m := 0; m < 12; m++ {
				for k := 0; k < 6; k++ {
					lx := yp + k
					ly := xp
					if rev_flg {
						pixelColor(lx, ly, ((uint16(HBuf[m*2])>>uint16(7-k))&0x01)^0x01)
					} else {
						pixelColor(lx, ly, (uint16(HBuf[m*2])>>uint16(7-k))&0x01)
					}
				}
				xp = (xp + 1) % 160
			}
			yp = (yp + 6) % 160
		} else {
			b1 := str[i]
			b2 := str[i+1]
			b1 -= 0xa0 //区码
			b2 -= 0xa0 //位码

			rec_offset := (94*(uint32(b1)-1) + (uint32(b2) - 1)) * 24 //- 0xb040;
			rec_offset += 96 * 24
			HBuf := HZbuf[rec_offset : rec_offset+24]
			xp := x
			for m := 0; m < 12; m++ {
				for k := 0; k < 8; k++ {
					lx := yp + k
					ly := xp
					if rev_flg {
						pixelColor(lx, ly, ((uint16(HBuf[m*2])>>uint16(7-k))&0x01)^0x01)
					} else {
						pixelColor(lx, ly, (uint16(HBuf[m*2])>>uint16(7-k))&0x01)
					}
				}
				for k := 0; k < 4; k++ {
					lx := yp + k + 8
					ly := xp
					if rev_flg {
						pixelColor(lx, ly, ((uint16(HBuf[m*2+1])>>uint16(7-k))&0x01)^0x01)
					} else {
						pixelColor(lx, ly, (uint16(HBuf[m*2+1])>>uint16(7-k))&0x01)
					}
				}
				xp = (xp + 1) % 160
			}
			yp = (yp + 12) % 160
			i += 1
		}
	}
	return
}
