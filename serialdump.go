package main

import "github.com/bsiegert/goserial"
import "fmt"
import "log"

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
	buf := make([]byte, 128)
	for {
		n, err = s.Read(buf)
		fmt.Printf("n=%d\n", n)
		if err != nil {
			log.Print(err)
		}
		log.Print("%q", buf[:n])
	}
}

func close() {
	fmt.Printf("closing\n")
}
