package Plugins

import (
	"../Misc"
	"../Parse"
	"github.com/stacktitan/smb/smb"
	"sync"
	"time"
	"context"
)

func Smb(info Misc.HostInfo,ch chan int,wg *sync.WaitGroup){
	tch:=make(chan int)
	defer func(){
		<-ch
		wg.Done()
	}()
	ctx,cancel:=context.WithTimeout(context.Background(),time.Second*time.Duration(info.Timeout))
	defer cancel()
	options := smb.Options{
		Host:        info.Host,
		Port:        info.Port,
		User:        info.Username,
		Domain:      "",
		Workstation: "",
		Password:    info.Password,
	}
	go func(){
		session, err := smb.NewSession(options, false)
		if err!=nil && info.ErrShow {
			info.PrintFail()
		}else if session.IsAuthenticated{
			success+=1
			session.Close()
			info.PrintSuccess()
			if info.Output!=""{
				info.OutputTXT()
			}
		}
		<-tch
	}()
	select{
	case <-ctx.Done():
		if info.ErrShow{
			Misc.ErrPrinter.Println(info.Url,"Time out")
		}
		return
	case <-tch:
		return
	}
}

func SmbConn(info *Misc.HostInfo,ch chan int){
	var hosts,usernames,passwords []string
	var err error
	var wg = sync.WaitGroup{}
	stime:=time.Now()
	if info.Ports == ""{
		info.Port=SMBPORT
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
				go Smb(*info,ch,&wg)
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