/* test xbee packet decode
   test against known inputs (packets) 
	   and outputs (check sums, types, lengths, payloads)
*/
package xbee

import (
	"encoding/hex"
	"fmt"
	"testing"
)

const actualPackets = "7E000C830001240001060001F7000058" +
	"7E000C830001240001060001E9000066" +
	"7E000C8300012400010600021000003E" +
	"7E000C830001240001060001E000006F" +
	"7E000C830001240001060001EB000064" +
	"7E000C830001240001060001F400005B" +
	"7E000C830001240001060001E9000066" +
	"7E000C830001240001060001F9000056" +
	"7E000C830001240001060001E500006A" +
	"7E000C830001240001060001E500006A" +
	"7E000C830001240001060001F9000056" +
	"7E000C830001240001060001E500006A" +
	"7E000C830001240001060001F6000059" +
	"7E000C830001240001060001DD000072" +
	"7E000C830001240001060001F500005A" +
	"7E000C830001240001060001F6000059" +
	"7E000C8300012400010600023A000014" +
	"7E000C830001240001060001FF000050" +
	"7E000C830001240001060001FC000053" +
	"7E000C830001240001060001FA000055" +
	"7E000C830001240001060001F8000057" +
	"7E000C8300012400010600023300001B" +
	"7E000C830001240001060001F9000056" +
	"7E000C83000124000106000209000045" +
	"7E000C830001240001060001F500005A" +
	"7E000C830001240001060001E9000066" +
	"7E000C830001240001060001F000005F" +
	"7E000C830001240001060001EF000060" +
	"7E000C830001240001060001F100005E" +
	"7E000C830001240001060001E400006B" +
	"7E000C830001240001060001F500005A" +
	"7E000C83000124000106000205000049" +
	"7E000C8300012400010600020200004C" +
	"7E000C830001240001060001EA000065" +
	"7E000C830001240001060001FF000050" +
	"7E000C8300012300010600020000004F" +
	"7E000C830001240001060001E7000068" +
	"7E000C830001240001060001F100005E" +
	"7E000C8300012400010600020A000044" +
	"7E000C830001240001060001E500006A" +
	"7E000C830001240001060001F8000057" +
	"7E000C83000124000106000206000048" +
	"7E000C830001240001060001E500006A"
const fiveAnalogSamplesTwoChannels =
      "7E001C83000133000506000201000001F6000001FC0000020100000201000040"
const helloPacket = "7E000A010150010048656C6C6FB8"
const seventeenSamplesTwoChannels = "7E004C8300012800110600021A000001FF0000021C000002200000021A0000021B0000022C00000236000002140000021D00000232000002250000020B0000021F0000022A000002010000022800002A"
//caseJ2a with quantity field missing resetting xbee xmitter this packet not sent again
//const caseJ2a = "7E004C83000133000600024F018D024F01DA024F01F4025001FE024F0201024F0202024F0201024F0202024F0202024F0202024F0202024F0203024F0203024F0202024F0203024F0203024F02033B"
//caseJ2b with quantity field present 0x11
const caseJ2b = "7E004C8300013300110600024F018D024F01DA024F01F4025001FE024F0201024F0202024F0201024F0202024F0202024F0202024F0202024F0203024F0203024F0202024F0203024F0203024F02033B"
// 7E 004C 83 0001 33 00 0600 024F018D024F01DA024F01F4025001FE024F0201024F0202024F0201024F0202024F0202024F0202024F0202024F0203024F0203024F0202024F0203024F0203024F02033B"
//len 004C
//type 83
// source 0001
// rssi 33
// options 00
// quantity -- missing? should be 0x11, which happens to be xoff?
// channel map 0600

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
const xoffPacket = "7E0002237D3113CB"

var frametests = []struct {
	frame       string
	apiPacketID uint
	name string
}{
	{"7E00028A066F", MdmStatus,"bob"},             //simplist good packet
	{"00028A066F7E00028A066F", MdmStatus,"joe"},   //wait for start (garbage in front)
	{"7E00028A066F102347921834", MdmStatus,"larry"}, //simplist good packet with trailing garbage
	{"7E00028A066F7E00028A066F", MdmStatus,"jill"}, //2 packets in a row
	{helloPacket, TXreq16,"susan"},
	{escapedPacket, 0x23,"frank"}, //don't know what type 0x23 is,maybe bogus
	//{xonPacket, 0x23,"barrack"}, // are xon and xoff sent by xbee  for rx only?
	//{xoffPacket, 0x23,"francine"},
 	{actualPackets, Input16,"nilla"},
  {fiveAnalogSamplesTwoChannels,Input16,"oreo"},
  {seventeenSamplesTwoChannels,Input16,"olando"},
 // {caseJ2a,Input16,"caseJ2a"},
  {caseJ2b,Input16,"caseJ2b"},
	}

//7E    : API Frame
//00 02 : Length
//8A    : Modem status
//06    : Coordinator started
//6F    : checksum FF – ((8A +06) & FF) = 6F
//
func TestMedianInt(t *testing.T){
	var onlyOne = []int {1}
	if MedianInt(onlyOne) != 1 { t.Error("array length of one fails")}
	var onlyTwo = []int{2,1}
	if MedianInt(onlyTwo) != 1 { t.Error("array length of two fails")}
  var threeOdd = []int{1,3,2}
	if MedianInt(threeOdd) != 2 { t.Error("array length of three fails")}
  var four = []int{1,2,3,4}
	if MedianInt(four) != 2 { t.Error("array length of four fails")}
	var sortaLikeRealDataOdd = []int{503,501,600,503,400,250,400,250,250,251,600,503,400,250,400,250,250}
	if i:=MedianInt(sortaLikeRealDataOdd);i != 400 { t.Errorf("sortaLikeRealDataOdd fails %d",i)}
	var sortaLikeRealDataEven = []int{503,503,501,600,503,400,250,400,250,250,251,600,503,400,250,400,250,250}
	if i:=MedianInt(sortaLikeRealDataEven);i != 400 { t.Errorf("sortaLikeRealDataOdd fails %d",i)}

}
var packet []byte

func TestFrames(t *testing.T) {
	for _, f := range frametests {
		packet, err := hex.DecodeString(f.frame) // convert test data hex string to byte sequence
		if err != nil {
			panic("bad test data\n")
		}
		var apiframe APIframe
		apiframe.Reset()
		found := false
		for _, b := range []byte(packet) {
			if apiframe.Add_byte(b) {
				packettype, sourceAddress, rssi, options, quantity, analogChannels, measurements, e := apiframe.Parse()
				found = true
				apiframe.Reset() //discard frame once parsed
				if e == nil {
					if packettype == Input16 {
						fmt.Printf("type: %X sourceAddress %d rssi %d options %b,quantity %d,analogChannels %b measurements %d\n",
							packettype, sourceAddress, rssi, options, quantity, analogChannels, measurements)
							for i := uint(0);i<quantity;i++ {
						fmt.Printf("measurement[%d] %f\n",i, ((float32(measurements[i])*(1500.0/1023.0)-500)/10.0)*1.8+32.0)
					}

				}
				} else {
					fmt.Printf("packet parse failed %v %s\n", e,f.name)
				}
			}
		}
    if !found {
			t.Errorf("no packet found %s",f.name)
		}
	}
}
