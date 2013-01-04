package main

import "fmt"

var datastore []string
const datastoresize = 1000

func addstring(s string) {
	if len(datastore) >= datastoresize {
		datastore = datastore[0:(len(datastore)-1)]
	}
	datastore = append(datastore, fmt.Sprintf("%d: %s",len(datastore),s))
}

func dumpdatabase(c  chan string)  {
	for _,s := range(datastore) {
   c <- s 
 }
}

