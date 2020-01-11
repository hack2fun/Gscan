package Plugins

import (
	"../Misc"
	"../Parse"
	"fmt"
	"github.com/jlaffaye/ftp"
	"sync"
	"time"
	//"os"
)

func Ftp(info Misc.HostInfo,ch chan int,wg *sync.WaitGroup){
	var err error
	addr := fmt.Sprintf("%s:%d",info.Host,info.Port)
	client, err := ftp.Dial(addr, ftp.DialWithTimeout(time.Duration(info.Timeout)*time.Second))
	if err!=nil && info.ErrShow{
		Misc.ErrPrinter.Println(err.Error())
	}else if err==nil{
		err = client.Login(info.Username, info.Password)
		if  err != nil && info.ErrShow{
			info.PrintFail()
		}else if err == nil{
			client.Quit()
			client.Logout()
			success+=1
			info.PrintSuccess()
			if info.Output!=""{
				info.OutputTXT()
			}
		}
	}
	wg.Done()
	<-ch
}

func FtpConn(info *Misc.HostInfo,ch chan int){
	var hosts,usernames,passwords []string
	var err error
	var wg = sync.WaitGroup{}
	stime:=time.Now()
	if info.Ports == ""{
		info.Port=FTPPORT
	}else{
		p,_:=Parse.ParsePort(info.Ports)
		info.Port = p[0]
	}

	hosts,err= Parse.ParseIP(info.Host)
	Misc.CheckErr(err)

	usernames,err= Parse.ParseUser(info)
	Misc.CheckErr(err)

	passwords,err= Parse.ParsePass(info)
	Misc.CheckErr(err)

	wg.Add(len(hosts)*len(usernames)*len(passwords))
	Misc.InfoPrinter.Println("Total length",len(hosts)*len(usernames)*len(passwords))
	for _,host:=range hosts{
		for _,user:=range usernames{
			for _,pass:=range passwords{
				info.Host = host
				info.Username = user
				info.Password = pass
				go Ftp(*info,ch,&wg)
				ch <- 1
			}
		}
	}
	wg.Wait()
	end:=time.Since(stime)
	Misc.InfoPrinter.Println("All Done")
	Misc.InfoPrinter.Println("Number of successes:",success)
	Misc.InfoPrinter.Println("Time consumed:",end)
}