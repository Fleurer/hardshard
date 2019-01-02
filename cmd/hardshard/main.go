package main

import (
	"github.com/Fleurer/hardshard/pkg/proxy"
	"github.com/siddontang/go-log/log"
)

func main() {
	s, err := proxy.NewServer("0.0.0.0:4001")
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	log.Info("Listen 0.0.0.0:4001..")
	s.Run()
}
