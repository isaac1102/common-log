package main

import "common-log/log"

func main() {
	log.Trace("this is trace test")
	log.Info("this is trace test")
	log.Debug("this is debug test")
	log.Error("this is error test")
	log.Fatal("this is fatal test")
}
