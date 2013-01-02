/*
read and decode xbee series 1 packets 
one byte at a time as they come in
  from: http://www.jsjf.demon.co.uk/xbee/xbee.pdf
  Input Line States how to decode 

	API id 0x82 or 0x83

	This packet (0x82 or 0x83) is used when a remote XBee and a base XBee have been conﬁgured 
  so that the remote will sample its inputs at
  set intervals, and transmit the results to the base. If the base XBee is conﬁgured
  to pass such data out through the UART,
  the base’s host will receive this packet.
   1. Byte: packet type id 0x82 for 64-bit source address, or 0x83 for 16-bit source address.
  2. Bytes: source address. 
     Two bytes for 16-bit source addressing, or eight for 64-bit source addressing.
   3. Byte: RSSI value.
   4. Byte: options.
     If bit 1 is set, this is an address broadcast.
     If bit 2 is set, it is a PAN broadcast. All other bits are   reserved.
  5. Byte: sample quantity. This is the number of full sets of samples in what follows.
  6. Word: 2-byte channel indicator msb ﬁrst.
     Bits 14–9 are a 6-bit mask, with 1 for each ADC channel in AD5–AD0
     respectively that will be reported.
     Bits 8–0 are for the digital lines D8–D0, showing which of them will be included
     in the values. Bit 15 is not used.
  7. Word: optional 16-bit bitﬁeld, with bits corresponding to lines D8–D0 
     as in the channel indicator. Where bits were
     set in the channel indicator, the corresponding bits here show the 
     state of the input. These bytes are not present if
     there are no lines enabled as digital inputs.
  8. Words: if any bits in the channel indicator were set among the A-D inputs, 
     those readings now follow. Each is a
     16-bit value with the A-D reading in the low-order 10 bits.
     A-D readings are given in order from AD0 to AD5.

*/
/*
0x7E  Start frame delimiter.
0x7D  Escape control character. Indicates that next byte is escaped.
0x11 0x13  These bytes are software flow control characters.
*/
// a packet with 2 analog measurements
//
//                  options
//                   ^
//                   | qual
//                   |  |        A0   A1
//7E 000C 83 0001 24 00 01 0601 01E9 0000 66
//   |     |   |  |         |             sum
//   |     |   |  |       channel
//   |_len |   |  |
//    type_|   |  |
//             |  |
//     source__|  |
//                |
//                |_rssi
//

package xbee

import (
	"errors"
	"sort"
	"os"
	"log"
)

const (
	//tx Packets
	TXreq64      = 0x00 // TX request with 64-bit destination address 5-4
	TXreq16      = 0x01 // TX request with 16-bit destination address 5-4
	LocalAT      = 0x08 // Local AT command, immediate action 5-2
	LocalATqueed = 0x09 // Local AT command, queued action 5-2
	RemoteAT64   = 0x17 // Remote AT command with 64-bit destination address 5-3
	//rx Packets 
	RX64       = 0x80 // RX with 64-bit source address 5-5
	RX16       = 0x81 // RX with 16-bit source address 5-5
	Input64    = 0x82 // Input line states with 64-bit source address 5-4
	Input16    = 0x83 // Input line states with 16-bit source address 5-4
	LocalATres = 0x88 // Local AT response 5-2
	TXres      = 0x89 // TX response 5-5
	MdmStatus  = 0x8a // Modem status packet 5-1
	RmtATres   = 0x97 // Remote AT response 5-3
)

// escapes
const (
	START_BYTE  = byte(0x7E)
	ESCAPE_BYTE = byte(0x7D)
	XON_BYTE    = byte(0x11)
	XOFF_BYTE   = byte(0x13)
)

// packet rx decoder states -- read_byte()
const (
	waitingForStart = iota
	waitingForLengthHi
	waitingForLengthLo
	waitingForData
	waitingForCheckSum
	waitingForEscape
	done
)

type APIframe struct {
	frame     []byte
	apiID     byte
	lengthHi  uint
	lengthLo  uint
	length    uint
	started   bool
	bytesLeft uint
	checkSum  uint8
	state     uint8
}
func init() {
	f, _ := os.Open("logfile.log")
	log.SetOutput(f)
}

func MedianInt(a [] int) int {
	var median int

  sort.Ints(a)	

	if i:=len(a); i % 2 == 0 {
      median = (a[i/2-1] + a[i/2])/2 
		} else {
      median = a[i/2] 
		}
		return median
}

func (f *APIframe) Reset() {
	f.frame = f.frame[:0]
	f.lengthHi = 0
	f.lengthLo = 0
	f.length = 0
	f.started = false
	f.bytesLeft = 0
	f.checkSum = 0
	f.state = waitingForStart
}

func bitCount(x uint) (n uint) {
	// n accumulates the total bits set in x, counting only set bits
	for ; x > 0; n++ {
		// clear the least significant bit set
		x &= x - 1
	}
	return n
}

func (f *APIframe) checksum() bool { return false }

func (f *APIframe) Add_byte(b uint8) bool {

	switch b {
		// maybe xon/xoff are used but scaped
//	case XON_BYTE, XOFF_BYTE:
//		if f.state == waitingForEscape {
//			break
//		} else {
//			return false //ignore xon xoff
//		}
	case ESCAPE_BYTE:
		if f.state == waitingForEscape {
			break
		} else {
			log.Printf("escape\n")
			//Todo: check sum fails on these packets
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
			log.Printf("Checksum %X != %X\n", f.checkSum, b)
			log.Printf("packet: %X,%X\n", f.frame, b)
			f.Reset()
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

func (f APIframe) Parse() (apiID uint, sourceAddress uint, rssi uint,
	options uint, quantity uint, analogChannels uint, analogMeasurements []int, e error) {
	apiID = uint(f.frame[0])
	if apiID != Input16 {
		return apiID, sourceAddress, rssi, options, quantity, analogChannels, analogMeasurements, e
	}
	if len(f.frame) > 3 {
		sourceAddress = uint(f.frame[1])<<8 + uint(f.frame[2])
	} else {
		e = errors.New("packet too short (sourceAddress")
		sourceAddress = 12345678
	}
	if len(f.frame) > 4 {
		rssi = uint(f.frame[3])
	} else {
		e = errors.New("packet too short (rssi)")
		rssi = 8321
	}
	if len(f.frame) > 5 {
		options = uint(f.frame[4])
	} else {
		e = errors.New("packet too short (options)")
		options = 8321
	}
	if len(f.frame) > 6 {
		quantity = uint(f.frame[5])
	} else {
		e = errors.New("packet too short (quantity)")
		quantity = 8321
	}
	if len(f.frame) > 8 {
		analogChannels = uint(uint(f.frame[6])<<8 + uint(f.frame[7]))
	} else {
		e = errors.New("packet too short (analogChanels)")
		analogChannels = 8321
	}

	channelCount := bitCount(analogChannels)
	analogMeasurements = make([]int, channelCount*quantity)

	if len(f.frame) > 9 {
		for q := uint(0); q < quantity;q++ {
			for i := uint(0); i < channelCount; i++ {
				measurement := int(int(f.frame[(8+i*2)+q*channelCount*2])<<8 + int(f.frame[(9+i*2)+q*channelCount*2]))
				analogMeasurements[i+q] = measurement
	  	}
	}
	}
	return apiID, sourceAddress, rssi, options, quantity, analogChannels, analogMeasurements, e
}
