package configparse

import (
	"github.com/Unknwon/goconfig"
	"strconv"
)

type Conf struct {
	items map[string]map[string]string
}

var GlobalConf Conf

var cfg *goconfig.ConfigFile

func (c *Conf) CfgInit(filename string) {
	c.items = make(map[string]map[string]string)
	secvalue := make(map[string]string)
	cfg, err := goconfig.LoadConfigFile(filename)
	if err != nil {
		panic("load config file failed " + filename)
	}
	cfgseclist := cfg.GetSectionList()
	for _, v := range cfgseclist {
		keys := cfg.GetKeyList(v)
		for _, b := range keys {
			secvalue[b], err = cfg.GetValue(v, b)
			if err != nil {
				panic("get key failed")
			}
			c.items[v] = secvalue
		}
	}
}

func (c *Conf) GetStr(section string, seckey string) string {
	return c.items[section][seckey]
}

func (c *Conf) GetBool(section string, seckey string) bool {
	return c.items[section][seckey] == "true" || c.items[section][seckey] == "1"
}

func (c *Conf) GetInt(section string, seckey string) int {
	intvalue, _ := strconv.Atoi(c.items[section][seckey])
	return intvalue
}
