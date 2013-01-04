package main

import (
  "fmt"
  "time"
  "github.com/tarm/goserial"
  "xbee"
  "log"
  "os"
  "math"
)

const (
  A = 8.271111E-4 // 108 probe betatherm 100K6A1B
  B = 2.088020E-4
  C = 8.059200E-8
)
var offset int = 0

func init() {
  f, e := os.OpenFile("monitor.log", os.O_WRONLY|os.O_CREATE| os.O_APPEND, 0640)
  if e != nil {
    log.Fatalln(e)
  }
  log.SetOutput(f)
  log.Printf("monitor start %v",time.Now())

}
func getdata() {
 var n int
  var apiframe xbee.APIframe

  c := &serial.Config{Name: "/dev/ttyUSB0", Baud: 9600}
  s, err := serial.OpenPort(c)
  if err != nil {
		for {
			h.broadcast <- fmt.Sprintf("%v",time.Now())
			s := fmt.Sprintf("%v",time.Now())
			addstring(s) 
			h.broadcast <- s
			time.Sleep(10e9)
		}
		log.Fatal(err)
  }

  buf := make([]byte, 256)
  for {
    n, err = s.Read(buf)
    if err != nil {
      time.Sleep(time.Second / 100)
      continue
    }
    for i := 0; i < n; i++ {
      if apiframe.Add_byte(buf[i]) {
        packettype, _, _, _, quantity, _, measurements, e := apiframe.Parse() // one index out of range thrown to here

        apiframe.Reset() //discard frame once parsed
        // buf = buf[:0]
        if e == nil {
          if packettype == xbee.Input16 {
            //fmt.Printf("raw: %d\n",measurements[0:quantity])
            median := xbee.MedianInt(measurements[0:quantity])
            t:=time.Now()
            mv := ((float64(median)*(2500.0/1023.0))) //~100k ohm thermistor
            r	 := 98000.0*(2500.0/mv -1)
            tk := 1/(A+B*math.Log(r) + C*(math.Pow(math.Log(r),3)))
            f :=(tk-273.15)*1.8+32
            s := fmt.Sprintf("%d %3.1f\n",t.UnixNano()/1e9 ,f)
            addstring(s)
            h.broadcast <- s
          }
        } else {
          fmt.Printf("packet parse failed %v\n", e)
        }
      }
    }
  }
  time.Sleep(1 * time.Second)
}



