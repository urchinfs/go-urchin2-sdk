package utils

import logging "github.com/ipfs/go-log"

func InitLog() {
	lvl, err := logging.LevelFromString("info")
	if err != nil {
		panic(err)
	}
	logging.SetAllLoggers(lvl)
}
