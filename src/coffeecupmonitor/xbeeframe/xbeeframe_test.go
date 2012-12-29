package xbeeframe

import (
	"testing"
	"fmt"
  "encoding/hex"
)

/*
0x7E  Start frame delimiter.

0x7D  Escape control character. Indicates that next byte is escaped.

0x11 0x13  These bytes are software flow control characters.
*/



var frametests = [] struct {
	frame string
	sum uint8
}{
	{"7E00028A066F",0x6f},
	{"00028A066F7E00028A066F",0x6f},
//	{"7E00028A066F7E00028A066F",0x6f},
//	{"028A066F7E00028A066F",0x6f},
//	{"066FA7E00028A066F",0x6f},
}



//7E    : API Frame
//00 02 : Length
//8A    : Modem status
//06    : Coordinator started
//6F    : checksum FF – ((8A +06) & FF) = 6F
//


var packet [256]byte

func TestFrames(t *testing.T) {
	for _,f := range(frametests) {
	  packet,err := hex.DecodeString(f.frame) // convert test data hex string to byte sequence, 2 characters per byte
		if err != nil {
			fmt.Print("oops\n")
			}
	  var apiframe APIframe 
		apiframe.init()
		for _,b := range([]byte(packet)) {
			apiframe.add_byte(b)
		}
//		res:=apiframe.parse()
//		 fmt.Printf("result %v\n",res)

	}
	//		fmt.Printf("len %v %X %v %X\n",i,packet,err,f.sum)

}



