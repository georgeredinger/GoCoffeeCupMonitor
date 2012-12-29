package xbee

import "fmt"
import "bytes"

/*
0x7E  Start frame delimiter.

0x7D  Escape control character. Indicates that next byte is escaped.

0x11 0x13  These bytes are software flow control characters.
*/

const (
	START_BYTE  = byte(0x7E)
	ESCAPE_BYTE = byte(0x7D)
	XON_BYTE    = byte(0x11)
	XOFF_BYTE   = byte(0x13)
)

var ESCAPE_BYTES = []byte{START_BYTE, ESCAPE_BYTE, XON_BYTE, XOFF_BYTE}

const (
	waitingForStart    = 0
	waitingForLengthHi = 1
	waitingForLengthLo = 2
	waitingForData     = 3
	waitingForCheckSum = 4
	waitingForEscape   = 5
	done               = 6
)

type APIframe struct {
	frame     []byte
	lengthHi  uint
	lengthLo  uint
	length    uint
	started   bool
	bytesLeft uint
	checkSum  uint8
	state     uint8
}

func isEsc(b byte) bool {
	return (bytes.IndexByte(ESCAPE_BYTES, b) != -1)
}

func (f *APIframe) init() {
		  f.frame = f.frame[:0]	
			f.lengthHi = 0
			f.lengthLo = 0
			f.length   = 0
			f.started  = false
			f.bytesLeft= 0
			f.checkSum  = 0
}

func (f *APIframe) checksum() bool { return false }

func (f *APIframe) add_byte(b uint8) bool {

	switch b {
	case XON_BYTE, XOFF_BYTE:
		if f.state == waitingForEscape {
			break
		} else {
			return false //ignore xon xoff
		}
	case ESCAPE_BYTE:
		if f.state == waitingForEscape {
			break
		} else {
			f.state = waitingForEscape
			return false
		}
	}

	switch f.state {

	case waitingForStart:
		if b == START_BYTE {
			f.state = waitingForLengthHi
			if f.bytesLeft != 0 {
				panic("new packet inside packet length must have been wrong")
			}
			return false
		}
		return false
	case waitingForLengthHi:
		f.state = waitingForLengthLo
		f.lengthHi = uint(b)
		return false
	case waitingForLengthLo:
		f.state = waitingForData
		f.lengthLo = uint(b)
		f.length = f.lengthHi*255 + f.lengthLo
		f.bytesLeft = f.length
		return false
	case waitingForEscape:
		b ^= byte(0x20)
		fallthrough
	case waitingForData:
		if f.length > uint(len(f.frame)) {
			f.frame = append(f.frame, b)
			f.bytesLeft--
		}
		if f.bytesLeft == 0 {
			f.state = waitingForCheckSum
		}
		return false
	case waitingForCheckSum:
		var sum uint8
		for _, c := range f.frame {
			sum += uint8(c)
		}
		f.checkSum = 0XFF - sum
		if b == f.checkSum {
			f.state = done
			return true
		} else {
			f.state = waitingForStart
			fmt.Printf("Checksum %X != %X\n", f.checkSum, b)
			fmt.Printf("packet: %X,%X\n", f.frame, b)
			f.init()
			return false // 
		}
	case done:
		return true
	default:
		panic("default frame state")
	}
	panic("impossible")
	return false
}

func (f APIframe) remaining_bytes() uint { return f.bytesLeft }

func (f APIframe) parse() bool {
	var good bool

	if f.frame[0] == byte(0x7E) {
		good = true
	} else {
		good = false
	}
	return good
}
