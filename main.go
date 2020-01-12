package main

import (
	"./Misc"
	"./Parse"
	"./Plugins"
	"flag"
)

func main() {
	flag.Parse()
	var params *Misc.HostInfo
	allparam := Parse.Param
	if allparam.INIFile != "" {
		ini, err := Parse.GetConfig(allparam.INIFile)
		Misc.CheckErr(err)
		params = ini
	} else {
		params = allparam
	}

	var ch = make(chan int, params.Thread)
	switch {
	case params.Help || flag.NFlag() == 0:
		flag.Usage()
	case params.Show:
		Plugins.Show()
	default:
		Plugins.Selector(params, ch)
	}
}
