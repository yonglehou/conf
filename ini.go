/*
import (
	"github.com/berkaroad/conf"
)

func main(){
	config := conf.LoadInitConfig("~/config1.ini")
	config.Get("command", "concurrent_num")
	config.Set("command", "concurrent_num", "3")
	if err := config.Reload(); err == nil {

	}
}
*/
package conf

import (
	"bufio"
	"errors"
	"io"
	"os"
	"strconv"
	"strings"
)

type iniConfig struct {
	filepath string
	loaded   bool
	conflist []map[string]map[string]string
}

func LoadIniConfig(filepath string) Config {
	c := new(iniConfig)
	c.filepath = filepath
	if err := c.Reload(); err == nil {
		if DEBUG {
			consoleLog.Printf("[info] Load file \"%s\" success.\n", c.filepath)
		}
	} else {
		consoleLog.Panicf("[error] Load file \"%s\" error:%s!\n", c.filepath, err.Error())
	}
	return c
}

func (self *iniConfig) Get(section, name string) string {
	if self.loaded {
		data := self.conflist
		for _, v := range data {
			for key, value := range v {
				if key == section {
					if DEBUG {
						consoleLog.Printf("[info] Get config:%s.%s.\n", section, name)
					}
					return value[name]
				}
			}
		}
	}
	return ""
}

func (self *iniConfig) GetInt(section, name string) int {
	str := self.Get(section, name)
	if val, err := strconv.Atoi(str); err == nil {
		return val
	} else {
		return 0
	}
}

func (self *iniConfig) Set(section, name, value string) {
	if self.loaded {
		data := self.conflist
		var ok bool
		var index = make(map[int]bool)
		var conf = make(map[string]map[string]string)
		for i, v := range data {
			_, ok = v[section]
			index[i] = ok
		}

		i, ok := func(m map[int]bool) (i int, v bool) {
			for i, v := range m {
				if v == true {
					return i, true
				}
			}
			return 0, false
		}(index)

		if ok {
			self.conflist[i][section][name] = value
		} else {
			conf[section] = make(map[string]string)
			conf[section][name] = value
			self.conflist = append(self.conflist, conf)
		}
		if DEBUG {
			consoleLog.Printf("[info] Set config:%s.%s=%s.\n", section, name, value)
		}
	}
}

func (self *iniConfig) Reload() error {
	file, err := os.OpenFile(self.filepath, os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	var data map[string]map[string]string
	var section string
	buf := bufio.NewReader(file)
	for {
		l, err := buf.ReadString('\n')
		line := strings.TrimSpace(l)
		if err != nil {
			if err != io.EOF {
				return err
			}
			if len(line) == 0 {
				break
			}
		}
		switch {
		case len(line) == 0 || line[0] == '#':
		case line[0] == '[' && len(line) > 2 && line[len(line)-1] == ']':
			section = strings.TrimSpace(line[1 : len(line)-1])
			data = make(map[string]map[string]string)
			data[section] = make(map[string]string)
		default:
			i := strings.IndexAny(line, "=")
			if i < 0 {
				return errors.New("Cann't find \"=\" in single line.")
			} else if i == 0 {
				return errors.New("Cann't find name before \"=\" in single line.")
			}
			value := strings.TrimSpace(line[i+1 : len(line)])
			data[section][strings.TrimSpace(line[0:i])] = value
			if self.uniqueAppend(section) == true {
				self.conflist = append(self.conflist, data)
			}
		}
	}
	self.loaded = true
	return nil
}

func (self *iniConfig) uniqueAppend(conf string) bool {
	for _, v := range self.conflist {
		for k, _ := range v {
			if k == conf {
				return false
			}
		}
	}
	return true
}
