package main

import "github.com/ziutek/serial"
import "fmt"
import "log"
import "time"


var	offset int  = 0
func main() {
	var n int

	s, err := serial.Open("/dev/ttyUSB0")
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("begin\n")
	buf := make([]byte, 16)
	for {
		n, err = s.Read(buf)
		if err != nil {
			continue
		}
		if n != 0 {
		}

			if buf[0] == byte(0x7E) {
				fmt.Printf("\n")
			}

			fmt.Printf("%X", buf[:n])
		
		time.Sleep(1 * time.Second)
	}
}

