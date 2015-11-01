package conf

import (
	"github.com/berkaroad/util"
	"log"
	"os"
	"strconv"
)

var consoleLog = log.New(os.Stdout, "[conf] ", log.LstdFlags)

type Config interface {
	Get(section, name string) string
	Set(section, name, value string)
	Load() error
}

var config = LoadIniConfig(util.GetExecFilePath() + ".ini")

func Get(section, name string) string {
	return config.Get(section, name)
}

func GetInt(section, name string) int {
	str := Get(section, name)
	if val, err := strconv.Atoi(str); err == nil {
		return val
	} else {
		return 0
	}
}

func Set(section, name, value string) {
	config.Set(section, name, value)
}

func Load() error {
	return config.Load()
}
