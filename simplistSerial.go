package main

import (
  "github.com/tarm/goserial"
  "log"
  "time"
)

func main() {

  log.Print("Open")
  c := &serial.Config{Name: "/dev/ttyUSB0", Baud: 9600}
  s, err := serial.OpenPort(c)
  if err != nil {
    log.Fatal(err)
  }

  for {
    var buf[] byte
    n, err := s.Read(buf)
    if err != nil {
      if n == 0 {
        time.Sleep(1 * time.Second)
        continue
      } else {
        log.Print(n)
      }
      log.Printf("%X",buf)
    }
  }
}
