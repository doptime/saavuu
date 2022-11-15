package logger

import (
	"log"
	"os"
)

var Lshortfile *log.Logger = log.New(os.Stdout, "", log.Lshortfile|log.Ldate|log.Ltime)
var Std = log.New(os.Stdout, "", log.Ldate|log.Ltime)
