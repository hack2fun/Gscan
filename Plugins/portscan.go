package Plugins

import (
	"../Misc"
	"../Parse"
	"fmt"
	"net"
	"sync"
	"time"
)

func PortScan(info Misc.HostInfo,ch chan int,wg *sync.WaitGroup){
	addr:=fmt.Sprintf("%s:%d",info.Host,info.Port)
	conn,err:=net.DialTimeout("tcp",addr,(time.Duration(info.Timeout))*time.Second)
	if err!=nil && info.ErrShow{
		info.PrintFailedPort()
	}else if err==nil{
		success+=1
		conn.Close()
		info.PrintSucceedPort()
		if info.Output!=""{
			info.OutputTXT()
		}
	}
	wg.Done()
	<- ch
}

func PortConn(info *Misc.HostInfo,ch chan int){
	var wg =sync.WaitGroup{}
	stime:=time.Now()
	ports,err:= Parse.ParsePort(info.Ports)
	Misc.CheckErr(err)

	hosts,err:=Parse.ParseIP(info.Host)
	Misc.CheckErr(err)
	wg.Add(len(hosts)*len(ports))
	Misc.InfoPrinter.Println("Total length",len(hosts)*len(ports))
	for _,host:=range hosts{
		for _,port:=range ports{
			info.Port=port
			info.Host = host
			go PortScan(*info,ch,&wg)
			ch<-1
		}
	}
	wg.Wait()
	end:=time.Since(stime)
	Misc.InfoPrinter.Println("All Done")
	Misc.InfoPrinter.Println("Number of successes:",success)
	Misc.InfoPrinter.Println("Time consumed:",end)
}
