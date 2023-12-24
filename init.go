package main

import (
	"io"
	"log"
	"os"
	"time"
)

func init() {
	// log handler to file and console out
	if err := os.MkdirAll("./logs", os.ModePerm); err != nil {
		panic(err)
	}
	logFile, err := os.OpenFile(time.Now().Format("./logs/2006-01-02.log"), os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile) // | log.Lmicroseconds  )
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
}
