package xbee

import (
	"encoding/hex"
//	"fmt"
	"testing"
)

/*
0x7E  Start frame delimiter.

0x7D  Escape control character. Indicates that next byte is escaped.

0x11 0x13  These bytes are software flow control characters.
*/

const helloPacket = "7E000A010150010048656C6C6FB8"

// helloPacket
// 7E // Start delimiter
// 00 0A // Length bytes
// 01 // API identifier
// 01 // API frame ID
// 50 01 // Destination address low
// 00 // Option byte
// 48 65 6C 6C 6F // Data packet(ASCII: "Hello")
// B8 // Checksum

const escapedPacket = "7E0002237D31CB"

// from page 58, http://ftp1.digi.com/support/documentation/90000982_B.pdf
// raw frame         0x7E 0x00 0x02 0x23 0x11 0xCB
//		 0x11 needs to be escaped which results in the following frame: 
// escaped frame		 0x7E 0x00 0x02 0x23 0x7D 0x31 0xCB
// Note: In the above example, the length of the raw data (excluding the checksum)
//   is 0x0002 and 
//   the checksum of the non-escaped data 
//   (excluding frame delimiter and length) is calculated as:
//   0xFF - (0x23 + 0x11) = (0xFF - 0x34) = 0xCB.
//  To test data integrity, a checksum is calculated and verified on non-escaped data.
//  To calculate: Not including frame delimiters and length, add all bytes
//    keeping only the lowest 8 
//  bits of the result and subtract from 0xFF.
//  To verify: Add all bytes (include checksum, but not the delimiter and length).
//  If the checksum is 
//    correct, the sum will equal 0xFF.
//

const xonPacket = "7E001102237D31CB"
const xoffPacket ="7E0002237D3113CB"

var frametests = []struct {
	frame string
	apiPacketID int
}{
	{"7E00028A066F", MdmStatus},             //simplist good packet
	{"00028A066F7E00028A066F", MdmStatus},   //wait for start (garbage in front)
	{"7E00028A066F102347921834", MdmStatus}, //simplist good packet with trailing garbage
	{"7E00028A066F7E00028A066F", MdmStatus}, //2 packets in a row
	{helloPacket, TXreq16},
	{escapedPacket, 0x23}, //don't know what type 0x23 is,maybe bogus
	{xonPacket, 0x23},
	{xoffPacket, 0x23},
}

//7E    : API Frame
//00 02 : Length
//8A    : Modem status
//06    : Coordinator started
//6F    : checksum FF – ((8A +06) & FF) = 6F
//

var packet [6]byte

func TestFrames(t *testing.T) {
	for _, f := range frametests {
		packet, err := hex.DecodeString(f.frame) // convert test data hex string to byte sequence
		if err != nil {
			panic("bad test data\n")
		}
		var apiframe APIframe
		apiframe.init()
		for _, b := range []byte(packet) {
			if apiframe.add_byte(b) {
				break
			}
		}
		packettype:=apiframe.parse()
		if packettype != f.apiPacketID {
			t.Errorf("packet %X, %X  not type %X\n",apiframe.frame,apiframe.frame[0],f.apiPacketID)
    }
	}
	//		fmt.Printf("len %v %X %v %X\n",i,packet,err,f.sum)
}
