package Plugins

import (
	"../Misc"
	"../Parse"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"sync"
	"time"
)

func WebScan(info Misc.HostInfo,ch chan int,wg *sync.WaitGroup,client *http.Client){
	res,err:=http.NewRequest("GET",info.Url,nil)
	if err!=nil && info.ErrShow{
		Misc.ErrPrinter.Println(err)
	}else if err==nil{
		res.Header.Add("Connection","keep-alive")
		res.Header.Add("User-agent","Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/28.0.1468.0 Safari/537.36")
		res.Header.Add("Accept","*/*")
		res.Header.Add("Accept-Language","en-us")
		res.Header.Add("Accept-Encoding","identity")
		if info.Cookie!=""{
			res.Header.Add("Cookie",info.Cookie)
		}
		if info.Header!=""{
			var header = make(map[string]string)
			err:=json.Unmarshal([]byte(info.Header),&header)
			if err!=nil{
				Misc.CheckErr(err)
			}
			for k,v:=range header{
				res.Header.Add(k,v)
			}
		}
		resp,err:=client.Do(res)
		if err!=nil && info.ErrShow{
			Misc.ErrPrinter.Println(err)
		}else if err==nil{
			UrlPrint(resp.StatusCode,info.Url,info.ErrShow,info)
			ioutil.ReadAll(resp.Body)
			resp.Body.Close()
		}
	}
	wg.Done()
	<-ch
}

func UrlConn(info *Misc.HostInfo,ch chan int){
	stime:=time.Now()
	var wg = sync.WaitGroup{}
	url,err:=Parse.ParseUrl(info.Url)
	Misc.CheckErr(err)

	var client = &http.Client{
		Transport:&http.Transport{
			DialContext:(&net.Dialer{
				Timeout:time.Second*time.Duration(info.Timeout),
			}).DialContext,
		},
		CheckRedirect:func(req *http.Request, via []*http.Request) error{
			return http.ErrUseLastResponse
		},
	}
	if info.UrlFile==""{
		info.UrlFile=URLFILE
	}
	pathes,err:=Parse.Readfile(info.UrlFile)
	Misc.CheckErr(err)
	wg.Add(len(pathes))
	Misc.InfoPrinter.Println("Total length",len(pathes))
	for _,path:=range pathes{
		info.Url = url+path
		go WebScan(*info,ch,&wg,client)
		ch<-1
	}
	wg.Wait()
	end:=time.Since(stime)
	Misc.InfoPrinter.Println("All Done")
	Misc.InfoPrinter.Println("Number of successes:",success)
	Misc.InfoPrinter.Println("Time consumed:",end)
}

func UrlPrint(code int,url string,ErrShow bool,info Misc.HostInfo){
	if code==200{
		urlpath:=fmt.Sprintf("%s\x1B[1;40;32m[%d]\x1B[0m",url,code)
		fmt.Println(urlpath)
		if info.Output!=""{
			info.OutputTXT()
		}
	}else if code!=404 && code!=200{
		urlpath:=fmt.Sprintf("%s\x1B[1;40;33m[%d]\x1B[0m",url,code)
		fmt.Println(urlpath)
		if info.Output!=""{
			info.OutputTXT()
		}
	}else if ErrShow && code==404{
		urlpath:=fmt.Sprintf("%s\x1B[0;40;31m[%d]\x1B[0m",url,code)
		fmt.Println(urlpath)
	}
}