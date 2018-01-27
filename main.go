package main

import (
	"log"

	"github.com/Fleurer/hardshard/proxy"
)

func main() {
	s := proxy.NewServer("0.0.0.0", 4001)
	log.Print("Listen 0.0.0.0:4001..")
	s.Run()
}
