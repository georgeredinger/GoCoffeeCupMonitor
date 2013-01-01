package main

import ( 
 "github.com/tarm/goserial"
 "xbee"
 "fmt"
 "log"
 "time"
)

var offset int = 0

func init() {

}


func main() {
	var n int
	var apiframe xbee.APIframe

	c := &serial.Config{Name: "/dev/ttyUSB0", Baud: 9600}
	s, err := serial.OpenPort(c)
	if err != nil {
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
				packettype, _, _, _, quantity, _, measurements, e := apiframe.Parse()
				apiframe.Reset() //discard frame once parsed
				// buf = buf[:0]
				if e == nil {
					if packettype == xbee.Input16 {
          //  sum := uint(0)
					//	for i := uint(0); i < quantity; i++ {
					//		sum += measurements[i]
					//	}
          //  mean := sum / quantity
					median := xbee.MedianInt(measurements[0:quantity])
						t:=time.Now()
						temp := ((float32(median)*(1500.0/1023.0)-500)/10.0)*1.8+32.0
						fmt.Printf("%d %3.1f\n",t.UnixNano()/1e9 ,temp)
					}
				} else {
					fmt.Printf("packet parse failed %v\n", e)
				}
			}
		}
	}
	time.Sleep(1 * time.Second)
}

