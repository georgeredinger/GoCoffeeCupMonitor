package main

import (
	"fmt"
	"time"
)

func getdata() {
	for {
   h.broadcast <- fmt.Sprintf("%v",time.Now())
	 time.Sleep(10e9)
  }

}

