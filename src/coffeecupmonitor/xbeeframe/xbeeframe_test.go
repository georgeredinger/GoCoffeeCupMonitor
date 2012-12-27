package xbeeframe

import (
  "testing"
	"fmt"
	"strconv"
)
/*
0x7E  Start frame delimiter.

0x7D  Escape control character. Indicates that next byte is escaped.

0x11 0x13  These bytes are software flow control characters.
*/




var frametests = [] struct {
	               frame  string
								 sum uint8
							  }{
									{"7E000x028A066F",0x6f},
									{"7E000x028A066F",0x6f},
									{"7E000x028A066F",0x6f},
									{"7E000x028A066F",0x6f},
									{"7E000x028A066F",0x6f},
						    }



//7E    : API Frame
//00 02 : Length
//8A    : Modem status
//06    : Coordinator started
//6F    : checksum FF – ((8A +06) & FF) = 6F
//

func TestFrames(t *testing.T) {
	bytes,err := strconv.ParseUint(frametests[0].frame[0:2],16,8)
	if err != nil {
		fmt.Printf("what %X",err)
	}
	fmt.Printf("frame[0] = %X\n",bytes)
	if bytes != 0x7E {
    t.Errorf("frame must start with 0x7f 0x%X",uint8(frametests[0].frame[0]))
  }
}

