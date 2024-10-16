package main

import (
	"fmt"
	"time"
)

func main() {
	t := time.Now()
	zone, offset := t.Zone()
	fmt.Println(zone, offset)

}
