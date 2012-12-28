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
									{"7E00028A066F",0x6f},
									{"7E00028A066F",0x6f},
									{"7E00028A066F",0x6f},
									{"7E00028A066F",0x6f},
						    }



//7E    : API Frame
//00 02 : Length
//8A    : Modem status
//06    : Coordinator started
//6F    : checksum FF – ((8A +06) & FF) = 6F
//

var packet [256]byte

func TestFrames(t *testing.T) {
	for i,f := range(frametests) {
		packet,err := hex.DecodeString(f.frame)
		fmt.Printf("len %v %X %v %X\n",i,packet,err,f.sum)
		  if err != nil {
				fmt.Print ("ooops\n")
			}
			fmt.Printf("first byte is %X\n",packet[0])

		}
}


