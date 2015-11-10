package conf

import (
	"log"
	"os"
)

var consoleLog = log.New(os.Stdout, "[conf] ", log.LstdFlags)

type Config interface {
	Get(section, name string) string
	GetInt(section, name string) int
	Set(section, name, value string)
	Reload() error
}
