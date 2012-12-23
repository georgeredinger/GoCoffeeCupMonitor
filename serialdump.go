package main

import "github.com/bsiegert/goserial"
import "fmt"
import "log"
import "time"


var	offset int  = 0
func main() {
	var n int
	defer close()
	fmt.Printf("open\n")
	c := &serial.Config{Name: "/dev/ttyUSB0", Baud: 9600}

	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("begin\n")

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
		fmt.Printf("n=%d\n", n)
		if err != nil {
			log.Print(err)
		}
		log.Print("%q", buf[:n])
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

func close() {
	fmt.Printf("closing\n")
}
