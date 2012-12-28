package xbeeframe

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

type APIFrame struct {

}

func checksum(f *APIFrame) bool { return false}
func add_byte( f *APIFrame) bool {
	return false

}
func remaining_bytes(f *APIFrame) int { return 0}
func parse(f *APIFrame) * map[string] string {
      var m = map[string] string {
        "demo0": "1234.56",
        "demo1": "6543.21" }
      return &m
}


