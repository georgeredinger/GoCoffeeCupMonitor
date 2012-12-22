package main

import (
	"github.com/tarm/goserial"
	"log"
	"runtime"
	"time"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	log.Print("Open")
	c := &serial.Config{Name: "/dev/ttyUSB0", Baud: 9600}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}

	type xbeeRecord struct {
		b   []byte
		err error
	}
	rc := make(chan *xbeeRecord)

	go func() {
		for {
			buf := make([]byte, 128)
			n, err := s.Read(buf)
			rc <- &xbeeRecord{buf[:n], nil}
			if err != nil {
        if n == 0 {
          continue
        }
        log.Print(err)
				//return
			}
		}
	}()

	for {
		timeout := time.NewTicker(1 * time.Second)
		defer timeout.Stop() //is this necessary?
		select {
		case got := <-rc:
			switch {
			case got.err != nil:
				//Catching an EOF error here can indicate the port was disconnected.
				// -- if using a USB to serial port, and the device is unplugged 
				//    while being read, we'll receive an EOF.
				log.Fatal("  error:" + got.err.Error())
			default:
				log.Print(got.b)
			}
		case <-timeout.C:
      log.Print(".")
		}
		time.Sleep(1 * time.Second) //stutter the infinite loop.
	}

}
