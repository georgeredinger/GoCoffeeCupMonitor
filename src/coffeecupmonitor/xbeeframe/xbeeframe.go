package xbeeframe
import "fmt"
/*
0x7E  Start frame delimiter.

0x7D  Escape control character. Indicates that next byte is escaped.

0x11 0x13  These bytes are software flow control characters.
*/

const (
	START_BYTE  = 0x7E
	ESCAPE_BYTE = 0x7D
	XON_BYTE    = 0x11
	XOFF_BYTE   = 0x13
)

var ESCAPE_BYTES = []int{START_BYTE, ESCAPE_BYTE, XON_BYTE, XOFF_BYTE}

type APIframe struct {
	frame []byte
	lengthHi uint
	lengthLo uint
	length uint
	started bool
	bytesLeft uint
	waitingForLengthLo bool
	waitingForLengthHi bool
	waitingForCheckSum bool
	checkSum uint8
	state uint8
//	const (
//		started,
//		waitingForLengthLo,
//		waitingForlenghtHi,
//	)
}

func (f *APIframe) init() {
	f.started = false
	f.waitingForCheckSum = false
	f.waitingForLengthLo  = false
	f.waitingForLengthHi = false
	
}
func (f *APIframe) checksum()  bool { return false}

func (f *APIframe) add_byte(b uint8 ) bool {
	
	if f.waitingForLengthHi{
		f.waitingForLengthHi = false
		f.lengthHi = uint(b)
		return false 
	} else {
		if f.waitingForLengthLo{
			f.waitingForLengthLo = false
			f.lengthLo = uint(b)
			f.length = f.lengthHi*255+f.lengthLo 	
			f.bytesLeft = f.length
			return false
		}
	}

	if b == byte(0x7E) {
		f.started = true
		f.waitingForCheckSum = false
		f.waitingForLengthLo,f.waitingForLengthHi = true,true
		if f.bytesLeft != 0{
			panic("new packet inside packet")
		}
		return false
	}

	if f.started {
		if f.bytesLeft > 0  {
			f.frame = append(f.frame,b)
			f.bytesLeft--
			return false 
		} else {
			fmt.Printf("packet body %X\n",f.frame)
			f.waitingForCheckSum = true
			var sum uint8 
			for _,c:= range(f.frame) {
				sum += uint8(c)
			}	
			f.checkSum = 0XFF - sum
			if(b == f.checkSum){
				fmt.Printf("Checksum %X == %X\n",f.checkSum,b)
				return true
			} else {
				fmt.Printf("Checksum %X != %X\n",f.checkSum,b)
				panic("checksum failure")
			}
		}
	}
return true // how do I get here?
}

func (f APIframe) remaining_bytes() uint { return f.bytesLeft}


func (f APIframe) parse()  bool {
	var good bool

	if f.frame[0] == byte(0x7E) {
		good = true
	}else{
		good = false
	}
	return good
}


