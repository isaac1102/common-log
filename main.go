package main

import "common-log/log"

var logger = log.GetLogger()

func main() {
	logger.Error().Msg("this is log")
}
