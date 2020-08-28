package utils

import (
	"log"
	"os"
)

var DefaultLog = new(log.Logger)

func init() {
	DefaultLog.SetOutput(os.Stdout)
}
