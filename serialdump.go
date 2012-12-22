package main
import "github.com/ziutek/serial"
import "fmt"
import "log"

func main() {
	var n int
  fmt.Printf("open\n")
	s, err := serial.Open("/dev/ttyUSB0")
	if err != nil {
		log.Fatal(err)
	}
  fmt.Printf("begin\n")
	buf := make([]byte, 128)
	for {
		n, err = s.Read(buf)
		if err != nil {
			fmt.Printf("%v\n",err)
		}
		fmt.Printf("%q", buf[:n])
	}
}

