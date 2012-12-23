package main

import (
	"github.com/tarm/goserial"
	"log"
	"runtime"
	"time"
//  "encoding/binary"
//  "unsafe"
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
    b []byte
    //  startbyte byte
    //  packetLenHi byte
    //  packetLenLo byte
	  //	data   [12]byte
	  //	sum    byte
		err error
	}
  const xbLen = 16
  //rc := make(chan(xbeeRecord),1 )
	rc := make(chan *xbeeRecord )


  //read serial port and put bytes on channel 
	go func() {
		for {
			//xb := new(xbeeRecord)
      xb := new(xbeeRecord)
      var buf[] byte
			n, err := s.Read(buf)
      log.Print(n)
      if n != 0 {
			  rc <- xb
      }
      if err != nil {
        if n == 0 {
          continue
        }
        log.Print(err)
			}
		}
	}()

	for {
		timeout := time.NewTicker(10 * time.Second)
		defer timeout.Stop() //is this necessary?
		select {
		case got := <-rc:
			switch {
			case got.err != nil:
  			log.Fatal("  error:" + got.err.Error())
			default:
				log.Printf("%X",got.b)
			}
		case <-timeout.C:
      log.Print(".")
		}
		time.Sleep(1 * time.Second) //stutter the infinite loop.
	}

}
