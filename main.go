package main

import (
	"github.com/Fleurer/hardshard/proxy"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("main")

func main() {
	s, err := proxy.NewServer("0.0.0.0:4001")
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	log.Info("Listen 0.0.0.0:4001..")
	s.Run()
}
