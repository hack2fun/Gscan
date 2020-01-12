package Parse

import (
	"../Misc"
	"errors"
	"github.com/go-ini/ini"
)

var ParseConfigErr = errors.New("An error occurred while parsing the configuration file, please check if your ini file format is correct or the file exists")

var CONFIG = &Misc.HostInfo{}

//Getting information of config file
func GetConfig(filename string) (*Misc.HostInfo, error) {
	conf, err := ini.Load(filename)
	if err != nil {
		return nil, ParseConfigErr
	}
	err = conf.Section("CONFIG").MapTo(CONFIG)
	if err != nil {
		return nil, ParseConfigErr
	}
	return CONFIG, nil
}
