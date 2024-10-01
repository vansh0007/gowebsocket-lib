package logger

import (
	"flag"
	"fmt"
	"go/build"
	"log"
	"os"
)

var (
	Log *log.Logger
)

func init() {
	// set location of log file

	var logpath = "/Users/gunvanshsiingh/GolandProjects/gowebsocket-lib/info.log"

	fmt.Println(build.Default.GOPATH)

	flag.Parse()
	var file, err1 = os.Create(logpath)

	if err1 != nil {
		panic(err1)
	}
	Log = log.New(file, "", log.LstdFlags|log.Lshortfile)
	Log.Println("LogFile : " + logpath)
}
