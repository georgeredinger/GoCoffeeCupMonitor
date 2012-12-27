package main

import "github.com/tarm/goserial"
import "fmt"
import "log"
import "time"
import "os"
import "./xbeeframe"

var offset int = 0
var packet[16] byte
func  packetdecode( packet [] byte) {

}

func  packetadd(b byte) {

		if offset == 0 && b != 0x7e {
      fmt.Printf("[%x]",b)
      os.Stdout.Sync()
			return
		}

    packet[offset] = b
		if offset == 15 {
			fmt.Printf("%X",packet[:offset])
			fmt.Printf("\n")

			fmt.Printf("len = %d \n",packet[1]*255+packet[2])
			fmt.Printf("API = %d \n",packet[3])
			fmt.Printf("frameid = %d \n",packet[4]*255+packet[5])
			fmt.Printf("destaddr = %d \n",packet[6]*255+packet[7])
			fmt.Printf("opt = %d \n",packet[8])
			fmt.Printf("ADC0 = %d \n",packet[9]*255+packet[10])
			fmt.Printf("ADC1 = %d \n",packet[11]*255+packet[12])
			sum := 0
			for i := 3;i<16;i++{
         fmt.Printf("+%d",int(packet[i]))
				 sum += int(packet[i])
			 }
       fmt.Print("=\n")
			 fmt.Printf("sum = %d\n",sum & 0x00ff)
			offset = -1
		}
		offset++
  }

func main() {
	fmt.Printf("begin\n")
	rc := make(chan byte)
	fmt.Printf("open\n")
	c := &serial.Config{Name: "/dev/ttyUSB0", Baud: 9600}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}


	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		log.Fatal(err)
	}

	go func() {
    fmt.Print("Begin read go func\n")
		buf := make([]byte, 16)
		for {
			n, err := s.Read(buf)
			if err != nil {
				time.Sleep(time.Second / 100)
				continue
			}

			if n != 0 {
  			rc <- buf[0]
			}
		}
	}()

  frame = new(APIFrame)

	for {
		timeout := time.NewTicker(10 * time.Second)
		defer timeout.Stop()
		select {
		case got := <-rc:
			switch {
			default:
				packetadd(got)
				frame.add_byte(got)
			}
		case <-timeout.C:
			log.Print(".")
		}
		time.Sleep(1 * time.Second/1000)
	}

}
