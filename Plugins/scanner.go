package Plugins

import (
	"../Misc"
	"errors"
	"os"
	"reflect"
)

//Reference：https://mikespook.com/2012/07/%E5%9C%A8-golang-%E4%B8%AD%E7%94%A8%E5%90%8D%E5%AD%97%E8%B0%83%E7%94%A8%E5%87%BD%E6%95%B0/
func Call_user_func(m map[string]interface{}, name string, params ...interface{}) (result []reflect.Value, err error) {
	f := reflect.ValueOf(m[name])
	if len(params) != f.Type().NumIn() {
		err = errors.New("The number of params is not adapted.")
		Misc.CheckErr(err)
	}
	in := make([]reflect.Value, len(params))
	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}
	result = f.Call(in)
	return result,nil
}

//根据-m参数选择要使用的功能
func Selector(info *Misc.HostInfo,ch chan int){
	_,ok:=PluginList[info.Scantype]
	if !ok{
		Misc.ErrPrinter.Println("The specified scan type does not exist, please use the -show parameter to view all supported scan types")
		os.Exit(0)
	}
	Call_user_func(PluginList,info.Scantype,info,ch)
}
