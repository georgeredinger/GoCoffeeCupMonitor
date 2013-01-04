package main

import (
	"fmt"
	"time"
)

func getdata() {
	for {
	 s := fmt.Sprintf("%v",time.Now())
   addstring(s) 
   h.broadcast <- s
	 time.Sleep(10e9)
  }

}

