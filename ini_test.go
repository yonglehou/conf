package conf

import (
	"testing"
)

var config Config = LoadIniConfig("~/config1.ini")

func Test_Get(t *testing.T) {
	config.Get("command", "concurrent_num")
	config.GetInt("command", "concurrent_num")
}

func Test_Set(t *testing.T) {
	config.Set("command", "concurrent_num", "2")
}

func Test_Reload(t *testing.T) {
	config.Reload()
}
