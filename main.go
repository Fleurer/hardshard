package main

import (
	"log"

	"github.com/Fleurer/hardshard/proxy"
)

func main() {
	s, err := proxy.NewServer("0.0.0.0:4001")
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	log.Print("Listen 0.0.0.0:4001..")
	s.Run()
}
