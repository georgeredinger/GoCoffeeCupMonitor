package xbee

import "fmt"
import "bytes"

/*
0x7E  Start frame delimiter.

0x7D  Escape control character. Indicates that next byte is escaped.

0x11 0x13  These bytes are software flow control characters.
*/
//Xbee API ids: (RX packets only
const (
		XBEE_IO64 = byte(0x82) // expect this
		XBEE_IO16 = byte(0x83) // or this
		//other API ids should throw an error
)
//  from: http://www.jsjf.demon.co.uk/xbee/xbee.pdf
//  Input Line States how to decode API id 0x82 or 0x83
//  
//  This packet is used when a remote XBee and a base XBee have been conﬁgured 
//  so that the remote will sample its inputs at
//  set intervals, and transmit the results to the base. If the base XBee is conﬁgured
//  to pass such data out through the UART,
//  the base’s host will receive this packet.
//  1. Byte: packet type id 0x82 for 64-bit source address, or 0x83 for 16-bit source address.
//  2. Bytes: source address. 
//     Two bytes for 16-bit source addressing, or eight for 64-bit source addressing.
//  3. Byte: RSSI value.
//  4. Byte: options.
//     If bit 1 is set, this is an address broadcast.
//     If bit 2 is set, it is a PAN broadcast. All other bits are   reserved.
//  5. Byte: sample quantity. This is the number of full sets of samples in what follows.
//  6. Word: 2-byte channel indicator msb ﬁrst.
//     Bits 14–9 are a 6-bit mask, with 1 for each ADC channel in AD5–AD0
//     respectively that will be reported.
//     Bits 8–0 are for the digital lines D8–D0, showing which of them will be included
//     in the values. Bit 15 is not used.
//  7. Word: optional 16-bit bitﬁeld, with bits corresponding to lines D8–D0 
//     as in the channel indicator. Where bits were
//     set in the channel indicator, the corresponding bits here show the 
//     state of the input. These bytes are not present if
//     there are no lines enabled as digital inputs.
//  8. Words: if any bits in the channel indicator were set among the A-D inputs, 
//     those readings now follow. Each is a
//     16-bit value with the A-D reading in the low-order 10 bits.
//     A-D readings are given in order from AD0 to AD5.
//  


// escapes
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
