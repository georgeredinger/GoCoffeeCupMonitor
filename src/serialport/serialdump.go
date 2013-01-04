package main

import "github.com/tarm/goserial"

import "fmt"
import "log"
import "time"

var offset int = 0

func init() {
	fmt.Printf("begin\n")
}

func main() {
	var n int
	fmt.Printf("open\n")
	c := &serial.Config{Name: "/dev/ttyUSB1", Baud: 9600}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}

	buf := make([]byte, 100)
	for {
		n, err = s.Read(buf)
		if err != nil {
			time.Sleep(time.Second / 100)
			continue
		}

		if n != 0 {
			if buf[0] == byte(0x7E) {
				fmt.Printf("\n")
			}
			fmt.Printf("%X", buf[:n])
		}
	}
	time.Sleep(1 * time.Second)
}
