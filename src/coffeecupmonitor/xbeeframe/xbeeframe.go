package xbeeframe
//import "fmt"
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
	length int
}

func checksum(f *APIframe) bool { return false}
func (f *APIframe) add_byte(b byte ) bool {
	//fmt.Printf(" add_byte: %2X\n",b)
	f.frame = append(f.frame,b)
	return true;
}
func (f APIframe) remaining_bytes() int { return 0}

func (f APIframe) parse()  bool {
  var good bool

	if f.frame[0] == byte(0x7E) {
    good = true
	}else{
		good = false
	}
	return good
}


