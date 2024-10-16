package main

import (
	"log"
	"time"
)

func main() {
	log.Println(time.Now())
	log.Println(time.Now())

	loc, _ := time.LoadLocation("Asia/Jakarta")
	log.Println(time.Now().In(loc))
}
